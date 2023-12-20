package main

import (
	"fmt"
	"ini-go/src"
)

type TTT struct {
	XX []string
}

func main() {
	iniObject, err := src.ParseFileToINIObject("test-input.ini")
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(iniObject)
		src.GenerateFileFromINIObject(iniObject, "test-output.ini")
	}
}
