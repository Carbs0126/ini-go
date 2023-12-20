package interfaces

import (
	"ini-go/src/position"
)

type INIContent interface {
	GetPosition() *position.INIPosition
	ToINIOutputString() (string, error)
}

type INIContentSlice []INIContent

func (cs INIContentSlice) Len() int {
	return len(cs)
}

func (cs INIContentSlice) Less(i, j int) bool {
	if cs[i] == nil || cs[j] == nil {
		return true
	}
	iniPositionA := cs[i].GetPosition()
	iniPositionB := cs[j].GetPosition()
	// 将 position 为空的元素排到最后
	if iniPositionA == nil {
		return false
	}
	if iniPositionB == nil {
		return true
	}
	lineNumberDelta := iniPositionA.LineNumber - iniPositionB.LineNumber
	if lineNumberDelta < 0 {
		return true
	}
	if lineNumberDelta > 0 {
		return false
	}
	return (iniPositionA.CharBegin - iniPositionB.CharBegin) < 0
}

func (cs INIContentSlice) Swap(i, j int) {
	cs[i], cs[j] = cs[j], cs[i]
}
