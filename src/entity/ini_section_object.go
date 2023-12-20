package entity

import (
	"ini-go/src/atom"
	"ini-go/src/atom/interfaces"
)

type INISectionObject struct {
	INISectionHeader *atom.INISectionHeader
	Comments         []*atom.INIComment
	EntryObjects     []*INIEntryObject
}

func NewINISectionObject() *INISectionObject {
	return &INISectionObject{
		Comments:     make([]*atom.INIComment, 0, 16),
		EntryObjects: make([]*INIEntryObject, 0, 16),
	}
}

func (iniSectionObject *INISectionObject) AddComment(comment *atom.INIComment) {
	iniSectionObject.Comments = append(iniSectionObject.Comments, comment)
}

func (iniSectionObject *INISectionObject) AddComments(comments []*atom.INIComment) {
	if len(comments) == 0 {
		return
	}
	iniSectionObject.Comments = append(iniSectionObject.Comments, comments...)
}

func (iniSectionObject *INISectionObject) AddEntryObject(entryObject *INIEntryObject) {
	if entryObject != nil {
		iniSectionObject.EntryObjects = append(iniSectionObject.EntryObjects, entryObject)
	}
}

func (iniSectionObject *INISectionObject) GetName() string {
	if iniSectionObject.INISectionHeader == nil {
		return ""
	}
	return iniSectionObject.INISectionHeader.Name
}

func (iniSectionObject *INISectionObject) GenerateContentLines() (contentLines []interfaces.INIContent) {
	for _, v := range iniSectionObject.Comments {
		contentLines = append(contentLines, v)
	}
	if iniSectionObject.INISectionHeader != nil {
		contentLines = append(contentLines, iniSectionObject.INISectionHeader)
	}
	for _, v := range iniSectionObject.EntryObjects {
		entryLines := v.GenerateContentLines()
		if len(entryLines) > 0 {
			contentLines = append(contentLines, entryLines...)
		}
	}
	return
}
