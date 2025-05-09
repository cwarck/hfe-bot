package spreadsheets

import (
	"context"
	"fmt"
	"strings"
	"time"

	"google.golang.org/api/sheets/v4"

	"hfe-go/pkg/config"
)

const (
	TimeFormat       = "1/2/2006"
	ValueInputOption = "USER_ENTERED"
)

// AddExpense adds an expense to the spreadsheet.
func AddExpense(sheetsCtx config.GoogleSheets, amount float64, category string, comment string) error {
	service, err := sheets.NewService(context.Background())
	if err != nil {
		return err
	}

	if strings.HasPrefix(comment, "https://") {
		comment = fmt.Sprintf(`=HYPERLINK("%s", "receipt")`, comment)
	}

	appendRange := fmt.Sprintf("%s!A1:D1", sheetsCtx.SheetName)
	appendValues := &sheets.ValueRange{
		Range: appendRange,
		Values: [][]any{{
			time.Now().Format(TimeFormat),
			category,
			amount,
			comment,
		}},
	}
	req := service.Spreadsheets.Values.Append(sheetsCtx.SpreadsheetId, appendRange, appendValues)
	_, err = req.ValueInputOption(ValueInputOption).Do()
	return err
}
