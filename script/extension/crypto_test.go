package extension

import (
	"testing"

	"github.com/dop251/goja"
	gethcmn "github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
)

const (
	AesGcmScriptTemplate = `
		const secret = new Uint8Array([1, 2, 3, 4, 8, 7, 6, 5, 3, 4, 5, 6, 91, 33, 255, 32]).buffer
		const nonce = new Uint8Array([1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16]).buffer
		const message = new Uint8Array([8, 8, 8, 8, 8, 8, 8, 8, 8, 8, 8, 8, 8, 8, 8, 8]).buffer
		const encryptedMsg = AesGcmEncrypt(secret, nonce, message)
		const [decryptedMsg, ok] = AesGcmDecrypt(secret, nonce, encryptedMsg)
	`

	PrivateKeyAndPublicKeyBasicScriptTemplate = `
		const privateKeyHex = 'c9cb992b13141bb3326d028020030f33b92ea9a64b6530291e7876938bd31479'
		const privateKeyBuf = HexToBuf(privateKeyHex)
		const privateKey = BufToPrivateKey(privateKeyBuf)
		
		const privateKeyToBuf = privateKey.Serialize()
		const privateKeyToHex = privateKey.Hex()

		const publicKey = privateKey.GetPublicKey()
		const publicKeyToCompressedBz = publicKey.SerializeCompressed()
		const publicKeyToUncompressedBz = publicKey.SerializeUncompressed()
	`

	PrivateKeyAndPublicKeyEncryptionScriptTemplate = `
		const privateKeyHex = 'c9cb992b13141bb3326d028020030f33b92ea9a64b6530291e7876938bd31479'
		const privateKeyBuf = HexToBuf(privateKeyHex)
		const privateKey = BufToPrivateKey(privateKeyBuf)

		const publicKey = privateKey.GetPublicKey()

		const message = new Uint8Array([8, 8, 8, 8, 8, 8, 8, 8, 8, 8, 8, 8, 8, 8, 8, 8]).buffer
		const nonce = new Uint8Array([1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16]).buffer
		const encryptedMsg = publicKey.Encrypt(message, nonce)

		const [decryptedMsg, ok] = privateKey.Decrypt(encryptedMsg)
	`

	PrivateKeyAndPublicKeySecretScriptTemplate = `
		const privateKeyHex1 = 'c9cb992b13141bb3326d028020030f33b92ea9a64b6530291e7876938bd31479'
		const privateKeyBuf1 = HexToBuf(privateKeyHex1)
		const privateKey1 = BufToPrivateKey(privateKeyBuf1)

		const privateKeyHex2 = '8cc74af1645ba43b2fa8e43de48773330589799260552d16081d985a8419a565'
		const privateKeyBuf2 = HexToBuf(privateKeyHex2)
		const privateKey2 = BufToPrivateKey(privateKeyBuf2)

		const publicKey1 = privateKey1.GetPublicKey()
		const publicKey2 = privateKey2.GetPublicKey()

		const sk1 = privateKey1.Encapsulate(publicKey2)
		const sk2 = publicKey1.Decapsulate(privateKey2)
	`

	PrivateKeyAndPublicKeyAddressScriptTemplate = `
		const privateKeyHex1 = 'c9cb992b13141bb3326d028020030f33b92ea9a64b6530291e7876938bd31479'
		const privateKeyBuf1 = HexToBuf(privateKeyHex1)
		const privateKey1 = BufToPrivateKey(privateKeyBuf1)

		const publicKey1 = privateKey1.GetPublicKey()

		const evmAddress = publicKey1.ToEvmAddress()
		const cashAddress = publicKey1.ToCashAddress()
	`

	VrfScriptTemplate = `
		const privateKeyHex1 = 'c9cb992b13141bb3326d028020030f33b92ea9a64b6530291e7876938bd31479'
		const privateKeyBuf1 = HexToBuf(privateKeyHex1)
		const privateKey1 = BufToPrivateKey(privateKeyBuf1)
		
		const alpha = new Uint8Array([1, 2, 3, 4, 5, 6, 7, 8]).buffer
		const [beta1, pi] = privateKey1.VrfProve(alpha)

		const publicKey1 = privateKey1.GetPublicKey()
		const beta2 = publicKey1.VrfVerify(alpha, pi)
	`

	SignatureScriptTemplate = `
		const privateKeyHex1 = 'c9cb992b13141bb3326d028020030f33b92ea9a64b6530291e7876938bd31479'
		const privateKeyBuf1 = HexToBuf(privateKeyHex1)
		const privateKey1 = BufToPrivateKey(privateKeyBuf1)

		const msg = new Uint8Array([1, 2, 3, 4, 5, 6, 7, 8]).buffer
		const digestHash = Keccak256(msg)
		const sig = privateKey1.Sign(digestHash)

		const publicKey1 = privateKey1.GetPublicKey()
		const publicKey1ToCompressedBz = publicKey1.SerializeCompressed()
		const ok = VerifySignature(publicKey1ToCompressedBz, digestHash, sig)

		const recoverPubKey = Ecrecover(digestHash, sig)
		const recoverAddress = recoverPubKey.ToEvmAddress()
		const address1 = publicKey1.ToEvmAddress()
	`

	GetEthSignedMessageScriptTemplate = `
		const msg = new Uint8Array([1, 2, 3, 4, 5, 6, 7, 8]).buffer
		const ethMsg = GetEthSignedMessage(msg)
	`
)

