package main

import (
	"bufio"
	"context"
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"golang.org/x/oauth2"
)

type chaseCreditItem struct {
	transactionDate string
	postedDate      string
	description     string
	category        string
	itemType        string
	amount          float64
	memo            string
}

type transaction struct {
	date        string
	category    string
	description string
	amount      float64
}

type keywords struct {
	Rent                []string `json:"rent"`
	Utilities           []string `json:"utilities"`
	GroceriesToiletries []string `json:"groceries_toiletries"`
	FoodDrinksOut       []string `json:"food_drinks_out"`
	Gas                 []string `json:"gas"`
	OtherNeed           []string `json:"other_need"`
	OtherWant           []string `json:"other_want"`
	GiftGiving          []string `json:"gift_giving"`
	Donations           []string `json:"donations"`
	Skip                []string `json:"skip"`
}

type keywordMap map[string]string

// Retrieve a token, saves the token, then returns the generated client.
func getClient(config *oauth2.Config) *http.Client {
	// The file token.json stores the user's access and refresh tokens, and is
	// created automatically when the authorization flow completes for the first
	// time.
	tokFile := "token.json"
	tok, err := tokenFromFile(tokFile)
	if err != nil {
		tok = getTokenFromWeb(config)
		saveToken(tokFile, tok)
	}
	return config.Client(context.Background(), tok)
}

// Request a token from the web, then returns the retrieved token.
func getTokenFromWeb(config *oauth2.Config) *oauth2.Token {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("Go to the following link in your browser then type the "+
		"authorization code: \n%v\n", authURL)

	var authCode string
	if _, err := fmt.Scan(&authCode); err != nil {
		log.Fatalf("Unable to read authorization code: %v", err)
	}

	tok, err := config.Exchange(context.TODO(), authCode)
	if err != nil {
		log.Fatalf("Unable to retrieve token from web: %v", err)
	}
	return tok
}

// Retrieves a token from a local file.
func tokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	tok := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(tok)
	return tok, err
}

// Saves a token to a file path.
func saveToken(path string, token *oauth2.Token) {
	fmt.Printf("Saving credential file to: %s\n", path)
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		log.Fatalf("Unable to cache oauth token: %v", err)
	}
	defer f.Close()
	json.NewEncoder(f).Encode(token)
}

