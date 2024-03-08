package extension

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/ecdsa"
	"crypto/sha256"
	"fmt"

	"github.com/dop251/goja"
	ecies "github.com/ecies/go/v2"
	gethcrypto "github.com/ethereum/go-ethereum/crypto"
	"github.com/gcash/bchutil"
	"github.com/vechain/go-ecvrf"

	"github.com/smartbch/egvm/script/utils"
)

// --------- AES ---------

func aesgcmEncrypt(secret, nonce, msg []byte) ([]byte, error) {
	block, err := aes.NewCipher(secret)
	if err != nil {
		return nil, fmt.Errorf("cannot create new aes block: %w", err)
	}

	aesgcm, err := cipher.NewGCMWithNonceSize(block, 16)
	if err != nil {
		return nil, fmt.Errorf("cannot create aes gcm: %w", err)
	}

	return aesgcm.Seal(nil, nonce, msg, nil), nil
}

func aesgcmDecrypt(secret, nonce, ciphertext []byte) ([]byte, error) {
	block, err := aes.NewCipher(secret)
	if err != nil {
		return nil, fmt.Errorf("cannot create new aes block: %w", err)
	}

	gcm, err := cipher.NewGCMWithNonceSize(block, 16)
	if err != nil {
		return nil, fmt.Errorf("cannot create gcm cipher: %w", err)
	}

	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, fmt.Errorf("cannot decrypt ciphertext: %w", err)
	}

	return plaintext, nil
}

func AesGcmEncrypt(f goja.FunctionCall, vm *goja.Runtime) goja.Value {
	secret, nonce, msg := utils.GetThreeArrayBuffers(f)
	bz, err := aesgcmEncrypt(secret, nonce, msg)
	if err != nil {
		panic(goja.NewSymbol("error in AesGcmEncrypt: " + err.Error()))
	}
	return vm.ToValue(vm.NewArrayBuffer(bz))
}

func AesGcmDecrypt(f goja.FunctionCall, vm *goja.Runtime) goja.Value {
	secret, nonce, ciphertext := utils.GetThreeArrayBuffers(f)
	bz, err := aesgcmDecrypt(secret, nonce, ciphertext)
	if err != nil {
		return vm.ToValue([2]any{nil, false})
	}
	return vm.ToValue([2]any{vm.NewArrayBuffer(bz), true})
}

// --------- Public-Key Cryptography ---------

type PrivateKey struct {
	key *ecies.PrivateKey
}

type PublicKey struct {
	key *ecies.PublicKey
}

func BufToPrivateKey(buf goja.ArrayBuffer) PrivateKey {
	key := ecies.NewPrivateKeyFromBytes(buf.Bytes())
	return PrivateKey{key: key}
}

func BufToPublicKey(buf goja.ArrayBuffer) PublicKey {
	key, _ := ecies.NewPublicKeyFromBytes(buf.Bytes())
	return PublicKey{key: key}
}

func (prv PrivateKey) GetPublicKey() PublicKey {
	return PublicKey{key: prv.key.PublicKey}
}

func (prv PrivateKey) ECDH(f goja.FunctionCall, vm *goja.Runtime) goja.Value {
	if len(f.Arguments) != 1 {
		panic(utils.IncorrectArgumentCount)
	}
	pub, ok := f.Arguments[0].Export().(PublicKey)
	if !ok {
		panic(goja.NewSymbol("The first argument must be PublicKey"))
	}
	bz, err := prv.key.ECDH(pub.key)
	if err != nil {
		panic(goja.NewSymbol("error in ECDH: " + err.Error()))
	}
	return vm.ToValue(vm.NewArrayBuffer(bz))
}

func (prv PrivateKey) Encapsulate(f goja.FunctionCall, vm *goja.Runtime) goja.Value {
	if len(f.Arguments) != 1 {
		panic(utils.IncorrectArgumentCount)
	}
	pub, ok := f.Arguments[0].Export().(PublicKey)
	if !ok {
		panic(goja.NewSymbol("The first argument must be PublicKey"))
	}
	bz, err := prv.key.Encapsulate(pub.key)
	if err != nil {
		panic(goja.NewSymbol("error in Encapsulate: " + err.Error()))
	}
	return vm.ToValue(vm.NewArrayBuffer(bz))
}

