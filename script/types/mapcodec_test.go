package types

import (
	"fmt"
	"testing"

	"github.com/dop251/goja"
	gethcmn "github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
)

const (
	SerializeAndDeserializeMapsScriptTemplate = `
		let im = NewOrderedIntMap()
		im.Set('a', 1)
		im.Set('b', 2)
		im.Set('c', 3)

		let sm = NewOrderedStrMap()
		sm.Set('d', 'a')
		sm.Set('e', 'b')
		sm.Set('f', 'c')

		let bm = NewOrderedBufMap()
		let buffer1 = new ArrayBuffer(8); // 8 bytes
		let buffer2 = new ArrayBuffer(8); // 8 bytes
		let buffer3 = new ArrayBuffer(8); // 8 bytes
		let view1 = new Uint8Array(buffer1);
		view1[7] = 1
		let view2 = new Uint8Array(buffer2);
		view2[7] = 2
		let view3 = new Uint8Array(buffer3);
		view3[7] = 3
		bm.Set('g', buffer1)
		bm.Set('h', buffer2)
		bm.Set('i', buffer3)

		const bz = SerializeMaps(im, sm, bm)

		let [m1, remainBz1] = DeserializeMap(bz, 0)
		const [v1, ok1] = m1.Get('a')
		const len1 = m1.Len()

		let [m2, remainBz2] = DeserializeMap(remainBz1, 1)
		const [v2, ok2] = m2.Get('f')
		const len2 = m2.Len()

		let [m3, remainBz3] = DeserializeMap(remainBz2, 2)
		const [v3, ok3] = m3.Get('h')
		const len3 = m3.Len()
	`

	OrderedMapReaderScriptTemplate = `
		let im = NewOrderedIntMap()
		im.Set('a', 1)
		im.Set('b', 2)
		im.Set('c', 3)

		let sm = NewOrderedStrMap()
		sm.Set('d', 'a')
		sm.Set('e', 'b')
		sm.Set('f', 'c')

		let bm = NewOrderedBufMap()
		let buffer1 = new ArrayBuffer(8); // 8 bytes
		let buffer2 = new ArrayBuffer(8); // 8 bytes
		let buffer3 = new ArrayBuffer(8); // 8 bytes
		let view1 = new Uint8Array(buffer1);
		view1[7] = 1
		let view2 = new Uint8Array(buffer2);
		view2[7] = 2
		let view3 = new Uint8Array(buffer3);
		view3[7] = 3
		bm.Set('g', buffer1)
		bm.Set('h', buffer2)
		bm.Set('i', buffer3)

		const bz = SerializeMaps(im, sm, bm)
		
		let mapReader = NewOrderedMapReader(bz)
		let m1 = mapReader.Read(0)
		const [v1, ok1] = m1.Get('a')
		const len1 = m1.Len()

		let m2 = mapReader.Read(1)
		const [v2, ok2] = m2.Get('f')
		const len2 = m2.Len()

		let m3 = mapReader.Read(2)
		const [v3, ok3] = m3.Get('h')
		const len3 = m3.Len()
	`
)

func setupGojaVmForCodec() *goja.Runtime {
	vm := goja.New()
	vm.Set("SerializeMaps", SerializeMaps)
	vm.Set("DeserializeMap", DeserializeMap)
	vm.Set("NewOrderedMapReader", NewOrderedMapReader)
	vm.Set("NewOrderedIntMap", NewOrderedIntMap)
	vm.Set("NewOrderedStrMap", NewOrderedStrMap)
	vm.Set("NewOrderedBufMap", NewOrderedBufMap)
	return vm
}

func TestSerializeAndDeserialize(t *testing.T) {
	vm := setupGojaVmForCodec()
	_, err := vm.RunString(SerializeAndDeserializeMapsScriptTemplate)
	require.NoError(t, err)

	bz := vm.Get("bz").Export().(goja.ArrayBuffer)
	fmt.Printf("bz: %v\n", bz.Bytes())

	ok1 := vm.Get("ok1").Export().(bool)
	v1 := vm.Get("v1").Export().(int64)
	len1 := vm.Get("len1").Export().(int64)
	require.True(t, ok1)
	require.EqualValues(t, 1, v1)
	require.EqualValues(t, 3, len1)

	ok2 := vm.Get("ok2").Export().(bool)
	v2 := vm.Get("v2").Export().(string)
	len2 := vm.Get("len2").Export().(int64)
	require.True(t, ok2)
	require.EqualValues(t, 'c', v2)
	require.EqualValues(t, 3, len2)

	ok3 := vm.Get("ok3").Export().(bool)
	v3 := vm.Get("v3").Export().(goja.ArrayBuffer)
	v3Hex := gethcmn.Bytes2Hex(v3.Bytes())
	len3 := vm.Get("len3").Export().(int64)
	require.True(t, ok3)
	require.EqualValues(t, "0000000000000002", v3Hex)
	require.EqualValues(t, 3, len3)

	remainBz3 := vm.Get("remainBz3").Export().(goja.ArrayBuffer)
	require.EqualValues(t, 0, len(remainBz3.Bytes()))
}

func TestOrderedMapReader(t *testing.T) {
	vm := setupGojaVmForCodec()
	_, err := vm.RunString(OrderedMapReaderScriptTemplate)
	require.NoError(t, err)

	bz := vm.Get("bz").Export().(goja.ArrayBuffer)
	fmt.Printf("bz: %v\n", bz.Bytes())

	ok1 := vm.Get("ok1").Export().(bool)
	v1 := vm.Get("v1").Export().(int64)
	len1 := vm.Get("len1").Export().(int64)
	require.True(t, ok1)
	require.EqualValues(t, 1, v1)
	require.EqualValues(t, 3, len1)

	ok2 := vm.Get("ok2").Export().(bool)
	v2 := vm.Get("v2").Export().(string)
	len2 := vm.Get("len2").Export().(int64)
	require.True(t, ok2)
	require.EqualValues(t, 'c', v2)
	require.EqualValues(t, 3, len2)

	ok3 := vm.Get("ok3").Export().(bool)
	v3 := vm.Get("v3").Export().(goja.ArrayBuffer)
	v3Hex := gethcmn.Bytes2Hex(v3.Bytes())
	len3 := vm.Get("len3").Export().(int64)
	require.True(t, ok3)
	require.EqualValues(t, "0000000000000002", v3Hex)
	require.EqualValues(t, 3, len3)

}
