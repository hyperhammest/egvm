package types

import (
	"fmt"
	"testing"

	"github.com/dop251/goja"
	gethcmn "github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"

	"github.com/smartbch/egvm/script/utils"
)

const (
	U256ScriptTemplate = `
		const a = U256(%v)
		a.String()
	`

	HexToU256ScriptTemplate = `
		const b = HexToU256('%v')
		b.String()
	`

	BufToU256ScriptTemplate = `
		let buffer = new ArrayBuffer(32); // 32 bytes
		let view = new Uint8Array(buffer);
		view[31] = %v
		const c = BufToU256(buffer)
		c.String()
	`

	ToS256ScriptTemplate = `
		const d = U256(%v)
		const sd = d.ToS256()
		sd.String()
	`

	U256AddScriptTemplate = `
		const a = U256(%v)
		const b = U256(%v)
		const ab = a.Add(b)
		ab.String()
	`

	U256UnsafeAddScriptTemplate = `
		const c = U256(%v)
		const d = U256(%v)
		const cd = c.UnsafeAdd(d)
		cd.String()
	`

	U256DivScriptTemplate = `
		const a = U256(%v)
		const b = U256(%v)
		const ab = a.Div(b)
		ab.String()
	`

	U256ModScriptTemplate = `
		const a = U256(%v)
		const b = U256(%v)
		const ab = a.Mod(b)
		ab.String()
	`

	U256DivModScriptTemplate = `
		const a = U256(%v)
		const b = U256(%v)
		const [z, m] = a.DivMod(b)
	`

	U256ExpScriptTemplate = `
		const a = U256(%v)
		const b = U256(%v)
		const ab = a.Exp(b)
		ab.String()
	`

	U256MulScriptTemplate = `
		const a = U256(%v)
		const b = U256(%v)
		const ab = a.Mul(b)
		ab.String()
	`

	U256UnsafeMulScriptTemplate = `
		const c = U256(%v)
		const d = U256(%v)
		const cd = c.UnsafeMul(d)
		cd.String()
	`

	U256AndScriptTemplate = `
		const a = HexToU256('%v')
		const b = HexToU256('%v')
		const ab = a.And(b)
		ab.String()
	`

	U256OrScriptTemplate = `
		const a = HexToU256('%v')
		const b = HexToU256('%v')
		const ab = a.Or(b)
		ab.String()
	`

	U256NotScriptTemplate = `
		const a = HexToU256('%v')
		const notA = a.Not()
		notA.String()
	`

	U256CompareScriptTemplate = `
		const a = U256(%v)
		const b = U256(%v)
		const gt = a.Gt(b)
		const gte = a.Gte(b)
		const lt = a.Lt(b)
		const lte = a.Lte(b)
		const eq = a.Equal(b)
		const isZero = a.IsZero()
	`

	U256CompareNumScriptTemplate = `
		const a = U256(%v)
		const b = 10
		const gtNum = a.GtNum(b)
		const gteNum = a.GteNum(b)
		const ltNum = a.LtNum(b)
		const lteNum = a.LteNum(b)
	`

	U256SqrtScriptTemplate = `
		const a = U256(%v)
		const sq = a.Sqrt()
		sq.String()
	`

	U256SubScriptTemplate = `
		const a = U256(%v)
		const b = U256(%v)
		const ab = a.Sub(b)
		ab.String()
	`

	U256UnsafeSubScriptTemplate = `
		const c = U256(%v)
		const d = U256(%v)
		const cd = c.UnsafeSub(d)
		cd.String()
	`

	U256ShiftScriptTemplate = `
		const a = HexToU256('%v')
		const l = %v
		const r = %v
		const lsh = a.Lsh(l)
		const rsh = a.Rsh(r)
	`

	U256SafeIntegerScriptTemplate = `
		const a = U256(%v)
		const isSafe = a.IsSafeInteger()
		const toSafe = a.ToSafeInteger()
	`

	U256ConversionScriptTemplate = `
		const a = U256(%v)
		const toBuf = a.ToBuf()
		const toHex = a.ToHex()
		const toString = a.String()
	`
)

func setupGojaVmForU256() *goja.Runtime {
	vm := goja.New()

	vm.Set("U256", U256)
	vm.Set("HexToU256", HexToU256)
	vm.Set("BufToU256", BufToU256)
	return vm
}

