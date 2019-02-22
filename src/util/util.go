package util

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

func GetScanner(fname string) (*os.File, *bufio.Scanner) {
	fmt.Println("File Name", fname)
	file, err := os.Open(fname)
	if err != nil {
		log.Fatal(err)
		file.Close()
	}

	scanner := bufio.NewScanner(file)
	return file, scanner
}
