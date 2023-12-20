package atom

import (
	"errors"
	"ini-go/src/position"
)

type INIKVPair struct {
	INIPosition *position.INIPosition
	Key         string
	Value       string
}

func NewINIKVPair(key string, value string, iniPosition *position.INIPosition) *INIKVPair {
	iniKVPair := INIKVPair{
		INIPosition: iniPosition,
		Key:         key,
		Value:       value,
	}
	return &iniKVPair
}

func (iniKVPair *INIKVPair) GetPosition() *position.INIPosition {
	return iniKVPair.INIPosition
}

func (iniKVPair *INIKVPair) ToINIOutputString() (string, error) {
	if len(iniKVPair.Key) == 0 {
		return "", errors.New("key of INIEntry should not be empty")
	}
	return iniKVPair.Key + "=" + iniKVPair.Value, nil
}
