package extension

import (
	"encoding/binary"
	"encoding/json"

	"github.com/dop251/goja"
	"github.com/dterei/gotsc"
	"github.com/klauspost/cpuid/v2"

	"github.com/smartbch/egvm/script/utils"
)

type cpuInfoOutput struct {
	cpuid.CPUInfo
	Features []string
	X64Level int
}

func GetCPUID() string {
	info := cpuInfoOutput{
		CPUInfo:  cpuid.CPU,
		Features: cpuid.CPU.FeatureSet(),
		X64Level: cpuid.CPU.X64Level(),
	}

	bz, err := json.Marshal(info)
	if err != nil {
		panic(goja.NewSymbol("Failed to marshal CPU info: " + err.Error()))
	}

	return string(bz)
}

func GetTSC(f goja.FunctionCall, vm *goja.Runtime) goja.Value {
	if len(f.Arguments) != 0 {
		// BIP-44
		// m / purpose' / coin' / account' / change / address_index
		panic(utils.IncorrectArgumentCount)
	}

	tsc := gotsc.TSCOverhead()
	var result [8]byte
	binary.BigEndian.PutUint64(result[:], tsc)
	return vm.ToValue(vm.NewArrayBuffer(result[:]))
}

func GetTSCBenchStart(f goja.FunctionCall, vm *goja.Runtime) goja.Value {
	if len(f.Arguments) != 0 {
		// BIP-44
		// m / purpose' / coin' / account' / change / address_index
		panic(utils.IncorrectArgumentCount)
	}

	start := gotsc.BenchStart()
	var result [8]byte
	binary.BigEndian.PutUint64(result[:], start)
	return vm.ToValue(vm.NewArrayBuffer(result[:]))
}

func GetTSCBenchEnd(f goja.FunctionCall, vm *goja.Runtime) goja.Value {
	if len(f.Arguments) != 0 {
		// BIP-44
		// m / purpose' / coin' / account' / change / address_index
		panic(utils.IncorrectArgumentCount)
	}

	end := gotsc.BenchEnd()
	var result [8]byte
	binary.BigEndian.PutUint64(result[:], end)
	return vm.ToValue(vm.NewArrayBuffer(result[:]))
}
