package middleware

type Middlewares[T any] []func(func(T) error) func(T) error

func New[T any]() Middlewares[T] {
	return make(Middlewares[T], 0)
}

func (mws Middlewares[T]) Wrap(hf func(T) error) func(T) error {
	for _, mw := range mws {
		hf = mw(hf)
	}
	return hf
}

func (mws Middlewares[T]) Use(f func(func(T) error) func(T) error) Middlewares[T] {
	newmws := make([]func(func(T) error) func(T) error, len(mws)+1)
	newmws[0] = f
	copy(newmws[1:], mws)
	return newmws
}
