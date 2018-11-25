package preprocess

import (
	"bufio"
	"image"
	"image/png"
	"os"
	"path/filepath"
	"strings"
)

func Sequential(fileName string) {
	file, _ := os.Open(fileName)
	directory := filepath.Dir(fileName)
	reader := bufio.NewReader(file)
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
		jobs := CreateSequentialWork(inFile, outFile, effectCodes...)
		for _, job := range jobs {
			job.execute()
		}
		jobs[len(jobs)-1].save()
	}
}

func CreateSequentialWork(fileName string, outFileName string, effects ...int) []work {
	file, _ := os.Open(fileName)
	img, _ := png.Decode(file)
	file.Close()

	inBuffer := padImage(img)
	imgHeight := img.Bounds().Max.Y
	jobs := make([]work, len(effects))

	for i, effect := range effects {

		outBuffer := image.NewRGBA(inBuffer.Rect)
		work := work{effect: effect, size: imgHeight, startRow: 1, inBuffer: inBuffer, outBuffer: outBuffer, outFileName: outFileName}
		jobs[i] = work
		inBuffer = outBuffer
	}
	return jobs
}
