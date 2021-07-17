package core

import (
	"io/ioutil"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

// FromYAML parses yaml file
func FromYAML(file string, dist interface{}) error {
	filename, _ := filepath.Abs(file)

	data, err := ioutil.ReadFile(filename)

	if err != nil {
		return err
	}

	return yaml.Unmarshal(data, dist)
}
