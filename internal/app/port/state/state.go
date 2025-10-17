package state

type State[T any] interface {
	Set(value T)
	Get() []T
}
