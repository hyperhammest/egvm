package extension

import (
	"encoding/hex"
	"fmt"
	"strings"
	"testing"

	"github.com/dop251/goja"
	ecies "github.com/ecies/go/v2"
	gethcmn "github.com/ethereum/go-ethereum/common"
	"github.com/gcash/bchd/bchec"
	"github.com/stretchr/testify/require"
)

const (
	TxAndSignScriptTemplate = `
		const txHex = '020000000147c8c5a1f4d7d5e3846a7e945daf634722340f617be0ff8736e668d7ee7d9fb402000000644128fd33544f9530b1a8ae03340bdfe9385324adf0ecefc39d53e6fddf9afdb64ccf3ef12bc692045d87e90380919429d3c5bafc29b51515aa0a992dd0d085663e412102dde6c067f5e1a641dedab654cbbd9c3b4c6f8adbf2aeb17c6500319d2c08f08e00000000030000000000000000666a04454754581456eb561cb6f98a985f80464fa99267a462c91bdb14e94358e473941de2d75d19fa330d607e05ffab4214efc507fb38cbcae3b32d1777e54593bc07eca5a1204ea5c508a6566e76240543f8feb06fd457777be300005af3107a40000000000110270000000000001976a9148097f6fbaa0dfdfe4f064bb650324c5e8018242088acca331e00000000001976a914307f40d73e01af33364901d82d5614e370f905d388ac00000000'
		const bchTx = ParseTxInHex(txHex)

		const privateKeyHex = 'c9cb992b13141bb3326d028020030f33b92ea9a64b6530291e7876938bd31479'
		const privateKeyBuf = HexToBuf(privateKeyHex)
		const privateKey = BufToPrivateKey(privateKeyBuf)

		const signedTxHex = SignTxAndSerialize(bchTx, privateKey)
	`

	MerkleProofToRootAndMatchesScriptTemplate = `
		const proofHex = '00e0a627f0ffea60563cc47c80d2a7f1854994158a7c8b0c10fed2000000000000000000368bfe8331d7c2223d4ecad3b3b80c0341792ba8294b6c43db3de3ff198b6c75d6373f64ec5e051871de02781b00000006fc8af02b8392ee24089818aea94069eacb75ecdc69f4e40b09a89ffd6e330a17082d2c4367aeb5a45b086c05f3bbae983d6c342aa47c8cca9cd4a78eeda2a20148f2a26ca8318a551117b927ced10eaab935fe55d5b93abb9d8485657c297e6622f48afe69608a59ea420a59081920af1199fa9793867fc05e168ad6e964deaa47f0917dbfc90531a7371834cc18a4e1c5de629b88f804076411be9b42c77e8ce68662e7bc7024408aba81dcfc114195a241265fab2ba584dc3d2eaefff2f0d9023f00'
		const proofBz = HexToBuf(proofHex)
		const [merkleRoot, txID1] = MerkleProofToRootAndMatches(proofBz)
	`
)

func setupGojaVmForBCH() *goja.Runtime {
	vm := goja.New()
	vm.Set("ParseTxInHex", ParseTxInHex)
	vm.Set("SignTxAndSerialize", SignTxAndSerialize)
	vm.Set("MerkleProofToRootAndMatches", MerkleProofToRootAndMatches)
	vm.Set("HexToBuf", HexToBuf)
	vm.Set("BufToPrivateKey", BufToPrivateKey)
	return vm
}

func TestTxAndSign(t *testing.T) {
	vm := setupGojaVmForBCH()
	_, err := vm.RunString(TxAndSignScriptTemplate)
	require.NoError(t, err)

	signedTxHex := vm.Get("signedTxHex").Export().(string)
	require.NotEmpty(t, signedTxHex)
}

func TestMerkleProofToRootAndMatches(t *testing.T) {
	vm := setupGojaVmForBCH()
	_, err := vm.RunString(MerkleProofToRootAndMatchesScriptTemplate)
	require.NoError(t, err)

	merkleRootBz := vm.Get("merkleRoot").Export().(goja.ArrayBuffer)
	merkleRootBzHex := gethcmn.Bytes2Hex(merkleRootBz.Bytes())
	require.EqualValues(t, "756c8b19ffe33ddb436c4b29a82b7941030cb8b3d3ca4e3d22c2d73183fe8b36", merkleRootBzHex)

	txID1Bz := vm.Get("txID1").Export().(goja.ArrayBuffer)
	txID1BzHex := gethcmn.Bytes2Hex(txID1Bz.Bytes())
	require.EqualValues(t, "170a336efd9fa8090be4f469dcec75cbea6940a9ae18980824ee92832bf08afc", txID1BzHex)
}

func TestGetPubkeyHashHex(t *testing.T) {
	pkScript, _ := hex.DecodeString("76a914307f40d73e01af33364901d82d5614e370f905d388ac")
	pbkHash := getPubkeyHashFromPkScript(pkScript)
	require.Equal(t, "307f40d73e01af33364901d82d5614e370f905d3", pbkHash)
}

func TestGetOpRetData(t *testing.T) {
	pkScript, _ := hex.DecodeString("6a04454754581456eb561cb6f98a985f80464fa99267a462c91bdb14e94358e473941de2d75d19fa330d607e05ffab4214efc507fb38cbcae3b32d1777e54593bc07eca5a1204ea5c508a6566e76240543f8feb06fd457777be300005af3107a400000000001")
	retData := getOpRetData(pkScript)
	require.Len(t, retData, 5)
	require.Equal(t, "45475458", retData[0])
	require.Equal(t, "56eb561cb6f98a985f80464fa99267a462c91bdb", retData[1])
	require.Equal(t, "e94358e473941de2d75d19fa330d607e05ffab42", retData[2])
	require.Equal(t, "efc507fb38cbcae3b32d1777e54593bc07eca5a1", retData[3])
	require.Equal(t, "4ea5c508a6566e76240543f8feb06fd457777be300005af3107a400000000001", retData[4])
}

