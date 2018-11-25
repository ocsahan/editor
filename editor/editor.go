package main

import (
	"os"
	"runtime"
	"strconv"
	"strings"

	"../preprocess"
)

func main() {
	if len(os.Args) == 2 {
		fileName := os.Args[1]
		preprocess.Sequential(fileName)
	} else if len(os.Args) == 3 {
		fileName := os.Args[1]
		if strings.Contains(os.Args[2], "-p") {
			var noOfThreads int
			if strings.Index(os.Args[2], "=") == 2 {
				noOfThreads, _ = strconv.Atoi(os.Args[2][3:])
			} else {
				noOfThreads = runtime.NumCPU()
			}
			preprocess.Parallel(fileName, noOfThreads)
		}
	}
}
