package observer

import (
	"bytes"
	"encoding/hex"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/wire"
	"github.com/ontio/btcrelayer/log"
	sdk "github.com/ontio/multi-chain-go-sdk"
	"time"
)

type BtcObConfig struct {
	NetType            string `json:"net_type"`
	BtcObLoopWaitTime  int64  `json:"btc_ob_loop_wait_time"`
	BtcObConfirmations uint32 `json:"btc_ob_confirmations"`
	BtcJsonRpcAddress  string `json:"btc_json_rpc_address"`
	User               string `json:"user"`
	Pwd                string `json:"pwd"`
	SleepTime int `json:"sleep_time"`
}

type BtcObserver struct {
	cli      *RestCli
	NetParam *chaincfg.Params
	conf     *BtcObConfig
}

func NewBtcObserver(conf *BtcObConfig, cli *RestCli) *BtcObserver {
	var param *chaincfg.Params
	switch conf.NetType {
	case "test":
		param = &chaincfg.TestNet3Params
	case "sim":
		param = &chaincfg.SimNetParams
	case "regtest":
		param = &chaincfg.RegressionNetParams
	default:
		param = &chaincfg.MainNetParams
	}
	var observer BtcObserver
	observer.cli = cli
	observer.NetParam = param
	observer.conf = conf

	return &observer
}

func (observer *BtcObserver) Listen(relaying chan *CrossChainItem) {
	top := btcCheckPoints[observer.NetParam.Name].Height
	log.Infof("[BtcObserver] get start height %d from checkpoint, check once %d seconds", top, observer.conf.BtcObLoopWaitTime)

	tick := time.NewTicker(time.Duration(observer.conf.BtcObLoopWaitTime) * time.Second)
	for {
		select {
		case <-tick.C:
			newTop, hash, err := observer.cli.GetCurrentHeightAndHash()
			if err != nil {
				log.Errorf("[BtcObserver] GetCurrentHeightAndHash failed, loop continue: %v", err)
				continue
			}
			log.Tracef("[BtcObserver] start observing from block %s at height %d", hash, newTop)

			if newTop <= top { // Prevent rollback
				log.Tracef("[BtcObserver] height not enough: now is %d, prev is %d", newTop, top)
				continue
			}
			for h := top - observer.conf.BtcObConfirmations + 2; h <= newTop-observer.conf.BtcObConfirmations + 1; h++ {
				txns, hash, err := observer.cli.GetTxsInBlockByHeight(h)
				if err != nil {
					log.Errorf("[BtcObserver] failed to check block %s, retry after 10 sec: %v", hash, err)
					h--
					time.Sleep(time.Second * time.Duration(observer.conf.SleepTime))
					continue
				}
				count := observer.SearchTxInBlock(txns, h, relaying)
				if count > 0 {
					log.Infof("[BtcObserver] %d tx found in block(height:%d) %s", count, h, hash)
				}
			}

			top = newTop
		}
	}
}

func (observer *BtcObserver) SearchTxInBlock(txns []*wire.MsgTx, height uint32, relaying chan *CrossChainItem) int {
	count := 0
	for i := 0; i < len(txns); i++ {
		if !checkIfCrossChainTx(txns[i], observer.NetParam) {
			continue
		}
		var buf bytes.Buffer
		err := txns[i].BtcEncode(&buf, wire.ProtocolVersion, wire.LatestEncoding)
		if err != nil {
			log.Errorf("[SearchTxInBlock] failed to encode transaction: %v", err)
			continue
		}
		txid := txns[i].TxHash()
		proof, err := observer.cli.GetProof([]string{txid.String()})
		if err != nil {
			switch err.(type) {
			case NetErr:
				log.Errorf("[SearchTxInBlock] post err when try to get proof for tx %s: %v", txid.String(), err)
				i--
				time.Sleep(time.Second * time.Duration(observer.conf.SleepTime))
			default:
				log.Errorf("[SearchTxInBlock] failed to get proof for tx %s: %v", txid.String(), err)
			}
			continue
		}
		proofBytes, _ := hex.DecodeString(proof)
		relaying <- &CrossChainItem {
			Proof:  proofBytes,
			Tx:     buf.Bytes(),
			Height: height,
			Txid:   txid,
		}
		log.Infof("[SearchTxInBlock] eligible transaction found, txid: %s", txid.String())
		count++
	}

	return count
}

type AllianceObConfig struct {
	AlliaObLoopWaitTime    int64  `json:"allia_ob_loop_wait_time"`
	WatchingKey            string `json:"watching_key"`
	AllianceJsonRpcAddress string `json:"alliance_json_rpc_address"`
	WalletFile             string `json:"wallet_file"`
	WalletPwd              string `json:"wallet_pwd"`
	SleepTime int `json:"sleep_time"`
}

type AllianceObserver struct {
	allia *sdk.MultiChainSdk
	conf  *AllianceObConfig
}

func NewAllianceObserver(allia *sdk.MultiChainSdk, conf *AllianceObConfig) *AllianceObserver {
	return &AllianceObserver{
		allia: allia,
		conf:  conf,
	}
}

func (observer *AllianceObserver) Listen(collecting chan *FromAllianceItem) {
	top := alliaCheckPoints["testnet"].Height //temp

	log.Infof("[AllianceObserver] get start height %d from checkpoint, check once %d seconds", top, observer.conf.AlliaObLoopWaitTime)
	tick := time.NewTicker(time.Duration(observer.conf.AlliaObLoopWaitTime) * time.Second)
	for {
		select {
		case <-tick.C:
			count := 0
			newTop, err := observer.allia.GetCurrentBlockHeight()
			if err != nil {
				log.Errorf("[AllianceObserver] failed to get current height, retry after 10 sec: %v", err)
				continue
			}
			log.Tracef("[AllianceObserver] start observing from height %d", newTop)

			if newTop-top == 0 {
				continue
			}

			h := top + 1
			for h <= newTop {
				events, err := observer.allia.GetSmartContractEventByBlock(h)
				if err != nil {
					log.Errorf("[AllianceObserver] GetSmartContractEventByBlock failed, retry after 10 sec: %v", err)
					time.Sleep(time.Second * time.Duration(observer.conf.SleepTime))
					continue
				}

				for _, e := range events {
					for _, n := range e.Notify {
						states, ok := n.States.([]interface{})
						if !ok {
							continue
						}
						name, ok := states[0].(string)
						if ok && name == observer.conf.WatchingKey {
							tx := states[1].(string)
							collecting <- &FromAllianceItem {
								Tx: tx,
							}
							count++
							log.Infof("[AllianceObserver] captured: %s when height is %d", tx, h)
						}
					}
				}

				h++
			}
			if count > 0 {
				log.Infof("[AllianceObserver] total %d transactions captured this time", count)
			}
			top = newTop
		}
	}
}
