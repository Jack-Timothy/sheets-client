package chase

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

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

var expectedColumnNames []string = []string{
	"Transaction Date",
	"Post Date",
	"Description",
	"Category",
	"Type",
	"Amount",
	"Memo",
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

func CsvContentsToStatement(csvContents [][]string) (statement []Transaction, err error) {
	if err := validateHeaderRow(csvContents[0]); err != nil {
		return nil, fmt.Errorf("failed to validate header row: %w", err)
	}
	csvContents = csvContents[1:]

	statement = make([]Transaction, 0, len(csvContents))
	for i, row := range csvContents {
		t, err := csvRowToTransaction(row)
		if err != nil {
			return nil, fmt.Errorf("failed to convert row %d to Chase transaction", i+2)
		}
		statement = append(statement, t)
	}
	return statement, nil
}

func validateHeaderRow(headerRow []string) error {
	if len(headerRow) != len(expectedColumnNames) {
		return fmt.Errorf("expected %d columns in header row but got %d", len(expectedColumnNames), len(headerRow))
	}
	for i, columnName := range headerRow {
		if expectedColumnNames[i] != columnName {
			return fmt.Errorf("expected column %d to be named %s but is named %s", i+1, expectedColumnNames[i], columnName)
		}
	}
	return nil
}

func csvRowToTransaction(row []string) (t Transaction, err error) {
	if len(row) != len(expectedColumnNames) {
		return t, fmt.Errorf("expected %d columns but got %d", len(expectedColumnNames), len(row))
	}
	amount, err := strconv.ParseFloat(row[5], 64)
	if err != nil {
		return t, fmt.Errorf("failed to parse 'amount' cell as float64: %w", err)
	}
	return Transaction{
		TransactionDate: row[0],
		PostedDate:      row[1],
		Description:     row[2],
		Category:        row[3],
		ItemType:        row[4],
		Amount:          amount,
		Memo:            row[6],
	}, nil
}

func StandardizeStatement(statement []Transaction, kwMap keywords.Map) (revisedStatement []standard.Transaction, err error) {
	revisedStatement = make([]standard.Transaction, 0)
	for i, t := range statement {
		st, skip, err := t.standardize(kwMap)
		if err != nil {
			return nil, fmt.Errorf("failed to standardize item %d: %w", i, err)
		}
		if skip {
			continue
		}
		revisedStatement = append(revisedStatement, st)
	}

	return revisedStatement, nil
}

func (item Transaction) standardize(kwMap keywords.Map) (t standard.Transaction, skip bool, err error) {
	t = standard.Transaction{
		Date:        item.TransactionDate,
		Description: item.Description,
		Amount:      item.Amount,
	}

	var foundMatch bool
	t.Category, foundMatch = kwMap.Search(item.Description)
	if !foundMatch {
		item.Print()

		fmt.Printf("Please provide a description of this transaction. Press Enter to accept the default description. Submit 'skip' to not include the transaction in the final statement.\n")
		reader := bufio.NewReader(os.Stdin)
		description, err := reader.ReadString('\n')
		if err != nil {
			return t, skip, fmt.Errorf("failed to scan user input for custom description: %w", err)
		}
		description = strings.TrimSuffix(description, "\n")
		description = strings.TrimSuffix(description, "\r")
		if description == "skip" {
			fmt.Printf("This transaction will be skipped.\n\n")
			return t, true, nil
		}
		if len(description) != 0 {
			t.Description = description
		}
		fmt.Printf("Received description: %s\n\n", t.Description)

		fmt.Printf("Please enter the enumeration of this transaction's category. Options are:\n")
		categories := []string{
			"Rent", "Utilities", "Groceries/Toiletries", "Food/Drinks Out", "Gas",
			"Other (Need)", "Other (Want)", "Gift Giving", "Donations",
		}
		for i, category := range categories {
			fmt.Printf("%d. %s ", i+1, category)
		}
		fmt.Printf("\n")
		var categoryEnum int
		fmt.Scanf("%d\n", &categoryEnum)
		if categoryEnum > len(categories) || categoryEnum < 1 {
			return t, skip, fmt.Errorf("received invalid category enumeration %d", categoryEnum)
		}
		t.Category = categories[categoryEnum-1]
		fmt.Printf("Received Category: %d. %s\n", categoryEnum, t.Category)
	} else if t.Category == "skip" {
		fmt.Printf("Skipping transaction %+v based on discovery of a skip keyword in the description.\n\n", item)
		return t, true, nil
	}

	log.Printf("Adding transaction %+v to statement.\n\n", t)
	return t, false, nil
}
