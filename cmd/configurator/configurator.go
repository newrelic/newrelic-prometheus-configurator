package main

import (
	"flag"
	"fmt"
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
	if inputPath == "" {
		return ioutil.ReadAll(os.Stdin)
	}
	fileReader, err := os.Open(inputPath)
	if err != nil {
		return nil, err
	}
	defer closeLoggingErr(fileReader)
	return ioutil.ReadAll(fileReader)
}

func writeOutput(outputPath string, output []byte) error {
	if outputPath == "" {
		_, err := os.Stdout.Write(output)
		return fmt.Errorf("error writing the output: %s", err)
	}
	fileWriter, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("error creating the output file: %s", err)
	}
	defer closeLoggingErr(fileWriter)
	_, err = fileWriter.Write(output)
	return fmt.Errorf("error writing output: %s", err)
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
	err = writeOutput(*outputFlag, output)
	if err != nil {
		logger.Errorf("Error writing the output configuration: %s", err)
		os.Exit(outputErrCode)
	}
}
