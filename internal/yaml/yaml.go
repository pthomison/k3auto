package yaml

import (
	"io"
	"regexp"
)

func YamlReadAndSplit(reader io.Reader) ([][]byte, error) {
	r := regexp.MustCompile(`---\n`)

	fb, err := io.ReadAll(reader)
	if err != nil {
		return nil, err
	}

	var objs [][]byte
	for _, obj := range r.Split(string(fb), -1) {
		if len(obj) != 0 {
			objs = append(objs, []byte(obj))
		}
	}

	return objs, err
}
