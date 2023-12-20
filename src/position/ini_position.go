package position

type INIPosition struct {
	FileLocation string
	LineNumber   int
	CharBegin    int
	CharEnd      int
}

func NewPosition(fileLocation string, lineNumber int, charBegin int, charEnd int) *INIPosition {
	iniPosition := INIPosition{
		FileLocation: fileLocation,
		LineNumber:   lineNumber,
		CharBegin:    charBegin,
		CharEnd:      charEnd,
	}
	return &iniPosition
}
