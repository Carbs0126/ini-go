package atom

import "ini-go/src/position"

type INIComment struct {
	INIPosition *position.INIPosition
	Comment     string
}

func NewINIComment(comment string, iniPosition *position.INIPosition) *INIComment {
	iniComment := INIComment{
		Comment:     comment,
		INIPosition: iniPosition,
	}
	return &iniComment
}

func (iniComment *INIComment) GetPosition() *position.INIPosition {
	return iniComment.INIPosition
}

func (iniComment *INIComment) ToINIOutputString() (string, error) {
	return iniComment.Comment, nil
}