func (prv PrivateKey) toECDSA() *ecdsa.PrivateKey {
	return &ecdsa.PrivateKey{
		D: prv.key.D,
		PublicKey: ecdsa.PublicKey{
			Curve: prv.key.PublicKey.Curve,
			X:     prv.key.PublicKey.X,
			Y:     prv.key.PublicKey.Y,
		},
	}
}

func (prv PrivateKey) Sign(f goja.FunctionCall, vm *goja.Runtime) goja.Value {
	msg := utils.GetOneArrayBuffer(f)
	ecdsaPrv := prv.toECDSA()
	sig, err := gethcrypto.Sign(msg, ecdsaPrv)
	if err != nil {
		panic(goja.NewSymbol("error in Sign: " + err.Error()))
	}
	return vm.ToValue(vm.NewArrayBuffer(sig))
}

func (prv PrivateKey) Equal(other PrivateKey) bool {
	return prv.key.Equals(other.key)
}

func (prv PrivateKey) Hex() string {
	return prv.key.Hex()
}

func (prv PrivateKey) Serialize(f goja.FunctionCall, vm *goja.Runtime) goja.Value {
	if len(f.Arguments) != 0 {
		panic(utils.IncorrectArgumentCount)
	}
	return vm.ToValue(vm.NewArrayBuffer(prv.key.Bytes()))
}

func (prv PrivateKey) Decrypt(f goja.FunctionCall, vm *goja.Runtime) goja.Value {
	buf := utils.GetOneArrayBuffer(f)
	bz, err := ecies.Decrypt(prv.key, buf)
	if err != nil {
		return vm.ToValue([2]any{nil, false})
	}
	return vm.ToValue([2]any{vm.NewArrayBuffer(bz), true})
}

func (pub PublicKey) Decapsulate(f goja.FunctionCall, vm *goja.Runtime) goja.Value {
	if len(f.Arguments) != 1 {
		panic(utils.IncorrectArgumentCount)
	}
	prv, ok := f.Arguments[0].Export().(PrivateKey)
	if !ok {
		panic(goja.NewSymbol("The first argument must be PrivateKey"))
	}
	bz, err := pub.key.Decapsulate(prv.key)
	if err != nil {
		panic(goja.NewSymbol("error in Decapsulate: " + err.Error()))
	}
	return vm.ToValue(vm.NewArrayBuffer(bz))
}

func (pub PublicKey) Equal(other PublicKey) bool {
	return pub.key.Equals(other.key)
}

func (pub PublicKey) Hex(compressed bool) string {
	return pub.key.Hex(compressed)
}

func (pub PublicKey) toBuf(f goja.FunctionCall, vm *goja.Runtime, compressed bool) goja.Value {
	if len(f.Arguments) != 0 {
		panic(utils.IncorrectArgumentCount)
	}
	return vm.ToValue(vm.NewArrayBuffer(pub.key.Bytes(compressed)))
}

func (pub PublicKey) SerializeCompressed(f goja.FunctionCall, vm *goja.Runtime) goja.Value {
	return pub.toBuf(f, vm, true)
}

func (pub PublicKey) SerializeUncompressed(f goja.FunctionCall, vm *goja.Runtime) goja.Value {
	return pub.toBuf(f, vm, false)
}

func (pub PublicKey) Encrypt(f goja.FunctionCall, vm *goja.Runtime) goja.Value {
	buf, entropy := utils.GetTwoArrayBuffers(f)
	bz, err := eciesEncrypt(pub.key, buf, entropy)
	if err != nil {
		panic(goja.NewSymbol("error in Encrypt: " + err.Error()))
	}
	return vm.ToValue(vm.NewArrayBuffer(bz))
}

func eciesEncrypt(pubkey *ecies.PublicKey, msg, entropy []byte) ([]byte, error) {
	var ct bytes.Buffer

	// Generate ephemeral key
	h := sha256.New()
	h.Write(entropy)
	h.Write(msg)
	ek := ecies.NewPrivateKeyFromBytes(h.Sum(nil))

	ct.Write(ek.PublicKey.Bytes(false))

	// Derive shared secret
	secret, err := ek.Encapsulate(pubkey)
	if err != nil {
		return nil, err
	}

	nonce := h.Sum(entropy)[:16]

	ct.Write(nonce)

	ciphertext, err := aesgcmEncrypt(secret, nonce, msg)
	if err != nil {
		return nil, err
	}

	tag := ciphertext[len(ciphertext)-16:]
	ct.Write(tag)
	ciphertext = ciphertext[:len(ciphertext)-len(tag)]
	ct.Write(ciphertext)

	return ct.Bytes(), nil
}

