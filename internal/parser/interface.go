package parser

type Fetcher[T any, S any] interface {
	GetById(id string) (T, error)
	GetAll(globalType int, page int) ([]S, error)
}
