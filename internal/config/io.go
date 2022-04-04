package config

import (
	"io/ioutil"
	"os"

	"gopkg.in/yaml.v3"
)

func StructToNode(in interface{}) (yaml.Node, error) {
	var node yaml.Node
	record, err := yaml.Marshal(in)
	if err != nil {
		return node, err
	}

	err = yaml.Unmarshal(record, &node)
	if err != nil {
		return node, err
	}
	return node, nil
}

func WriteString(data string, file string) error {
	err := ioutil.WriteFile(file, []byte(data), os.ModePerm)
	if err != nil {
		return err
	}

	return nil
}

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

func Parse(file string, out interface{}) error {
	data, err := ioutil.ReadFile(file)
	if err != nil {
		return err
	}

	err = yaml.Unmarshal(data, out)
	if err != nil {
		return err
	}
	return nil
}
