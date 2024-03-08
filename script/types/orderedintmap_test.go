package types

import (
	"testing"

	"github.com/dop251/goja"
	"github.com/stretchr/testify/require"
)

const (
	OrderedIntMapCRUDScriptTemplate = `
		let m = NewOrderedIntMap()
		const [v1, ok1] = m.Get('a')
		const len1 = m.Len()

		m.Set('a', 1)
		m.Set('b', 2)
		m.Set('c', 3)
		m.Set('c', 4)
		const [v2, ok2] = m.Get('c')
		const len2 = m.Len()

		m.Delete('b')
		const [v3, ok3] = m.Get('b')
		const len3 = m.Len()
	`

	OrderedIntMapGetOrDeleteEmptyKeyScriptTemplate = `
		let m = NewOrderedIntMap()
		const [v1, ok1] = m.Get('')
		const len1 = m.Len()

		m.Delete('')
		const len2 = m.Len()
	`

	OrderedIntMapSetEmptyKeyScriptTemplate = `
		let m = NewOrderedIntMap()
		m.Set('', 1)
	`

	OrderedIntMapSeekScriptTemplate = `
		let m = NewOrderedIntMap()
		m.Set('a', 1)
		m.Set('b', 2)
		m.Set('c', 3)
		m.Set('d', 4)
		m.Set('e', 5)

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

	OrderedIntMapSeekFirstAndLastScriptTemplate = `
		let m = NewOrderedIntMap()
		m.Set('a', 1)
		m.Set('b', 2)
		m.Set('c', 3)
		m.Set('d', 4)
		m.Set('e', 5)

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

	OrderedIntMapClearScriptTemplate = `
		let m = NewOrderedIntMap()
		m.Set('a', 1)
		m.Set('b', 2)
		m.Set('c', 3)
		m.Set('d', 4)
		m.Set('e', 5)
		m.Clear()
		const len1 = m.Len()

		m.Set('e', 5)
		m.Clear()
		const len2 = m.Len()
	`
)

func setupGojaVmForOrderedIntMap() *goja.Runtime {
	vm := goja.New()

	vm.Set("NewOrderedIntMap", NewOrderedIntMap)
	return vm
}

func TestOrderedIntMapCRUD(t *testing.T) {
	vm := setupGojaVmForOrderedIntMap()
	_, err := vm.RunString(OrderedIntMapCRUDScriptTemplate)
	require.NoError(t, err)

	// 1. get non-existed key
	ok1 := vm.Get("ok1").Export().(bool)
	v1 := vm.Get("v1").Export().(int64)
	len1 := vm.Get("len1").Export().(int64)
	require.False(t, ok1)
	require.EqualValues(t, 0, v1)
	require.EqualValues(t, 0, len1)

	// 2. set and get keys
	ok2 := vm.Get("ok2").Export().(bool)
	v2 := vm.Get("v2").Export().(int64)
	len2 := vm.Get("len2").Export().(int64)
	require.True(t, ok2)
	require.EqualValues(t, 4, v2)
	require.EqualValues(t, 3, len2)

	// 3. delete key
	ok3 := vm.Get("ok3").Export().(bool)
	v3 := vm.Get("v3").Export().(int64)
	len3 := vm.Get("len3").Export().(int64)
	require.False(t, ok3)
	require.EqualValues(t, 0, v3)
	require.EqualValues(t, 2, len3)
}

func TestOrderedIntMapGetOrDeleteEmptyKey(t *testing.T) {
	vm := setupGojaVmForOrderedIntMap()
	_, err := vm.RunString(OrderedIntMapGetOrDeleteEmptyKeyScriptTemplate)
	require.NoError(t, err)

	len1 := vm.Get("len1").Export().(int64)
	len2 := vm.Get("len2").Export().(int64)
	require.EqualValues(t, 0, len1)
	require.EqualValues(t, 0, len2)
}

func TestOrderedIntMapSetEmptyKey(t *testing.T) {
	vm := setupGojaVmForOrderedIntMap()
	_, err := vm.RunString(OrderedIntMapSetEmptyKeyScriptTemplate)
	require.Error(t, err, "Empty key string")
}

