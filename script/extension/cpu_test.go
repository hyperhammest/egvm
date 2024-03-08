package extension

import (
	"fmt"
	"testing"

	"github.com/dop251/goja"
	"github.com/stretchr/testify/require"

	"github.com/smartbch/egvm/script/types"
)

const (
	CPUIDScriptTemplate = `
		const cpuInfo = GetCPUID()
		const cpuInfoObj = JSON.parse(cpuInfo)
	`

	TSCScriptTemplate = `
		const n = 100
		const tscBz = GetTSC()
		const tsc = BufToU256(tscBz)

  		const startBz = GetTSCBenchStart()
		for (let i = 0; i < n; i++) {
			// code to evaluate
		}
		const endBz = GetTSCBenchEnd()

		const start = BufToU256(startBz)
		const end = BufToU256(endBz)
		const avg = end.Sub(start).Sub(tsc).Div(U256(n))
	`
)

func setupGojaVmForCPU() *goja.Runtime {
	vm := goja.New()
	vm.Set("GetCPUID", GetCPUID)
	vm.Set("GetTSC", GetTSC)
	vm.Set("GetTSCBenchStart", GetTSCBenchStart)
	vm.Set("GetTSCBenchEnd", GetTSCBenchEnd)
	vm.Set("U256", types.U256)
	vm.Set("BufToU256", types.BufToU256)
	return vm
}

func TestGetCPUID(t *testing.T) {
	vm := setupGojaVmForCPU()
	_, err := vm.RunString(CPUIDScriptTemplate)
	require.NoError(t, err)

	cpuInfo := vm.Get("cpuInfo").Export().(string)
	require.NotEmpty(t, cpuInfo)
	fmt.Printf("cpuInfo: %s\n", cpuInfo)

	cpuInfoObj := vm.Get("cpuInfoObj").Export().(map[string]interface{})
	require.NotNil(t, cpuInfoObj)
	fmt.Printf("cpuInfoObj: %+v\n", cpuInfoObj)
}

func TestTSC(t *testing.T) {
	vm := setupGojaVmForCPU()
	_, err := vm.RunString(TSCScriptTemplate)
	require.NoError(t, err)

	tsc := vm.Get("tsc").Export().(types.Uint256)
	start := vm.Get("start").Export().(types.Uint256)
	end := vm.Get("end").Export().(types.Uint256)
	avg := vm.Get("avg").Export().(types.Uint256)
	require.NotZero(t, tsc)
	require.NotZero(t, start)
	require.NotZero(t, end)
	require.NotZero(t, avg)
	require.True(t, end.Gt(start))

	fmt.Printf("tsc: %v\n", tsc.ToSafeInteger())
	fmt.Printf("start: %v\n", start.ToSafeInteger())
	fmt.Printf("end: %v\n", end.ToSafeInteger())
	fmt.Printf("avg: %v\n", avg.ToSafeInteger())
}
