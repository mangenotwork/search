package api

import (
	"container/list"
	"github.com/mangenotwork/search/entity"
	"sync"
)

// DocLru 默认缓存 10w个文档， 如果文档大小 1kb 则是 100000kb
var DocLru = NewLru(10 * 1000)

// Lru 文档缓存， LRU算法
// key theme_docId
// value doc
type Lru interface {
	Get(key string) (value *entity.Doc, ok bool)
	GetKeyFromValue(value *entity.Doc) (key string, ok bool)
	Put(key string, value *entity.Doc)
}

type lru struct {
	capacity         int
	doubleLinkedList *list.List
	keyToElement     *sync.Map
	valueToElement   *sync.Map
	mu               *sync.Mutex
}

type lruElement struct {
	key   string
	value *entity.Doc
}

// NewLru initializes a lru cache
func NewLru(cap int) Lru {
	return &lru{
		capacity:         cap,
		doubleLinkedList: list.New(),
		keyToElement:     new(sync.Map),
		valueToElement:   new(sync.Map),
		mu:               new(sync.Mutex),
	}
}

func (l *lru) Get(key string) (value *entity.Doc, ok bool) {
	l.mu.Lock()
	defer l.mu.Unlock()
	if v, ok := l.keyToElement.Load(key); ok {
		element := v.(*list.Element)
		l.doubleLinkedList.MoveToFront(element)
		return element.Value.(*lruElement).value, true
	}
	return nil, false
}

func (l *lru) GetKeyFromValue(value *entity.Doc) (key string, ok bool) {
	l.mu.Lock()
	defer l.mu.Unlock()
	if k, ok := l.valueToElement.Load(value); ok {
		element := k.(*list.Element)
		l.doubleLinkedList.MoveToFront(element)
		return element.Value.(*lruElement).key, true
	}
	return "", false
}

func (l *lru) Put(key string, value *entity.Doc) {
	l.mu.Lock()
	e := &lruElement{"", value}
	if v, ok := l.keyToElement.Load(key); ok {
		element := v.(*list.Element)
		element.Value = e
		l.doubleLinkedList.MoveToFront(element)
	} else {
		element := l.doubleLinkedList.PushFront(e)
		l.keyToElement.Store(key, element)
		l.valueToElement.Store(value, element)
		if l.doubleLinkedList.Len() > l.capacity {
			toBeRemove := l.doubleLinkedList.Back()
			l.doubleLinkedList.Remove(toBeRemove)
			l.keyToElement.Delete(toBeRemove.Value.(*lruElement).key)
			l.valueToElement.Delete(toBeRemove.Value.(*lruElement).value)
		}
	}
	l.mu.Unlock()
}
