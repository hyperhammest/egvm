package request

import (
	"testing"

	"github.com/dop251/goja"
	"github.com/stretchr/testify/require"
)

const (
	AttestScriptTemplate = `
		const signerID = '8c83745f1d946d0b9ab8b2233d63449f4274504fc6f67870598ef90f663187df'
		const uniqueID = 'a30551d6f49a81fb52b74ebe6743e8ace4e141e5d89d59ea51eedcfdc74d8654'
		const [ok, reason] = AttestEnclaveServer('https://elfincdn111.paralinker.io', signerID, uniqueID)
	`
)

func setupGojaVmForEnclave() *goja.Runtime {
	vm := goja.New()
	vm.Set("AttestEnclaveServer", AttestEnclaveServer)
	return vm
}

// Note: test it with CGO_FLAGS and CGO_LDFLAGS
func TestAttestEnclaveServer(t *testing.T) {
	vm := setupGojaVmForEnclave()
	tlsConfig, _ = loadTlsConfigForTest(TrustedCertsPathForTest)
	_, err := vm.RunString(AttestScriptTemplate)
	require.NoError(t, err)

	ok := vm.Get("ok").Export().(bool)
	reason := vm.Get("reason").Export().(string)
	require.True(t, ok)
	require.Empty(t, reason)
}
