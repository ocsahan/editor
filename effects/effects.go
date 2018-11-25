package effects

import (
	"image"
)

func convolutePixel(kernel [9]int, denominator int, inBuffer *image.RGBA, outBuffer *image.RGBA, index int) {
	rowLength := inBuffer.Rect.Max.X

	for i := 0; i < 3; i++ {
		sum := 0
		pixel := index*4 + i
		up := pixel - rowLength*4
		down := pixel + rowLength*4

		sum += kernel[8] * int(inBuffer.Pix[up-4])
		sum += kernel[7] * int(inBuffer.Pix[up])
		sum += kernel[6] * int(inBuffer.Pix[up+4])
		sum += kernel[5] * int(inBuffer.Pix[pixel-4])
		sum += kernel[4] * int(inBuffer.Pix[pixel])
		sum += kernel[3] * int(inBuffer.Pix[pixel+4])
		sum += kernel[2] * int(inBuffer.Pix[down-4])
		sum += kernel[1] * int(inBuffer.Pix[down])
		sum += kernel[0] * int(inBuffer.Pix[down+4])

		result := sum / denominator
		if result > 255 {
			outBuffer.Pix[pixel] = 255
		} else if result < 0 {
			outBuffer.Pix[pixel] = 0
		} else {
			outBuffer.Pix[pixel] = uint8(result)
		}
	}
}

func convoluteImage(kernel [9]int, denominator int, inBuffer *image.RGBA, outBuffer *image.RGBA, startRow int, noOfRows int) {
	rowLength := inBuffer.Rect.Max.X

	for i := startRow; i < startRow+noOfRows; i++ {
		for j := rowLength*i + 1; j < rowLength*(i+1)-1; j++ {
			convolutePixel(kernel, denominator, inBuffer, outBuffer, j)
		}
	}
}

func grayScaleImage(inBuffer *image.RGBA, outBuffer *image.RGBA, startRow int, noOfRows int) {
	rowLength := inBuffer.Rect.Max.X

	for i := startRow; i < startRow+noOfRows; i++ {
		for j := (rowLength*i + 1) * 4; j < (rowLength*(i+1)-1)*4; j += 4 {
			sum := 0
			sum += int(inBuffer.Pix[j])
			sum += int(inBuffer.Pix[j+1])
			sum += int(inBuffer.Pix[j+2])
			avg := uint8(sum / 3)
			outBuffer.Pix[j] = avg
			outBuffer.Pix[j+1] = avg
			outBuffer.Pix[j+2] = avg
		}
	}
}

func AddEffect(effect int, inBuffer *image.RGBA, outBuffer *image.RGBA, startRow int, noOfRows int) {
	var kernel [9]int
	var denominator int

	switch effect {
	case 0:
		kernel = [9]int{0, -1, 0 - 1, 5, -1, 0, -1, 0}
		denominator = 1
		convoluteImage(kernel, denominator, inBuffer, outBuffer, startRow, noOfRows)
	case 1:
		kernel = [9]int{-1, -1, -1, -1, 8, -1, -1, -1, -1}
		denominator = 1
		convoluteImage(kernel, denominator, inBuffer, outBuffer, startRow, noOfRows)
	case 2:
		kernel = [9]int{1, 1, 1, 1, 1, 1, 1, 1, 1}
		denominator = 9
		convoluteImage(kernel, denominator, inBuffer, outBuffer, startRow, noOfRows)
	case 3:
		grayScaleImage(inBuffer, outBuffer, startRow, noOfRows)
	}
}
