package datastore

import (
	"bufio"
	"os"
)

func WriteFile(string string, filePath string) error {
	// open file for writing simple moving average values
	f, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	w := bufio.NewWriter(f)
	w.WriteString(string)
	w.Flush()
	return nil
}
