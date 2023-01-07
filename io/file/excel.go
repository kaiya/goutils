package file

import (
	"fmt"

	"github.com/pkg/errors"
	"github.com/xuri/excelize"
)

func ReadExcelRecords(fileName, sheetName string) ([][]string, error) {
	f, err := excelize.OpenFile(fileName)
	if err != nil {
		return nil, errors.Wrap(err, "openFile")
	}
	defer func() {
		// Close the spreadsheet.
		if err := f.Close(); err != nil {
			fmt.Println(err)
		}
	}()
	return f.GetRows(sheetName)
}