func TestCreateU256(t *testing.T) {
	vm := setupGojaVmForU256()

	// U256
	v1, err := vm.RunString(fmt.Sprintf(U256ScriptTemplate, 1))
	require.NoError(t, err)
	require.EqualValues(t, "0x1", v1.String())

	// HexToU256
	v2, err := vm.RunString(fmt.Sprintf(HexToU256ScriptTemplate, "0x2"))
	require.NoError(t, err)
	require.EqualValues(t, "0x2", v2.String())

	// BufToU256
	v3, err := vm.RunString(fmt.Sprintf(BufToU256ScriptTemplate, 3))
	require.NoError(t, err)
	require.EqualValues(t, "0x3", v3.String())
}

func TestToS256(t *testing.T) {
	// positive
	vm := setupGojaVmForU256()
	v1, err := vm.RunString(fmt.Sprintf(ToS256ScriptTemplate, 100))
	require.NoError(t, err)
	require.EqualValues(t, "0x64", v1.String())

	// zero
	vm = setupGojaVmForU256()
	v2, err := vm.RunString(fmt.Sprintf(ToS256ScriptTemplate, 0))
	require.NoError(t, err)
	require.EqualValues(t, "0x0", v2.String())
}

func TestU256Add(t *testing.T) {
	vm := setupGojaVmForU256()

	v1, err := vm.RunString(fmt.Sprintf(U256AddScriptTemplate, 5, 8))
	require.NoError(t, err)
	require.EqualValues(t, "0xd", v1.String())

	v2, err := vm.RunString(fmt.Sprintf(U256UnsafeAddScriptTemplate, 5, 8))
	require.NoError(t, err)
	require.EqualValues(t, "0xd", v2.String())

	// overflow
	vm = setupGojaVmForU256()
	v3, err := vm.RunString(fmt.Sprintf(U256UnsafeAddScriptTemplate, utils.MaxSafeInteger, utils.MaxSafeInteger))
	require.NoError(t, err)
	require.EqualValues(t, "0x3ffffffffffffe", v3.String())
}

func TestU256DivAndMod(t *testing.T) {
	// Div
	vm := setupGojaVmForU256()
	v1, err := vm.RunString(fmt.Sprintf(U256DivScriptTemplate, 10, 3))
	require.NoError(t, err)
	require.EqualValues(t, "0x3", v1.String())

	// Mod
	vm = setupGojaVmForU256()
	v2, err := vm.RunString(fmt.Sprintf(U256ModScriptTemplate, 10, 3))
	require.NoError(t, err)
	require.EqualValues(t, "0x1", v2.String())

	// DivMod
	vm = setupGojaVmForU256()
	_, err = vm.RunString(fmt.Sprintf(U256DivModScriptTemplate, 10, 3))
	require.NoError(t, err)

	z := vm.Get("z").Export().(Uint256)
	require.EqualValues(t, "0x3", z.X.String())
	m := vm.Get("m").Export().(Uint256)
	require.EqualValues(t, "0x1", m.X.String())
}

func TestU256Exp(t *testing.T) {
	vm := setupGojaVmForU256()
	v1, err := vm.RunString(fmt.Sprintf(U256ExpScriptTemplate, 10, 2))
	require.NoError(t, err)
	require.EqualValues(t, "0x64", v1.String())
}

func TestU256Mul(t *testing.T) {
	vm := setupGojaVmForU256()

	v1, err := vm.RunString(fmt.Sprintf(U256MulScriptTemplate, 2, 5))
	require.NoError(t, err)
	require.EqualValues(t, "0xa", v1.String())

	v2, err := vm.RunString(fmt.Sprintf(U256UnsafeMulScriptTemplate, 2, 5))
	require.NoError(t, err)
	require.EqualValues(t, "0xa", v2.String())

	// overflow
	vm = setupGojaVmForU256()
	v3, err := vm.RunString(fmt.Sprintf(U256UnsafeMulScriptTemplate, 2, utils.MaxSafeInteger))
	require.NoError(t, err)
	require.EqualValues(t, "0x3ffffffffffffe", v3.String())
}

func TestU256And(t *testing.T) {
	vm := setupGojaVmForU256()

	v1, err := vm.RunString(fmt.Sprintf(U256AndScriptTemplate, "0xff", "0xf0"))
	require.NoError(t, err)
	require.EqualValues(t, "0xf0", v1.String())
}

func TestU256Or(t *testing.T) {
	vm := setupGojaVmForU256()

	v1, err := vm.RunString(fmt.Sprintf(U256OrScriptTemplate, "0xf", "0xf0"))
	require.NoError(t, err)
	require.EqualValues(t, "0xff", v1.String())
}

