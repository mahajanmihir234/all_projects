package system

import "errors"

type History struct {
	undoStack []Action
	redoStack []Action
}

func NewHistory() *History {
	return &History{
		undoStack: []Action{},
		redoStack: []Action{},
	}
}

func (h *History) Execute(action Action) error {
	if err := action.Apply(); err != nil {
		return err
	}

	h.undoStack = append(h.undoStack, action)
	h.redoStack = []Action{}
	return nil
}

func (h *History) Undo() error {
	if len(h.undoStack) == 0 {
		return errors.New("nothing to undo")
	}

	action := h.undoStack[len(h.undoStack)-1]
	h.undoStack = h.undoStack[:len(h.undoStack)-1]

	if err := action.Revert(); err != nil {
		return err
	}

	h.redoStack = append(h.redoStack, action)
	return nil
}

func (h *History) Redo() error {
	if len(h.redoStack) == 0 {
		return errors.New("nothing to redo")
	}

	action := h.redoStack[len(h.redoStack)-1]
	h.redoStack = h.redoStack[:len(h.redoStack)-1]

	if err := action.Apply(); err != nil {
		return err
	}

	h.undoStack = append(h.undoStack, action)
	return nil
}
