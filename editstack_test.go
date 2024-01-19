package undoredo

import (
	"slices"
	"testing"
)

func TestUndoFailsWhenEmpty(t *testing.T) {
	u := NewEditStack[int]()
	if _, err := u.Undo(); err == nil {
		t.Error("expected error on empty stack")
	}
}

func TestRedoFailsWhenEmpty(t *testing.T) {
	u := NewEditStack[int]()
	if _, err := u.Redo(); err == nil {
		t.Error("expected error on empty stack")
	}
}

func TestCannotPeekWhenEmpty(t *testing.T) {
	u := NewEditStack[int]()
	if _, err := u.Peek(); err == nil {
		t.Error("expected error on empty stack")
	}
}

func TestCanPeekAfterPush(t *testing.T) {
	u := NewEditStack[int]()
	u.Push(1)
	if _, err := u.Peek(); err != nil {
		t.Error("expected to not be able to push")
	}
}

func TestCanUndo(t *testing.T) {
	u := NewEditStack[int]()
	u.Push(1)
	if !u.CanUndo() {
		t.Error("expected to be able to undo")
	}
}

func TestCanRedo(t *testing.T) {
	u := NewEditStack[int]()
	u.Push(1)
	u.Undo()
	if !u.CanRedo() {
		t.Error("expected to be able to redo")
	}
}

func TestUndoSizeIncreasesWithPush(t *testing.T) {
	u := NewEditStack[int]()
	for i := 0; i < 10; i++ {
		u.Push(1)
	}
	if u.UndoSize() != 10 {
		t.Error("expected undo size to be 11")
	}
}

func TestRedoSizeIncreasesWithUndo(t *testing.T) {
	u := NewEditStack[int]()
	for i := 0; i < 10; i++ {
		u.Push(1)
	}
	if u.RedoSize() != 0 {
		t.Error("expected redo size to be 0 when no undo")
	}

	for i := 0; i < 10; i++ {
		u.Undo()
	}
	if u.RedoSize() != 10 {
		t.Error("expected redo size to be 10")
	}
}

func TestUndoNUndoesLastNChanges(t *testing.T) {
	u := NewEditStack[int]()
	for i := 0; i < 10; i++ {
		u.Push(1)
	}
	res, _ := u.UndoN(5)
	if !slices.Equal(res, []int{1, 1, 1, 1, 1}) {
		t.Error("expected to undo 5 changes")
	}
	if u.UndoSize() != 5 {
		t.Error("expected remaining undo to be 5")
	}
}

func TestUndoAllChanges(t *testing.T) {
	u := NewEditStack[int]()
	for i := 0; i < 10; i++ {
		u.Push(1)
	}
	_, _ = u.UndoN(u.UndoSize())
	if u.UndoSize() != 0 {
		t.Error("expected remaining undo to be 0")
	}
}

func TestRedoNRedoesLastNUndos(t *testing.T) {
	u := NewEditStack[int]()
	for i := 0; i < 10; i++ {
		u.Push(1)
	}
	u.UndoN(5)
	res, _ := u.RedoN(5)
	if len(res) != 5 {
		t.Error("expected to redo 5 undos")
	}
}

func TestRedoAllUndos(t *testing.T) {
	u := NewEditStack[int]()
	for i := 0; i < 10; i++ {
		u.Push(1)
	}
	u.UndoN(u.UndoSize())
	res, _ := u.RedoN(u.RedoSize())
	if len(res) != 10 {
		t.Error("expected to redo 5 undos")
	}
}

func TestLenIncreasesWithPush(t *testing.T) {
	u := NewEditStack[int]()
	u.Push(1)
	if u.Len() != 1 {
		t.Error("expected undo size to be 1")
	}
	for i := 0; i < 10; i++ {
		u.Push(1)
	}
	if u.Len() != 11 {
		t.Error("expected undo size to be 11")
	}
}

func TestLenDoesNotChangeWithUndoOrRedo(t *testing.T) {
	u := NewEditStack[int]()
	for i := 0; i < 10; i++ {
		u.Push(1)
	}
	if u.Len() != 10 {
		t.Error("expected stack to have size 10")
	}
	for i := 0; i < 10; i++ {
		u.Undo()
	}
	if u.Len() != 10 {
		t.Error("expected stack size to remain at 10")
	}
}
