package errmsg

import "errors"

var (
	ErrUndoMax    = errors.New("maximum undo reached")
	ErrUndoNMax   = errors.New("n exceeds maximum undo")
	ErrRedoMax    = errors.New("maximum redo reached")
	ErrRedoNMax   = errors.New("n exceeds maximum redo")
	ErrEmptyStack = errors.New("stack is empty")
)