func TestU256Not(t *testing.T) {
	vm := setupGojaVmForU256()

	v1, err := vm.RunString(fmt.Sprintf(U256NotScriptTemplate, "0xf"))
	require.NoError(t, err)
	require.EqualValues(t, "0xfffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0", v1.String())
}

func TestU256Compare(t *testing.T) {
	vm := setupGojaVmForU256()

	_, err := vm.RunString(fmt.Sprintf(U256CompareScriptTemplate, 5, 2))
	require.NoError(t, err)

	gt := vm.Get("gt").Export().(bool)
	gte := vm.Get("gte").Export().(bool)
	lt := vm.Get("lt").Export().(bool)
	lte := vm.Get("lte").Export().(bool)
	eq := vm.Get("eq").Export().(bool)
	isZero := vm.Get("isZero").Export().(bool)

	require.True(t, gt)
	require.True(t, gte)
	require.False(t, lt)
	require.False(t, lte)
	require.False(t, eq)
	require.False(t, isZero)
}

func TestU256CompareNum(t *testing.T) {
	vm := setupGojaVmForU256()

	_, err := vm.RunString(fmt.Sprintf(U256CompareNumScriptTemplate, 5))
	require.NoError(t, err)

	gtNum := vm.Get("gtNum").Export().(bool)
	gteNum := vm.Get("gteNum").Export().(bool)
	ltNum := vm.Get("ltNum").Export().(bool)
	lteNum := vm.Get("lteNum").Export().(bool)

	require.False(t, gtNum)
	require.False(t, gteNum)
	require.True(t, ltNum)
	require.True(t, lteNum)
}

func TestU256Sqrt(t *testing.T) {
	vm := setupGojaVmForU256()
	v1, err := vm.RunString(fmt.Sprintf(U256SqrtScriptTemplate, 10))
	require.NoError(t, err)
	require.EqualValues(t, "0x3", v1.String())
}

func TestU256Sub(t *testing.T) {
	// Sub
	vm := setupGojaVmForU256()
	v1, err := vm.RunString(fmt.Sprintf(U256SubScriptTemplate, 10, 5))
	require.NoError(t, err)
	require.EqualValues(t, "0x5", v1.String())

	// UnsafeSub
	vm = setupGojaVmForU256()
	v2, err := vm.RunString(fmt.Sprintf(U256UnsafeSubScriptTemplate, 10, 5))
	require.NoError(t, err)
	require.EqualValues(t, "0x5", v2.String())

	// overflow
	vm = setupGojaVmForU256()
	v3, err := vm.RunString(fmt.Sprintf(U256UnsafeSubScriptTemplate, 0, 1))
	require.NoError(t, err)
	require.EqualValues(t, "0xffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff", v3.String())
}

func TestU256Shift(t *testing.T) {
	vm := setupGojaVmForU256()
	_, err := vm.RunString(fmt.Sprintf(U256ShiftScriptTemplate, "0xf", 2, 2))
	require.NoError(t, err)

	lsh := vm.Get("lsh").Export().(Uint256)
	rsh := vm.Get("rsh").Export().(Uint256)
	require.EqualValues(t, "0x3c", lsh.X.String())
	require.EqualValues(t, "0x3", rsh.X.String())
}

func TestU256SafeInteger(t *testing.T) {
	vm := setupGojaVmForU256()
	_, err := vm.RunString(fmt.Sprintf(U256SafeIntegerScriptTemplate, 100))
	require.NoError(t, err)

	isSafe := vm.Get("isSafe").Export().(bool)
	toSafe := vm.Get("toSafe").Export().(int64)
	require.True(t, isSafe)
	require.EqualValues(t, int64(100), toSafe)
}

func TestU256Conversion(t *testing.T) {
	vm := setupGojaVmForU256()
	_, err := vm.RunString(fmt.Sprintf(U256ConversionScriptTemplate, 100))
	require.NoError(t, err)

	toBuf := vm.Get("toBuf").Export().(goja.ArrayBuffer)
	toHex := vm.Get("toHex").Export().(string)
	toString := vm.Get("toString").Export().(string)
	require.EqualValues(t, gethcmn.FromHex("0x0000000000000000000000000000000000000000000000000000000000000064"), toBuf.Bytes())
	require.EqualValues(t, "0x64", toHex)
	require.EqualValues(t, "0x64", toString)
}
