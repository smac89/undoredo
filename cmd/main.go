package main

import (
	"errors"

	"github.com/smac89/undoredo"
	"github.com/smac89/undoredo/pkg/errmsg"
)

func main() {
	udrd := undoredo.NewUndoRedoStack[int]()
	udrd.Push(1)
	udrd.Push(2)
	udrd.Push(3)
	var undo, redo, data int
	_, _ = udrd.Undo()
	undo, _ = udrd.Undo()
	redo, _ = udrd.Redo()
	println(undo, redo) // 2 2
	redo, _ = udrd.Redo()
	println(undo, redo) // 2 3
	udrd.Push(4)
	data, _ = udrd.Peek()
	println(data) // 4
	udrd.Undo()
	udrd.Undo()
	udrd.Undo()
	println(udrd.CanRedo()) // true
	udrd.Undo()
	println(udrd.CanUndo()) // false
	udrd.Push(5)
	println(udrd.CanRedo()) // false
	data, _ = udrd.Peek()
	println(data) // 5
	_, err := udrd.Redo()
	println(errors.Is(err, errmsg.ErrRedoMax)) // true
}
