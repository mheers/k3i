package helpers

import (
	"encoding/json"
	"fmt"

	"github.com/common-nighthawk/go-figure"
	"github.com/gocarina/gocsv"
	"gopkg.in/yaml.v3"
)

// PrintInfo print Info
func PrintInfo() {
	f := figure.NewColorFigure("k3i", "big", "red", true)
	figletStr := f.String()
	fmt.Println(figletStr)
	fmt.Println()
}

func PrintJSON(obj interface{}) error {
	b, err := json.MarshalIndent(obj, "", "  ")
	if err != nil {
		return err
	}
	fmt.Println(string(b))
	return nil
}

func PrintYAML(obj interface{}) error {
	b, err := yaml.Marshal(obj)
	if err != nil {
		return err
	}
	fmt.Println(string(b))
	return nil
}

func PrintCSV(obj interface{}) error {
	csv, err := gocsv.MarshalBytes(obj)
	if err != nil {
		return err
	}
	fmt.Println(string(csv))
	return nil
}
