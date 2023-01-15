package standard

import "github.com/Jack-Timothy/sheets-client/cleanprint"

type Statement []Transaction

func (s Statement) Print() {
	statementStrings := [][]string{
		{"Date", "Category", "Description", "Amount"},
	}
	for _, t := range s {
		statementStrings = append(statementStrings, t.makePrintableLine())
	}
	cleanprint.Print(statementStrings)
}
