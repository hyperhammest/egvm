package types

import (
	"github.com/dop251/goja"
	"github.com/holiman/uint256"

	"github.com/smartbch/egvm/script/utils"
)

type Sint256 struct {
	x *uint256.Int
}

func HexToS256(hex string) Sint256 {
	x, err := uint256.FromHex(hex)
	if err != nil {
		panic(goja.NewSymbol("invalid hex string"))
	}
	return Sint256{x: x}
}

func BufToS256(buf goja.ArrayBuffer) Sint256 {
	return Sint256{x: uint256.NewInt(0).SetBytes(buf.Bytes())}
}

func S256(v int64) Sint256 {
	if v > int64(utils.MaxSafeInteger) || -v > int64(utils.MaxSafeInteger) {
		panic(utils.LargerThanMaxInteger)
	}
	if v >= 0 {
		return Sint256{x: uint256.NewInt(uint64(v))}
	}
	tmp := uint256.NewInt(uint64(-v))
	return Sint256{x: uint256.NewInt(0).Neg(tmp)}
}

func (s Sint256) ToU256() Uint256 {
	return Uint256{X: s.x.Clone()}
}

func (s Sint256) Abs() Sint256 {
	if s.x.Eq(utils.MinNegValue) {
		panic(goja.NewSymbol("overflow in Abs"))
	}
	result := uint256.NewInt(0).Abs(s.x)
	return Sint256{x: result}
}

func (s Sint256) Neg() Sint256 {
	if s.x.Eq(utils.MinNegValue) {
		panic(goja.NewSymbol("overflow in Neg"))
	}
	result := uint256.NewInt(0).Neg(s.x)
	return Sint256{x: result}
}

func (s Sint256) Add(v Sint256) Sint256 {
	result, overflow := (*uint256.Int)(nil), false
	if s.x.Sign() >= 0 && v.x.Sign() >= 0 {
		result, overflow = uint256.NewInt(0).AddOverflow(s.x, v.x)
	} else {
		result = uint256.NewInt(0).Add(s.x, v.x)
		overflow = s.x.Sign() < 0 && v.x.Sign() < 0 && result.Sign() >= 0
	}
	if overflow {
		panic(goja.NewSymbol("overflow in signed addition"))
	}
	return Sint256{x: result}
}

func (s Sint256) Equal(v Sint256) bool {
	return s.x.Eq(v.x)
}

func (s Sint256) IsZero() bool {
	return s.x.IsZero()
}

func (s Sint256) Sign() int {
	return s.x.Sign()
}

func (s Sint256) Sub(v Sint256) Sint256 {
	result := uint256.NewInt(0).Sub(s.x, v.x)
	uSign, vSign, rSign := s.x.Sign(), v.x.Sign(), result.Sign()
	overflow := (uSign > 0 && vSign < 0 && rSign < 0) || // pos - neg < 0
		(uSign < 0 && vSign > 0 && rSign > 0) // neg - pos > 0
	if overflow {
		panic(goja.NewSymbol("overflow in signed substraction"))
	}
	return Sint256{x: result}
}

func (s Sint256) Mul(v Sint256) Sint256 {
	if s.x.Eq(utils.MinNegValue) || v.x.Eq(utils.MinNegValue) {
		panic(goja.NewSymbol("overflow in signed multiplication"))
	}
	p, q := uint256.NewInt(0).Abs(s.x), uint256.NewInt(0).Abs(v.x)
	result, overflow := uint256.NewInt(0).MulOverflow(p, q)
	if overflow {
		panic(goja.NewSymbol("overflow in signed multiplication"))
	}
	if s.x.Sign()*v.x.Sign() == -1 {
		result = uint256.NewInt(0).Neg(result)
	}
	return Sint256{x: result}
}

