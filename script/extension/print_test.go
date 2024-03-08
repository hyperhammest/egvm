package extension

import (
	"testing"

	"github.com/dop251/goja"
	"github.com/stretchr/testify/require"
)

const (
	PrintlnScriptTemplate = `
		const owner1 = '1234abcd1'
		const owner2 = '1234abcd2'
		const owner3 = '1234abcd3'
		const ownerList = [owner1, owner2, owner3]
		Println(ownerList)
	`

	PrintfScriptTemplate = `
		const str1 = 'abc'
		const str2 = '1234'
		Printf("%v-%v\n", str1, str2)
	`
)

func setupGojaVmForPrint() *goja.Runtime {
	vm := goja.New()
	vm.Set("Println", Println)
	vm.Set("Printf", Printf)
	return vm
}

func TestPrintln(t *testing.T) {
	vm := setupGojaVmForPrint()
	_, err := vm.RunString(PrintlnScriptTemplate)
	require.NoError(t, err)
}

func TestPrintf(t *testing.T) {
	vm := setupGojaVmForPrint()
	_, err := vm.RunString(PrintfScriptTemplate)
	require.NoError(t, err)
}
