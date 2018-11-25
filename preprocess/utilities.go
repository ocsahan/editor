package preprocess

import (
	"image"
	"image/draw"
	"image/png"
	"os"

	"../effects"
)

type work struct {
	effect      int
	size        int
	startRow    int
	inBuffer    *image.RGBA
	outBuffer   *image.RGBA
	outFileName string
}

func parallelExecute(workChannel <-chan parallelWork) {
	for {
		parallelWork, isOpen := <-workChannel
		if !isOpen {
			break
		}
		parallelWork.work.execute()
		parallelWork.wg.Done()

		parallelWork.wg.Wait()
	}
}

func (work *work) execute() {
	effects.AddEffect(work.effect, work.inBuffer, work.outBuffer, work.startRow, work.size)
}

func (work *work) save() {
	outFile, _ := os.Create(work.outFileName)
	for i := 3; i < len(work.outBuffer.Pix); i += 4 {
		work.outBuffer.Pix[i] = 255
	}
	oldMax := work.outBuffer.Rect.Max
	newMax := image.Point{oldMax.X - 1, oldMax.Y - 1}
	work.outBuffer.Rect.Max = newMax
	work.outBuffer.Rect.Min = image.Point{1, 1}
	png.Encode(outFile, work.outBuffer)
	outFile.Close()
}

func padImage(img image.Image) *image.RGBA {
	destBounds := img.Bounds()
	destBounds.Max.X += 2
	destBounds.Max.Y += 2

	drawArea := img.Bounds()
	drawArea.Min.X++
	drawArea.Min.Y++
	drawArea.Max.X++
	drawArea.Max.Y++

	rgba := image.NewRGBA(destBounds)
	draw.Draw(rgba, drawArea, img, img.Bounds().Min, draw.Src)

	xLen := destBounds.Max.X
	yLen := destBounds.Max.Y
	bottomLeft := (xLen * 4 * (yLen - 1))
	topRight := (xLen - 1) * 4

	for i := 3; i < (xLen * 4); i += 4 {
		rgba.Pix[i-3] = 255
		rgba.Pix[i-2] = 255
		rgba.Pix[i-1] = 255
		rgba.Pix[i] = 255

		rgba.Pix[bottomLeft+i-3] = 255
		rgba.Pix[bottomLeft+i-2] = 255
		rgba.Pix[bottomLeft+i-1] = 255
		rgba.Pix[bottomLeft+i] = 255
	}

	for i := 3; i < (xLen * 4 * yLen); i += (xLen * 4) {
		rgba.Pix[i] = 255
		rgba.Pix[i-3] = 255
		rgba.Pix[i-2] = 255
		rgba.Pix[i-1] = 255

		rgba.Pix[i+topRight] = 255
		rgba.Pix[i-3+topRight] = 255
		rgba.Pix[i-2+topRight] = 255
		rgba.Pix[i-1+topRight] = 255
	}

	return rgba
}