func main() {
	// ctx := context.Background()
	// b, err := os.ReadFile("credentials.json")
	// if err != nil {
	// 	log.Fatalf("Unable to read client secret file: %v", err)
	// }

	// // If modifying these scopes, delete your previously saved token.json.
	// // For full list of scopes: https://developers.google.com/identity/protocols/oauth2/scopes#sheets.
	// config, err := google.ConfigFromJSON(b, "https://www.googleapis.com/auth/spreadsheets")
	// if err != nil {
	// 	log.Fatalf("Unable to parse client secret file to config: %v", err)
	// }
	// client := getClient(config)

	// srv, err := sheets.NewService(ctx, option.WithHTTPClient(client))
	// if err != nil {
	// 	log.Fatalf("Unable to retrieve Sheets client: %v", err)
	// }

	// // Prints the data in a test spreadsheet:
	// // https://docs.google.com/spreadsheets/d/15KWFkIY-RW81leDLXqahARB0gtSnWAIGDg-lkx2g04Q/edit
	// spreadsheetId := "15KWFkIY-RW81leDLXqahARB0gtSnWAIGDg-lkx2g04Q"
	// readRange := "Sheet1!A2:B3"
	// resp, err := srv.Spreadsheets.Values.Get(spreadsheetId, readRange).Do()
	// if err != nil {
	// 	log.Fatalf("Unable to retrieve data from sheet: %v", err)
	// }

	// if len(resp.Values) == 0 {
	// 	fmt.Println("No data found.")
	// } else {
	// 	fmt.Println("x, y:")
	// 	for _, row := range resp.Values {
	// 		// Print columns A and B, which correspond to indices 0 and 1.
	// 		fmt.Printf("%s, %s\n", row[0], row[1])
	// 	}
	// }

	// fmt.Println("Checking for expense sheets...")
	// respSpreadsheet, err := srv.Spreadsheets.Get(spreadsheetId).Do()
	// if err != nil {
	// 	log.Fatalf("Failed to get spreadsheet: %v", err)
	// }

	// var missingMonths []string = []string{
	// 	"January", "February", "March", "April", "May", "June",
	// 	"July", "August", "September", "October", "November", "December",
	// }
	// for _, sheet := range respSpreadsheet.Sheets {
	// 	for monthIndex, monthName := range missingMonths {
	// 		if sheet.Properties.Title == monthName+"Expenses" {
	// 			missingMonths = append(missingMonths[:monthIndex], missingMonths[monthIndex+1:]...)
	// 			break
	// 		}
	// 	}
	// }

	// if len(missingMonths) > 0 {
	// 	fmt.Println("Creating expense sheets for months: ", missingMonths)
	// 	var batchAddSheetsReq sheets.BatchUpdateSpreadsheetRequest
	// 	for _, monthName := range missingMonths {
	// 		batchAddSheetsReq.Requests = append(batchAddSheetsReq.Requests, &sheets.Request{
	// 			AddSheet: &sheets.AddSheetRequest{
	// 				Properties: &sheets.SheetProperties{
	// 					Title: monthName + "Expenses",
	// 				},
	// 			},
	// 		})
	// 	}
	// 	_, err = srv.Spreadsheets.BatchUpdate(spreadsheetId, &batchAddSheetsReq).Do()
	// 	if err != nil {
	// 		log.Fatalf("Failed to add sheets: %v", err)
	// 	}
	// } else {
	// 	fmt.Println("All months already have expense sheets.")
	// }

	// fmt.Println("Writing...")
	// writeRange := readRange
	// newValues := &sheets.ValueRange{
	// 	MajorDimension: "ROWS",
	// 	Range:          writeRange,
	// 	Values: [][]interface{}{
	// 		{2, 4},
	// 		{6, 8},
	// 	},
	// }
	// _, err = srv.Spreadsheets.Values.Update(spreadsheetId, writeRange, newValues).ValueInputOption("USER_ENTERED").Do()
	// if err != nil {
	// 	log.Fatalf("Unable to write data to sheet: %v", err)
	// }

	statementFile, err := os.Open("sample-statement.csv")
	if err != nil {
		log.Fatalf("Failed to open CSV file: %v", err)
	}
	defer statementFile.Close()

	csvReader := csv.NewReader(statementFile)
	statementData, err := csvReader.ReadAll()
	if err != nil {
		log.Fatalf("Failed to read CSV file: %v", err)
	}
	if len(statementData) == 0 {
		log.Fatalf("Empty CSV file.")
	}

	headerRow := statementData[0]
	statementData = statementData[1:]

	expectedColumnNames := []string{
		"Transaction Date",
		"Post Date",
		"Description",
		"Category",
		"Type",
		"Amount",
		"Memo",
	}
	if len(headerRow) != len(expectedColumnNames) {
		log.Fatalf("Statement format error. Expected %d columns in header row but got %d.", len(expectedColumnNames), len(headerRow))
	}
	for i, columnName := range headerRow {
		if expectedColumnNames[i] != columnName {
			log.Fatalf("Statement format error. Expected column %d to be named %s but instead is named %s.", i+1, expectedColumnNames[i], columnName)
		}
	}

	longestDescription := 0
	longestCategory := 0
	longestType := 0

	statement := make([]chaseCreditItem, 0, len(statementData))
	for i, row := range statementData {
		if len(row) != len(expectedColumnNames) {
			log.Fatalf("Statement format error. Expected %d columns in row %d but got %d.", len(expectedColumnNames), i+2, len(row))
		}
		amount, err := strconv.ParseFloat(row[5], 64)
		if err != nil {
			log.Fatalf("Could not parse amount as float64 from row %d: %v", i+2, err)
		}
		if len(row[2]) > longestDescription {
			longestDescription = len(row[2])
		}
		if len(row[3]) > longestCategory {
			longestCategory = len(row[3])
		}
		if len(row[4]) > longestType {
			longestType = len(row[4])
		}
		statement = append(statement, chaseCreditItem{
			transactionDate: row[0],
			postedDate:      row[1],
			description:     row[2],
			category:        row[3],
			itemType:        row[4],
			amount:          amount,
			memo:            row[6],
		})
	}

	keywordsFile, err := os.Open("keywords.json")
	if err != nil {
		log.Fatalf("Failed to open keywords.json: %v", err)
	}
	defer keywordsFile.Close()

	keywordsFileBytes, err := io.ReadAll(keywordsFile)
	if err != nil {
		log.Fatalf("Failed to read keywords.json: %v", err)
	}
	var kw keywords
	err = json.Unmarshal(keywordsFileBytes, &kw)
	if err != nil {
		log.Fatalf("Failed to unmarshal keywords: %v", err)
	}
	kwMap, err := buildKeywordMap(kw)
	if err != nil {
		log.Fatalf("Failed to validate keywords: %v", err)
	}

	revisedStatement := make([]transaction, 0)
	for _, item := range statement {
		t := transaction{
			date:        item.transactionDate,
			description: item.description,
			amount:      item.amount,
		}

		var foundMatch bool
		t.category, foundMatch = keywordSearch(kwMap, item.description)
		if !foundMatch {
			fmt.Printf("Transaction Date\tPost Date\t%s\t%s\t%s\tAmount\t\tMemo\n", "Description"+strings.Repeat(" ", longestDescription-len("Description")), "Category"+strings.Repeat(" ", longestCategory-len("Category")), "Type"+strings.Repeat(" ", longestType-len("Type")))
			fmt.Printf("%s\t\t%s\t%s\t%s\t%s\t%f\t%s\n\n", item.transactionDate, item.postedDate, item.description+strings.Repeat(" ", longestDescription-len(item.description)), item.category+strings.Repeat(" ", longestCategory-len(item.category)), item.itemType+strings.Repeat(" ", longestType-len(item.itemType)), item.amount, item.memo)

			fmt.Printf("Please provide a description of this transaction. Press Enter to accept the default description. Submit 'skip' to not include the transaction in the final statement.\n")
			reader := bufio.NewReader(os.Stdin)
			description, err := reader.ReadString('\n')
			if err != nil {
				log.Fatalf("Failed to scan description: %v", err)
			}
			description = strings.TrimSuffix(description, "\n")
			description = strings.TrimSuffix(description, "\r")
			if description == "skip" {
				fmt.Printf("This transaction will be skipped.\n\n")
				continue
			}
			if len(description) != 0 {
				t.description = description
			}
			fmt.Printf("Received description: %s\n\n", t.description)

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
				log.Fatalf("Received invalid category enumeration %d.", categoryEnum)
			}
			t.category = categories[categoryEnum-1]
			fmt.Printf("Received Category: %d. %s\n", categoryEnum, t.category)
		} else if t.category == "skip" {
			fmt.Printf("Skipping transaction %+v based on discovery of a skip keyword in the description.\n\n", item)
			continue
		}

		log.Printf("Adding transaction %+v to statement.\n\n", t)
		revisedStatement = append(revisedStatement, t)
	}

	fmt.Printf("Full revised statement:\n")
	for _, t := range revisedStatement {
		fmt.Printf("%+v\n", t)
	}
}

