package src

import (
	"bufio"
	"errors"
	"ini-go/src/atom"
	"ini-go/src/atom/interfaces"
	"ini-go/src/entity"
	"ini-go/src/position"
	"io"
	"os"
	"strconv"
	"strings"
)

func ParseFileToINIObject(iniFileAbsolutePath string) (*entity.INIObject, error) {
	file, err := os.Open(iniFileAbsolutePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := bufio.NewReader(file)

	content := make([]string, 0, 32)

	for {
		line, err := reader.ReadString('\n')
		if err == io.EOF {
			content = append(content, line)
			break
		}
		content = append(content, line)
	}

	var iniLines []interfaces.INIContent
	lineNumber := 0
	fileName := iniFileAbsolutePath

	for _, originLine := range content {
		trimmedLine := strings.TrimSpace(originLine)
		if strings.HasPrefix(trimmedLine, ";") {
			// comment
			originLineTrimmedRight := strings.TrimRight(originLine, "\n")
			iniComment := atom.NewINIComment(originLineTrimmedRight, position.NewPosition(fileName, lineNumber, 0, len(originLineTrimmedRight)))
			appendLineContentIntoLineList(iniComment, &iniLines)
		} else if strings.HasPrefix(trimmedLine, "[") {
			// section header
			rightSquareBracketsPosition := strings.Index(trimmedLine, "]")
			if rightSquareBracketsPosition < 2 {
				return nil, errors.New("Right square bracket's position should be greater than 1, now it is " +
					strconv.Itoa(rightSquareBracketsPosition))
			}
			sectionName := trimmedLine[0 : rightSquareBracketsPosition+1]
			if strings.Contains(sectionName, ";") {
				return nil, errors.New("Section's name should not contain ';' symbol")
			}
			charBegin := strings.Index(originLine, "[")
			charEnd := strings.Index(originLine, "]")
			sectionHeader := atom.NewINISectionHeader(sectionName, position.NewPosition(fileName, lineNumber, charBegin, charEnd))
			appendLineContentIntoLineList(sectionHeader, &iniLines)
			checkSemicolon(originLine, charEnd+1, &iniLines, fileName, lineNumber)
		} else if len(trimmedLine) == 0 {
			iniEmpty := atom.NewINIEmpty(position.NewPosition(fileName, lineNumber, 0, 0))
			appendLineContentIntoLineList(iniEmpty, &iniLines)
		} else {
			// kv
			indexOfEqualInTrimmedString := strings.Index(trimmedLine, "=")
			if indexOfEqualInTrimmedString < 1 {
				return nil, errors.New("Equal's position should be greater than 0")
			}
			indexOfEqualInOriginString := strings.Index(originLine, "=")
			keyName := strings.TrimSpace(trimmedLine[0:indexOfEqualInTrimmedString])
			rightStringOfEqual := trimmedLine[indexOfEqualInTrimmedString+1:]
			var valueNameSB strings.Builder
			length := len(rightStringOfEqual)
			if length > 0 {
				// 0: 过滤前面的空格，还未找到value
				// 1: 正在记录value
				// 2: value结束
				stat := 0
				i := 0
				for i = 0; i < length; i++ {
					c := string([]rune(rightStringOfEqual)[i])
					if stat == 0 {
						// 过滤前面的空格
						if c == " " || c == "\t" {
							continue
						} else {
							stat = 1
							valueNameSB.WriteString(c)
						}
					} else if stat == 1 {
						// 正在记录value
						// value中允许有空格
						if c == ";" {
							// 记录 value 结束
							stat = 2
							break
						} else {
							stat = 1
							valueNameSB.WriteString(c)
						}
					}
				}
				valueName := valueNameSB.String()
				charBegin := strings.Index(originLine, keyName)
				charEnd := indexOfEqualInOriginString + 1 + i
				inikvPair := atom.NewINIKVPair(keyName, valueName, position.NewPosition(fileName, lineNumber, charBegin, charEnd))
				appendLineContentIntoLineList(inikvPair, &iniLines)
				if i != length {
					// 没有到结尾，检测是不是有分号
					checkSemicolon(originLine, indexOfEqualInOriginString+1+i, &iniLines, fileName, lineNumber)
				}
			}
		}
		lineNumber++
	}
	// 最终解析为一个实体
	iniObject := entity.NewINIObject()
	// 收集 section 或者 kv 的 comments
	var commentsCache []*atom.INIComment
	// 解析完当前的 section ，一次存入
	var currentSectionObject *entity.INISectionObject = nil
	// 解析当前的 kvPair
	var currentEntryObject *entity.INIEntryObject = nil

	// 0 解析 section 阶段，还是没有解析到 section
	// 1 已经解析出 sectionName 阶段，（刚刚解析完 sectionHeader）还没有解析到下一个 section
	parseState := 0
	var preINIContent interfaces.INIContent = nil
	var curINIContent interfaces.INIContent = nil

	for _, iniContent := range iniLines {
		_, ok := iniContent.(*atom.INIEmpty)
		if ok {
			continue
		}
		curINIContent = iniContent
		switch cur := curINIContent.(type) {
		case *atom.INIComment:
			var iniComment = cur
			if parseState == 0 {
				// 还没解析到 section
				commentsCache = append(commentsCache, iniComment)
			} else {
				switch preINIContent.(type) {
				case *atom.INISectionHeader:
					if checkSameLine(preINIContent, curINIContent) {
						// 当前 comment 属于 section
						commentsCache = append(commentsCache, iniComment)
						if currentSectionObject == nil {
							currentSectionObject = entity.NewINISectionObject()
						}
						currentSectionObject.AddComments(commentsCache)
						commentsCache = commentsCache[0:0]
						// 当前 section 的所有 comment 已经结束
					} else {
						// 当前 comment 属于当前 section 的 kv 或者下一个 section 的 section
						if currentSectionObject == nil {
							currentSectionObject = entity.NewINISectionObject()
						}
						currentSectionObject.AddComments(commentsCache)
						commentsCache = commentsCache[0:0]
						commentsCache = append(commentsCache, iniComment)
					}
				case *atom.INIComment:
					// comment 累加
					commentsCache = append(commentsCache, iniComment)
				case *atom.INIKVPair:
					if checkSameLine(preINIContent, curINIContent) {
						// 当前 comment 属于 kv
						commentsCache = append(commentsCache, iniComment)
						if currentEntryObject == nil {
							// 不走这里
							currentEntryObject = entity.NewINIEntryObject()
						}
						currentEntryObject.AddComments(commentsCache)
						if currentSectionObject == nil {
							currentSectionObject = entity.NewINISectionObject()
						}
						currentSectionObject.AddEntryObject(currentEntryObject)
						currentEntryObject = nil
						commentsCache = commentsCache[0:0]
						// 当前 kv 收尾
					} else {
						// 当前comments 属于当前 section 的下一个 kv 或者下一个 section 的 section
						commentsCache = commentsCache[0:0]
						commentsCache = append(commentsCache, iniComment)
					}
				}
			}
		case *atom.INISectionHeader:
			var iniSectionHeader = cur
			if parseState == 0 {
				// 解析到第一个 section
				parseState = 1
				currentSectionObject = entity.NewINISectionObject()
				currentSectionObject.INISectionHeader = iniSectionHeader
			} else {
				switch preINIContent.(type) {
				case *atom.INISectionHeader:
					// 连着两个 section header
					// 收尾上一个 section header
					if currentSectionObject != nil {
						currentSectionObject.AddComments(commentsCache)
						commentsCache = commentsCache[0:0]
						iniObject.AddSection(currentSectionObject)
					}
					// 新建 section header
					currentSectionObject = entity.NewINISectionObject()
					currentSectionObject.INISectionHeader = iniSectionHeader
				case *atom.INIComment:
					if len(commentsCache) == 0 {
						// 说明上一个 comment 和其之前的元素是一行，需要收尾上一个 section
						if currentSectionObject != nil {
							iniObject.AddSection(currentSectionObject)
						}
						currentSectionObject = entity.NewINISectionObject()
						currentSectionObject.INISectionHeader = iniSectionHeader
					} else {
						currentSectionObject = entity.NewINISectionObject()
						currentSectionObject.INISectionHeader = iniSectionHeader
						currentSectionObject.AddComments(commentsCache)
						commentsCache = commentsCache[0:0]
					}
				case *atom.INIKVPair:
					// 说明上一个section结束了，需要收尾
					if currentSectionObject != nil {
						if currentEntryObject != nil {
							currentSectionObject.AddEntryObject(currentEntryObject)
							currentEntryObject = nil
						}
						iniObject.AddSection(currentSectionObject)
					}
					currentSectionObject = entity.NewINISectionObject()
					currentSectionObject.INISectionHeader = iniSectionHeader
				}
			}
		case *atom.INIKVPair:
			inikvPair := cur
			if parseState == 0 {
				// 没有 section，就出现了 kv，说明格式出错
				return nil, errors.New("There should be a section header before key-value pairs")
			} else {

				switch preINIContent.(type) {
				case *atom.INISectionHeader:
					currentEntryObject = entity.NewINIEntryObject()
					currentEntryObject.INIKVPair = inikvPair
				case *atom.INIComment:
					if len(commentsCache) == 0 {
						// 说明上一行中，comment 是右边的注释，还包含左边的元素
						// 当上一行的左侧是 section 时，不需要关心 section
						// 当上一行的左侧是 kv 时，不需要关心当前 section 或者上一个 kv
						currentEntryObject = entity.NewINIEntryObject()
						currentEntryObject.INIKVPair = inikvPair
					} else {
						currentEntryObject = entity.NewINIEntryObject()
						currentEntryObject.INIKVPair = inikvPair
					}
				case *atom.INIKVPair:
					// 把前一个 kv 收尾到 section 中
					if currentEntryObject != nil {
						currentEntryObject.AddComments(commentsCache)
						commentsCache = commentsCache[0:0]
						if currentSectionObject != nil {
							currentSectionObject.AddEntryObject(currentEntryObject)
						}
					}
					currentEntryObject = entity.NewINIEntryObject()
					currentEntryObject.INIKVPair = inikvPair
				}
			}
		}
		preINIContent = curINIContent
	}
	// 最后一个元素
	if currentEntryObject != nil {
		currentEntryObject.AddComments(commentsCache)
		commentsCache = commentsCache[0:0]
	}
	if currentSectionObject != nil {
		currentSectionObject.AddComments(commentsCache)
		commentsCache = commentsCache[0:0]
		if currentEntryObject != nil {
			currentSectionObject.AddEntryObject(currentEntryObject)
			currentEntryObject = nil
		}
		iniObject.AddSection(currentSectionObject)
	}
	return iniObject, nil
}

func appendLineContentIntoLineList(iniContent interfaces.INIContent, iniLines *[]interfaces.INIContent) {
	*iniLines = append(*iniLines, iniContent)
}

func checkSemicolon(originString string, charBegin int, iniLines *[]interfaces.INIContent, fileLocation string, lineNumber int) (*atom.INIComment, error) {
	remainStr := originString[charBegin:]
	trimmedRemainStr := strings.TrimSpace(remainStr)
	if len(trimmedRemainStr) > 0 {
		if strings.HasPrefix(trimmedRemainStr, ";") {
			iniComment := atom.NewINIComment(trimmedRemainStr,
				position.NewPosition(fileLocation, lineNumber, strings.Index(originString, ";"), len(originString)))
			appendLineContentIntoLineList(iniComment, iniLines)
			return iniComment, nil
		} else {
			return nil, errors.New("Need ';' symbol, but find " + string([]rune(trimmedRemainStr)[0]) + " instead")
		}
	}
	return nil, nil
}

func checkSameLine(preINIContent interfaces.INIContent, curINIContent interfaces.INIContent) bool {
	if preINIContent != nil && curINIContent != nil {
		return preINIContent.GetPosition().LineNumber == curINIContent.GetPosition().LineNumber
	}
	return false
}
