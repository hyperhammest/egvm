package extension

import (
	"fmt"
	"math"

	"github.com/dop251/goja"
	ecies "github.com/ecies/go/v2"
	"github.com/tyler-smith/go-bip32"

	"github.com/smartbch/egvm/script/utils"
	"github.com/smartbch/egvm/keygrantor"
)

// ===============

type Bip32Key struct {
	key *bip32.Key
}

// Only for golang
func NewBip32Key(key *bip32.Key) Bip32Key {
	return Bip32Key{key: key}
}

func B58ToBip32Key(data string) Bip32Key {
	key, err := bip32.B58Deserialize(data)
	if err != nil {
		panic(goja.NewSymbol("error in B58ToBip32Key: " + err.Error()))
	}
	return Bip32Key{key: key}
}

func BufToBip32Key(f goja.FunctionCall, vm *goja.Runtime) goja.Value {
	buf := utils.GetOneArrayBuffer(f)
	key, err := bip32.Deserialize(buf)
	if err != nil {
		panic(goja.NewSymbol("error in BufToBip32Key: " + err.Error()))
	}
	return vm.ToValue(Bip32Key{key: key})
}

func GenerateRandomBip32Key() Bip32Key {
	return Bip32Key{key: keygrantor.GetRandomExtPrivKey()}
}

func (key Bip32Key) B58Serialize() string {
	return key.key.B58Serialize()
}

func (key Bip32Key) NewChildKey(childIdx uint32) Bip32Key {
	k, err := key.key.NewChildKey(childIdx)
	if err != nil {
		panic(goja.NewSymbol("error in NewChildKey: " + err.Error()))
	}
	return Bip32Key{key: k}
}

func (key Bip32Key) PublicKey() Bip32Key {
	return Bip32Key{key: key.key.PublicKey()}
}

func (key Bip32Key) Serialize(f goja.FunctionCall, vm *goja.Runtime) goja.Value {
	bz, _ := key.key.Serialize() // impossible to generate error
	return vm.ToValue(vm.NewArrayBuffer(bz))
}

func (key Bip32Key) IsPrivate() bool {
	return key.key.IsPrivate
}

func (key Bip32Key) ToPrivateKey() PrivateKey {
	if !key.IsPrivate() {
		panic(goja.NewSymbol("It's not private"))
	}
	return PrivateKey{key: ecies.NewPrivateKeyFromBytes(key.key.Key)}
}

func (key Bip32Key) DeriveWithBytes32(f goja.FunctionCall, vm *goja.Runtime) goja.Value {
	var hash [32]byte
	bz := utils.GetOneArrayBuffer(f)
	copy(hash[:], bz)
	bip32Key := keygrantor.DeriveKey(key.key, hash)
	return vm.ToValue(Bip32Key{key: bip32Key})
}

func (key Bip32Key) Derive(f goja.FunctionCall, vm *goja.Runtime) goja.Value {
	if len(f.Arguments) != 5 {
		// BIP-44
		// m / purpose' / coin' / account' / change / address_index
		panic(utils.IncorrectArgumentCount)
	}

	child := key.key
	for i, arg := range f.Arguments {
		n, ok := arg.Export().(int64)
		if !ok || n < 0 || n > math.MaxUint32 {
			panic(goja.NewSymbol(fmt.Sprintf("The argument #%d is not uint32", i)))
		}

		var err error
		child, err = child.NewChildKey(uint32(n))
		if err != nil {
			panic(vm.ToValue("Error in Derive: " + err.Error()))
		}
	}
	return vm.ToValue(Bip32Key{key: child})
}