func buildKeywordMap(kw keywords) (keywordMap, error) {
	categories := map[string][]string{
		"Rent":                 kw.Rent,
		"Utilities":            kw.Utilities,
		"Groceries/Toiletries": kw.GroceriesToiletries,
		"Food/Drinks Out":      kw.FoodDrinksOut,
		"Gas":                  kw.Gas,
		"Other (Need)":         kw.OtherNeed,
		"Other (Want)":         kw.OtherWant,
		"Gift Giving":          kw.GiftGiving,
		"Donations":            kw.Donations,
		"skip":                 kw.Skip,
	}
	kwMap := keywordMap{}
	for category, categoryWords := range categories {
		err := kwMap.add(category, categoryWords)
		if err != nil {
			return nil, fmt.Errorf("failed to add keywords for %s to keyword map: %v", category, err)
		}
	}
	return kwMap, nil
}

func (kwMap keywordMap) add(category string, words []string) error {
	if kwMap == nil {
		return errors.New("keywordMap is nil")
	}
	for _, word := range words {
		if _, ok := kwMap[word]; ok {
			return fmt.Errorf("found duplicate word: %s", word)
		}
		kwMap[word] = category
	}
	return nil
}

func keywordSearch(kwMap keywordMap, description string) (category string, foundMatch bool) {
	for word, associatedCategory := range kwMap {
		if strings.Contains(strings.ToLower(description), strings.ToLower(word)) {
			return associatedCategory, true
		}
	}
	return "", false
}
