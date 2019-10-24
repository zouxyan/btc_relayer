package db

import (
	"os"
	"testing"
)

var (
	txArr = []string{
		"01000000019f074c07f34ffdcac88f76aa403e0725a90870b974c777a7236d6db067481ff2020000006b483045022100c5647452812dd245de91536de723d35239cbd49bb4dd924a5b6376b099a8a716022078938060af6771a44913893eaf0b091de365ee3a7a6ecefa830bf5d4caf6c996012103128a2c4525179e47f38cf3fefca37a61548ca4610255b3fb4ee86de2d3e80c0fffffffff03204e00000000000017a91487a9652e9b396545598c0fc72cb5a98848bf93d3870000000000000000276a256600000000000000020000000000000000f3b8a17f1f957f60c88f105e32ebff3f022e56a4a8ae0800000000001976a91428d2e8cee08857f569e5a1b147c5d5e87339e08188ac00000000",
		"01000000014d369f59ba828a4996f47b03229cbb976a0d0ed841c8c2f7a8e843289b15e631020000006a47304402202676093014919f3aa5dd5237566e4690dfc3503809246c3601cac38c4bd7636202207277e47ec95a28cf6eb30bceb36252704086e2b04c13e35797b0862741d44528012103128a2c4525179e47f38cf3fefca37a61548ca4610255b3fb4ee86de2d3e80c0fffffffff03204e00000000000017a91487a9652e9b396545598c0fc72cb5a98848bf93d3870000000000000000276a256600000000000000020000000000000000f3b8a17f1f957f60c88f105e32ebff3f022e56a4d8230900000000001976a91428d2e8cee08857f569e5a1b147c5d5e87339e08188ac00000000",
		"01000000012de9998067a46027dff12f8c114cb67f9788f52af731caa22a4c0b68babc58d0020000006a4730440220014f7b9c643ce47275552583fe3785f1d72307207b9e95ec293b16cbb15e10bc02205b6452982da43ebc2ddba531cba43a88baad003f2fcaef265e4ec30e8c077b0a012103128a2c4525179e47f38cf3fefca37a61548ca4610255b3fb4ee86de2d3e80c0fffffffff03204e00000000000017a91487a9652e9b396545598c0fc72cb5a98848bf93d3870000000000000000276a256600000000000000020000000000000000f3b8a17f1f957f60c88f105e32ebff3f022e56a4e84a0900000000001976a91428d2e8cee08857f569e5a1b147c5d5e87339e08188ac00000000",
		"0100000001e94f4feea73d3a44ad41cc52ea9ba5ccba9094dbd351310497eafff8233e2bec020000006b483045022100b1d0ec22de2404bf25d620e2c4488d2489d0eb46beca29f5dc219ee464ecd9e3022068a3f0c045b86e14fb99c72173ff9dd8d65dd784fc3af2050085a7c8e9b67cdb012103128a2c4525179e47f38cf3fefca37a61548ca4610255b3fb4ee86de2d3e80c0fffffffff03204e00000000000017a91487a9652e9b396545598c0fc72cb5a98848bf93d3870000000000000000276a256600000000000000020000000000000000f3b8a17f1f957f60c88f105e32ebff3f022e56a418c00900000000001976a91428d2e8cee08857f569e5a1b147c5d5e87339e08188ac00000000",
		"0100000001cbf0d729953a6227098a45752ed871c621f938260a53bfc6c998a347e2fd62e4020000006b483045022100f2e9590d6a9fcd5abce3775b0d263194547890f825350e81927468d807c43dd70220165a5edbca9f1cd42a6197a20ac4d97985dcde02f303b061541c1f86e026c398012103128a2c4525179e47f38cf3fefca37a61548ca4610255b3fb4ee86de2d3e80c0fffffffff03204e00000000000017a91487a9652e9b396545598c0fc72cb5a98848bf93d3870000000000000000276a256600000000000000020000000000000000f3b8a17f1f957f60c88f105e32ebff3f022e56a430390a00000000001976a91428d2e8cee08857f569e5a1b147c5d5e87339e08188ac00000000",
		"01000000017b60622809f978db6b95b09cec9789e2d68539b5eb7b2015ae2f7d79028b320c020000006b48304502210088a67aa604db724f34b2410e813927d26ad3489655cfb252934b25713890a17b02203bbe698595467da469fd20a09049c5920e8f70cae591605375d16bc801253c1f012103128a2c4525179e47f38cf3fefca37a61548ca4610255b3fb4ee86de2d3e80c0fffffffff03204e00000000000017a91487a9652e9b396545598c0fc72cb5a98848bf93d3870000000000000000276a256600000000000000020000000000000000f3b8a17f1f957f60c88f105e32ebff3f022e56a430390a00000000001976a91428d2e8cee08857f569e5a1b147c5d5e87339e08188ac00000000",
		"010000000140b1e5b05e26757f799ebf9f9018f970e77b0e6f7123d19fabaf373a32c3d0a0020000006a47304402204fb2d2d3edebd65675aa8b27f24d2b2398cec0fbf9cac2cfcb7def6df9b7cd6202205d6881aae4b4d65d63efbc1ed1cbb669dc493811b94f2c1629573a5ef5568706012103128a2c4525179e47f38cf3fefca37a61548ca4610255b3fb4ee86de2d3e80c0fffffffff03204e00000000000017a91487a9652e9b396545598c0fc72cb5a98848bf93d3870000000000000000276a256600000000000000020000000000000000f3b8a17f1f957f60c88f105e32ebff3f022e56a430390a00000000001976a91428d2e8cee08857f569e5a1b147c5d5e87339e08188ac00000000",
	}
)

