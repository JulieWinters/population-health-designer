package config

import (
	"io/ioutil"
	"os"

	"gopkg.in/yaml.v2"
)

func Write(data interface{}, file string) error {

	bytes, err := yaml.Marshal(data)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(file, bytes, os.ModePerm)
	if err != nil {
		return err
	}

	return nil
}
