package atom

import (
	"errors"
	"ini-go/src/position"
)

type INISectionHeader struct {
	INIPosition *position.INIPosition
	Name        string
}

func NewINISectionHeader(name string, iniPosition *position.INIPosition) *INISectionHeader {
	iniKVPair := INISectionHeader{
		INIPosition: iniPosition,
		Name:        name,
	}
	return &iniKVPair
}

func (iniSectionHeader *INISectionHeader) GetPosition() *position.INIPosition {
	return iniSectionHeader.INIPosition
}

func (iniSectionHeader *INISectionHeader) ToINIOutputString() (string, error) {
	if len(iniSectionHeader.Name) == 0 {
		return "", errors.New("key of INISectionHeader should not be empty")
	}
	return iniSectionHeader.Name, nil
}
