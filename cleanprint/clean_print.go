package cleanprint

import (
	"fmt"
	"strings"
)

type cleanLines struct {
	line1 []string
	line2 []string
}

type altCleanLines [][]string

func PrintLines(lines [][]string) {
	acl := makeAltCleanLines(lines)
	outputLines := make([]string, 0)
	for _, cleanLine := range acl {
		var outputLine string
		for _, str := range cleanLine {
			outputLine += str + " "
		}
		if len(outputLine) > 0 {
			outputLine = outputLine[:len(outputLine)-1]
		}
		outputLines = append(outputLines, outputLine)
	}
	for _, outputLine := range outputLines {
		fmt.Println(outputLine)
	}
}

func Print(rawFirstLine, rawSecondLine []string) {
	cl := makeCleanLines(rawFirstLine, rawSecondLine)
	var firstLine string
	for _, s := range cl.line1 {
		firstLine += s + " "
	}
	if len(firstLine) > 0 {
		firstLine = firstLine[:len(firstLine)-1]
	}
	var dataLine string
	for _, s := range cl.line2 {
		dataLine += s + " "
	}
	if len(dataLine) > 0 {
		dataLine = dataLine[:len(dataLine)-1]
	}
	fmt.Println(firstLine)
	fmt.Println(dataLine)
}

func makeCleanLines(l1, l2 []string) cleanLines {
	cl := cleanLines{
		line1: l1,
		line2: l2,
	}
	cl.equalizeLineLengths()
	cl.equalizeStringLengths()
	return cl
}

func makeAltCleanLines(inputStrings [][]string) altCleanLines {
	acl := altCleanLines(inputStrings)
	acl.equalizeLineLengths()
	acl.equalizeStringLengths()
	return acl
}

func (cl *cleanLines) equalizeLineLengths() {
	n := len(cl.line1) - len(cl.line2)
	if n > 0 {
		cl.line2 = extendSlice(cl.line2, n)
	} else if n < 0 {
		cl.line1 = extendSlice(cl.line1, -1*n)
	}
}

func (acl *altCleanLines) equalizeLineLengths() {
	var maxLineLength int
	for _, line := range *acl {
		lineLength := len(line)
		if lineLength > maxLineLength {
			maxLineLength = lineLength
		}
	}
	for i := 0; i < len(*acl); i++ {
		(*acl)[i] = altExtendSlice((*acl)[i], maxLineLength)
	}
}

func extendSlice(s []string, n int) []string {
	for i := 0; i < n; i++ {
		s = append(s, "")
	}
	return s
}

func altExtendSlice(s []string, size int) []string {
	if len(s) >= size {
		return s
	}
	numToAdd := size - len(s)
	for i := 0; i < numToAdd; i++ {
		s = append(s, "")
	}
	return s
}

func (cl *cleanLines) equalizeStringLengths() {
	for i := 0; i < len(cl.line1); i++ {
		cl.line1[i], cl.line2[i] = equalizeStringLengths(cl.line1[i], cl.line2[i])
	}
}

func equalizeStringLengths(str1, str2 string) (string, string) {
	n := len(str1) - len(str2)
	if n > 0 {
		return str1, str2 + strings.Repeat(" ", n)
	} else if n < 0 {
		return str1 + strings.Repeat(" ", -1*n), str2
	}
	return str1, str2
}

func (acl *altCleanLines) equalizeStringLengths() {
	for i := 0; i < len(*acl); i++ {
		acl.equalizeStringLengthsForColumn(i)
	}
}

func (acl *altCleanLines) equalizeStringLengthsForColumn(columnIndex int) {
	var maxStringLength int
	for _, line := range *acl {
		stringLength := len(line[columnIndex])
		if stringLength > maxStringLength {
			maxStringLength = stringLength
		}
	}
	for i := 0; i < len(*acl); i++ {
		acl.extendStringLength(i, columnIndex, maxStringLength)
	}
}

func (acl *altCleanLines) extendStringLength(rowIndex, columnIndex, size int) {
	targetString := (*acl)[rowIndex][columnIndex]
	stringLength := len(targetString)
	if stringLength >= size {
		return
	}
	(*acl)[rowIndex][columnIndex] = targetString + strings.Repeat(" ", size - len(targetString))
}
