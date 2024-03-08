package extension

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"hash"

	"github.com/dop251/goja"
	"golang.org/x/crypto/sha3"

	"github.com/smartbch/egvm/script/utils"
)

func verifyMerkleProof(f goja.FunctionCall, vm *goja.Runtime, h hash.Hash) goja.Value {
	root, proof, leaf := utils.GetThreeArrayBuffers(f)
	if len(leaf) != 32 {
		panic(goja.NewSymbol(fmt.Sprintf("Invalid merkle tree leaf size %d", len(leaf))))
	}

	if len(proof)%32 != 0 {
		panic(goja.NewSymbol(fmt.Sprintf("Invalid merkle proof size %d", len(proof))))
	}

	computedHash := leaf
	for offset := 0; offset < len(proof); offset += 32 {
		h.Reset()
		node := proof[offset : offset+32]
		h.Write(combineTwoHash(node, computedHash))
		computedHash = h.Sum(nil)
	}

	ok := bytes.Equal(root, computedHash)
	return vm.ToValue(ok)
}

func combineTwoHash(a, b []byte) []byte {
	bf := bytes.NewBuffer(nil)
	if bytes.Compare(a, b) < 0 {
		bf.Write(a)
		bf.Write(b)
		return bf.Bytes()
	}

	bf.Write(b)
	bf.Write(a)
	return bf.Bytes()
}

func VerifyMerkleProofSha256(f goja.FunctionCall, vm *goja.Runtime) goja.Value {
	return verifyMerkleProof(f, vm, sha256.New())
}

func VerifyMerkleProofKeccak256(f goja.FunctionCall, vm *goja.Runtime) goja.Value {
	return verifyMerkleProof(f, vm, sha3.NewLegacyKeccak256())
}