func setupGojaVmForCrypto() *goja.Runtime {
	vm := goja.New()
	vm.Set("HexToBuf", HexToBuf)
	vm.Set("AesGcmDecrypt", AesGcmDecrypt)
	vm.Set("AesGcmEncrypt", AesGcmEncrypt)
	vm.Set("BufToPrivateKey", BufToPrivateKey)
	vm.Set("BufToPublicKey", BufToPublicKey)
	vm.Set("Keccak256", Keccak256)
	vm.Set("GetEthSignedMessage", GetEthSignedMessage)
	vm.Set("VerifySignature", VerifySignature)
	vm.Set("Ecrecover", Ecrecover)
	return vm
}

func TestAesGcm(t *testing.T) {
	vm := setupGojaVmForCrypto()
	_, err := vm.RunString(AesGcmScriptTemplate)
	require.NoError(t, err)

	encryptedMsg := vm.Get("encryptedMsg").Export().(goja.ArrayBuffer)
	encryptedMsgHex := gethcmn.Bytes2Hex(encryptedMsg.Bytes())
	decryptedMsg := vm.Get("decryptedMsg").Export().(goja.ArrayBuffer)
	decryptedMsgHex := gethcmn.Bytes2Hex(decryptedMsg.Bytes())
	ok := vm.Get("ok").Export().(bool)
	require.True(t, ok)
	require.EqualValues(t, "fe5695ffe3ce8fb998bb9e35f529507b2c990e7de5c5d99e32d1f546130c6ef4", encryptedMsgHex)
	require.EqualValues(t, "08080808080808080808080808080808", decryptedMsgHex)

}

func TestPrivateKeyAndPublicKeyBasic(t *testing.T) {
	vm := setupGojaVmForCrypto()
	_, err := vm.RunString(PrivateKeyAndPublicKeyBasicScriptTemplate)
	require.NoError(t, err)

	privateKeyToBuf := vm.Get("privateKeyToBuf").Export().(goja.ArrayBuffer)
	privateKeyToBufHex := gethcmn.Bytes2Hex(privateKeyToBuf.Bytes())
	privateKeyToHex := vm.Get("privateKeyToHex").Export().(string)
	require.EqualValues(t, "c9cb992b13141bb3326d028020030f33b92ea9a64b6530291e7876938bd31479", privateKeyToBufHex)
	require.EqualValues(t, "c9cb992b13141bb3326d028020030f33b92ea9a64b6530291e7876938bd31479", privateKeyToHex)

	publicKeyToCompressedBz := vm.Get("publicKeyToCompressedBz").Export().(goja.ArrayBuffer)
	publicKeyToCompressedBzHex := gethcmn.Bytes2Hex(publicKeyToCompressedBz.Bytes())
	publicKeyToUncompressedBz := vm.Get("publicKeyToUncompressedBz").Export().(goja.ArrayBuffer)
	publicKeyToUncompressedBzHex := gethcmn.Bytes2Hex(publicKeyToUncompressedBz.Bytes())
	require.EqualValues(t, "02dde6c067f5e1a641dedab654cbbd9c3b4c6f8adbf2aeb17c6500319d2c08f08e", publicKeyToCompressedBzHex)
	require.EqualValues(t, "04dde6c067f5e1a641dedab654cbbd9c3b4c6f8adbf2aeb17c6500319d2c08f08ec14990f4beca82491ca7134e7c637a35073df4e4ef8ee4251f5b4570f222777a",
		publicKeyToUncompressedBzHex)
}

