package extension

import (
	"testing"

	"github.com/dop251/goja"
	gethcmn "github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
)

const (
	ZstdCompressScriptTemplate = `
		const buffer1 = new Uint8Array([1, 2, 3, 4, 8, 7, 6, 5]).buffer
		const compressedBz = ZstdCompress(buffer1)
		const decompressedBz = ZstdDecompress(compressedBz)
	`
)

func setupGojaVmForCompress() *goja.Runtime {
	vm := goja.New()
	vm.Set("ZstdCompress", ZstdCompress)
	vm.Set("ZstdDecompress", ZstdDecompress)
	return vm
}

func TestZstdCompressAndDecompress(t *testing.T) {
	vm := setupGojaVmForCompress()
	_, err := vm.RunString(ZstdCompressScriptTemplate)
	require.NoError(t, err)

	compressedBz := vm.Get("compressedBz").Export().(goja.ArrayBuffer)
	compressedBzHex := gethcmn.Bytes2Hex(compressedBz.Bytes())
	decompressedBz := vm.Get("decompressedBz").Export().(goja.ArrayBuffer)
	decompressedBzHex := gethcmn.Bytes2Hex(decompressedBz.Bytes())
	require.EqualValues(t, "28b52ffd04004100000102030408070605a9c79424", compressedBzHex)
	require.EqualValues(t, "0102030408070605", decompressedBzHex)
}
