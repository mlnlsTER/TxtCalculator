package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"io/fs"
	"os"
	"regexp"
	"strconv"
)

func calculate(expression string) (int, error) {
	re := regexp.MustCompile(`(\d+)\s*([\+\-\*/])\s*(\d+)`)
	match := re.FindStringSubmatch(expression)

	if len(match) != 4 {
		return 0, fmt.Errorf("invalid expression: %s", expression)
	}

	num1, err := strconv.Atoi(match[1])
	if err != nil {
		return 0, fmt.Errorf("invalid number: %s", match[1])
	}

	num2, err := strconv.Atoi(match[3])
	if err != nil {
		return 0, fmt.Errorf("invalid number: %s", match[3])
	}

	operator := match[2]

	switch operator {
	case "+":
		return num1 + num2, nil
	case "-":
		return num1 - num2, nil
	case "*":
		return num1 * num2, nil
	case "/":
		if num2 == 0 {
			return 0, fmt.Errorf("division by zero")
		}
		return num1 / num2, nil
	default:
		return 0, fmt.Errorf("unsupported operator: %s", operator)
	}
}

func processFile(inputFileName string, outputFileName string) error {
	inputData, err := func() ([]byte, error) {
		f, err := os.OpenFile(inputFileName, os.O_RDONLY, 0)
		if err != nil {
			return nil, err
		}
		defer f.Close()
		var size int
		if info, err := f.Stat(); err == nil {
			size64 := info.Size()
			if int64(int(size64)) == size64 {
				size = int(size64)
			}
		}
		size++
		if size < 512 {
			size = 512
		}
		data := make([]byte, 0, size)
		for {
			if len(data) >= cap(data) {
				d := append(data[:cap(data)], 0)
				data = d[:len(data)]
			}
			n, err := f.Read(data[len(data):cap(data)])
			data = data[:len(data)+n]
			if err != nil {
				if err == io.EOF {
					err = nil
				}
				return data, err
			}
		}
	}()
	if err != nil {
		return err
	}

	lines := bufio.NewScanner(bytes.NewReader(inputData))

	outputFile, err := os.Create(outputFileName)
	if err != nil {
		return err
	}
	defer outputFile.Close()

	writer := bufio.NewWriter(outputFile)

	for lines.Scan() {
		line := lines.Text()

		line = line[:len(line)-1]

		result, err := calculate(line)
		if err == nil {
			output := fmt.Sprintf("%s=%d\n", line[:len(line)-1], result)
			writer.WriteString(output)
		}
	}

	writer.Flush()

	return nil
}

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: program input.txt output.txt")
		os.Exit(1)
	}

	inputFileName := os.Args[1]
	outputFileName := os.Args[2]

	if _, err := os.Stat(outputFileName); err == nil {
		err := os.WriteFile(outputFileName, []byte(nil), fs.FileMode(0644))
		if err != nil {
			fmt.Println("Error clearing output file:", err)
			os.Exit(1)
		}
	}

	err := processFile(inputFileName, outputFileName)
	if err != nil {
		fmt.Println("Error processing file:", err)
		os.Exit(1)
	}

	fmt.Println("Calculation results written to", outputFileName)
}
