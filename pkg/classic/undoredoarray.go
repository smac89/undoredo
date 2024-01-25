package classic

import (
	"github.com/smac89/undoredo/pkg/errmsg"
)

type UndoRedoArray[T any] struct {
	currentPos int
	data       []T
}

func (u *UndoRedoArray[T]) Init() *UndoRedoArray[T] {
	u.currentPos = 0
	u.data = make([]T, 0)
	return u
}

func (u *UndoRedoArray[T]) IsEmpty() bool {
	return u == nil || u.currentPos < 1
}

func (u *UndoRedoArray[T]) Clear() {
	u.Init()
}

func (u *UndoRedoArray[T]) CanUndo() bool {
	return !u.IsEmpty()
}

func (u *UndoRedoArray[T]) CanRedo() bool {
	return u.currentPos < len(u.data)
}

func (u *UndoRedoArray[T]) UndoSize() int {
	return u.currentPos
}

func (u *UndoRedoArray[T]) RedoSize() int {
	return len(u.data) - u.currentPos
}

func (u *UndoRedoArray[T]) Len() int {
	return len(u.data)
}

func (u *UndoRedoArray[T]) Peek() (T, error) {
	if u.IsEmpty() {
		return *new(T), errmsg.ErrEmptyStack
	}
	return u.data[u.currentPos-1], nil
}

func (u *UndoRedoArray[T]) Push(change T) {
	u.data = append(u.data[:u.currentPos], change)
	u.currentPos = len(u.data)
}

func (u *UndoRedoArray[T]) Pop() (T, error) {
	return u.Undo()
}

func (u *UndoRedoArray[T]) Undo() (T, error) {
	if !u.CanUndo() {
		return *new(T), errmsg.ErrEmptyStack
	}
	u.currentPos--
	return u.data[u.currentPos], nil
}

func (u *UndoRedoArray[T]) Redo() (T, error) {
	if !u.CanRedo() {
		return *new(T), errmsg.ErrEmptyStack
	}
	res := u.data[u.currentPos]
	u.currentPos++
	return res, nil
}

func (u *UndoRedoArray[T]) UndoN(n int) ([]T, error) {
	var (
		undos      = make([]T, 0)
		currentPos = u.currentPos
	)
	if n > u.UndoSize() {
		return undos, errmsg.ErrUndoMax
	} else if n < 0 {
		return undos, nil
	}
	for i := 0; i < n; i++ {
		if u.CanUndo() {
			if undo, err := u.Undo(); err != nil {
				u.currentPos = currentPos
				return []T{}, err
			} else {
				undos = append(undos, undo)
			}
		} else {
			u.currentPos = currentPos
			break
		}
	}
	return undos, nil
}

func (u *UndoRedoArray[T]) RedoN(n int) ([]T, error) {
	var (
		redos      = make([]T, 0)
		currentPos = u.currentPos
	)
	if n > u.RedoSize() {
		return redos, errmsg.ErrRedoMax
	} else if n < 0 {
		return redos, nil
	}
	for i := 0; i < n; i++ {
		if u.CanRedo() {
			if redo, err := u.Redo(); err != nil {
				u.currentPos = currentPos
				return []T{}, err
			} else {
				redos = append(redos, redo)
			}
		} else {
			u.currentPos = currentPos
			break
		}
	}
	return redos, nil
}

func (u *UndoRedoArray[T]) IterUndos(f func(T) bool) error {
	if !u.CanUndo() {
		return errmsg.ErrUndoMax
	}
	undo := *u
	for undo.CanUndo() {
		val, _ := undo.Undo()
		if !f(val) {
			return nil
		}
	}
	return nil
}

func (u *UndoRedoArray[T]) IterRedos(f func(T) bool) error {
	if !u.CanRedo() {
		return errmsg.ErrRedoMax
	}
	redo := *u
	for redo.CanRedo() {
		val, _ := redo.Redo()
		if !f(val) {
			return nil
		}
	}
	return nil
}
