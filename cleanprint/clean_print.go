package cleanprint

import (
	"fmt"
	"strings"
)

func Print(inputStrings [][]string) {
	cl := equalizeLinesAndStrings(inputStrings)
	outputLines := make([]string, 0)
	for _, cleanLine := range cl {
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
	fmt.Printf("\n")
}

func equalizeLinesAndStrings(linesOfStrings [][]string) [][]string {
	linesOfStrings = equalizeLineLengths(linesOfStrings)
	linesOfStrings = equalizeStringLengths(linesOfStrings)
	return linesOfStrings
}

func equalizeLineLengths(linesOfStrings [][]string) [][]string {
	var maxLineLength int
	for _, line := range linesOfStrings {
		lineLength := len(line)
		if lineLength > maxLineLength {
			maxLineLength = lineLength
		}
	}
	for i := 0; i < len(linesOfStrings); i++ {
		linesOfStrings[i] = extendSlice(linesOfStrings[i], maxLineLength)
	}
	return linesOfStrings
}

func extendSlice(s []string, size int) []string {
	if len(s) >= size {
		return s
	}
	numToAdd := size - len(s)
	for i := 0; i < numToAdd; i++ {
		s = append(s, "")
	}
	return s
}

func equalizeStringLengths(linesOfStrings [][]string) [][]string {
	if len(linesOfStrings) == 0 {
		return linesOfStrings
	}
	for i := 0; i < len(linesOfStrings[0]); i++ {
		linesOfStrings = equalizeStringLengthsForColumn(i, linesOfStrings)
	}
	return linesOfStrings
}

func equalizeStringLengthsForColumn(columnIndex int, linesOfStrings [][]string) [][]string {
	var maxStringLength int
	for _, line := range linesOfStrings {
		stringLength := len(line[columnIndex])
		if stringLength > maxStringLength {
			maxStringLength = stringLength
		}
	}
	for i, line := range linesOfStrings {
		linesOfStrings[i][columnIndex] += strings.Repeat(" ", maxStringLength-len(line[columnIndex]))
	}
	return linesOfStrings
}