func TestOrderedIntMapSeek(t *testing.T) {
	vm := setupGojaVmForOrderedIntMap()
	_, err := vm.RunString(OrderedIntMapSeekScriptTemplate)
	require.NoError(t, err)

	ok1 := vm.Get("ok1").Export().(bool)
	ok2 := vm.Get("ok2").Export().(bool)
	require.True(t, ok1)
	require.True(t, ok2)

	k1 := vm.Get("k1").Export().(string)
	v1 := vm.Get("v1").Export().(int64)
	k2 := vm.Get("k2").Export().(string)
	v2 := vm.Get("v2").Export().(int64)
	k3 := vm.Get("k3").Export().(string)
	v3 := vm.Get("v3").Export().(int64)
	k4 := vm.Get("k4").Export().(string)
	v4 := vm.Get("v4").Export().(int64)
	require.EqualValues(t, "c", k1)
	require.EqualValues(t, 3, v1)
	require.EqualValues(t, "b", k2)
	require.EqualValues(t, 2, v2)
	require.EqualValues(t, "a", k3)
	require.EqualValues(t, 1, v3)
	require.EqualValues(t, "", k4)
	require.EqualValues(t, 0, v4)

	k5 := vm.Get("k5").Export().(string)
	v5 := vm.Get("v5").Export().(int64)
	k6 := vm.Get("k6").Export().(string)
	v6 := vm.Get("v6").Export().(int64)
	k7 := vm.Get("k7").Export().(string)
	v7 := vm.Get("v7").Export().(int64)
	k8 := vm.Get("k8").Export().(string)
	v8 := vm.Get("v8").Export().(int64)
	require.EqualValues(t, "c", k5)
	require.EqualValues(t, 3, v5)
	require.EqualValues(t, "d", k6)
	require.EqualValues(t, 4, v6)
	require.EqualValues(t, "e", k7)
	require.EqualValues(t, 5, v7)
	require.EqualValues(t, "", k8)
	require.EqualValues(t, 0, v8)
}

func TestOrderedIntMapSeekFirstAndLast(t *testing.T) {
	vm := setupGojaVmForOrderedIntMap()
	_, err := vm.RunString(OrderedIntMapSeekFirstAndLastScriptTemplate)
	require.NoError(t, err)

	k1 := vm.Get("k1").Export().(string)
	v1 := vm.Get("v1").Export().(int64)
	k2 := vm.Get("k2").Export().(string)
	v2 := vm.Get("v2").Export().(int64)
	k3 := vm.Get("k3").Export().(string)
	v3 := vm.Get("v3").Export().(int64)
	k4 := vm.Get("k4").Export().(string)
	v4 := vm.Get("v4").Export().(int64)
	k5 := vm.Get("k5").Export().(string)
	v5 := vm.Get("v5").Export().(int64)
	require.EqualValues(t, "a", k1)
	require.EqualValues(t, 1, v1)
	require.EqualValues(t, "b", k2)
	require.EqualValues(t, 2, v2)
	require.EqualValues(t, "c", k3)
	require.EqualValues(t, 3, v3)
	require.EqualValues(t, "d", k4)
	require.EqualValues(t, 4, v4)
	require.EqualValues(t, "e", k5)
	require.EqualValues(t, 5, v5)

	k6 := vm.Get("k6").Export().(string)
	v6 := vm.Get("v6").Export().(int64)
	k7 := vm.Get("k7").Export().(string)
	v7 := vm.Get("v7").Export().(int64)
	k8 := vm.Get("k8").Export().(string)
	v8 := vm.Get("v8").Export().(int64)
	k9 := vm.Get("k9").Export().(string)
	v9 := vm.Get("v9").Export().(int64)
	k10 := vm.Get("k10").Export().(string)
	v10 := vm.Get("v10").Export().(int64)
	require.EqualValues(t, "e", k6)
	require.EqualValues(t, 5, v6)
	require.EqualValues(t, "d", k7)
	require.EqualValues(t, 4, v7)
	require.EqualValues(t, "c", k8)
	require.EqualValues(t, 3, v8)
	require.EqualValues(t, "b", k9)
	require.EqualValues(t, 2, v9)
	require.EqualValues(t, "a", k10)
	require.EqualValues(t, 1, v10)

}

func TestOrderedIntMapClear(t *testing.T) {
	vm := setupGojaVmForOrderedIntMap()
	_, err := vm.RunString(OrderedIntMapClearScriptTemplate)
	require.NoError(t, err)

	len1 := vm.Get("len1").Export().(int64)
	len2 := vm.Get("len2").Export().(int64)
	require.EqualValues(t, 0, len1)
	require.EqualValues(t, 0, len2)
}
