package main

import (
	"flag"
	"io"
	"io/ioutil"
	"os"

	"newrelic-prometheus/configurator"

	log "github.com/sirupsen/logrus"
)

const (
	inputErrCode = iota + 1
	outputErrCode
	parseErrCode
)

var logger = log.StandardLogger()

func readInput(inputPath string) ([]byte, error) {
	var reader io.Reader
	if inputPath == "" {
		reader = os.Stdin
	} else {
		fileReader, err := os.Open(inputPath)
		if err != nil {
			return nil, err
		}
		defer closeLoggingErr(fileReader)
		reader = fileReader
	}
	return ioutil.ReadAll(reader)
}

func writeOutput(outputPath string, output []byte) error {
	var writer io.Writer
	if outputPath == "" {
		writer = os.Stdout
	} else {
		fileWriter, err := os.Create(outputPath)
		if err != nil {
			return err
		}
		defer closeLoggingErr(fileWriter)
		writer = fileWriter
	}
	_, err := writer.Write(output)
	return err
}

func closeLoggingErr(f *os.File) {
	if err := f.Close(); err != nil {
		logger.Errorf("Fail closing file: %s", err)
	}
}

func main() {
	inputFlag := flag.String("input", "", "Input file to load the configuration from, defaults to stdin.")
	outputFlag := flag.String("output", "", "Output file to use as output, defaults to stdout.")
	flag.Parse()

	input, err := readInput(*inputFlag)
	if err != nil {
		logger.Errorf("Error loading the input: %s", err)
		os.Exit(inputErrCode)
	}
	output, err := configurator.Parse(input)
	if err != nil {
		logger.Errorf("Error parsing the configuration: %s", err)
		os.Exit(parseErrCode)
	}
	if err := writeOutput(*outputFlag, output); err != nil {
		logger.Errorf("Error parsing the configuration: %s", err)
		os.Exit(outputErrCode)
	}
}
