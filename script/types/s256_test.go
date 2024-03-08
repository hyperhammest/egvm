package types

import (
	"fmt"
	"testing"

	"github.com/dop251/goja"
	gethcmn "github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
)

const (
	S256ScriptTemplate = `
		const a = S256(%v)
		a.String()
	`

	HexToS256ScriptTemplate = `
		const b = HexToS256('%v')
		b.String()
	`

	BufToS256ScriptTemplate = `
		let buffer = new ArrayBuffer(32); // 32 bytes
		let view = new Uint8Array(buffer);
		view[31] = %v
		const c = BufToS256(buffer)
		c.String()
	`

	ToU256ScriptTemplate = `
		const d = S256(%v)
		const ud = d.ToU256()
		ud.String()
	`

	S256AddScriptTemplate = `
		const a = S256(%v)
		const b = S256(%v)
		const ab = a.Add(b)
		ab.String()
	`

	S256SubScriptTemplate = `
		const a = S256(%v)
		const b = S256(%v)
		const ab = a.Sub(b)
		ab.String()
	`

	S256DivScriptTemplate = `
		const a = S256(%v)
		const b = S256(%v)
		const ab = a.Div(b)
		ab.String()
	`

	S256MulScriptTemplate = `
		const a = S256(%v)
		const b = S256(%v)
		const ab = a.Mul(b)
		ab.String()
	`

	S256NegScriptTemplate = `
		const a = S256(%v)
		const na = a.Neg()
	`

	S256AbsScriptTemplate = `
		const a = S256(%v)
		const aa = a.Abs()
	`

	S256SignScriptTemplate = `
		const a = S256(%v)
		const b = S256(%v)
		const sa = a.Sign()
		const sb = b.Sign()
	`

	S256CompareScriptTemplate = `
		const a = S256(%v)
		const b = S256(%v)
		const gt = a.Gt(b)
		const gte = a.Gte(b)
		const lt = a.Lt(b)
		const lte = a.Lte(b)
		const eq = a.Equal(b)
		const isZero = a.IsZero()
	`

	S256CompareNumScriptTemplate = `
		const a = S256(%v)
		const b = 10
		const gtNum = a.GtNum(b)
		const gteNum = a.GteNum(b)
		const ltNum = a.LtNum(b)
		const lteNum = a.LteNum(b)
	`

	S256ShiftScriptTemplate = `
		const a = HexToS256('%v')
		const l = %v
		const r = %v
		const lsh = a.Lsh(l)
		const rsh = a.Rsh(r)
	`

	S256SafeIntegerScriptTemplate = `
		const a = S256(%v)
		const isSafe = a.IsSafeInteger()
		const toSafe = a.ToSafeInteger()
	`

	S256ConversionScriptTemplate = `
		const a = S256(%v)
		const toBuf = a.ToBuf()
		const toHex = a.ToHex()
		const toString = a.String()
	`
)

func setupGojaVmForS256() *goja.Runtime {
	vm := goja.New()

	vm.Set("S256", S256)
	vm.Set("HexToS256", HexToS256)
	vm.Set("BufToS256", BufToS256)
	return vm
}

func TestCreateS256(t *testing.T) {
	vm := setupGojaVmForS256()

	// S256
	v1, err := vm.RunString(fmt.Sprintf(S256ScriptTemplate, -1))
	require.NoError(t, err)
	require.EqualValues(t, "-0x1", v1.String())

	// HexToS256
	v2, err := vm.RunString(fmt.Sprintf(HexToS256ScriptTemplate, "0x2"))
	require.NoError(t, err)
	require.EqualValues(t, "0x2", v2.String())

	// BufToS256
	v3, err := vm.RunString(fmt.Sprintf(BufToS256ScriptTemplate, 3))
	require.NoError(t, err)
	require.EqualValues(t, "0x3", v3.String())
}

func TestToU256(t *testing.T) {
	// positive
	vm := setupGojaVmForS256()
	v1, err := vm.RunString(fmt.Sprintf(ToU256ScriptTemplate, 100))
	require.NoError(t, err)
	require.EqualValues(t, "0x64", v1.String())

	// zero
	vm = setupGojaVmForS256()
	v2, err := vm.RunString(fmt.Sprintf(ToU256ScriptTemplate, 0))
	require.NoError(t, err)
	require.EqualValues(t, "0x0", v2.String())

	// negative
	vm = setupGojaVmForS256()
	v3, err := vm.RunString(fmt.Sprintf(ToU256ScriptTemplate, -100))
	require.NoError(t, err)
	require.EqualValues(t, "0xffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff9c", v3.String())
}

