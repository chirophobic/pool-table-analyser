package main

import (
	"github.com/anthonynsimon/bild/effect"
	"github.com/anthonynsimon/bild/imgio"
	"image"
)

func main() {
	onImage("res/images/one.jpg", "one")
	onImage("res/images/two.jpg", "two")
	onImage("res/images/three.jpg", "three")
}

func onImage(filename string, ext string) {
	img, err := imgio.Open(filename)
	if err != nil {
		panic(err)
	}

	operate(img, edgeDetect, "image-output/edge-detect-" + ext)
	operate(img, sobel, "image-output/sobel-" + ext)
	operate(img, invert, "image-output/invert-" + ext)
}

func operate(img image.Image, f func(image.Image) *image.RGBA, filename string) {
	output := f(img)
	if err := imgio.Save(filename, output, imgio.JPEG); err != nil {
		panic(err)
	}
}

func edgeDetect(img image.Image) *image.RGBA {
	return effect.EdgeDetection(img, 1.5)
}

func sobel(img image.Image) *image.RGBA {
	return effect.Sobel(img)
}

func invert(img image.Image) *image.RGBA {
	return effect.Invert(img)
}
