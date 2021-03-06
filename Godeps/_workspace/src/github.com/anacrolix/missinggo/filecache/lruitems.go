package filecache

import (
	"container/list"
	"io"

	"gx/ipfs/QmVgtwPh2NNoZTSyYkr4Y3epaYACBKf26r8hV6EFA7xS6c/b"
)

type Iterator interface {
	Next() Iterator
	Value() interface{}
}

type listElementIterator struct {
	le *list.Element
}

func (me listElementIterator) Next() Iterator {
	e := me.le.Next()
	if e == nil {
		return nil
	}
	return listElementIterator{e}
}

func (me listElementIterator) Value() interface{} {
	return me.le.Value
}

func newLRUItems() *lruItems {
	return &lruItems{b.TreeNew(func(_a, _b interface{}) int {
		a := _a.(ItemInfo)
		b := _b.(ItemInfo)
		if a.Accessed != b.Accessed {
			if a.Accessed.Before(b.Accessed) {
				return -1
			} else {
				return 1
			}
		}
		if a.Path == b.Path {
			return 0
		}
		if a.Path < b.Path {
			return -1
		}
		return 1
	})}
}

type lruItems struct {
	tree *b.Tree
}

type bEnumeratorIterator struct {
	e *b.Enumerator
	v ItemInfo
}

func (me bEnumeratorIterator) Next() Iterator {
	_, v, err := me.e.Next()
	if err == io.EOF {
		return nil
	}
	return bEnumeratorIterator{me.e, v.(ItemInfo)}
}

func (me bEnumeratorIterator) Value() interface{} {
	return me.v
}

func (me *lruItems) Front() Iterator {
	e, _ := me.tree.SeekFirst()
	if e == nil {
		return nil
	}
	return bEnumeratorIterator{
		e: e,
	}.Next()
}

func (me *lruItems) LRU() ItemInfo {
	_, v := me.tree.First()
	return v.(ItemInfo)
}

func (me *lruItems) Insert(ii ItemInfo) {
	me.tree.Set(ii, ii)
}

func (me *lruItems) Remove(ii ItemInfo) bool {
	return me.tree.Delete(ii)
}
