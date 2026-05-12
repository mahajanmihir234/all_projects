package system

import "errors"

type Document struct {
	content string
}

func NewDocument() *Document {
	return &Document{content: ""}
}

func (d *Document) Append(text string) {
	d.content += text
}

func (d *Document) DeleteLast(n int) error {
	if n < 0 || n > len(d.content) {
		return errors.New("delete length exceeds document content")
	}

	d.content = d.content[:len(d.content)-n]
	return nil
}

func (d *Document) Content() string {
	return d.content
}