func (pub PublicKey) toECDSA() *ecdsa.PublicKey {
	return &ecdsa.PublicKey{
		Curve: pub.key.Curve,
		X:     pub.key.X,
		Y:     pub.key.Y,
	}
}

func (pub PublicKey) ToEvmAddress(f goja.FunctionCall, vm *goja.Runtime) goja.Value {
	if len(f.Arguments) != 0 {
		panic(utils.IncorrectArgumentCount)
	}
	key := pub.toECDSA()
	addr := gethcrypto.PubkeyToAddress(*key)
	return vm.ToValue(vm.NewArrayBuffer(addr[:]))
}

func (pub PublicKey) ToCashAddress(f goja.FunctionCall, vm *goja.Runtime) goja.Value {
	if len(f.Arguments) != 0 {
		panic(utils.IncorrectArgumentCount)
	}
	pubKeyHash := bchutil.Hash160(pub.key.Bytes(true))
	return vm.ToValue(vm.NewArrayBuffer(pubKeyHash[:]))
}

func (prv PrivateKey) VrfProve(f goja.FunctionCall, vm *goja.Runtime) goja.Value {
	alpha := utils.GetOneArrayBuffer(f)
	beta, pi, err := ecvrf.Secp256k1Sha256Tai.Prove(prv.toECDSA(), alpha)
	if err != nil {
		panic(goja.NewSymbol("error in VrfProve: " + err.Error()))
	}
	return vm.ToValue([2]goja.ArrayBuffer{vm.NewArrayBuffer(beta), vm.NewArrayBuffer(pi)})
}

func (pub PublicKey) VrfVerify(f goja.FunctionCall, vm *goja.Runtime) goja.Value {
	alpha, pi := utils.GetTwoArrayBuffers(f)
	beta, err := ecvrf.Secp256k1Sha256Tai.Verify(pub.toECDSA(), alpha, pi)
	if err != nil {
		panic(goja.NewSymbol("error in VrfVerify: " + err.Error()))
	}
	return vm.ToValue(vm.NewArrayBuffer(beta[:]))
}

// --------- Signature ---------

func GetEthSignedMessage(f goja.FunctionCall, vm *goja.Runtime) goja.Value {
	msg := utils.GetOneArrayBuffer(f)
	ethMsg := fmt.Sprintf("\x19Ethereum Signed Message:\n%d%s", len(msg[:]), msg[:])
	return vm.ToValue(vm.NewArrayBuffer([]byte(ethMsg)))
}

func VerifySignature(f goja.FunctionCall, vm *goja.Runtime) goja.Value {
	pubkey, digestHash, signature := utils.GetThreeArrayBuffers(f)
	if len(pubkey) != 33 && len(pubkey) != 65 {
		panic(goja.NewSymbol("invalid public key length, must be 33 or 65"))
	}

	if len(digestHash) != 32 {
		panic(goja.NewSymbol("invalid digestHash length, must be 32"))
	}

	if len(signature) != 65 {
		panic(goja.NewSymbol("invalid sig length, must be 65"))
	}
	result := gethcrypto.VerifySignature(pubkey, digestHash, signature[:64])
	return vm.ToValue(result)
}

func Ecrecover(f goja.FunctionCall, vm *goja.Runtime) goja.Value {
	hash, sig := utils.GetTwoArrayBuffers(f)
	if len(sig) != 65 {
		panic(goja.NewSymbol("invalid sig length, must be 65"))
	}

	bz, err := gethcrypto.Ecrecover(hash, sig)
	if err != nil {
		panic(goja.NewSymbol("error in Ecrecover: " + err.Error()))
	}

	key, err := ecies.NewPublicKeyFromBytes(bz)
	if err != nil {
		panic(goja.NewSymbol("error in NewPublicKeyFromBytes: " + err.Error()))
	}
	return vm.ToValue(PublicKey{key: key})
}
