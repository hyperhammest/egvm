package types

import (
	"github.com/dop251/goja"
	"github.com/holiman/uint256"

	"github.com/smartbch/egvm/script/utils"
)

type Uint256 struct {
	X *uint256.Int
}

func HexToU256(hex string) Uint256 {
	x, err := uint256.FromHex(hex)
	if err != nil {
		panic(goja.NewSymbol("invalid hex string"))
	}
	return Uint256{X: x}
}

func BufToU256(buf goja.ArrayBuffer) Uint256 {
	return Uint256{X: uint256.NewInt(0).SetBytes(buf.Bytes())}
}

func U256(v uint64) Uint256 {
	if v > utils.MaxSafeInteger {
		panic(utils.LargerThanMaxInteger)
	}
	return Uint256{X: uint256.NewInt(v)}
}

func (u Uint256) ToS256() Sint256 {
	return Sint256{x: u.X.Clone()}
}

func (u Uint256) Incr() Uint256 {
	result, overflow := uint256.NewInt(0).AddOverflow(u.X, uint256.NewInt(1))
	if overflow {
		panic(goja.NewSymbol("overflow in addition"))
	}
	return Uint256{X: result}
}

func (u Uint256) Add(v Uint256) Uint256 {
	result, overflow := uint256.NewInt(0).AddOverflow(u.X, v.X)
	if overflow {
		panic(goja.NewSymbol("overflow in addition"))
	}
	return Uint256{X: result}
}

func (u Uint256) UnsafeAdd(v Uint256) Uint256 {
	result := uint256.NewInt(0).Add(u.X, v.X)
	return Uint256{X: result}
}

func (u Uint256) And(v Uint256) Uint256 {
	result := uint256.NewInt(0).And(u.X, v.X)
	return Uint256{X: result}
}

func (u Uint256) Div(v Uint256) Uint256 {
	if v.X.IsZero() {
		panic(goja.NewSymbol("divide by zero"))
	}
	result := uint256.NewInt(0).Div(u.X, v.X)
	return Uint256{X: result}
}

func (u Uint256) Mod(v Uint256) Uint256 {
	if v.X.IsZero() {
		panic(goja.NewSymbol("divide by zero"))
	}
	result := uint256.NewInt(0).Mod(u.X, v.X)
	return Uint256{X: result}
}

func (u Uint256) DivMod(v Uint256) [2]Uint256 {
	if v.X.IsZero() {
		panic(goja.NewSymbol("divide by zero"))
	}
	quo, rem := uint256.NewInt(0).DivMod(u.X, v.X, uint256.NewInt(0))
	return [2]Uint256{{X: quo}, {X: rem}}
}

func (u Uint256) Exp(v Uint256) Uint256 {
	result := uint256.NewInt(0).Exp(u.X, v.X)
	return Uint256{X: result}
}

func (u Uint256) Gt(v Uint256) bool {
	return u.X.Gt(v.X)
}

func (u Uint256) Gte(v Uint256) bool {
	return !v.X.Gt(u.X)
}

func (u Uint256) GtNum(v uint64) bool {
	if v > utils.MaxSafeInteger {
		panic(utils.LargerThanMaxInteger)
	}
	return u.X.GtUint64(v)
}

func (u Uint256) GteNum(v uint64) bool {
	if v > utils.MaxSafeInteger {
		panic(utils.LargerThanMaxInteger)
	}
	return !U256(v).X.Gt(u.X)
}

func (u Uint256) IsZero() bool {
	return u.X.IsZero()
}

func (u Uint256) Equal(v Uint256) bool {
	return u.X.Eq(v.X)
}

func (u Uint256) Lt(v Uint256) bool {
	return u.X.Lt(v.X)
}

func (u Uint256) Lte(v Uint256) bool {
	return !v.X.Lt(u.X)
}

func (u Uint256) LtNum(v uint64) bool {
	if v > utils.MaxSafeInteger {
		panic(utils.LargerThanMaxInteger)
	}
	return u.X.LtUint64(v)
}

func (u Uint256) LteNum(v uint64) bool {
	if v > utils.MaxSafeInteger {
		panic(utils.LargerThanMaxInteger)
	}
	return !U256(v).X.Lt(u.X)
}

func (u Uint256) Mul(v Uint256) Uint256 {
	result, overflow := uint256.NewInt(0).MulOverflow(u.X, v.X)
	if overflow {
		panic(utils.OverflowInSigned)
	}
	return Uint256{X: result}
}

func (u Uint256) UnsafeMul(v Uint256) Uint256 {
	result := uint256.NewInt(0).Mul(u.X, v.X)
	return Uint256{X: result}
}

func (u Uint256) Not() Uint256 {
	result := uint256.NewInt(0).Not(u.X)
	return Uint256{X: result}
}

func (u Uint256) Or(v Uint256) Uint256 {
	result := uint256.NewInt(0).Or(u.X, v.X)
	return Uint256{X: result}
}

func (u Uint256) Lsh(v uint) Uint256 {
	result := uint256.NewInt(0).Lsh(u.X, v)
	return Uint256{X: result}
}

func (u Uint256) Rsh(v uint) Uint256 {
	result := uint256.NewInt(0).Rsh(u.X, v)
	return Uint256{X: result}
}

func (u Uint256) Sqrt() Uint256 {
	result := uint256.NewInt(0).Sqrt(u.X)
	return Uint256{X: result}
}

func (u Uint256) Sub(v Uint256) Uint256 {
	if u.X.Lt(v.X) {
		panic(goja.NewSymbol("Overflow in substraction"))
	}
	result := uint256.NewInt(0).Sub(u.X, v.X)
	return Uint256{X: result}
}

func (u Uint256) UnsafeSub(v Uint256) Uint256 {
	result := uint256.NewInt(0).Sub(u.X, v.X)
	return Uint256{X: result}
}

func (u Uint256) IsSafeInteger() bool {
	u64, overflow := u.X.Uint64WithOverflow()
	return u64 <= utils.MaxSafeInteger && !overflow
}

func (u Uint256) ToSafeInteger() int64 {
	u64, overflow := u.X.Uint64WithOverflow()
	safe := u64 <= utils.MaxSafeInteger && !overflow
	if !safe {
		panic(goja.NewSymbol("Overflow in ToSafeInteger"))
	}
	return int64(u64)
}

func (u Uint256) ToBuf(f goja.FunctionCall, vm *goja.Runtime) goja.Value {
	if len(f.Arguments) != 0 {
		panic(vm.ToValue("ToBuf has no arguments."))
	}
	var dest [32]byte
	u.X.WriteToArray32(&dest)
	return vm.ToValue(vm.NewArrayBuffer(dest[:]))
}

func (u Uint256) ToHex(f goja.FunctionCall, vm *goja.Runtime) goja.Value {
	if len(f.Arguments) != 0 {
		panic(vm.ToValue("ToHex has no arguments."))
	}
	return vm.ToValue(u.X.Hex())
}

// Stringer interface
func (u Uint256) String() string {
	return u.X.String()
}
