package src

import (
	"bufio"
	"errors"
	"ini-go/src/entity"
	"os"
)

func GenerateFileFromINIObject(iniObject *entity.INIObject, fileAbsolutePath string) (*os.File, error) {
	if iniObject == nil {
		return nil, errors.New("IniObject should not be null")
	}
	strLines := iniObject.GenerateStringLines()
	if strLines == nil {
		return nil, errors.New("IniObject is empty")
	}
	file, err := os.OpenFile(fileAbsolutePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		return file, err
	}
	defer file.Close()
	writer := bufio.NewWriter(file)
	length := len(strLines)
	for index, str := range strLines {
		writer.WriteString(str)
		if index < length-1 {
			writer.WriteString("\n")
		}
	}
	writer.Flush()
	return file, err
}
