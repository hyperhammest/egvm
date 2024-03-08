package types

import (
	"testing"

	"github.com/dop251/goja"
	gethcmn "github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
)

const (
	OrderedBufMapCRUDScriptTemplate = `
		let m = NewOrderedBufMap()
		const [v1, ok1] = m.Get('a')
		const len1 = m.Len()

		let buffer1 = new ArrayBuffer(8); // 8 bytes
		let view1 = new Uint8Array(buffer1);
		view1[7] = 15
		
		let buffer2 = new ArrayBuffer(8); // 8 bytes
		let view2 = new Uint8Array(buffer2);
		view2[7] = 10
		
		m.Set('b', buffer1)
		m.Set('c', buffer2)
		const [v2, ok2] = m.Get('c')
		const len2 = m.Len()
		
		m.Delete('b')
		const [v3, ok3] = m.Get('b')
		const len3 = m.Len()
	`

	OrderedBufMapGetOrDeleteEmptyKeyScriptTemplate = `
		let m = NewOrderedBufMap()
		const [v1, ok1] = m.Get('')
		const len1 = m.Len()

		m.Delete('')
		const len2 = m.Len()
	`

	OrderedBufMapSetEmptyKeyScriptTemplate = `
		let m = NewOrderedBufMap()
		let buffer1 = new ArrayBuffer(8); // 8 bytes
		let view1 = new Uint8Array(buffer1);
		view1[7] = 15
		
		m.Set('', buffer1)
	`

	OrderedBufMapSeekScriptTemplate = `
		let m = NewOrderedBufMap()

		let buffer1 = new ArrayBuffer(8); // 8 bytes
		let buffer2 = new ArrayBuffer(8); // 8 bytes
		let buffer3 = new ArrayBuffer(8); // 8 bytes
		let buffer4 = new ArrayBuffer(8); // 8 bytes
		let buffer5 = new ArrayBuffer(8); // 8 bytes
		let view1 = new Uint8Array(buffer1);
		view1[7] = 1
		let view2 = new Uint8Array(buffer2);
		view2[7] = 2
		let view3 = new Uint8Array(buffer3);
		view3[7] = 3
		let view4 = new Uint8Array(buffer4);
		view4[7] = 4
		let view5 = new Uint8Array(buffer5);
		view5[7] = 5


		m.Set('a', buffer1)
		m.Set('b', buffer2)
		m.Set('c', buffer3)
		m.Set('d', buffer4)
		m.Set('e', buffer5)

		const [it1, ok1] = m.Seek('c')
		const [k1, v1] = it1.Prev()
		const [k2, v2] = it1.Prev()
		const [k3, v3] = it1.Prev()
		const [k4, v4] = it1.Prev()
		it1.Close()

		const [it2, ok2] = m.Seek('c')
		const [k5, v5] = it2.Next()
		const [k6, v6] = it2.Next()
		const [k7, v7] = it2.Next()
		const [k8, v8] = it2.Next()
		it2.Close()
	`

	OrderedBufMapSeekFirstAndLastScriptTemplate = `
		let m = NewOrderedBufMap()

		let buffer1 = new ArrayBuffer(8); // 8 bytes
		let buffer2 = new ArrayBuffer(8); // 8 bytes
		let buffer3 = new ArrayBuffer(8); // 8 bytes
		let buffer4 = new ArrayBuffer(8); // 8 bytes
		let buffer5 = new ArrayBuffer(8); // 8 bytes
		let view1 = new Uint8Array(buffer1);
		view1[7] = 1
		let view2 = new Uint8Array(buffer2);
		view2[7] = 2
		let view3 = new Uint8Array(buffer3);
		view3[7] = 3
		let view4 = new Uint8Array(buffer4);
		view4[7] = 4
		let view5 = new Uint8Array(buffer5);
		view5[7] = 5


		m.Set('a', buffer1)
		m.Set('b', buffer2)
		m.Set('c', buffer3)
		m.Set('d', buffer4)
		m.Set('e', buffer5)

		const it1 = m.SeekFirst()
		const [k1, v1] = it1.Next()
		const [k2, v2] = it1.Next()
		const [k3, v3] = it1.Next()
		const [k4, v4] = it1.Next()
		const [k5, v5] = it1.Next()
		it1.Close()

		const it2 = m.SeekLast()
		const [k6, v6] = it2.Prev()
		const [k7, v7] = it2.Prev()
		const [k8, v8] = it2.Prev()
		const [k9, v9] = it2.Prev()
		const [k10, v10] = it2.Prev()
		it2.Close()
	`

	OrderedBufMapClearScriptTemplate = `
		let m = NewOrderedBufMap()

		let buffer1 = new ArrayBuffer(8); // 8 bytes
		let buffer2 = new ArrayBuffer(8); // 8 bytes
		let buffer3 = new ArrayBuffer(8); // 8 bytes
		let buffer4 = new ArrayBuffer(8); // 8 bytes
		let buffer5 = new ArrayBuffer(8); // 8 bytes
		let view1 = new Uint8Array(buffer1);
		view1[7] = 1
		let view2 = new Uint8Array(buffer2);
		view2[7] = 2
		let view3 = new Uint8Array(buffer3);
		view3[7] = 3
		let view4 = new Uint8Array(buffer4);
		view4[7] = 4
		let view5 = new Uint8Array(buffer5);
		view5[7] = 5


		m.Set('a', buffer1)
		m.Set('b', buffer2)
		m.Set('c', buffer3)
		m.Set('d', buffer4)
		m.Set('e', buffer5)


		m.Clear()
		const len1 = m.Len()

		m.Set('e', buffer5)
		m.Clear()
		const len2 = m.Len()
	`
)

