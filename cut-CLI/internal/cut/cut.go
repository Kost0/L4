package cut

import (
	"bufio"
	"errors"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
	"unicode"
)

type CutOptions struct {
	Fields    []int
	Delimiter string
	Separated bool
}

func ReadLines(inputFile string) ([]string, error) {
	var reader io.Reader

	if inputFile == "-" || inputFile == "" {
		reader = os.Stdin
	} else {
		file, err := os.Open(inputFile)
		if err != nil {
			return nil, err
		}
		defer func() {
			err = file.Close()
			if err != nil {
				log.Fatal(err)
			}
		}()
		reader = file
	}

	var lines []string
	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			break
		}
		lines = append(lines, line)
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return lines, nil
}

func ParseFields(fields string) ([]int, error) {
	res := make([]int, 0)
	prevNum := ""
	curNum := ""
	for i := range fields {
		if unicode.IsDigit(rune(fields[i])) {
			curNum += string(fields[i])
		} else if fields[i] == ',' {
			if prevNum == "" {
				digit, err := strconv.Atoi(curNum)
				if err != nil {
					return nil, err
				}
				res = append(res, digit-1)
			} else {
				prevDigit, err := strconv.Atoi(prevNum)
				if err != nil {
					return nil, err
				}
				lastDigit, err := strconv.Atoi(curNum)
				if err != nil {
					return nil, err
				}

				for j := prevDigit; j <= lastDigit; j++ {
					res = append(res, j-1)
				}
			}
			prevNum = ""
			curNum = ""
		} else if fields[i] == '-' {
			prevNum = curNum
			curNum = ""
		} else {
			return nil, errors.New("invalid field")
		}
	}

	if prevNum == "" {
		digit, err := strconv.Atoi(curNum)
		if err != nil {
			return nil, err
		}
		res = append(res, digit-1)
	} else {
		prevDigit, err := strconv.Atoi(prevNum)
		if err != nil {
			return nil, err
		}
		lastDigit, err := strconv.Atoi(curNum)
		if err != nil {
			return nil, err
		}

		for j := prevDigit; j <= lastDigit; j++ {
			res = append(res, j-1)
		}
	}

	if len(res) == 0 {
		return nil, errors.New("need fields")
	}

	return res, nil
}

func CutLines(lines []string, opts *CutOptions) [][]string {
	res := make([][]string, 0)

	for _, line := range lines {
		if strings.Contains(line, opts.Delimiter) {
			parts := strings.Split(line, opts.Delimiter)
			resParts := make([]string, 0)

			for _, field := range opts.Fields {
				if field >= len(parts) {
					break
				}
				resParts = append(resParts, parts[field])
			}

			res = append(res, resParts)
		} else {
			if !opts.Separated {
				res = append(res, []string{line})
			}
		}
	}

	return res
}
