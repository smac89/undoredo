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

type undoRedo[T any] struct {
	currentData T
	previous    *undoRedo[T]
	next        *undoRedo[T]
	hasData     bool
	size        *undoRedoSize
}

var (
	MaxUndoError = errors.New("maximum undo reached")
	MaxRedoError = errors.New("maximum redo reached")
	EmptyStack   = errors.New("stack is empty")
)

func NewEditStack[T any]() EditStack[T] {
	edStack := &undoRedo[T]{size: &undoRedoSize{}}
	return edStack
}

func (u *undoRedo[T]) IsEmpty() bool {
	return u == nil || !u.hasData
}

func (u *undoRedo[T]) Clear() {
	newStk := NewEditStack[T]()
	*u = *(newStk.(*undoRedo[T]))
}

func (u *undoRedo[T]) CanUndo() bool {
	return u.previous != nil
}

func (u *undoRedo[T]) CanRedo() bool {
	return u.next != nil
}

func (u *undoRedo[T]) UndoSize() int {
	return u.size.undoSize
}

func (u *undoRedo[T]) RedoSize() int {
	return u.size.redoSize
}

func (u *undoRedo[T]) Len() int {
	return u.RedoSize() + u.UndoSize()
}

func (u *undoRedo[T]) Peek() (T, error) {
	if u.IsEmpty() {
		return *new(T), EmptyStack
	}
	return u.currentData, nil
}

func (u *undoRedo[T]) Push(change T) {
	prev := *u
	*u = undoRedo[T]{previous: &prev, currentData: change, hasData: true, size: prev.size}
	u.size.undoSize++
	u.size.redoSize = 0
	prev.next = u
}

func (u *undoRedo[T]) Pop() (T, error) {
	return u.Undo()
}

func (u *undoRedo[T]) Undo() (T, error) {
	if !u.CanUndo() {
		return *new(T), MaxUndoError
	}

	u.size.undoSize--
	u.size.redoSize++
	res := *u
	*u = *res.previous
	u.next = &res
	res.previous = u
	return u.Peek()
}

func (u *undoRedo[T]) Redo() (T, error) {
	if !u.CanRedo() {
		return *new(T), MaxRedoError
	}

	u.size.redoSize--
	u.size.undoSize++
	res := *u
	*u = *res.next
	u.previous = &res
	res.next = u
	return u.Peek()
}

func (u *undoRedo[T]) UndoN(n int) ([]T, error) {
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
