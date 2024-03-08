package extension

import (
	"github.com/dop251/goja"
	"github.com/klauspost/compress/zstd"

	"github.com/smartbch/egvm/script/utils"
)

// ================================

func ZstdDecompress(f goja.FunctionCall, vm *goja.Runtime) goja.Value {
	src := utils.GetOneArrayBuffer(f)
	var decoder, _ = zstd.NewReader(nil, zstd.WithDecoderConcurrency(0))
	bz, err := decoder.DecodeAll(src, nil)
	if err != nil {
		panic(goja.NewSymbol("error in ZstdDecompress: " + err.Error()))
	}
	return vm.ToValue(vm.NewArrayBuffer(bz))
}

func ZstdCompress(f goja.FunctionCall, vm *goja.Runtime) goja.Value {
	src := utils.GetOneArrayBuffer(f)
	var encoder, _ = zstd.NewWriter(nil)
	bz := encoder.EncodeAll(src, make([]byte, 0, len(src)))
	return vm.ToValue(vm.NewArrayBuffer(bz))
}

// =================== non-deterministic =============

func ND_ReadTsc() uint64 {
	return 0 //TODO
}

func ND_GetEphemeralID() string {
	return "" //TODO
}
