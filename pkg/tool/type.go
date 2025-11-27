package tool

import (
	"sort"
)

type null struct{}

type Array[T comparable] struct {
	values []T
}

type ArrayMap[K, T comparable] struct {
	keys   []K
	values map[K]T
}

type HashSet struct {
	values map[string]null
}

type ArraySet struct {
	keys []string
	hash map[string]bool
}

type HashMap[K, T comparable] struct {
	values map[K]T
}

func NewArray[T comparable]() *Array[T] {
	return &Array[T]{values: make([]T, 0)}
}

func (m *Array[T]) Add(value T) {
	m.values = append(m.values, value)
}

func (m *Array[T]) Adds(value ...T) {
	for _, val := range value {
		m.values = append(m.values, val)
	}
}

func (m *Array[T]) Get(index int) (T, bool) {
	if index < len(m.values) {
		return m.values[index], true
	}
	var n T
	return n, false
}

func (m *Array[T]) Len() int {
	return len(m.values)
}

func (m *Array[T]) Values() []T {
	return m.values
}

func NewHashSet() *HashSet {
	return &HashSet{values: make(map[string]null)}
}

func NewHashSetByValues(values ...string) *HashSet {
	result := &HashSet{values: make(map[string]null)}
	result.Add(values...)
	return result
}

func (m *HashSet) Values() (values []string) {
	for key, _ := range m.values {
		values = append(values, key)
	}
	return
}

func (m *HashSet) Contains(key string) bool {
	if _, ok := m.values[key]; ok {
		return ok
	}
	return false
}

func (m *HashSet) ContainsPrefix(key string) bool {
	l := len(key)
	for tmpKey, _ := range m.values {
		tl := len(tmpKey)
		if l >= tl {
			if key[:tl] == tmpKey {
				return true
			}
		}
	}
	return false
}

func (m *HashSet) Add(value ...string) {
	for _, v := range value {
		m.values[v] = null{}
	}
}

func (m *HashSet) Iterate(fn func(key string) bool) {
	for key, _ := range m.values {
		if fn(key) {
			break
		}
	}
}

func NewArraySet() *ArraySet {
	return &ArraySet{keys: make([]string, 0), hash: make(map[string]bool)}
}

func (m *ArraySet) ContainsKey(key string) bool {
	if _, ok := m.hash[key]; ok {
		return true
	}
	return false
}

func (m *ArraySet) Add(key string) {
	if _, ok := m.hash[key]; !ok {
		m.keys = append(m.keys, key)
		m.hash[key] = true
	}
}

func (m *ArraySet) Adds(keys ...string) {
	for _, key := range keys {
		if _, ok := m.hash[key]; !ok {
			m.keys = append(m.keys, key)
			m.hash[key] = true
		}
	}
}

func (m *ArraySet) Values() (result []string) {
	return m.keys
}

func (m *ArraySet) Iterate(fn func(key string) bool) {
	if fn == nil {
		return
	}
	for _, key := range m.keys {
		if fn(key) {
			break
		}
	}
}

func (m *ArraySet) Sort() {
	sort.Strings(m.keys)
}

func (m *ArrayMap[K, T]) Del(key K) {
	if _, ok := m.values[key]; ok {
		delete(m.values, key)
		for idx, tmpKey := range m.keys {
			if key == tmpKey {
				m.keys = append(m.keys[:idx], m.keys[idx+1:]...)
				break
			}
		}
	}
}

func NewArrayMap[K, T comparable]() *ArrayMap[K, T] {
	return &ArrayMap[K, T]{keys: make([]K, 0), values: make(map[K]T)}
}

func (m *ArrayMap[K, T]) ContainsKey(key K) bool {
	if _, ok := m.values[key]; ok {
		return true
	}
	return false
}

func (m *ArrayMap[K, T]) ContainsValue(value T) bool {
	for _, val := range m.values {
		if val == value {
			return true
		}
	}
	return false
}

func (m *ArrayMap[K, T]) Add(key K, value T) {
	if m.values == nil {
		m.values = make(map[K]T)
	}
	if _, ok := m.values[key]; !ok {
		m.keys = append(m.keys, key)
	}
	m.values[key] = value
}

func (m *ArrayMap[K, T]) Get(key K) T {
	return m.values[key]
}

func (m *ArrayMap[K, T]) Keys() []K {
	return m.keys
}

func (m *ArrayMap[K, T]) Values() (result []T) {
	for _, key := range m.keys {
		result = append(result, m.values[key])
	}
	return
}
func (m *ArrayMap[K, T]) Iterate(fn func(key K, value T) bool) {
	if fn == nil {
		return
	}
	for _, key := range m.keys {
		if fn(key, m.values[key]) {
			break
		}
	}
}

func NewHashMap[K, T comparable]() *HashMap[K, T] {
	return &HashMap[K, T]{values: make(map[K]T)}
}

func (m *HashMap[K, T]) ContainsKey(key K) bool {
	if _, ok := m.values[key]; ok {
		return true
	}
	return false
}

func (m *HashMap[K, T]) ContainsValue(value T) bool {
	for _, val := range m.values {
		if val == value {
			return true
		}
	}
	return false
}

func (m *HashMap[K, T]) Add(key K, value T) {
	if m.values == nil {
		m.values = make(map[K]T)
	}
	m.values[key] = value
}

func (m *HashMap[K, T]) Clear() {
	m.values = make(map[K]T)
}

func (m *HashMap[K, T]) Get(key K) T {
	return m.values[key]
}

func (m *HashMap[K, T]) Values() map[K]T {
	return m.values
}

func (m *HashMap[K, T]) Iterate(fn func(key K, value T) bool) {
	if fn == nil {
		return
	}
	for key, val := range m.values {
		if fn(key, val) {
			break
		}
	}
}
