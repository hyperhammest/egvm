package extension

import (
	"crypto/sha256"
	"io"
	"unicode"

	xxh32 "github.com/OneOfOne/xxhash"
	"github.com/cespare/xxhash/v2"
	"github.com/dop251/goja"
	"github.com/zeebo/xxh3"
	"golang.org/x/crypto/ripemd160"
	"golang.org/x/crypto/sha3"

	"github.com/smartbch/egvm/script/types"
)

// ===============

func isASCII(s string) bool {
	for _, c := range s {
		if c > unicode.MaxASCII {
			return false
		}
	}
	return true
}

func hashFunc(f goja.FunctionCall, vm *goja.Runtime, h io.Writer) {
	var buf [32]byte
	for _, arg := range f.Arguments {
		switch v := arg.Export().(type) {
		case string:
			if !isASCII(v) {
				panic(vm.ToValue("Non-ascii string is not supported for hash"))
			}
			h.Write([]byte(v))
		case goja.ArrayBuffer:
			h.Write(v.Bytes())
		case types.Uint256:
			v.X.WriteToArray32(&buf)
			h.Write(buf[:])
		default:
			panic(vm.ToValue("Unsupported type for hash"))
		}
	}
}

func Keccak256(f goja.FunctionCall, vm *goja.Runtime) goja.Value {
	h := sha3.NewLegacyKeccak256()
	hashFunc(f, vm, h)
	return vm.ToValue(vm.NewArrayBuffer(h.Sum(nil)))
}

func Sha256(f goja.FunctionCall, vm *goja.Runtime) goja.Value {
	h := sha256.New()
	hashFunc(f, vm, h)
	return vm.ToValue(vm.NewArrayBuffer(h.Sum(nil)))
}

func Ripemd160(f goja.FunctionCall, vm *goja.Runtime) goja.Value {
	h := ripemd160.New()
	hashFunc(f, vm, h)
	return vm.ToValue(vm.NewArrayBuffer(h.Sum(nil)))
}

func XxHash32(f goja.FunctionCall, vm *goja.Runtime) goja.Value {
	h := xxh32.New32()
	hashFunc(f, vm, h)
	return vm.ToValue(vm.NewArrayBuffer(h.Sum(nil)))
}

func XxHash64(f goja.FunctionCall, vm *goja.Runtime) goja.Value {
	h := xxhash.New()
	hashFunc(f, vm, h)
	return vm.ToValue(vm.NewArrayBuffer(h.Sum(nil)))
}

func XxHash128(f goja.FunctionCall, vm *goja.Runtime) goja.Value {
	h := xxh3.New()
	hashFunc(f, vm, h)
	hash128 := h.Sum128().Bytes()
	return vm.ToValue(vm.NewArrayBuffer(hash128[:]))
}

func XxHash32Int(f goja.FunctionCall, vm *goja.Runtime) goja.Value {
	h := xxh32.New32()
	hashFunc(f, vm, h)
	return vm.ToValue(h.Sum32())
}
