package pcommon

type Else[T any] interface {
	ElseDo(fn func() T) T
	Else(value T) T
}

type Then[T any] interface {
	ThenDo(fn func() T) Else[T]
	Then(value T) Else[T]
}

type Condition[T any] struct {
	condition bool
	thenValue T
	thenFn    func() T
}

func When[T any](condition bool) Then[T] {
	return &Condition[T]{condition: condition}
}

func (c *Condition[T]) ThenDo(fn func() T) Else[T] {
	c.thenFn = fn
	return c
}

func (c *Condition[T]) Then(value T) Else[T] {
	c.thenValue = value
	return c
}

func (c *Condition[T]) ElseDo(fn func() T) T {
	if c.condition {
		return c.then()
	}

	return fn()
}

func (c *Condition[T]) Else(value T) T {
	if c.condition {
		return c.then()
	}

	return value
}

func (c *Condition[T]) then() T {
	if c.thenFn != nil {
		return c.thenFn()
	}
	return c.thenValue
}
