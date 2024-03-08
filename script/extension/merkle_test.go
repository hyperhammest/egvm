package extension

import (
	"testing"

	"github.com/dop251/goja"
	"github.com/stretchr/testify/require"
)

const (
	VerifyMerkleProofKeccak256ScriptTemplate = `
		const rootHex = '0xb54f0bff950c0ffa22e64267f5d655f51ea54c1cd4cad1773d9388fc2a69aceb'
		const leafHex = '0xb001d570adb732e5fa562757f7124c2ff02f4f2f6ea93bba4349918789518d22'
		const proof1Hex = '0x01845ca97f435a23e9a16639b8d061bb8c1f90e74209da8d8e9cbfe424974059'
		const proof2Hex = '0x6a431fe5154748a6a8111d139994ab1f62f086a27b778ec467f14fbe318f35a2'
		const proof3Hex = '0xb0cf8fa6c34e1b6a2a168171f50310c8ba2cb34d877e9fecdc5179c96321d6e2'

		const rootBz = HexToBuf(rootHex)
		const leafBz = HexToBuf(leafHex)
		const proofBz = BufConcat(HexToBuf(proof1Hex), HexToBuf(proof2Hex), HexToBuf(proof3Hex))

		const ok = VerifyMerkleProofKeccak256(rootBz, proofBz, leafBz)
	`

	VerifyMerkleProofSha256ScriptTemplate = `
		const rootHex = '0x5d4ed29f119592a57e6b17a48c67e658ab4fe4f54d203f760571a45911186e0f'
		const leafHex = '0x72b4fc7fcc23336eeb3d0dd35ed777c39c619138da22fdd772eae05292cef106'
		const proof1Hex = '0x245843abef9e72e7efac30138a994bf6301e7e1d7d7042a33d42e863d2638811'
		const proof2Hex = '0xf07039046ff89a76fabadfc0411b6401433869236f76c47a5a3e4cdadcd1118d'
		const proof3Hex = '0x67db80cd1029a8918c45cd6eca91469c711a19c161dc8e23af30e39f5e9e0681'

		const rootBz = HexToBuf(rootHex)
		const leafBz = HexToBuf(leafHex)
		const proofBz = BufConcat(HexToBuf(proof1Hex), HexToBuf(proof2Hex), HexToBuf(proof3Hex))

		const ok = VerifyMerkleProofSha256(rootBz, proofBz, leafBz)
	`
)

func setupGojaVmForMerkle() *goja.Runtime {
	vm := goja.New()
	vm.Set("VerifyMerkleProofSha256", VerifyMerkleProofSha256)
	vm.Set("VerifyMerkleProofKeccak256", VerifyMerkleProofKeccak256)
	vm.Set("HexToBuf", HexToBuf)
	vm.Set("BufConcat", BufConcat)
	return vm
}

func TestVerifyMerkleProofKeccak256(t *testing.T) {
	vm := setupGojaVmForMerkle()
	_, err := vm.RunString(VerifyMerkleProofKeccak256ScriptTemplate)
	require.NoError(t, err)

	ok := vm.Get("ok").Export().(bool)
	require.True(t, ok)
}

func TestVerifyMerkleProofSha256(t *testing.T) {
	vm := setupGojaVmForMerkle()
	_, err := vm.RunString(VerifyMerkleProofSha256ScriptTemplate)
	require.NoError(t, err)

	ok := vm.Get("ok").Export().(bool)
	require.True(t, ok)
}
