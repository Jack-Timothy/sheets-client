package main

import (
	"context"
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/Jack-Timothy/sheets-client/chase"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)

// Retrieves a token, saves the token, then returns the generated client.
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
	csvFileName := "sample-statement.csv"
	csvContents, err := getCsvContents(csvFileName)
	if err != nil {
		log.Fatalf("Error getting contents of %s: %v", csvFileName, err)
	}

	chaseStatement, err := chase.CsvContentsToStatement(csvContents)
	if err != nil {
		log.Fatalf("Error converting csv contents to Chase statement: %v", err)
	}

	standardStatement, err := chaseStatement.Standardize()
	if err != nil {
		log.Fatalf("Error standardizing Chase statement: %v", err)
	}

	if err = standardStatement.AcceptUserEdits(); err != nil {
		log.Fatalf("Error during user edits of statement: %v", err)
	}

	b, err := os.ReadFile("credentials.json")
	if err != nil {
		log.Fatalf("Unable to read client secret file: %v", err)
	}

	// If modifying these scopes, delete your previously saved token.json.
	// For full list of scopes: https://developers.google.com/identity/protocols/oauth2/scopes#sheets.
	config, err := google.ConfigFromJSON(b, "https://www.googleapis.com/auth/spreadsheets")
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}
	client := getClient(config)

	srv, err := sheets.NewService(context.Background(), option.WithHTTPClient(client))
	if err != nil {
		log.Fatalf("Unable to retrieve Sheets client: %v", err)
	}

	// Prints the data in a test spreadsheet:
	// https://docs.google.com/spreadsheets/d/15KWFkIY-RW81leDLXqahARB0gtSnWAIGDg-lkx2g04Q/edit
	spreadsheetId := "1dnKqyF20h90PT1ualeQHZJkJVH1LpnhZMRNH4-kwxws"
	readRange := "Sheet1!A:D"
	resp, err := srv.Spreadsheets.Values.Get(spreadsheetId, readRange).Do()
	if err != nil {
		log.Fatalf("Unable to retrieve data from sheet: %v", err)
	}

	if len(resp.Values) == 0 {
		fmt.Println("No data found.")
		return
	}

	numRows := len(resp.Values) - 1
	fmt.Printf("Found %d rows of data.\n", numRows)
	for _, row := range resp.Values {
		output := ""
		for _, column := range row {
			output += column.(string) + ", "
		}
		output = output[:len(output)-2]
		fmt.Println(output)
	}

	fmt.Println("Writing...")
	rowNumToStartWrite := 2 + numRows
	rowNumToEndWrite := rowNumToStartWrite + len(standardStatement)
	writeRange := fmt.Sprintf("Sheet1!A%d:D%d", rowNumToStartWrite, rowNumToEndWrite)
	newValues := &sheets.ValueRange{
		MajorDimension: "ROWS",
		Range:          writeRange,
		Values:         standardStatement.GetRawData(),
	}
	_, err = srv.Spreadsheets.Values.Update(spreadsheetId, writeRange, newValues).ValueInputOption("USER_ENTERED").Do()
	if err != nil {
		log.Fatalf("Unable to write data to sheet: %v", err)
	}
}

func getCsvContents(fileName string) (csvContents [][]string, err error) {
	csvFile, err := os.Open(fileName)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer csvFile.Close()

	csvReader := csv.NewReader(csvFile)
	csvContents, err = csvReader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}
	if len(csvContents) == 0 {
		return nil, errors.New("file is empty")
	}
	return csvContents, nil
}
