package extension

import (
	"bytes"
	"encoding/base64"
	"encoding/binary"
	"encoding/hex"
	"github.com/dop251/goja"
	gethcmn "github.com/ethereum/go-ethereum/common"

	"github.com/smartbch/egvm/script/utils"
)

func BufConcat(f goja.FunctionCall, vm *goja.Runtime) goja.Value {
	var data [][]byte
	totalLen := 0
	for _, arg := range f.Arguments {
		switch v := arg.Export().(type) {
		case goja.ArrayBuffer:
			data = append(data, v.Bytes())
			totalLen += len(v.Bytes())
		default:
			panic(vm.ToValue("Unsupported type for BufConcat"))
		}
	}

	result := make([]byte, 0, totalLen)
	for _, bz := range data {
		result = append(result, bz...)
	}
	return vm.ToValue(vm.NewArrayBuffer(result))
}

func B64ToBuf(f goja.FunctionCall, vm *goja.Runtime) goja.Value {
	if len(f.Arguments) != 1 {
		panic(utils.IncorrectArgumentCount)
	}
	str, ok := f.Arguments[0].Export().(string)
	if !ok {
		panic(goja.NewSymbol("The first argument must be string"))
	}
	data, err := base64.StdEncoding.DecodeString(str)
	if err != nil {
		panic(goja.NewSymbol("error in B64ToBuf: " + err.Error()))
	}
	return vm.ToValue(vm.NewArrayBuffer(data))
}

func HexToBuf(f goja.FunctionCall, vm *goja.Runtime) goja.Value {
	if len(f.Arguments) != 1 {
		panic(utils.IncorrectArgumentCount)
	}
	str, ok := f.Arguments[0].Export().(string)
	if !ok {
		panic(goja.NewSymbol("The first argument must be string"))
	}

	data := gethcmn.FromHex(str)
	return vm.ToValue(vm.NewArrayBuffer(data))
}

func UTF8StrToBuf(f goja.FunctionCall, vm *goja.Runtime) goja.Value {
	if len(f.Arguments) != 1 {
		panic(utils.IncorrectArgumentCount)
	}
	str, ok := f.Arguments[0].Export().(string)
	if !ok {
		panic(goja.NewSymbol("The first argument must be string"))
	}

	return vm.ToValue(vm.NewArrayBuffer([]byte(str)))
}

// HexToPaddingBuf encodes a hex string to a padding buffer in big-endian
func HexToPaddingBuf(f goja.FunctionCall, vm *goja.Runtime) goja.Value {
	if len(f.Arguments) != 2 {
		panic(utils.IncorrectArgumentCount)
	}
	str, ok := f.Arguments[0].Export().(string)
	if !ok {
		panic(goja.NewSymbol("The first argument must be string"))
	}

	n, ok := f.Arguments[1].Export().(int64)
	if !ok {
		panic(goja.NewSymbol("The second argument must be int"))
	}

	if n <= 0 || n > 256 {
		panic(goja.NewSymbol("The second argument must be between 0 and 256"))
	}

	data := gethcmn.FromHex(str)
	paddingData := gethcmn.LeftPadBytes(data, int(n))
	return vm.ToValue(vm.NewArrayBuffer(paddingData))
}

func BufToB64(f goja.FunctionCall, vm *goja.Runtime) goja.Value {
	buf := utils.GetOneArrayBuffer(f)
	str := base64.StdEncoding.EncodeToString(buf)
	return vm.ToValue(str)
}

func BufToHex(f goja.FunctionCall, vm *goja.Runtime) goja.Value {
	buf := utils.GetOneArrayBuffer(f)
	str := hex.EncodeToString(buf)
	return vm.ToValue(str)
}

func BufEqual(f goja.FunctionCall, vm *goja.Runtime) goja.Value {
	a, b := utils.GetTwoArrayBuffers(f)
	return vm.ToValue(bytes.Equal(a, b))
}

func BufCompare(f goja.FunctionCall, vm *goja.Runtime) goja.Value {
	a, b := utils.GetTwoArrayBuffers(f)
	return vm.ToValue(bytes.Compare(a, b))
}

func BufReverse(f goja.FunctionCall, vm *goja.Runtime) goja.Value {
	a := utils.GetOneArrayBuffer(f)
	return vm.ToValue(vm.NewArrayBuffer(bytesReverse(a)))
}

func BufToU32BE(f goja.FunctionCall, vm *goja.Runtime) goja.Value {
	bz := utils.GetOneArrayBuffer(f)
	return vm.ToValue(binary.BigEndian.Uint32(bz))
}

func BufToU32LE(f goja.FunctionCall, vm *goja.Runtime) goja.Value {
	bz := utils.GetOneArrayBuffer(f)
	return vm.ToValue(binary.LittleEndian.Uint32(bz))
}

func U64ToBufBE(f goja.FunctionCall, vm *goja.Runtime) goja.Value {
	u64 := utils.GetOneUint64(f)
	var buf [8]byte
	binary.BigEndian.PutUint64(buf[:], u64)
	return vm.ToValue(vm.NewArrayBuffer(buf[:]))
}

func U64ToBufLE(f goja.FunctionCall, vm *goja.Runtime) goja.Value {
	u64 := utils.GetOneUint64(f)
	var buf [8]byte
	binary.LittleEndian.PutUint64(buf[:], u64)
	return vm.ToValue(vm.NewArrayBuffer(buf[:]))
}

func U32ToBufBE(f goja.FunctionCall, vm *goja.Runtime) goja.Value {
	u64 := utils.GetOneUint64(f)
	var buf [4]byte
	binary.BigEndian.PutUint32(buf[:], uint32(u64))
	return vm.ToValue(vm.NewArrayBuffer(buf[:]))
}

func U32ToBufLE(f goja.FunctionCall, vm *goja.Runtime) goja.Value {
	u64 := utils.GetOneUint64(f)
	var buf [4]byte
	binary.LittleEndian.PutUint32(buf[:], uint32(u64))
	return vm.ToValue(vm.NewArrayBuffer(buf[:]))
}

func bytesReverse(bz []byte) []byte {
	rBz := make([]byte, 0, len(bz))
	for i := range bz {
		rBz = append(rBz, bz[len(bz)-1-i])
	}
	return rBz
}