func afterTest() {
	os.RemoveAll("./retry.bin")
}

func TestNewRetryDB(t *testing.T) {
	defer afterTest()
	_, err := NewRetryDB("./", 5, 1)
	if err != nil {
		t.Fatal(err)
	}
}

func TestRetryDB_Put(t *testing.T) {
	defer afterTest()
	db, _ := NewRetryDB("./", 5, 1)
	err := db.Put(txArr[0])
	if err != nil {
		t.Fatal(err)
	}

	val, err := db.GetAll()
	if err != nil {
		t.Fatal(err)
	}
	if len(val) != 1 || val[0] != txArr[0] {
		t.Fatal("not right val")
	}

}

func TestRetryDB_GetAll(t *testing.T) {
	defer afterTest()
	db, _ := NewRetryDB("./", 2, 1)
	for _, tx := range txArr[:3] {
		db.Put(tx)
	}

	vals, err := db.GetAll()
	if err != nil {
		t.Fatal(err)
	}
	if len(vals) != 3 {
		t.Fatal("not right length 3")
	}

	for _, tx := range txArr[3:] {
		db.Put(tx)
	}
	vals, err = db.GetAll()
	if err != nil {
		t.Fatal(err)
	}
	if len(vals) != 7 {
		t.Fatal("not right length 7 ")
	}
	vals, err = db.GetAll()
	if err != nil {
		t.Fatal(err)
	}
	if len(vals) != 4 {
		t.Fatal("not right length 4")
	}

	_, err = db.GetAll()
	if err == nil {
		t.Fatal("err should not be nil")
	}
}

func TestRetryDB_GetAll2(t *testing.T) {
	defer afterTest()
	db, _ := NewRetryDB("./", 0, 1)
	for _, tx := range txArr[:3] {
		db.Put(tx)
	}
	vals, err := db.GetAll()
	if err != nil {
		t.Fatal(err)
	}
	if len(vals) != 3 {
		t.Fatal("not right length 3")
	}
	vals, err = db.GetAll()
	if err != nil {
		t.Fatal(err)
	}
	if len(vals) != 3 {
		t.Fatal("length shouldn't change")
	}
}

func TestRetryDB_Del(t *testing.T) {
	defer afterTest()
}