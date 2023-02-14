package standard

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"

	"github.com/Jack-Timothy/sheets-client/cleanprint"
)

type Statement []Transaction

var columnTitles []string = []string{
	"Index", "Date", "Category", "Description", "Amount",
}

func (s Statement) Print(withIndex bool) {
	headings := columnTitles
	if !withIndex {
		headings = headings[1:]
	}
	statementStrings := [][]string{
		headings,
	}
	for i, t := range s {
		statementStrings = append(statementStrings, t.makePrintableLine(i, withIndex))
	}
	cleanprint.Print(statementStrings)
}

func getUserInput() (userInput string, err error) {
	reader := bufio.NewReader(os.Stdin)
	userInput, err = reader.ReadString('\n')
	if err != nil {
		return "", fmt.Errorf("failed to read user input: %w", err)
	}
	userInput = strings.TrimSuffix(userInput, "\n")
	userInput = strings.TrimSuffix(userInput, "\r")
	return userInput, nil
}

func (s *Statement) AcceptUserEdits() error {
	fmt.Println("Statement:")
	s.Print(true)
	for {
		fmt.Println("Please select one of the following actions:")
		fmt.Println("- Enter 'ok' to accept statement.")
		fmt.Println("- Enter 'add' to add a new transaction.")
		fmt.Println("- Enter 'delete <TRANSACTION_INDEX>' to delete a transaction.")
		fmt.Println("- Enter 'edit <TRANSACTION_INDEX>' to edit a transaction.")

		selectedAction, err := getUserInput()
		if err != nil {
			log.Printf("Error getting action selection: %v", err)
			continue
		}

		if selectedAction == "ok" {
			return nil
		}

		err = s.editBasedOnUserInput(selectedAction)
		if err != nil {
			log.Printf("Error editing based on user input: %v", err)
			continue
		}

		fmt.Println("Updated statement:")
		s.Print(true)
	}
}

func (s *Statement) editBasedOnUserInput(input string) error {
	frags := strings.Split(input, " ")
	if len(frags) == 0 {
		return errors.New("user input is empty")
	}
	selectedAction := frags[0]

	switch selectedAction {
	case "add":
		if err := s.handleUserAddingTransaction(); err != nil {
			return fmt.Errorf("failed to handle user adding transaction: %w", err)
		}
	case "delete":
		if err := s.handleUserDeletingTransaction(input); err != nil {
			return fmt.Errorf("failed to handle user deleting transaction: %w", err)
		}
	case "edit":
		if err := s.handleUserEditingTransaction(input); err != nil {
			return fmt.Errorf("failed to handle user editing transaction: %w", err)
		}
	default:
		return fmt.Errorf("'%s' is not a valid action", selectedAction)
	}
	return nil
}

func (s *Statement) handleUserAddingTransaction() error {
	t, err := getSingleTransactionFromUser()
	if err != nil {
		return fmt.Errorf("failed to get single transaction from user: %w", err)
	}

	err = s.addTransaction(t)
	if err != nil {
		return fmt.Errorf("failed to add transaction to statement: %w", err)
	}
	return nil
}

func (s *Statement) addTransaction(t Transaction) error {
	*s = append(*s, t)
	err := s.sort()
	if err != nil {
		return fmt.Errorf("failed to sort statement: %w", err)
	}
	return nil
}

func (s *Statement) handleUserDeletingTransaction(input string) error {
	input = strings.TrimPrefix(input, "delete")
	input = strings.TrimSpace(input)
	indexToDelete, err := strconv.ParseUint(input, 10, bitsPerWord)
	if err != nil {
		return fmt.Errorf("failed to parse unsigned integer from user input: %w", err)
	}
	if err := s.deleteTransactionIndex(int(indexToDelete)); err != nil {
		return fmt.Errorf("failed to delete transaction with index %d: %w", indexToDelete, err)
	}
	return nil
}

func (s *Statement) deleteTransactionIndex(index int) error {
	if index < 0 {
		return fmt.Errorf("given index %d is negative", index)
	}
	if index >= len(*s) {
		return fmt.Errorf("%d exceeds the bounds of statement which has %d transactions", index, len(*s))
	}
	*s = append((*s)[:index], (*s)[index+1:]...)
	return nil
}

func (s *Statement) getTransactionWithIndex(index int) (*Transaction, error) {
	if index < 0 {
		return nil, fmt.Errorf("given index %d is negative", index)
	}
	if index >= len(*s) {
		return nil, fmt.Errorf("%d exceeds the bounds of statement which has %d transactions", index, len(*s))
	}
	return &(*s)[index], nil
}

func (s *Statement) handleUserEditingTransaction(input string) error {
	input = strings.TrimPrefix(input, "edit")
	input = strings.TrimSpace(input)
	indexToEdit, err := strconv.ParseUint(input, 10, bitsPerWord)
	if err != nil {
		return fmt.Errorf("failed to parse unsigned integer from user input: %w", err)
	}

	tr, err := s.getTransactionWithIndex(int(indexToEdit))
	if err != nil {
		return fmt.Errorf("failed to get transaction with index %d: %w", indexToEdit, err)
	}

	fmt.Println("Editing the following transaction:")
	tr.printWithHeadings()

	// edit Date
	fmt.Println("Enter a new Date or press Enter to leave it the same.")
	dateInput, err := getUserInput()
	if err != nil {
		return fmt.Errorf("failed to get user input for Date: %w", err)
	}
	if dateInput != "" {
		if err = validateDateString(dateInput); err != nil {
			return fmt.Errorf("invalid date format: %w", err)
		}
		tr.Date = dateInput
	}

	// edit Category
	if err = tr.getCategoryFromUser(); err != nil {
		return fmt.Errorf("failed to get category from user: %w", err)
	}

	// edit Description
	fmt.Println("Enter a new Description or press Enter to leave it the same.")
	descriptionInput, err := getUserInput()
	if err != nil {
		return fmt.Errorf("failed to get user input for Description: %w", err)
	}
	if descriptionInput != "" {
		tr.Description = descriptionInput
	}

	// edit Amount
	if err = tr.getAmountFromUser(); err != nil {
		return fmt.Errorf("failed to get amount from user: %w", err)
	}

	fmt.Println("Resulting transaction data after edits:")
	tr.printWithHeadings()
	return nil
}

func (s *Statement) sort() error {
	for i, t := range *s {
		if err := validateDateString(t.Date); err != nil {
			return fmt.Errorf("transaction with index %d has invalid date string", i)
		}
	}
	sort.Slice(*s, func(x, y int) bool {
		return isDateXBeforeDateY((*s)[x].Date, (*s)[y].Date)
	})
	return nil
}

func BuildTestStatement(numTransactions int) Statement {
	s := make(Statement, 0, numTransactions)
	for i := 0; i < numTransactions; i++ {
		s = append(s, buildTestTransaction(i))
	}
	return s
}
