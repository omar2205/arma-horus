package utils

import (
	"errors"
	"os"
)

func ReadlastLine(fileName string) (string, error) {
	file, err := os.Open(fileName)
	if err != nil {
		return "", err
	}
	defer file.Close()

	stat, err := file.Stat()
	if err != nil {
		return "", err
	}

	if stat.IsDir() {
		return "", errors.New("Not a file")
	}

	fileSize := stat.Size()
	var lastLine string
	buf := make([]byte, 1)

	for i := fileSize - 1; i > 0; i-- {
		file.ReadAt(buf, i-1)
		if buf[0] == '\n' {
			break
		}

		lastLine = string(buf) + lastLine
	}

	return lastLine, nil
}

func Contains[T string | int | bool](s []T, value T) bool {
	for _, v := range s {
		if v == value {
			return true
		}
	}

	return false
}
