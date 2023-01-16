package standard

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"

	"github.com/Jack-Timothy/sheets-client/cleanprint"
)

type Statement []Transaction

func (s Statement) Print() {
	statementStrings := [][]string{
		{"Index", "Date", "Category", "Description", "Amount"},
	}
	for i, t := range s {
		statementStrings = append(statementStrings, t.makePrintableLine(i))
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
	s.Print()
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
		s.Print()
	}
}

func (s *Statement) editBasedOnUserInput(selectedAction string) error {
	if selectedAction == "add" {
		err := s.handleUserAddingTransaction()
		if err != nil {
			return fmt.Errorf("failed to handle user adding transaction: %w", err)
		}
		return nil
	}

	if strings.HasPrefix(selectedAction, "delete") {
		err := s.handleUserDeletingTransaction(selectedAction)
		if err != nil {
			return fmt.Errorf("failed to handle user deleting transaction: %w", err)
		}
		return nil
	}

	fmt.Printf("'%s' is not a valid action.", selectedAction)
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

func (s *Statement) handleUserDeletingTransaction(userInput string) error {
	userInput = strings.TrimPrefix(userInput, "delete")
	userInput = strings.TrimSpace(userInput)
	indexOfTransactionToDelete, err := strconv.ParseInt(userInput, 10, bitsPerWord)
	if err != nil {
		return fmt.Errorf("failed to parse integer from user input: %w", err)
	}
	if indexOfTransactionToDelete < 0 {
		return fmt.Errorf("received index %d but need non-negative index", indexOfTransactionToDelete)
	}
	if int(indexOfTransactionToDelete) >= len(*s) {
		return fmt.Errorf("%d exceeds the bounds of statement which has %d transactions", indexOfTransactionToDelete, len(*s))
	}
	*s = append((*s)[:indexOfTransactionToDelete], (*s)[indexOfTransactionToDelete+1:]...)
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
