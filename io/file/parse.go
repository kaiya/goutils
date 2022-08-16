package file

import (
	"bufio"
	"io"
	"os"

	"github.com/goccy/go-json"
	"github.com/pkg/errors"
)

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
