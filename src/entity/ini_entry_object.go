package entity

import (
	"ini-go/src/atom"
	"ini-go/src/atom/interfaces"
)

type INIEntryObject struct {
	Comments  []*atom.INIComment
	INIKVPair *atom.INIKVPair
}

func NewINIEntryObject() *INIEntryObject {
	return &INIEntryObject{
		Comments: make([]*atom.INIComment, 0, 16),
	}
}

func (iniEntryObject *INIEntryObject) AddComments(comments []*atom.INIComment) {
	if comments == nil || len(comments) == 0 {
		return
	}
	iniEntryObject.Comments = append(iniEntryObject.Comments, comments...)
}

func (iniEntryObject *INIEntryObject) AddComment(comment *atom.INIComment) {
	if comment == nil {
		return
	}
	iniEntryObject.Comments = append(iniEntryObject.Comments, comment)
}

func (iniEntryObject *INIEntryObject) GenerateContentLines() (contentLines []interfaces.INIContent) {
	if iniEntryObject.Comments != nil {
		for _, v := range iniEntryObject.Comments {
			contentLines = append(contentLines, v)
		}
	}
	contentLines = append(contentLines, iniEntryObject.INIKVPair)
	return
}
