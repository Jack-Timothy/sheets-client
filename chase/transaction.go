package chase

import (
	"fmt"
	"log"

	"github.com/Jack-Timothy/sheets-client/cleanprint"
	"github.com/Jack-Timothy/sheets-client/keywords"
	"github.com/Jack-Timothy/sheets-client/standard"
)

type Transaction struct {
	TransactionDate string
	PostedDate      string
	Description     string
	Category        string
	ItemType        string
	Amount          float64
	Memo            string
}

func (t Transaction) Print() {
	dataStrings := []string{
		t.TransactionDate,
		t.PostedDate,
		t.Description,
		t.Category,
		t.ItemType,
		fmt.Sprintf("%f", t.Amount),
		t.Memo,
	}
	cleanprint.Print(expectedColumnNames, dataStrings)
}

func (t Transaction) standardize(kwMap keywords.Map) (st standard.Transaction, skip bool, err error) {
	st = standard.Transaction{
		Date:        t.TransactionDate,
		Description: t.Description,
		Amount:      t.Amount,
	}

	var foundMatch bool
	st.Category, foundMatch = kwMap.Search(t.Description)
	if foundMatch {
		skip = st.Category == "skip"
	} else {
		t.Print()
		skip, err = st.GetDescriptionAndCategoryFromUser()
		if err != nil {
			return st, false, fmt.Errorf("failed to get description or category from user: %w", err)
		}
	}
	if skip {
		fmt.Printf("Skipping transaction %+v.\n\n", t)
		return st, true, nil
	}

	log.Printf("Adding transaction %+v to statement.\n\n", st)
	return st, false, nil
}
