package lru

type keyValue struct{
	key string
	value Valuer
}

type Valuer interface {
	Len() int
}

