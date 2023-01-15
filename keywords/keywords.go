package keywords

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
)

type Map map[string]string

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

func (kwMap Map) Search(description string) (category string, foundMatch bool) {
	for word, associatedCategory := range kwMap {
		if strings.Contains(strings.ToLower(description), strings.ToLower(word)) {
			return associatedCategory, true
		}
	}
	return "", false
}

func MapFromFile(fileName string) (kwMap Map, err error) {
	keywordsFile, err := os.Open(fileName)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer keywordsFile.Close()

	keywordsFileBytes, err := io.ReadAll(keywordsFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}
	var kw keywords
	err = json.Unmarshal(keywordsFileBytes, &kw)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal keywords: %w", err)
	}
	kwMap, err = buildKeywordMap(kw)
	if err != nil {
		return nil, fmt.Errorf("failed to build keyword map: %w", err)
	}
	return kwMap, err
}

func buildKeywordMap(kw keywords) (Map, error) {
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
	kwMap := Map{}
	for category, categoryWords := range categories {
		err := kwMap.add(category, categoryWords)
		if err != nil {
			return nil, fmt.Errorf("failed to add keywords for %s to keyword map: %v", category, err)
		}
	}
	return kwMap, nil
}

func (kwMap Map) add(category string, words []string) error {
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
