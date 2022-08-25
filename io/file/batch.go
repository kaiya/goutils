package file

func Slice2Records(cells []string, batchSize int) [][]string {
	ret := make([][]string, 0)
	rowCount := 0
	rowCells := make([]string, 0, batchSize)
	for _, cell := range cells {
		rowCells = append(rowCells, cell)
		rowCount++
		if rowCount%batchSize == 0 {
			ret = append(ret, rowCells)
			rowCells = make([]string, 0)
			rowCount = 0
		}
	}
	if rowCount != 0 {
		ret = append(ret, rowCells)
	}
	return ret
}

func SingleColumn2Table(records [][]string, idx int) []string {
	ret := make([]string, 0, len(records))
	for _, record := range records {
		ret = append(ret, record[idx])
	}
	return ret
}
