package entity

import (
	"fmt"
	"ini-go/src/atom/interfaces"
	"sort"
	"strings"
)

type INIObject struct {
	OrderedSectionsName []string
	SectionsMap         map[string]*INISectionObject
}

func NewINIObject() *INIObject {
	return &INIObject{
		OrderedSectionsName: make([]string, 0, 16),
		SectionsMap:         make(map[string]*INISectionObject),
	}
}

func (iniObject *INIObject) AddSection(section *INISectionObject) {
	if section == nil {
		return
	}
	sectionName := section.GetName()
	iniObject.OrderedSectionsName = append(iniObject.OrderedSectionsName, sectionName)
	iniObject.SectionsMap[sectionName] = section
}

func (iniObject *INIObject) GetSection(sectionName string) *INISectionObject {
	return iniObject.SectionsMap[sectionName]
}

func (iniObject *INIObject) GenerateStringLines() (strLines []string) {

	var iniContentLines []interfaces.INIContent

	for _, v := range iniObject.OrderedSectionsName {
		if len(v) > 0 {
			sectionObj, ok := iniObject.SectionsMap[v]
			if ok {
				oneSectionLines := sectionObj.GenerateContentLines()
				if len(oneSectionLines) > 0 {
					iniContentLines = append(iniContentLines, oneSectionLines...)
				}
			}
		}
	}

	sort.Sort(interfaces.INIContentSlice(iniContentLines))

	var sbOneLine strings.Builder

	preLineNumber := -1
	curLineNumber := -1

	for _, iniContent := range iniContentLines {
		if iniContent == nil {
			continue
		}
		curINIPosition := iniContent.GetPosition()
		if curINIPosition == nil {
			if sbOneLine.Len() > 0 {
				strLines = append(strLines, sbOneLine.String())
				sbOneLine.Reset()
			}
			str, err := iniContent.ToINIOutputString()
			if err != nil {
				fmt.Println(err)
			} else {
				strLines = append(strLines, str)
			}
			continue
		}

		curLineNumber = curINIPosition.LineNumber
		if preLineNumber != curLineNumber {
			if preLineNumber > -1 {
				strLines = append(strLines, sbOneLine.String())
				sbOneLine.Reset()
			}
			lineDelta := curLineNumber - preLineNumber
			if lineDelta > 1 {
				// 中间有空行
				for i := 0; i < lineDelta-1; i++ {
					strLines = append(strLines, "")
				}
			}
			str, err := iniContent.ToINIOutputString()
			if err != nil {
				fmt.Println(err)
			} else {
				sbOneLine.WriteString(str)
			}
		} else {
			str, err := iniContent.ToINIOutputString()
			if err != nil {
				fmt.Println(err)
			} else {
				sbOneLine.WriteString(str)
			}
		}
		preLineNumber = curLineNumber
	}
	strLines = append(strLines, sbOneLine.String())
	return
}