func TestParseTxInHex(t *testing.T) {
	// https://blockchair.com/bitcoin-cash/transaction/c1d33d4c03c81f8dee2f176b944bc05eb08524163698d4d397a8fb9ca7cd2651
	txHex := "020000000147c8c5a1f4d7d5e3846a7e945daf634722340f617be0ff8736e668d7ee7d9fb402000000644128fd33544f9530b1a8ae03340bdfe9385324adf0ecefc39d53e6fddf9afdb64ccf3ef12bc692045d87e90380919429d3c5bafc29b51515aa0a992dd0d085663e4121031c60b05831b6f3c31739856575cde27d97d9fe926a63d51abce4a0c16b4108be00000000030000000000000000666a04454754581456eb561cb6f98a985f80464fa99267a462c91bdb14e94358e473941de2d75d19fa330d607e05ffab4214efc507fb38cbcae3b32d1777e54593bc07eca5a1204ea5c508a6566e76240543f8feb06fd457777be300005af3107a40000000000110270000000000001976a9148097f6fbaa0dfdfe4f064bb650324c5e8018242088acca331e00000000001976a914307f40d73e01af33364901d82d5614e370f905d388ac00000000"

	tx, err := parseTxInHex(txHex)
	require.NoError(t, err)
	require.Equal(t, int32(2), tx.Version)
	require.Equal(t, uint32(0), tx.LockTime)
	require.Equal(t, "c1d33d4c03c81f8dee2f176b944bc05eb08524163698d4d397a8fb9ca7cd2651", tx.HexTxID)
	require.Len(t, tx.TxIn, 1)
	require.Len(t, tx.TxOut, 3)
	require.Equal(t, TxIn{
		PreviousOutPoint: OutPoint{
			HexTxID: "b49f7deed768e63687ffe07b610f34224763af5d947e6a84e3d5d7f4a1c5c847",
			Index:   2,
		},
		HexPubkey: "031c60b05831b6f3c31739856575cde27d97d9fe926a63d51abce4a0c16b4108be",
		Sequence:  0,
	}, tx.TxIn[0])
	require.Equal(t, TxOut{
		Value:         0,
		HexPubkeyHash: "",
		HexDataElements: []string{
			"45475458",
			"56eb561cb6f98a985f80464fa99267a462c91bdb",
			"e94358e473941de2d75d19fa330d607e05ffab42",
			"efc507fb38cbcae3b32d1777e54593bc07eca5a1",
			"4ea5c508a6566e76240543f8feb06fd457777be300005af3107a400000000001",
		},
	}, tx.TxOut[0])
	require.Equal(t, TxOut{
		Value:           10000,
		HexPubkeyHash:   "8097f6fbaa0dfdfe4f064bb650324c5e80182420",
		HexDataElements: nil,
	}, tx.TxOut[1])
	require.Equal(t, TxOut{
		Value:           1979338,
		HexPubkeyHash:   "307f40d73e01af33364901d82d5614e370f905d3",
		HexDataElements: nil,
	}, tx.TxOut[2])
}

func TestSignTxAndSerialize(t *testing.T) {
	eciesKey, err := ecies.GenerateKey()
	require.NoError(t, err)

	privKey := PrivateKey{key: eciesKey}
	bchecKey := (*bchec.PrivateKey)(privKey.toECDSA())

	txHex := "020000000147c8c5a1f4d7d5e3846a7e945daf634722340f617be0ff8736e668d7ee7d9fb402000000644128fd33544f9530b1a8ae03340bdfe9385324adf0ecefc39d53e6fddf9afdb64ccf3ef12bc692045d87e90380919429d3c5bafc29b51515aa0a992dd0d085663e4121031c60b05831b6f3c31739856575cde27d97d9fe926a63d51abce4a0c16b4108be00000000030000000000000000666a04454754581456eb561cb6f98a985f80464fa99267a462c91bdb14e94358e473941de2d75d19fa330d607e05ffab4214efc507fb38cbcae3b32d1777e54593bc07eca5a1204ea5c508a6566e76240543f8feb06fd457777be300005af3107a40000000000110270000000000001976a9148097f6fbaa0dfdfe4f064bb650324c5e8018242088acca331e00000000001976a914307f40d73e01af33364901d82d5614e370f905d388ac00000000"
	tx, err := parseTxInHex(strings.ReplaceAll(txHex,
		"031c60b05831b6f3c31739856575cde27d97d9fe926a63d51abce4a0c16b4108be",
		hex.EncodeToString(bchecKey.PubKey().SerializeCompressed())))
	require.NoError(t, err)

	signedTxHex, err := signTxAndSerialize(*tx, privKey)
	require.NoError(t, err)
	fmt.Println(signedTxHex)

	signedTx, err := parseTxInHex(signedTxHex)
	require.NoError(t, err)
	require.Equal(t, int32(2), signedTx.Version)
	require.Equal(t, uint32(0), signedTx.LockTime)
	require.Len(t, signedTx.TxIn, 1)
	require.Len(t, signedTx.TxOut, 3)
	require.Len(t, signedTx.TxOut[0].HexDataElements, 5)
	require.Equal(t, "8097f6fbaa0dfdfe4f064bb650324c5e80182420", signedTx.TxOut[1].HexPubkeyHash)
}
