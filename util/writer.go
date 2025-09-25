package util

import (
	"os"
)

// WriteResults writes results from a channel into a file (append mode)
func WriteResults(filename string, results <-chan string) error {
	f, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	for r := range results {
		_, _ = f.WriteString(r)
	}
	return nil
}
