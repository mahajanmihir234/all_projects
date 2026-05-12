package main

import (
	"fmt"
	"log"

	"main/text_editor/system"
)

func main() {
	editor := system.NewTextEditor()

	if err := editor.Write("Hello"); err != nil {
		log.Fatal(err)
	}
	if err := editor.Write(", World"); err != nil {
		log.Fatal(err)
	}
	fmt.Println(editor.Content())

	for range 5 {
		if err := editor.Undo(); err != nil {
			log.Fatal(err)
		}
		fmt.Println(editor.Content())

		if err := editor.Redo(); err != nil {
			log.Fatal(err)
		}
		fmt.Println(editor.Content())
	}
}
