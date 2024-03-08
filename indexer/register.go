package indexer

import (
	"github.com/dop251/goja"

	"github.com/smartbch/egvm/script/extension"
	"github.com/smartbch/egvm/script/request"
	"github.com/smartbch/egvm/script/types"
)

func registerFunctions(vm *goja.Runtime) {
	// ---------- types ----------
	// uint256
	vm.Set("U256", types.U256)
	vm.Set("HexToU256", types.HexToU256)
	vm.Set("BufToU256", types.BufToU256)

	// signed int256
	vm.Set("S256", types.S256)
	vm.Set("HexToS256", types.HexToS256)
	vm.Set("BufToS256", types.BufToS256)

	// ordered map
	vm.Set("SerializeMaps", types.SerializeMaps)
	vm.Set("DeserializeMap", types.DeserializeMap)
	vm.Set("NewOrderedMapReader", types.NewOrderedMapReader)
	vm.Set("NewOrderedIntMap", types.NewOrderedIntMap)
	vm.Set("NewOrderedStrMap", types.NewOrderedStrMap)
	vm.Set("NewOrderedBufMap", types.NewOrderedBufMap)

	// buffer builder
	vm.Set("NewBufBuilder", types.NewBufBuilder)

	// ---------- extension functions ----------
	// bch
	vm.Set("ParseTxInHex", extension.ParseTxInHex)
	vm.Set("SignTxAndSerialize", extension.SignTxAndSerialize)
	vm.Set("MerkleProofToRootAndMatches", extension.MerkleProofToRootAndMatches)

	// aes encryption
	vm.Set("AesGcmDecrypt", extension.AesGcmDecrypt)
	vm.Set("AesGcmEncrypt", extension.AesGcmEncrypt)

	// public-key encryption
	vm.Set("BufToPrivateKey", extension.BufToPrivateKey)
	vm.Set("BufToPublicKey", extension.BufToPublicKey)

	// hash functions
	vm.Set("Keccak256", extension.Keccak256)
	vm.Set("Sha256", extension.Sha256)
	vm.Set("Ripemd160", extension.Ripemd160)
	vm.Set("XxHash32", extension.XxHash32)
	vm.Set("XxHash64", extension.XxHash64)
	vm.Set("XxHash128", extension.XxHash128)
	vm.Set("XxHash32Int", extension.XxHash32Int)

	// buffer functions
	vm.Set("BufConcat", extension.BufConcat)
	vm.Set("HexToBuf", extension.HexToBuf)
	vm.Set("UTF8StrToBuf", extension.UTF8StrToBuf)
	vm.Set("HexToPaddingBuf", extension.HexToPaddingBuf)
	vm.Set("B64ToBuf", extension.B64ToBuf)
	vm.Set("BufToB64", extension.BufToB64)
	vm.Set("BufToHex", extension.BufToHex)
	vm.Set("BufEqual", extension.BufEqual)
	vm.Set("BufCompare", extension.BufCompare)
	vm.Set("BufReverse", extension.BufReverse)
	vm.Set("BufToU32BE", extension.BufToU32BE)
	vm.Set("BufToU32LE", extension.BufToU32LE)
	vm.Set("U64ToBufBE", extension.U64ToBufBE)
	vm.Set("U64ToBufLE", extension.U64ToBufLE)
	vm.Set("U32ToBufBE", extension.U32ToBufBE)
	vm.Set("U32ToBufLE", extension.U32ToBufLE)

	// signature
	vm.Set("GetEthSignedMessage", extension.GetEthSignedMessage)
	vm.Set("VerifySignature", extension.VerifySignature)
	vm.Set("Ecrecover", extension.Ecrecover)

	// compress
	vm.Set("ZstdCompress", extension.ZstdCompress)
	vm.Set("ZstdDecompress", extension.ZstdDecompress)

	// merkle tree
	vm.Set("VerifyMerkleProofSha256", extension.VerifyMerkleProofSha256)
	vm.Set("VerifyMerkleProofKeccak256", extension.VerifyMerkleProofKeccak256)

	// bip32 key
	vm.Set("GenerateRandomBip32Key", extension.GenerateRandomBip32Key)
	vm.Set("B58ToBip32Key", extension.B58ToBip32Key)
	vm.Set("BufToBip32Key", extension.BufToBip32Key)

	// cpu
	vm.Set("GetCPUID", extension.GetCPUID)
	vm.Set("GetTSC", extension.GetTSC)
	vm.Set("GetTSCBenchStart", extension.GetTSCBenchStart)
	vm.Set("GetTSCBenchEnd", extension.GetTSCBenchEnd)

	// debug
	vm.Set("Printf", extension.Printf)
	vm.Set("Println", extension.Println)

	// system
	vm.Set("Sleep", extension.Sleep)
	vm.Set("SleepMs", extension.SleepMs)

	// ---------- http(s) request ----------
	vm.Set("HttpsRequest", request.HttpsRequest)
	vm.Set("AttestEnclaveServer", request.AttestEnclaveServer)

}