func setupGojaVmForOrderedBufMap() *goja.Runtime {
	vm := goja.New()

	vm.Set("NewOrderedBufMap", NewOrderedBufMap)
	return vm
}

func TestOrderedBufMapCRUD(t *testing.T) {
	vm := setupGojaVmForOrderedBufMap()
	_, err := vm.RunString(OrderedBufMapCRUDScriptTemplate)
	require.NoError(t, err)

	// 1. get non-existed key
	ok1 := vm.Get("ok1").Export().(bool)
	v1 := vm.Get("v1").Export().(goja.ArrayBuffer)
	v1Hex := gethcmn.Bytes2Hex(v1.Bytes())
	len1 := vm.Get("len1").Export().(int64)
	require.False(t, ok1)
	require.EqualValues(t, "", v1Hex)
	require.EqualValues(t, 0, len1)

	// 2. set and get keys
	ok2 := vm.Get("ok2").Export().(bool)
	v2 := vm.Get("v2").Export().(goja.ArrayBuffer)
	v2Hex := gethcmn.Bytes2Hex(v2.Bytes())
	len2 := vm.Get("len2").Export().(int64)
	require.True(t, ok2)
	require.EqualValues(t, "000000000000000a", v2Hex)
	require.EqualValues(t, 2, len2)

	// 3. delete key
	ok3 := vm.Get("ok3").Export().(bool)
	v3 := vm.Get("v3").Export().(goja.ArrayBuffer)
	v3Hex := gethcmn.Bytes2Hex(v3.Bytes())
	len3 := vm.Get("len3").Export().(int64)
	require.False(t, ok3)
	require.EqualValues(t, "", v3Hex)
	require.EqualValues(t, 1, len3)
}

func TestOrderedBufMapGetOrDeleteEmptyKey(t *testing.T) {
	vm := setupGojaVmForOrderedBufMap()
	_, err := vm.RunString(OrderedBufMapGetOrDeleteEmptyKeyScriptTemplate)
	require.NoError(t, err)

	len1 := vm.Get("len1").Export().(int64)
	len2 := vm.Get("len2").Export().(int64)
	require.EqualValues(t, 0, len1)
	require.EqualValues(t, 0, len2)
}

func TestOrderedBufMapSetEmptyKey(t *testing.T) {
	vm := setupGojaVmForOrderedBufMap()
	_, err := vm.RunString(OrderedBufMapSetEmptyKeyScriptTemplate)
	require.Error(t, err, "Empty key string")
}

func TestOrderedBufMapSeek(t *testing.T) {
	vm := setupGojaVmForOrderedBufMap()
	_, err := vm.RunString(OrderedBufMapSeekScriptTemplate)
	require.NoError(t, err)

	ok1 := vm.Get("ok1").Export().(bool)
	ok2 := vm.Get("ok2").Export().(bool)
	require.True(t, ok1)
	require.True(t, ok2)

	k1 := vm.Get("k1").Export().(string)
	v1 := vm.Get("v1").Export().(goja.ArrayBuffer)
	v1Hex := gethcmn.Bytes2Hex(v1.Bytes())
	k2 := vm.Get("k2").Export().(string)
	v2 := vm.Get("v2").Export().(goja.ArrayBuffer)
	v2Hex := gethcmn.Bytes2Hex(v2.Bytes())
	k3 := vm.Get("k3").Export().(string)
	v3 := vm.Get("v3").Export().(goja.ArrayBuffer)
	v3Hex := gethcmn.Bytes2Hex(v3.Bytes())
	k4 := vm.Get("k4").Export().(string)
	_, ok4 := vm.Get("v4").Export().(goja.ArrayBuffer)
	require.EqualValues(t, "c", k1)
	require.EqualValues(t, "0000000000000003", v1Hex)
	require.EqualValues(t, "b", k2)
	require.EqualValues(t, "0000000000000002", v2Hex)
	require.EqualValues(t, "a", k3)
	require.EqualValues(t, "0000000000000001", v3Hex)
	require.EqualValues(t, "", k4)
	require.False(t, ok4)

	k5 := vm.Get("k5").Export().(string)
	v5 := vm.Get("v5").Export().(goja.ArrayBuffer)
	v5Hex := gethcmn.Bytes2Hex(v5.Bytes())
	k6 := vm.Get("k6").Export().(string)
	v6 := vm.Get("v6").Export().(goja.ArrayBuffer)
	v6Hex := gethcmn.Bytes2Hex(v6.Bytes())
	k7 := vm.Get("k7").Export().(string)
	v7 := vm.Get("v7").Export().(goja.ArrayBuffer)
	v7Hex := gethcmn.Bytes2Hex(v7.Bytes())
	k8 := vm.Get("k8").Export().(string)
	_, ok8 := vm.Get("v8").Export().(goja.ArrayBuffer)
	require.EqualValues(t, "c", k5)
	require.EqualValues(t, "0000000000000003", v5Hex)
	require.EqualValues(t, "d", k6)
	require.EqualValues(t, "0000000000000004", v6Hex)
	require.EqualValues(t, "e", k7)
	require.EqualValues(t, "0000000000000005", v7Hex)
	require.EqualValues(t, "", k8)
	require.False(t, ok8)
}

