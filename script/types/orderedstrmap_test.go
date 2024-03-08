package types

import (
	"testing"

	"github.com/dop251/goja"
	"github.com/stretchr/testify/require"
)

const (
	OrderedStrMapCRUDScriptTemplate = `
		let m = NewOrderedStrMap()
		const [v1, ok1] = m.Get('a')
		const len1 = m.Len()

		m.Set('a', 'a')
		m.Set('b', 'b')
		m.Set('c', 'c')
		m.Set('c', 'd')
		const [v2, ok2] = m.Get('c')
		const len2 = m.Len()

		m.Delete('b')
		const [v3, ok3] = m.Get('b')
		const len3 = m.Len()
	`

	OrderedStrMapGetOrDeleteEmptyKeyScriptTemplate = `
		let m = NewOrderedStrMap()
		const [v1, ok1] = m.Get('')
		const len1 = m.Len()

		m.Delete('')
		const len2 = m.Len()
	`

	OrderedStrMapSetEmptyKeyScriptTemplate = `
		let m = NewOrderedStrMap()
		m.Set('', 'a')
	`

	OrderedStrMapSeekScriptTemplate = `
		let m = NewOrderedStrMap()
		m.Set('a', 'a')
		m.Set('b', 'b')
		m.Set('c', 'c')
		m.Set('d', 'd')
		m.Set('e', 'e')

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

	OrderedStrMapSeekFirstAndLastScriptTemplate = `
		let m = NewOrderedStrMap()
		m.Set('a', 'a')
		m.Set('b', 'b')
		m.Set('c', 'c')
		m.Set('d', 'd')
		m.Set('e', 'e')

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

	OrderedStrMapClearScriptTemplate = `
		let m = NewOrderedStrMap()
		m.Set('a', 'a')
		m.Set('b', 'b')
		m.Set('c', 'c')
		m.Set('d', 'd')
		m.Set('e', 'e')

		m.Clear()
		const len1 = m.Len()

		m.Set('e', 'e')
		m.Clear()
		const len2 = m.Len()
	`
)

func setupGojaVmForOrderedStrMap() *goja.Runtime {
	vm := goja.New()

	vm.Set("NewOrderedStrMap", NewOrderedStrMap)
	return vm
}

func TestOrderedStrMapCRUD(t *testing.T) {
	vm := setupGojaVmForOrderedStrMap()
	_, err := vm.RunString(OrderedStrMapCRUDScriptTemplate)
	require.NoError(t, err)

	// 1. get non-existed key
	ok1 := vm.Get("ok1").Export().(bool)
	v1 := vm.Get("v1").Export().(string)
	len1 := vm.Get("len1").Export().(int64)
	require.False(t, ok1)
	require.EqualValues(t, "", v1)
	require.EqualValues(t, 0, len1)

	// 2. set and get keys
	ok2 := vm.Get("ok2").Export().(bool)
	v2 := vm.Get("v2").Export().(string)
	len2 := vm.Get("len2").Export().(int64)
	require.True(t, ok2)
	require.EqualValues(t, "d", v2)
	require.EqualValues(t, 3, len2)

	// 3. delete key
	ok3 := vm.Get("ok3").Export().(bool)
	v3 := vm.Get("v3").Export().(string)
	len3 := vm.Get("len3").Export().(int64)
	require.False(t, ok3)
	require.EqualValues(t, "", v3)
	require.EqualValues(t, 2, len3)
}

func TestOrderedStrMapGetOrDeleteEmptyKey(t *testing.T) {
	vm := setupGojaVmForOrderedStrMap()
	_, err := vm.RunString(OrderedStrMapGetOrDeleteEmptyKeyScriptTemplate)
	require.NoError(t, err)

	len1 := vm.Get("len1").Export().(int64)
	len2 := vm.Get("len2").Export().(int64)
	require.EqualValues(t, 0, len1)
	require.EqualValues(t, 0, len2)
}

func TestOrderedStrMapSetEmptyKey(t *testing.T) {
	vm := setupGojaVmForOrderedStrMap()
	_, err := vm.RunString(OrderedStrMapSetEmptyKeyScriptTemplate)
	require.Error(t, err, "Empty key string")
}