func (s Sint256) Div(v Sint256) Sint256 {
	if s.x.Eq(utils.MinNegValue) || v.x.Eq(utils.MinNegValue) {
		panic(goja.NewSymbol("overflow in division"))
	}
	if v.x.IsZero() {
		panic(goja.NewSymbol("divide by zero"))
	}
	p, q := uint256.NewInt(0).Abs(s.x), uint256.NewInt(0).Abs(v.x)
	result := uint256.NewInt(0).Div(p, q)
	if s.x.Sign()*v.x.Sign() == -1 {
		result = uint256.NewInt(0).Neg(result)
	}
	return Sint256{x: result}
}

func (s Sint256) Lsh(v uint) Sint256 {
	result := uint256.NewInt(0).Lsh(s.x, v)
	return Sint256{x: result}
}

func (s Sint256) Rsh(v uint) Sint256 {
	result := uint256.NewInt(0).SRsh(s.x, v)
	return Sint256{x: result}
}

func (s Sint256) Gt(v Sint256) bool {
	return s.x.Sgt(v.x)
}

func (s Sint256) Gte(v Sint256) bool {
	return !v.x.Sgt(s.x)
}

func (s Sint256) GtNum(v int64) bool {
	if v > int64(utils.MaxSafeInteger) || -v > int64(utils.MaxSafeInteger) {
		panic(utils.LargerThanMaxInteger)
	}
	return s.x.Sgt(S256(v).x)
}

func (s Sint256) GteNum(v int64) bool {
	if v > int64(utils.MaxSafeInteger) || -v > int64(utils.MaxSafeInteger) {
		panic(utils.LargerThanMaxInteger)
	}
	return !S256(v).x.Sgt(s.x)
}

func (s Sint256) Lt(v Sint256) bool {
	return s.x.Slt(v.x)
}

func (s Sint256) Lte(v Sint256) bool {
	return !v.x.Slt(s.x)
}

func (s Sint256) LtNum(v int64) bool {
	if v > int64(utils.MaxSafeInteger) || -v > int64(utils.MaxSafeInteger) {
		panic(utils.LargerThanMaxInteger)
	}
	return s.x.Slt(S256(v).x)
}

func (s Sint256) LteNum(v int64) bool {
	if v > int64(utils.MaxSafeInteger) {
		panic(utils.LargerThanMaxInteger)
	}
	return !S256(v).x.Slt(s.x)
}

func (s Sint256) IsSafeInteger() bool {
	u64, overflow := uint64(0), false
	if s.x.Sign() < 0 {
		y := uint256.NewInt(0).Neg(s.x)
		u64, overflow = y.Uint64WithOverflow()
	}
	u64, overflow = s.x.Uint64WithOverflow()
	return u64 <= utils.MaxSafeInteger && !overflow
}

func (s Sint256) ToSafeInteger() int64 {
	u64, overflow := uint64(0), false
	if s.x.Sign() < 0 {
		y := uint256.NewInt(0).Neg(s.x)
		u64, overflow = y.Uint64WithOverflow()
	}
	u64, overflow = s.x.Uint64WithOverflow()
	safe := u64 <= utils.MaxSafeInteger && !overflow
	if !safe {
		panic(goja.NewSymbol("Overflow in ToSafeInteger"))
	}
	if s.x.Sign() < 0 {
		return -int64(u64)
	}
	return int64(u64)
}

func (s Sint256) ToBuf(f goja.FunctionCall, vm *goja.Runtime) goja.Value {
	if len(f.Arguments) != 0 {
		panic(vm.ToValue("ToBuf has no arguments."))
	}
	var dest [32]byte
	s.x.WriteToArray32(&dest)
	return vm.ToValue(vm.NewArrayBuffer(dest[:]))
}

func (s Sint256) ToHex(f goja.FunctionCall, vm *goja.Runtime) goja.Value {
	if len(f.Arguments) != 0 {
		panic(vm.ToValue("ToHex has no arguments."))
	}
	return vm.ToValue(s.x.Hex())
}

// Stringer interface
func (s Sint256) String() string {
	if s.x.Sign() < 0 {
		return "-" + uint256.NewInt(0).Neg(s.x).String()
	}
	return s.x.String()
}