func TestS256Add(t *testing.T) {
	// positive
	vm := setupGojaVmForS256()
	v1, err := vm.RunString(fmt.Sprintf(S256AddScriptTemplate, 5, 8))
	require.NoError(t, err)
	require.EqualValues(t, "0xd", v1.String())

	// negative
	vm = setupGojaVmForS256()
	v2, err := vm.RunString(fmt.Sprintf(S256AddScriptTemplate, -5, -8))
	require.NoError(t, err)
	require.EqualValues(t, "-0xd", v2.String())
}

func TestS256Sub(t *testing.T) {
	// Sub
	vm := setupGojaVmForS256()
	v1, err := vm.RunString(fmt.Sprintf(S256SubScriptTemplate, -10, 5))
	require.NoError(t, err)
	require.EqualValues(t, "-0xf", v1.String())

}

func TestS256Div(t *testing.T) {
	// Div
	vm := setupGojaVmForS256()
	v1, err := vm.RunString(fmt.Sprintf(S256DivScriptTemplate, 10, -3))
	require.NoError(t, err)
	require.EqualValues(t, "-0x3", v1.String())
}

func TestS256Mul(t *testing.T) {
	vm := setupGojaVmForS256()

	v1, err := vm.RunString(fmt.Sprintf(S256MulScriptTemplate, 2, -5))
	require.NoError(t, err)
	require.EqualValues(t, "-0xa", v1.String())
}

func TestS256Neg(t *testing.T) {
	vm := setupGojaVmForS256()

	_, err := vm.RunString(fmt.Sprintf(S256NegScriptTemplate, 5))
	require.NoError(t, err)
	na := vm.Get("na").Export().(Sint256)

	require.EqualValues(t, "0xfffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffb", na.x.String())
}

func TestS256Abs(t *testing.T) {
	vm := setupGojaVmForS256()

	_, err := vm.RunString(fmt.Sprintf(S256AbsScriptTemplate, -5))
	require.NoError(t, err)
	aa := vm.Get("aa").Export().(Sint256)
	require.EqualValues(t, "0x5", aa.x.String())
}

func TestS256Sign(t *testing.T) {
	vm := setupGojaVmForS256()

	_, err := vm.RunString(fmt.Sprintf(S256SignScriptTemplate, 2, -5))
	require.NoError(t, err)
	sa := vm.Get("sa").Export().(int64)
	sb := vm.Get("sb").Export().(int64)

	require.EqualValues(t, 1, sa)
	require.EqualValues(t, -1, sb)
}

func TestS256Compare(t *testing.T) {
	vm := setupGojaVmForS256()

	_, err := vm.RunString(fmt.Sprintf(S256CompareScriptTemplate, 5, -5))
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

func TestS256CompareNum(t *testing.T) {
	vm := setupGojaVmForS256()

	_, err := vm.RunString(fmt.Sprintf(S256CompareNumScriptTemplate, -5))
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

func TestS256Shift(t *testing.T) {
	vm := setupGojaVmForS256()
	_, err := vm.RunString(fmt.Sprintf(S256ShiftScriptTemplate, "0xf", 2, 2))
	require.NoError(t, err)

	lsh := vm.Get("lsh").Export().(Sint256)
	rsh := vm.Get("rsh").Export().(Sint256)
	require.EqualValues(t, "0x3c", lsh.x.String())
	require.EqualValues(t, "0x3", rsh.x.String())
}

func TestS256SafeInteger(t *testing.T) {
	vm := setupGojaVmForS256()
	_, err := vm.RunString(fmt.Sprintf(S256SafeIntegerScriptTemplate, 100))
	require.NoError(t, err)

	isSafe := vm.Get("isSafe").Export().(bool)
	toSafe := vm.Get("toSafe").Export().(int64)
	require.True(t, isSafe)
	require.EqualValues(t, int64(100), toSafe)
}

func TestS256Conversion(t *testing.T) {
	vm := setupGojaVmForS256()
	_, err := vm.RunString(fmt.Sprintf(S256ConversionScriptTemplate, -100))
	require.NoError(t, err)

	toBuf := vm.Get("toBuf").Export().(goja.ArrayBuffer)
	toHex := vm.Get("toHex").Export().(string)
	toString := vm.Get("toString").Export().(string)
	require.EqualValues(t, gethcmn.FromHex("0xffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff9c"), toBuf.Bytes())
	require.EqualValues(t, "0xffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff9c", toHex)
	require.EqualValues(t, "-0x64", toString)
}
