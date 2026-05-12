package system

type TextEditor struct {
	document *Document
	history  *History
}

func NewTextEditor() *TextEditor {
	return &TextEditor{
		document: NewDocument(),
		history:  NewHistory(),
	}
}

func (t *TextEditor) Write(text string) error {
	return t.history.Execute(NewWriteAction(t.document, text))
}

func (t *TextEditor) Undo() error {
	return t.history.Undo()
}

func (t *TextEditor) Redo() error {
	return t.history.Redo()
}

func (t *TextEditor) Content() string {
	return t.document.Content()
}
