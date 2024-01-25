package undoredo

import "github.com/smac89/undoredo/pkg/classic"

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
	// iterate over undos
	IterUndos(f func(T) bool) error
	// iterate over redos
	IterRedos(f func(T) bool) error
}

var (
	_ EditStack[any] = (*classic.UndoRedoArray[any])(nil)
	_ EditStack[any] = (*classic.UndoRedoStack[any])(nil)
)

func NewUndoRedoArray[T any]() *classic.UndoRedoArray[T] {
	arr := &classic.UndoRedoArray[T]{}
	return arr.Init()
}

func NewUndoRedoStack[T any]() *classic.UndoRedoStack[T] {
	edStack := &classic.UndoRedoStack[T]{}
	return edStack.Init()
}
