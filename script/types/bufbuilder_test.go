package types

import (
	"testing"

	"github.com/dop251/goja"
	gethcmn "github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
)

const (
	BufBuilderScriptTemplate = `
		let bb = NewBufBuilder()

		const buffer1 = new Uint8Array([1, 2, 3, 4]).buffer
		const buffer2 = new Uint8Array([5, 6, 7, 8]).buffer
		const write1 = bb.Write(buffer1)
		const write2 = bb.Write(buffer2)
		const len1 = bb.Len()
	
		const bz = bb.ToBuf()
		bb.Reset()
		const len2 = bb.Len()
	`
)

func setupGojaVmForBuffer() *goja.Runtime {
	vm := goja.New()
	vm.Set("NewBufBuilder", NewBufBuilder)
	return vm
}

func TestBufBuilder(t *testing.T) {
	vm := setupGojaVmForBuffer()
	_, err := vm.RunString(BufBuilderScriptTemplate)
	require.NoError(t, err)

	write1 := vm.Get("write1").Export().(int64)
	write2 := vm.Get("write2").Export().(int64)
	len1 := vm.Get("len1").Export().(int64)
	len2 := vm.Get("len2").Export().(int64)
	require.EqualValues(t, 4, write1)
	require.EqualValues(t, 4, write2)
	require.EqualValues(t, 8, len1)
	require.EqualValues(t, 0, len2)

	bz := vm.Get("bz").Export().(goja.ArrayBuffer)
	bzHex := gethcmn.Bytes2Hex(bz.Bytes())
	require.EqualValues(t, "0102030405060708", bzHex)

}
