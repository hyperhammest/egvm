package types

import (
	"strings"

	"github.com/tinylib/msgp/msgp"
	"modernc.org/b/v2"

	"github.com/smartbch/egvm/script/utils"
)

type OrderedIntMapIter struct {
	e *b.Enumerator[string, int64]
}

func (iter OrderedIntMapIter) Close() {
	iter.e.Close()
}

func (iter OrderedIntMapIter) Next() (string, int64) {
	k, v, err := iter.e.Next()
	if err != nil {
		return "", 0
	}
	return k, v
}

func (iter OrderedIntMapIter) Prev() (string, int64) {
	k, v, err := iter.e.Prev()
	if err != nil {
		return "", 0
	}
	return k, v
}

type OrderedIntMap struct {
	estimatedSize int
	tree          *b.Tree[string, int64]
}

func NewOrderedIntMap() OrderedIntMap {
	return OrderedIntMap{tree: b.TreeNew[string, int64](func(a, b string) int {
		return strings.Compare(a, b)
	})}
}

func (m *OrderedIntMap) loadFrom(b []byte) ([]byte, error) {
	m.tree.Clear()
	initSize := len(b)
	count, b, err := msgp.ReadIntBytes(b)
	if err != nil {
		return nil, err
	}
	for i := 0; i < count; i++ {
		k, v := "", int64(0)
		k, b, err = msgp.ReadStringBytes(b)
		if err != nil {
			return nil, err
		}
		v, b, err = msgp.ReadInt64Bytes(b)
		if err != nil {
			return nil, err
		}
		m.tree.Set(k, v)
	}
	m.estimatedSize = initSize - len(b)
	return b, nil
}

func (m *OrderedIntMap) dumpTo(b []byte) []byte {
	b = msgp.AppendInt(b, m.tree.Len())
	if m.tree.Len() == 0 {
		return b
	}
	e, _ := m.tree.SeekFirst()
	defer e.Close()

	k, v, err := e.Next()
	for err == nil && k != "" {
		b = msgp.AppendString(b, k)
		b = msgp.AppendInt64(b, v)
		k, v, err = e.Next()
	}
	return b
}

func (m *OrderedIntMap) Clear() {
	m.tree.Clear()
	m.estimatedSize = 0
}

func (m *OrderedIntMap) Delete(k string) {
	existed := m.tree.Delete(k)
	if existed {
		m.estimatedSize -= len(k)
	}
}

func (m *OrderedIntMap) Get(k string) (int64, bool) {
	return m.tree.Get(k)
}

func (m *OrderedIntMap) Len() int {
	return m.tree.Len()
}

func (m *OrderedIntMap) Set(k string, v int64) {
	if len(k) == 0 {
		panic(utils.EmptyKeyString)
	}

	m.tree.Put(k, func(_ int64, exists bool) (int64, bool) {
		if !exists {
			m.estimatedSize += 10 + len(k)
		}
		return v, true
	})
}

func (m *OrderedIntMap) Seek(k string) (OrderedIntMapIter, bool) {
	e, ok := m.tree.Seek(k)
	return OrderedIntMapIter{e: e}, ok
}

func (m *OrderedIntMap) SeekFirst() (OrderedIntMapIter, error) {
	e, err := m.tree.SeekFirst()
	return OrderedIntMapIter{e: e}, err
}

func (m *OrderedIntMap) SeekLast() (OrderedIntMapIter, error) {
	e, err := m.tree.SeekLast()
	return OrderedIntMapIter{e: e}, err
}