func TestOrderedBufMapSeekFirstAndLast(t *testing.T) {
	vm := setupGojaVmForOrderedBufMap()
	_, err := vm.RunString(OrderedBufMapSeekFirstAndLastScriptTemplate)
	require.NoError(t, err)

	k1 := vm.Get("k1").Export().(string)
	v1 := vm.Get("v1").Export().(goja.ArrayBuffer)
	v1Hex := gethcmn.Bytes2Hex(v1.Bytes())
	k2 := vm.Get("k2").Export().(string)
	v2 := vm.Get("v2").Export().(goja.ArrayBuffer)
	v2Hex := gethcmn.Bytes2Hex(v2.Bytes())
	k3 := vm.Get("k3").Export().(string)
	v3 := vm.Get("v3").Export().(goja.ArrayBuffer)
	v3Hex := gethcmn.Bytes2Hex(v3.Bytes())
	k4 := vm.Get("k4").Export().(string)
	v4 := vm.Get("v4").Export().(goja.ArrayBuffer)
	v4Hex := gethcmn.Bytes2Hex(v4.Bytes())
	k5 := vm.Get("k5").Export().(string)
	v5 := vm.Get("v5").Export().(goja.ArrayBuffer)
	v5Hex := gethcmn.Bytes2Hex(v5.Bytes())
	require.EqualValues(t, "a", k1)
	require.EqualValues(t, "0000000000000001", v1Hex)
	require.EqualValues(t, "b", k2)
	require.EqualValues(t, "0000000000000002", v2Hex)
	require.EqualValues(t, "c", k3)
	require.EqualValues(t, "0000000000000003", v3Hex)
	require.EqualValues(t, "d", k4)
	require.EqualValues(t, "0000000000000004", v4Hex)
	require.EqualValues(t, "e", k5)
	require.EqualValues(t, "0000000000000005", v5Hex)

	k6 := vm.Get("k6").Export().(string)
	v6 := vm.Get("v6").Export().(goja.ArrayBuffer)
	v6Hex := gethcmn.Bytes2Hex(v6.Bytes())
	k7 := vm.Get("k7").Export().(string)
	v7 := vm.Get("v7").Export().(goja.ArrayBuffer)
	v7Hex := gethcmn.Bytes2Hex(v7.Bytes())
	k8 := vm.Get("k8").Export().(string)
	v8 := vm.Get("v8").Export().(goja.ArrayBuffer)
	v8Hex := gethcmn.Bytes2Hex(v8.Bytes())
	k9 := vm.Get("k9").Export().(string)
	v9 := vm.Get("v9").Export().(goja.ArrayBuffer)
	v9Hex := gethcmn.Bytes2Hex(v9.Bytes())
	k10 := vm.Get("k10").Export().(string)
	v10 := vm.Get("v10").Export().(goja.ArrayBuffer)
	v10Hex := gethcmn.Bytes2Hex(v10.Bytes())
	require.EqualValues(t, "e", k6)
	require.EqualValues(t, "0000000000000005", v6Hex)
	require.EqualValues(t, "d", k7)
	require.EqualValues(t, "0000000000000004", v7Hex)
	require.EqualValues(t, "c", k8)
	require.EqualValues(t, "0000000000000003", v8Hex)
	require.EqualValues(t, "b", k9)
	require.EqualValues(t, "0000000000000002", v9Hex)
	require.EqualValues(t, "a", k10)
	require.EqualValues(t, "0000000000000001", v10Hex)
}

func TestOrderedBufMapClear(t *testing.T) {
	vm := setupGojaVmForOrderedBufMap()
	_, err := vm.RunString(OrderedBufMapClearScriptTemplate)
	require.NoError(t, err)

	len1 := vm.Get("len1").Export().(int64)
	len2 := vm.Get("len2").Export().(int64)
	require.EqualValues(t, 0, len1)
	require.EqualValues(t, 0, len2)
}
