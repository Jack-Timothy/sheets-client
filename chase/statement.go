package chase

import (
	"fmt"

	"github.com/Jack-Timothy/sheets-client/keywords"
	"github.com/Jack-Timothy/sheets-client/standard"
)

type Statement []Transaction

func (s Statement) Standardize() (ss standard.Statement, err error) {
	kwMap, err := keywords.MapFromFile("keywords.json")
	if err != nil {
		return ss, fmt.Errorf("failed to make keyword map from file: %w", err)
	}

	ss = make([]standard.Transaction, 0)
	for i, t := range s {
		st, skip, err := t.standardize(kwMap)
		if err != nil {
			return nil, fmt.Errorf("failed to standardize item %d: %w", i, err)
		}
		if skip {
			continue
		}
		ss = append(ss, st)
	}
	return ss, nil
}
