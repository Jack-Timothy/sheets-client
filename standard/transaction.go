package standard

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

type Transaction struct {
	Date        string
	Category    string
	Description string
	Amount      float64
}

func (t *Transaction) GetDescriptionAndCategoryFromUser() (skip bool, err error) {
	skip, err = t.getDescriptionFromUserWithOptions()
	if err != nil {
		return false, fmt.Errorf("failed to get description from user: %w", err)
	}
	if skip {
		return skip, nil
	}

	err = t.getCategoryFromUser()
	if err != nil {
		return false, fmt.Errorf("failed to get category from user: %w", err)
	}
	return false, nil
}

func (t *Transaction) getDescriptionFromUserWithOptions() (skip bool, err error) {
	fmt.Printf("Please provide a description of this transaction. Press Enter to accept the default description. Submit 'skip' to not include the transaction in the final statement.\n")

	description, err := getUserInput()
	if err != nil {
		return false, err
	}

	if description == "skip" {
		fmt.Printf("This transaction will be skipped.\n\n")
		return true, nil
	}
	if len(description) != 0 {
		t.Description = description
	}
	return false, nil
}

func (t *Transaction) getDescriptionFromUser() error {
	fmt.Printf("Please enter a description of the transaction.\n")

	description, err := getUserInput()
	if err != nil {
		return fmt.Errorf("failed to get user input: %w", err)
	}
	if len(description) == 0 {
		return errors.New("received empty description")
	}
	t.Description = description
	return nil
}

func (t *Transaction) getAmountFromUser() error {
	fmt.Printf("Please enter the amount of the transaction.\n")
	_, err := fmt.Scanf("%f\n", &t.Amount)
	return err
}

func (t *Transaction) getCategoryFromUser() error {
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
	_, err := fmt.Scanf("%d\n", &categoryEnum)
	if err != nil {
		return fmt.Errorf("failed to scan integer from user: %w", err)
	}
	if categoryEnum > len(categories) || categoryEnum < 1 {
		return fmt.Errorf("received invalid category enumeration %d", categoryEnum)
	}

	t.Category = categories[categoryEnum-1]
	return nil
}

func (t *Transaction) getDateFromUser() error {
	fmt.Printf("Please enter the date of the transaction with the format MM/DD/YYYY.\n")
	date, err := getUserInput()
	if err != nil {
		return fmt.Errorf("failed to get user input: %w", err)
	}
	if err = validateDateString(date); err != nil {
		return fmt.Errorf("invalid date format: %w", err)
	}
	t.Date = date
	return nil
}

const bitsPerWord = 32 << (^uint(0) >> 63)

func validateDateString(date string) error {
	dateElements := strings.Split(date, "/")
	if len(dateElements) != 3 {
		return fmt.Errorf("expected 2 slashes, but got %d", len(dateElements)-1)
	}
	for i, dateElement := range dateElements {
		num, err := strconv.ParseInt(dateElement, 10, bitsPerWord)
		if err != nil {
			return fmt.Errorf("failed to parse number with index %d as an integer: %w", i, err)
		}
		switch i {
		case 0: // month must be between 1 and 12
			if num < 1 || num > 12 {
				return fmt.Errorf("received %d for month but it must be between 1 and 12", num)
			}
		case 1: // day must be between 1 and 31
			if num < 1 || num > 31 {
				return fmt.Errorf("received %d for day but it must be between 1 and 31", num)
			}
		default: // year just must be positive
			if num < 1 {
				return fmt.Errorf("received %d for year but it must be positive", num)
			}
		}
	}
	return nil
}

func isDateXBeforeDateY(x, y string) bool {
	xElements := strings.Split(x, "/")
	yElements := strings.Split(y, "/")
	// year is first comparison
	xYear, _ := strconv.ParseInt(xElements[2], 10, bitsPerWord)
	yYear, _ := strconv.ParseInt(yElements[2], 10, bitsPerWord)
	if xYear < yYear {
		return true
	}
	if xYear > yYear {
		return false
	}
	// month is next comparison
	xMonth, _ := strconv.ParseInt(xElements[0], 10, bitsPerWord)
	yMonth, _ := strconv.ParseInt(yElements[0], 10, bitsPerWord)
	if xMonth < yMonth {
		return true
	}
	if xMonth > yMonth {
		return false
	}
	// day is final comparison
	xDay, _ := strconv.ParseInt(xElements[1], 10, bitsPerWord)
	yDay, _ := strconv.ParseInt(yElements[1], 10, bitsPerWord)
	return xDay < yDay
}

func (tr *Transaction) printWithHeadings() {
	statementCopy := make(Statement, 0, 1)
	statementCopy = append(statementCopy, *tr)
	statementCopy.Print()
}

func (t *Transaction) makePrintableLine(index int) []string {
	return []string{
		fmt.Sprintf("%d", index),
		t.Date,
		t.Category,
		t.Description,
		fmt.Sprintf("%f", t.Amount),
	}
}

func buildTestTransaction(i int) (t Transaction) {
	now := time.Now()
	t.Date = fmt.Sprintf("%d/%d/%d", now.Month(), now.Day(), now.Year())
	t.Category = fmt.Sprintf("Test Category %d", i)
	t.Description = fmt.Sprintf("Test Description %d", i)
	t.Amount = float64(20 + i)
	return t
}

func getSingleTransactionFromUser() (t Transaction, err error) {
	if err = t.getDateFromUser(); err != nil {
		return t, fmt.Errorf("failed to get date from user: %w", err)
	}
	if err = t.getCategoryFromUser(); err != nil {
		return t, fmt.Errorf("failed to get category from user: %w", err)
	}
	if err = t.getDescriptionFromUser(); err != nil {
		return t, fmt.Errorf("failed to get description from user: %w", err)
	}
	if err = t.getAmountFromUser(); err != nil {
		return t, fmt.Errorf("failed to get amount from user: %w", err)
	}
	return t, nil
}