func TestOrderedStrMapSeek(t *testing.T) {
	vm := setupGojaVmForOrderedStrMap()
	_, err := vm.RunString(OrderedStrMapSeekScriptTemplate)
	require.NoError(t, err)

	ok1 := vm.Get("ok1").Export().(bool)
	ok2 := vm.Get("ok2").Export().(bool)
	require.True(t, ok1)
	require.True(t, ok2)

	k1 := vm.Get("k1").Export().(string)
	v1 := vm.Get("v1").Export().(string)
	k2 := vm.Get("k2").Export().(string)
	v2 := vm.Get("v2").Export().(string)
	k3 := vm.Get("k3").Export().(string)
	v3 := vm.Get("v3").Export().(string)
	k4 := vm.Get("k4").Export().(string)
	v4 := vm.Get("v4").Export().(string)
	require.EqualValues(t, "c", k1)
	require.EqualValues(t, "c", v1)
	require.EqualValues(t, "b", k2)
	require.EqualValues(t, "b", v2)
	require.EqualValues(t, "a", k3)
	require.EqualValues(t, "a", v3)
	require.EqualValues(t, "", k4)
	require.EqualValues(t, "", v4)

	k5 := vm.Get("k5").Export().(string)
	v5 := vm.Get("v5").Export().(string)
	k6 := vm.Get("k6").Export().(string)
	v6 := vm.Get("v6").Export().(string)
	k7 := vm.Get("k7").Export().(string)
	v7 := vm.Get("v7").Export().(string)
	k8 := vm.Get("k8").Export().(string)
	v8 := vm.Get("v8").Export().(string)
	require.EqualValues(t, "c", k5)
	require.EqualValues(t, "c", v5)
	require.EqualValues(t, "d", k6)
	require.EqualValues(t, "d", v6)
	require.EqualValues(t, "e", k7)
	require.EqualValues(t, "e", v7)
	require.EqualValues(t, "", k8)
	require.EqualValues(t, "", v8)
}

func TestOrderedStrMapSeekFirstAndLast(t *testing.T) {
	vm := setupGojaVmForOrderedStrMap()
	_, err := vm.RunString(OrderedStrMapSeekFirstAndLastScriptTemplate)
	require.NoError(t, err)

	k1 := vm.Get("k1").Export().(string)
	v1 := vm.Get("v1").Export().(string)
	k2 := vm.Get("k2").Export().(string)
	v2 := vm.Get("v2").Export().(string)
	k3 := vm.Get("k3").Export().(string)
	v3 := vm.Get("v3").Export().(string)
	k4 := vm.Get("k4").Export().(string)
	v4 := vm.Get("v4").Export().(string)
	k5 := vm.Get("k5").Export().(string)
	v5 := vm.Get("v5").Export().(string)
	require.EqualValues(t, "a", k1)
	require.EqualValues(t, "a", v1)
	require.EqualValues(t, "b", k2)
	require.EqualValues(t, "b", v2)
	require.EqualValues(t, "c", k3)
	require.EqualValues(t, "c", v3)
	require.EqualValues(t, "d", k4)
	require.EqualValues(t, "d", v4)
	require.EqualValues(t, "e", k5)
	require.EqualValues(t, "e", v5)

	k6 := vm.Get("k6").Export().(string)
	v6 := vm.Get("v6").Export().(string)
	k7 := vm.Get("k7").Export().(string)
	v7 := vm.Get("v7").Export().(string)
	k8 := vm.Get("k8").Export().(string)
	v8 := vm.Get("v8").Export().(string)
	k9 := vm.Get("k9").Export().(string)
	v9 := vm.Get("v9").Export().(string)
	k10 := vm.Get("k10").Export().(string)
	v10 := vm.Get("v10").Export().(string)
	require.EqualValues(t, "e", k6)
	require.EqualValues(t, "e", v6)
	require.EqualValues(t, "d", k7)
	require.EqualValues(t, "d", v7)
	require.EqualValues(t, "c", k8)
	require.EqualValues(t, "c", v8)
	require.EqualValues(t, "b", k9)
	require.EqualValues(t, "b", v9)
	require.EqualValues(t, "a", k10)
	require.EqualValues(t, "a", v10)
}

func TestOrderedStrMapClear(t *testing.T) {
	vm := setupGojaVmForOrderedStrMap()
	_, err := vm.RunString(OrderedStrMapClearScriptTemplate)
	require.NoError(t, err)

	len1 := vm.Get("len1").Export().(int64)
	len2 := vm.Get("len2").Export().(int64)
	require.EqualValues(t, 0, len1)
	require.EqualValues(t, 0, len2)
}
