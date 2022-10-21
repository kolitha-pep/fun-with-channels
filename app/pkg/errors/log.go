package errors

import (
	"bufio"
	"os"
)

func Log(err error) {
	// open file for writing simple moving average values
	f, _ := os.OpenFile("errors.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	defer f.Close()

	w := bufio.NewWriter(f)
	w.WriteString(err.Error())
	w.Flush()
}
