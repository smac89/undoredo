# undoredo
[![Go Reference](https://pkg.go.dev/badge/github.com/smac89/undoredo.svg)](https://pkg.go.dev/github.com/smac89/undoredo)

![stacks](./docs/stacks.jpg)

Robust implementation of a data structure which allows for recording changes and replaying them similar to how a text editor undo/redo history works

### Rationale
We all love `Ctrl + Z` and `Ctrl + Shift + Z`, but do we know how they work internally? Atleast on the level of data stuctures. Inspired by a recent @ThePrimeagen's [video](https://youtu.be/yeatOU5vVsA), I decided to implement such a datastructure to see if I can reason my way to a working solution.

### EditStack[T]
I present to you: `EditStack`

| Operation        | Description |
| ---              | ---         |
| `Undo() (T, error)` | undo the last change to the stack |
| `Redo() (T, error)` | redo the last undo |
| `Push(change T)` | push a change into the stack |
| `Pop() (T, error)` | alias for Undo |
| `Peek() (T, error)` | get the latest change at the top of the stack |
| `UndoN(n int) ([]T, error)` | undo the last N changes to the stack |
| `RedoN(n int) ([]T, error)` | redo the last N undos |
| `CanUndo() bool` | check if we can undo |
| `CanRedo() bool` | check if we can redo |
| `UndoSize() int` | get the maximum number of undos that can be done |
| `RedoSize() int` | get the maximum number of redos that can be done |
| `Len() int` | get the number of changes in the stack (undos + redos) |
| `Clear()` | clear the edit stack. no undo or redo can be done after this |
| `IterUndos(func (T) bool) error` | iterate the remaining undos |
| `IterRedos(func (T) bool) error` | iterate the remaining redos |

Other operations which can be generalized from the above:

- `UndoAll() ([]T, error)` can be achieved by `UndoN(UndoSize())`
- `RedoAll() ([]T, error)` can be achieved by `RedoN(RedoSize())`

### Sample usage
```go
import (
	"errors"

	"github.com/smac89/undoredo"
)

func main() {
	udrd := undoredo.NewUndoRedo[int]()
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
	println(errors.Is(err, undoredo.ErrRedoMax)) // true
}
```

### Future improvements
- 🚀 Allow setting maximum stack size and implement eviction strategy
- 🚀 Consider an array-based implementation which may help with eviction and saves memory
