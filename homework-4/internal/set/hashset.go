package set

type Set[T comparable] interface {
	Contains(item T) bool
	Add(item T)
	Remove(item T)
	RemoveAll()
	Size() uint32
}

type HashSet[T comparable] map[T]struct{}

func NewHashSet[T comparable]() HashSet[T] {
	return make(HashSet[T])
}

func (set HashSet[T]) Contains(item T) bool {
	_, found := set[item]
	return found
}

func (set HashSet[T]) Add(item T) {
	set[item] = struct{}{}
}

func (set HashSet[T]) Remove(item T) {
	delete(set, item)
}

func (set HashSet[T]) RemoveAll() {
	for item := range set {
		delete(set, item)
	}
}

func (set HashSet[T]) Size() uint32 {
	return uint32(len(set))
}
