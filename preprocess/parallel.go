package preprocess

import (
	"bufio"
	"image"
	"image/png"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

type parallelWork struct {
	wg   *sync.WaitGroup
	work work
}

func Parallel(fileName string, noOfThreads int) {
	file, _ := os.Open(fileName)
	directory := filepath.Dir(fileName)
	reader := bufio.NewReader(file)
	maxBufferSize := noOfThreads * 4
	workChannel := make(chan parallelWork, maxBufferSize)

	for i := 1; i < noOfThreads; i++ {
		go parallelExecute(workChannel)
	}

	for {
		line, _, err := reader.ReadLine()
		if err != nil || len(line) == 0 {
			break
		}

		command := strings.Split(string(line), ",")
		inFile := directory + "/" + command[0]
		outFile := directory + "/" + command[1]
		effects := command[2:]

		effectCodes := make([]int, len(effects))

		for i, effect := range effects {
			switch effect {
			case "S":
				effectCodes[i] = 0
			case "E":
				effectCodes[i] = 1
			case "B":
				effectCodes[i] = 2
			case "G":
				effectCodes[i] = 3
			}
		}

		CreateParallelWork(inFile, noOfThreads, workChannel, outFile, effectCodes...)
		var pWork parallelWork
		for i := 0; i < len(effects); i++ {
			pWork = <-workChannel
			pWork.work.execute()
			pWork.wg.Done()

			pWork.wg.Wait()
		}
		pWork.work.save()
	}
	close(workChannel)
}

func CreateParallelWork(fileName string, threads int, workChannel chan<- parallelWork, outFileName string, effects ...int) {
	file, _ := os.Open(fileName)
	img, _ := png.Decode(file)
	file.Close()

	inBuffer := padImage(img)
	imgHeight := img.Bounds().Max.Y
	canvasHeight := imgHeight + 2

	blockSize := (imgHeight / threads)

	for _, effect := range effects {

		outBuffer := image.NewRGBA(inBuffer.Rect)
		rowCounter := 1
		leftovers := imgHeight - blockSize*threads
		var wg sync.WaitGroup

		for rowCounter < canvasHeight-1 {
			workSize := blockSize

			if leftovers > 0 {
				workSize++
				leftovers--
			}

			wg.Add(1)
			work := work{effect: effect, size: workSize, startRow: rowCounter, inBuffer: inBuffer, outBuffer: outBuffer, outFileName: outFileName}
			workChannel <- parallelWork{work: work, wg: &wg}

			rowCounter += workSize
		}
		inBuffer = outBuffer
	}
}
