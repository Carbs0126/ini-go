package atom

import "ini-go/src/position"

type INIEmpty struct {
	INIPosition *position.INIPosition
}

func NewINIEmpty(iniPosition *position.INIPosition) *INIEmpty {
	iniEmpty := INIEmpty{
		INIPosition: iniPosition,
	}
	return &iniEmpty
}

func (iniEmpty *INIEmpty) GetPosition() *position.INIPosition {
	return iniEmpty.INIPosition
}

func (iniEmpty *INIEmpty) ToINIOutputString() (string, error) {
	return "", nil
}
