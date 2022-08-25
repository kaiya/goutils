package file

import (
	"bufio"
	"encoding/csv"
	"io"
	"os"

	"github.com/goccy/go-json"
	"github.com/pkg/errors"
)

func ReadCsvRecordsFromFile(filename string) ([][]string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, errors.Wrap(err, "open file")
	}
	csvReader := csv.NewReader(file)
	csvReader.FieldsPerRecord = -1
	records, err := csvReader.ReadAll()
	if err != nil {
		return nil, errors.Wrap(err, "read csv")
	}
	return records, nil
}

func ParseJsonFile(path string, out interface{}) error {
	file, err := os.Open(path)
	if err != nil {
		return errors.Wrap(err, "open file")
	}
	bytes, err := io.ReadAll(bufio.NewReader(file))
	if err != nil {
		return errors.Wrap(err, "readall")
	}
	return json.Unmarshal(bytes, out)
}

// json format {"data": ["a","b"]}
func DiffJsonFiles(allPath, setPath string) ([]string, error) {
	setDocs := struct {
		Data []string `json:"data"`
	}{}
	err := ParseJsonFile(setPath, &setDocs)
	if err != nil {
		return nil, err
	}
	allDocs := struct {
		Data []string `json:"data"`
	}{}
	err = ParseJsonFile(allPath, &allDocs)
	if err != nil {
		return nil, err
	}
	return DiffSlices(allDocs.Data, setDocs.Data), nil
}

// get not exists of all in set
// TODO replace in place
func DiffSlices(all, set []string) []string {
	setMap := make(map[string]struct{})
	for _, s := range set {
		setMap[s] = struct{}{}
	}
	ret := make([]string, 0)
	for _, a := range all {
		if _, ok := setMap[a]; !ok {
			ret = append(ret, a)
		}
	}
	return ret
}
