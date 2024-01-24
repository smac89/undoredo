package undoredo

import "errors"

type EditStack[T any] interface {
	// push a change into the stack
	Push(change T)
	// alias for Undo
	Pop() (T, error)
	// undo the last change to the stack
	Undo() (T, error)
	// undo the last N changes to the stack
	UndoN(n int) ([]T, error)
	// redo the last undo
	Redo() (T, error)
	// redo the last N undos
	RedoN(n int) ([]T, error)
	// check if we can redo
	CanRedo() bool
	// check if we can undo
	CanUndo() bool
	// get the maximum number of undos that can be done
	UndoSize() int
	// get the maximum number of redos that can be done
	RedoSize() int
	// get the number of changes in the stack
	Len() int
	// clear the edit stack. no undo or redo can be done after this
	Clear()
	// get the result of the last change to the stack
	Peek() (T, error)
}

type undoRedoSize struct {
	undoSize int
	redoSize int
}

type UndoRedo[T any] struct {
	currentData T
	previous    *UndoRedo[T]
	next        *UndoRedo[T]
	hasData     bool
	size        *undoRedoSize
}

var (
	ErrUndoMax    = errors.New("maximum undo reached")
	ErrRedoMax    = errors.New("maximum redo reached")
	ErrEmptyStack = errors.New("stack is empty")
)

var (
	_ EditStack[any] = (*UndoRedo[any])(nil)
)

func NewUndoRedo[T any]() *UndoRedo[T] {
	edStack := &UndoRedo[T]{size: &undoRedoSize{}}
	return edStack
}

func (u *UndoRedo[T]) IsEmpty() bool {
	return u == nil || !u.hasData
}

func (u *UndoRedo[T]) Clear() {
	newStk := NewUndoRedo[T]()
	*u = *newStk
}

func (u *UndoRedo[T]) CanUndo() bool {
	return !u.IsEmpty()
}

func (u *UndoRedo[T]) CanRedo() bool {
	return u.next != nil
}

func (u *UndoRedo[T]) UndoSize() int {
	return u.size.undoSize
}

func (u *UndoRedo[T]) RedoSize() int {
	return u.size.redoSize
}

func (u *UndoRedo[T]) Len() int {
	return u.RedoSize() + u.UndoSize()
}

func (u *UndoRedo[T]) Peek() (T, error) {
	if u.IsEmpty() {
		return *new(T), ErrEmptyStack
	}
	return u.currentData, nil
}

func (u *UndoRedo[T]) Push(change T) {
	res := *u
	if u.previous != nil {
		u.previous.next = &res
	}
	*u = UndoRedo[T]{previous: &res, currentData: change, hasData: true, size: res.size}
	u.size.undoSize++
	u.size.redoSize = 0
	res.next = u
}

func (u *UndoRedo[T]) Pop() (T, error) {
	return u.Undo()
}

func (u *UndoRedo[T]) Undo() (T, error) {
	if !u.CanUndo() {
		return *new(T), ErrUndoMax
	}

	u.size.undoSize--
	u.size.redoSize++
	res := *u
	if res.previous != nil {
		res.previous.next = &res
	}
	if res.next != nil {
		res.next.previous = &res
	}
	*u = *res.previous
	return res.currentData, nil
}

func (u *UndoRedo[T]) Redo() (T, error) {
	if !u.CanRedo() {
		return *new(T), ErrRedoMax
	}

	u.size.redoSize--
	u.size.undoSize++
	res := *u
	if res.previous != nil {
		res.previous.next = &res
	}
	if res.next != nil {
		res.next.previous = &res
	}
	*u = *res.next
	return u.currentData, nil
}

func (u *UndoRedo[T]) UndoN(n int) ([]T, error) {
	var (
		undo  = u
		undos = make([]T, 0)
	)
	for i := 0; i < n-1; i++ {
		if undo.CanUndo() {
			undos = append(undos, undo.currentData)
			undo.size.redoSize++
			undo.size.undoSize--
			undo = undo.previous
		} else {
			break
		}
	}
	if u.previous != nil {
		u.previous.next = u.next
	}
	if u.next != nil {
		u.next.previous = u.previous
	}
	*u = *undo
	d, err := u.Undo()
	if err == nil {
		undos = append(undos, d)
	}
	return undos, err
}

func (u *undoRedo[T]) RedoN(n int) ([]T, error) {
	var (
		redo  = u
		redos = make([]T, 0)
	)
	for i := 0; i < n-1; i++ {
		if redo.CanRedo() {
			redos = append(redos, redo.currentData)
			redo.size.redoSize--
			redo.size.undoSize++
			redo = redo.next
		} else {
			break
		}
	}
	if u.previous != nil {
		u.previous.next = u.next
	}
	if u.next != nil {
		u.next.previous = u.previous
	}
	*u = *redo
	d, err := u.Redo()
	if err == nil {
		redos = append(redos, d)
	}
	return redos, err
}