func TestPrivateKeyAndPublicKeyEncryption(t *testing.T) {
	vm := setupGojaVmForCrypto()
	_, err := vm.RunString(PrivateKeyAndPublicKeyEncryptionScriptTemplate)
	require.NoError(t, err)

	encryptedMsg := vm.Get("encryptedMsg").Export().(goja.ArrayBuffer)
	encryptedMsgHex := gethcmn.Bytes2Hex(encryptedMsg.Bytes())
	decryptedMsg := vm.Get("decryptedMsg").Export().(goja.ArrayBuffer)
	decryptedMsgHex := gethcmn.Bytes2Hex(decryptedMsg.Bytes())
	ok := vm.Get("ok").Export().(bool)
	require.True(t, ok)
	require.EqualValues(t, "040447c2dcdd17d5a65b95a697304212da4b8064486d0862eb66090632bb320f71286f9d8fa99affd10be414bec4be093a0a50c70e78fe4d439ef8e1a8762a87010102030405060708090a0b0c0d0e0f10f1bf1dd69b1e370effb15b45d0c9814429fbaaad106c9d4e0fa5bdef4c24ff6d",
		encryptedMsgHex)
	require.EqualValues(t, "08080808080808080808080808080808", decryptedMsgHex)
}

func TestPrivateKeyAndPublicKeySecret(t *testing.T) {
	vm := setupGojaVmForCrypto()
	_, err := vm.RunString(PrivateKeyAndPublicKeySecretScriptTemplate)
	require.NoError(t, err)

	sk1 := vm.Get("sk1").Export().(goja.ArrayBuffer)
	sk1Hex := gethcmn.Bytes2Hex(sk1.Bytes())
	sk2 := vm.Get("sk2").Export().(goja.ArrayBuffer)
	sk2Hex := gethcmn.Bytes2Hex(sk2.Bytes())
	require.EqualValues(t, sk1Hex, sk2Hex)
}

func TestPrivateKeyAndPublicKeyAddress(t *testing.T) {
	vm := setupGojaVmForCrypto()
	_, err := vm.RunString(PrivateKeyAndPublicKeyAddressScriptTemplate)
	require.NoError(t, err)

	evmAddress := vm.Get("evmAddress").Export().(goja.ArrayBuffer)
	evmAddressHex := gethcmn.Bytes2Hex(evmAddress.Bytes())
	cashAddress := vm.Get("cashAddress").Export().(goja.ArrayBuffer)
	cashAddressHex := gethcmn.Bytes2Hex(cashAddress.Bytes())
	require.EqualValues(t, "1df57cdcfdaffc7595f5104bb2732e6b98dd58c0", evmAddressHex)
	require.True(t, gethcmn.IsHexAddress(evmAddressHex))
	require.EqualValues(t, "f3101cbb3bed3115fbbe06cb15134f7ce2d89489", cashAddressHex)
}

func TestVrf(t *testing.T) {
	vm := setupGojaVmForCrypto()
	_, err := vm.RunString(VrfScriptTemplate)
	require.NoError(t, err)

	beta1 := vm.Get("beta1").Export().(goja.ArrayBuffer)
	beta1Hex := gethcmn.Bytes2Hex(beta1.Bytes())
	beta2 := vm.Get("beta2").Export().(goja.ArrayBuffer)
	beta2Hex := gethcmn.Bytes2Hex(beta2.Bytes())
	require.EqualValues(t, beta2Hex, beta1Hex)
}

func TestSignature(t *testing.T) {
	vm := setupGojaVmForCrypto()
	_, err := vm.RunString(SignatureScriptTemplate)
	require.NoError(t, err)

	ok := vm.Get("ok").Export().(bool)
	require.True(t, ok)

	recoverAddress := vm.Get("recoverAddress").Export().(goja.ArrayBuffer)
	recoverAddressHex := gethcmn.Bytes2Hex(recoverAddress.Bytes())
	address1 := vm.Get("address1").Export().(goja.ArrayBuffer)
	address1Hex := gethcmn.Bytes2Hex(address1.Bytes())
	require.EqualValues(t, address1Hex, recoverAddressHex)
}

func TestEthSignedMessage(t *testing.T) {
	vm := setupGojaVmForCrypto()
	_, err := vm.RunString(GetEthSignedMessageScriptTemplate)
	require.NoError(t, err)

	ethMsg := vm.Get("ethMsg").Export().(goja.ArrayBuffer)
	ethMsgStr := string(ethMsg.Bytes())
	require.EqualValues(t, "\x19Ethereum Signed Message:\n8\x01\x02\x03\x04\x05\x06\a\b", ethMsgStr)
}
