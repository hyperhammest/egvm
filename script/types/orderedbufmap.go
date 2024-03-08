package types

import (
	"strings"

	"github.com/dop251/goja"
	"github.com/tinylib/msgp/msgp"
	"modernc.org/b/v2"

	"github.com/smartbch/egvm/script/utils"
)

type OrderedBufMapIter struct {
	e *b.Enumerator[string, []byte]
}

func (iter OrderedBufMapIter) Close() {
	iter.e.Close()
}

func (iter OrderedBufMapIter) Next(f goja.FunctionCall, vm *goja.Runtime) goja.Value {
	if len(f.Arguments) != 0 {
		panic(utils.IncorrectArgumentCount)
	}
	var result [2]any
	k, v, err := iter.e.Next()
	if err != nil {
		result = [2]any{"", nil}
	} else {
		result = [2]any{k, vm.NewArrayBuffer(v)}
	}
	return vm.ToValue(result)
}

func (iter OrderedBufMapIter) Prev(f goja.FunctionCall, vm *goja.Runtime) goja.Value {
	if len(f.Arguments) != 0 {
		panic(utils.IncorrectArgumentCount)
	}
	var result [2]any
	k, v, err := iter.e.Prev()
	if err != nil {
		result = [2]any{"", nil}
	} else {
		result = [2]any{k, vm.NewArrayBuffer(v)}
	}
	return vm.ToValue(result)
}

type OrderedBufMap struct {
	estimatedSize int
	tree          *b.Tree[string, []byte]
}

func NewOrderedBufMap() OrderedBufMap {
	return OrderedBufMap{tree: b.TreeNew[string, []byte](func(a, b string) int {
		return strings.Compare(a, b)
	})}
}

func (m *OrderedBufMap) loadFrom(b []byte) ([]byte, error) {
	m.tree.Clear()
	initSize := len(b)
	count, b, err := msgp.ReadIntBytes(b)
	if err != nil {
		return nil, err
	}
	for i := 0; i < count; i++ {
		var v []byte
		k := ""
		k, b, err = msgp.ReadStringBytes(b)
		if err != nil {
			return nil, err
		}
		v, b, err = msgp.ReadBytesBytes(b, nil)
		if err != nil {
			return nil, err
		}
		m.tree.Set(k, v)
	}
	m.estimatedSize = initSize - len(b)
	return b, nil
}

func (m *OrderedBufMap) dumpTo(b []byte) []byte {
	b = msgp.AppendInt(b, m.tree.Len())
	if m.tree.Len() == 0 {
		return b
	}
	e, _ := m.tree.SeekFirst()
	defer e.Close()

	k, v, err := e.Next()
	for err == nil && k != "" {
		b = msgp.AppendString(b, k)
		b = msgp.AppendBytes(b, v)
		k, v, err = e.Next()
	}
	return b
}

func (m *OrderedBufMap) Clear() {
	m.tree.Clear()
	m.estimatedSize = 0
}

func (m *OrderedBufMap) Delete(k string) {
	existed := m.tree.Delete(k)
	if existed {
		m.estimatedSize -= len(k)
	}
}

func (m *OrderedBufMap) Get(f goja.FunctionCall, vm *goja.Runtime) goja.Value {
	if len(f.Arguments) != 1 {
		panic(utils.IncorrectArgumentCount)
	}
	k, ok := f.Arguments[0].Export().(string)
	if !ok {
		panic(goja.NewSymbol("The first argument must be string"))
	}
	v, ok := m.tree.Get(k)
	return vm.ToValue([2]any{vm.NewArrayBuffer(v), ok})
}

func (m *OrderedBufMap) Len() int {
	return m.tree.Len()
}

func (m *OrderedBufMap) Set(k string, buf goja.ArrayBuffer) {
	if len(k) == 0 {
		panic(utils.EmptyKeyString)
	}

	v := buf.Bytes()
	m.tree.Put(k, func(oldV []byte, exists bool) (newV []byte, write bool) {
		if exists {
			m.estimatedSize += len(v) - len(oldV)
		} else {
			m.estimatedSize += 10 + len(k) + len(v)
		}
		return v, true
	})
}

func (m *OrderedBufMap) Seek(k string) (OrderedBufMapIter, bool) {
	e, ok := m.tree.Seek(k)
	return OrderedBufMapIter{e: e}, ok
}

func (m *OrderedBufMap) SeekFirst() (OrderedBufMapIter, error) {
	e, err := m.tree.SeekFirst()
	return OrderedBufMapIter{e: e}, err
}

func (m *OrderedBufMap) SeekLast() (OrderedBufMapIter, error) {
	e, err := m.tree.SeekLast()
	return OrderedBufMapIter{e: e}, err
}
