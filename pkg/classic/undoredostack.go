package classic

import (
	"github.com/smac89/undoredo/pkg/errmsg"
)

type undoRedoSize struct {
	undoSize int
	redoSize int
}

type UndoRedoStack[T any] struct {
	currentData T
	previous    *UndoRedoStack[T]
	next        *UndoRedoStack[T]
	hasData     bool
	size        *undoRedoSize
}

func (u *UndoRedoStack[T]) Init() *UndoRedoStack[T] {
	u.size = &undoRedoSize{}
	return u
}

func (u *UndoRedoStack[T]) IsEmpty() bool {
	return u == nil || !u.hasData
}

func (u *UndoRedoStack[T]) Clear() {
	*u = UndoRedoStack[T]{size: &undoRedoSize{}}
}

func (u *UndoRedoStack[T]) CanUndo() bool {
	return !u.IsEmpty()
}

func (u *UndoRedoStack[T]) CanRedo() bool {
	return u.next != nil
}

func (u *UndoRedoStack[T]) UndoSize() int {
	return u.size.undoSize
}

func (u *UndoRedoStack[T]) RedoSize() int {
	return u.size.redoSize
}

func (u *UndoRedoStack[T]) Len() int {
	return u.RedoSize() + u.UndoSize()
}

func (u *UndoRedoStack[T]) Peek() (T, error) {
	if u.IsEmpty() {
		return *new(T), errmsg.ErrEmptyStack
	}
	return u.currentData, nil
}

func (u *UndoRedoStack[T]) Push(change T) {
	res := *u
	if u.previous != nil {
		u.previous.next = &res
	}
	*u = UndoRedoStack[T]{previous: &res, currentData: change, hasData: true, size: res.size}
	u.size.undoSize++
	u.size.redoSize = 0
	res.next = u
}

func (u *UndoRedoStack[T]) Pop() (T, error) {
	return u.Undo()
}

func (u *UndoRedoStack[T]) Undo() (T, error) {
	if !u.CanUndo() {
		return *new(T), errmsg.ErrUndoMax
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

func (u *UndoRedoStack[T]) Redo() (T, error) {
	if !u.CanRedo() {
		return *new(T), errmsg.ErrRedoMax
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

func (u *UndoRedoStack[T]) UndoN(n int) ([]T, error) {
	var (
		undo  = u
		undos = make([]T, 0)
	)
	if n > u.UndoSize() {
		return undos, errmsg.ErrUndoNMax
	} else if n < 0 {
		return undos, nil
	}
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
	res := *u
	if res.previous != nil {
		res.previous.next = &res
	}
	if res.next != nil {
		res.next.previous = &res
	}
	d, err := undo.Undo()
	if err == nil {
		undos = append(undos, d)
		*u = *undo
	}
	return undos, err
}

func (u *UndoRedoStack[T]) RedoN(n int) ([]T, error) {
	var (
		redo  = u
		redos = make([]T, 0)
	)
	if n > u.RedoSize() {
		return redos, errmsg.ErrRedoNMax
	} else if n < 0 {
		return redos, nil
	}
	for i := 0; i < n-1; i++ {
		if redo.CanRedo() {
			redo.size.redoSize--
			redo.size.undoSize++
			redo = redo.next
			redos = append(redos, redo.currentData)
		} else {
			break
		}
	}
	res := *u
	if res.previous != nil {
		res.previous.next = &res
	}
	if res.next != nil {
		res.next.previous = &res
	}
	d, err := redo.Redo()
	if err == nil {
		redos = append(redos, d)
		*u = *redo
	}
	return redos, err
}

func (u *UndoRedoStack[T]) IterUndos(f func(T) bool) error {
	if !u.CanUndo() {
		return errmsg.ErrUndoMax
	}
	undo := u
	for undo.CanUndo() {
		if !f(undo.currentData) {
			return nil
		}
		undo = undo.previous
	}
	return nil
}

func (u *UndoRedoStack[T]) IterRedos(f func(T) bool) error {
	if !u.CanRedo() {
		return errmsg.ErrRedoMax
	}
	redo := u
	for redo.CanRedo() {
		redo = redo.next
		if !f(redo.currentData) {
			return nil
		}
	}
	return nil
}
