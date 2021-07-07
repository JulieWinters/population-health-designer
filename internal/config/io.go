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

func Parse(file string, out interface{}) {

	data, err := ioutil.ReadFile(file)
	if err != nil {
		panic(err)
	}

	err = yaml.Unmarshal(data, out)
	if err != nil {
		panic(err)
	}
}
