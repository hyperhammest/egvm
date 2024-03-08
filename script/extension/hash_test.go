package extension

import (
	"testing"

	"github.com/dop251/goja"
	gethcmn "github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
)

const (
	HashBufScriptTemplate = `
		const buffer1 = new Uint8Array([1, 2, 3, 4, 8, 7, 6, 5]).buffer
		const keccackHash = Keccak256(buffer1)
		const sha256Hash = Sha256(buffer1)
		const ripemdHash = Ripemd160(buffer1)
		const xxhash32Hash = XxHash32(buffer1)
		const xxhash64Hash = XxHash64(buffer1)
		const xxhash128Hash = XxHash128(buffer1)
		const xxhash32Int = XxHash32Int(buffer1)
	`
)

func setupGojaVmForHash() *goja.Runtime {
	vm := goja.New()
	vm.Set("Keccak256", Keccak256)
	vm.Set("Sha256", Sha256)
	vm.Set("Ripemd160", Ripemd160)
	vm.Set("XxHash32", XxHash32)
	vm.Set("XxHash64", XxHash64)
	vm.Set("XxHash128", XxHash128)
	vm.Set("XxHash32Int", XxHash32Int)
	return vm
}

func TestHashFunctions(t *testing.T) {
	vm := setupGojaVmForHash()
	_, err := vm.RunString(HashBufScriptTemplate)
	require.NoError(t, err)

	keccackHash := vm.Get("keccackHash").Export().(goja.ArrayBuffer)
	keccackHashHex := gethcmn.Bytes2Hex(keccackHash.Bytes())
	sha256Hash := vm.Get("sha256Hash").Export().(goja.ArrayBuffer)
	sha256HashHex := gethcmn.Bytes2Hex(sha256Hash.Bytes())
	ripemdHash := vm.Get("ripemdHash").Export().(goja.ArrayBuffer)
	ripemdHashHex := gethcmn.Bytes2Hex(ripemdHash.Bytes())
	xxhash32Hash := vm.Get("xxhash32Hash").Export().(goja.ArrayBuffer)
	xxhash32HashHex := gethcmn.Bytes2Hex(xxhash32Hash.Bytes())
	xxhash64Hash := vm.Get("xxhash64Hash").Export().(goja.ArrayBuffer)
	xxhash64HashHex := gethcmn.Bytes2Hex(xxhash64Hash.Bytes())
	xxhash128Hash := vm.Get("xxhash128Hash").Export().(goja.ArrayBuffer)
	xxhash128HashHex := gethcmn.Bytes2Hex(xxhash128Hash.Bytes())
	xxhash32Int := vm.Get("xxhash32Int").Export().(int64)

	require.EqualValues(t, "2691f4fdfe3aa541af9ba914f133fe37517a37ed32a42c39c715e735afb7e94d", keccackHashHex)
	require.EqualValues(t, "258d410a3aa33a094daeba2a99bae7dc416d45ebc02a5e8d7a7974110285dd87", sha256HashHex)
	require.EqualValues(t, "a69e12499849e30bd564770692d185109913d416", ripemdHashHex)
	require.EqualValues(t, "e0d5f344", xxhash32HashHex)
	require.EqualValues(t, "b1005c2e2494c7a9", xxhash64HashHex)
	require.EqualValues(t, "3110682add3b6d6c9223d40cc1d4797e", xxhash128HashHex)
	require.EqualValues(t, 3772117828, xxhash32Int)
}
