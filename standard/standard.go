package standard

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type Transaction struct {
	Date        string
	Category    string
	Description string
	Amount      float64
}

type Statement []Transaction

func (t *Transaction) GetDescriptionAndCategoryFromUser() (skip bool, err error) {
	skip, err = t.getDescriptionFromUser()
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

func (t *Transaction) getDescriptionFromUser() (skip bool, err error) {
	fmt.Printf("Please provide a description of this transaction. Press Enter to accept the default description. Submit 'skip' to not include the transaction in the final statement.\n")

	reader := bufio.NewReader(os.Stdin)
	description, err := reader.ReadString('\n')
	if err != nil {
		return false, fmt.Errorf("failed to scan user input for custom description: %w", err)
	}
	description = strings.TrimSuffix(description, "\n")
	description = strings.TrimSuffix(description, "\r")

	if description == "skip" {
		fmt.Printf("This transaction will be skipped.\n\n")
		return true, nil
	}
	if len(description) != 0 {
		t.Description = description
	}
	fmt.Printf("Received description: %s\n\n", t.Description)
	return false, nil
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
	fmt.Scanf("%d\n", &categoryEnum)
	if categoryEnum > len(categories) || categoryEnum < 1 {
		return fmt.Errorf("received invalid category enumeration %d", categoryEnum)
	}

	t.Category = categories[categoryEnum-1]
	fmt.Printf("Received Category: %d. %s\n", categoryEnum, t.Category)
	return nil
}
