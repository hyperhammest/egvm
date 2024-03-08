package extension

import (
	"testing"

	"github.com/dop251/goja"
	"github.com/stretchr/testify/require"
)

const (
	Bip32KeyBasicScriptTemplate = `
		const randomBip32Key = GenerateRandomBip32Key()
		const b58ToBip32Key = B58ToBip32Key('xprv9s21ZrQH143K3QTDL4LXw2F7HEK3wJUD2nW2nRk4stbPy6cq3jPPqjiChkVvvNKmPGJxWUtg6LnF5kejMRNNU3TGtRBeJgk33yuGBxrMPHi')

		const b58Str = b58ToBip32Key.B58Serialize()
		const isPrivate = b58ToBip32Key.IsPrivate()
		const publicKey = b58ToBip32Key.PublicKey()
		const pubKeyB58Str = publicKey.B58Serialize()

		const privateKey = b58ToBip32Key.ToPrivateKey()
	`

	Bip32KeyDeriveScriptTemplate = `
		const b58ToBip32Key = B58ToBip32Key('xprv9s21ZrQH143K3QTDL4LXw2F7HEK3wJUD2nW2nRk4stbPy6cq3jPPqjiChkVvvNKmPGJxWUtg6LnF5kejMRNNU3TGtRBeJgk33yuGBxrMPHi')
		const deriveKey = b58ToBip32Key.Derive(0, 1, 0, 0, 0)
		const deriveKeyB58Str = deriveKey.B58Serialize()

		const rand32 = new Uint8Array([1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31, 32]).buffer
		const deriveWith32Key = b58ToBip32Key.DeriveWithBytes32(rand32)
		const deriveWith32KeyB58Str = deriveWith32Key.B58Serialize()

	`
)

func setupGojaVmForBip32Key() *goja.Runtime {
	vm := goja.New()
	vm.Set("GenerateRandomBip32Key", GenerateRandomBip32Key)
	vm.Set("B58ToBip32Key", B58ToBip32Key)
	vm.Set("BufToBip32Key", BufToBip32Key)
	return vm
}

func TestBip32KeyBasic(t *testing.T) {
	vm := setupGojaVmForBip32Key()
	_, err := vm.RunString(Bip32KeyBasicScriptTemplate)
	require.NoError(t, err)

	b58Str := vm.Get("b58Str").Export().(string)
	require.EqualValues(t, "xprv9s21ZrQH143K3QTDL4LXw2F7HEK3wJUD2nW2nRk4stbPy6cq3jPPqjiChkVvvNKmPGJxWUtg6LnF5kejMRNNU3TGtRBeJgk33yuGBxrMPHi", b58Str)
	isPrivate := vm.Get("isPrivate").Export().(bool)
	require.True(t, isPrivate)
	pubKeyB58Str := vm.Get("pubKeyB58Str").Export().(string)
	require.EqualValues(t, "xpub661MyMwAqRbcFtXgS5sYJABqqG9YLmC4Q1Rdap9gSE8NqtwybGhePY2gZ29ESFjqJoCu1Rupje8YtGqsefD265TMg7usUDFdp6W1EGMcet8", pubKeyB58Str)
}

func TestBip32KeyDerive(t *testing.T) {
	vm := setupGojaVmForBip32Key()
	_, err := vm.RunString(Bip32KeyDeriveScriptTemplate)
	require.NoError(t, err)

	deriveKeyB58Str := vm.Get("deriveKeyB58Str").Export().(string)
	require.EqualValues(t, "xprvA3pian7QYLdwyi9yxXpZVuNpuc78cFGDekYE71G1rJCBRKbnHYv8YyxAaadptJG6hRvUTVMhDXdVRS61afUWE6FYCjNdJ5qetKG7prRGrrq", deriveKeyB58Str)
	deriveWith32KeyB58Str := vm.Get("deriveWith32KeyB58Str").Export().(string)
	require.EqualValues(t, "xprvABK1GV7D4nrPbWbHK7AsSRkWD83x6dCi22mKe9o3TcUS1W9VbPuNLZUQGWaWdiEbWUsNvCTwUCRhLTVea9WUcLRD11GN9gZbAPi2kqBsf6g", deriveWith32KeyB58Str)
}
