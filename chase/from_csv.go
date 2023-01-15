package chase

import (
	"fmt"
	"strconv"
)

var expectedColumnNames []string = []string{
	"Transaction Date",
	"Post Date",
	"Description",
	"Category",
	"Type",
	"Amount",
	"Memo",
}

func CsvContentsToStatement(csvContents [][]string) (s Statement, err error) {
	if err := validateHeaderRow(csvContents[0]); err != nil {
		return nil, fmt.Errorf("failed to validate header row: %w", err)
	}
	csvContents = csvContents[1:]

	s = make([]Transaction, 0, len(csvContents))
	for i, row := range csvContents {
		t, err := csvRowToTransaction(row)
		if err != nil {
			return nil, fmt.Errorf("failed to convert row %d to Chase transaction", i+2)
		}
		s = append(s, t)
	}
	return s, nil
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
