package system

type Action interface {
	Apply() error
	Revert() error
}

type WriteAction struct {
	document *Document
	text     string
}

func NewWriteAction(document *Document, text string) *WriteAction {
	return &WriteAction{
		document: document,
		text:     text,
	}
}

func (w *WriteAction) Apply() error {
	w.document.Append(w.text)
	return nil
}

func (w *WriteAction) Revert() error {
	return w.document.DeleteLast(len(w.text))
}
