package main

import (
	"flag"
	"image/color"
	"os"
	"path/filepath"

	"github.com/fogleman/gg"
)

func main() {
	// Create output directory "Data" to store frames.
	outputDir := filepath.Join(".", "Data")
	os.MkdirAll(outputDir, os.ModePerm)

	// Parse command line args.
	fontsize := flag.Int("font", 12, "Font Size for ASCII Graphics.")
	ascii := flag.Bool("ascii", false, "Use ASCII Graphics.")
	hascolor := flag.Bool("color", false, "Include color with ASCII.")
	grayscale := flag.Bool("grayscale", false, "Grayscale the image.")
	keep := flag.Bool("keep", false, "Keep frames used for GIF.")
	//background := flag.Bool("b", false, "White background for ASCII.")
	delay := flag.Int("delay", 100, "GIF Frame Delay in 1/100 of a second.")
	// background
	// contrast

	flag.Parse()

	args := flag.Args()

	filename := args[0]

	//fmt.Println(*grayscale, *keep, *ascii, *hascolor, *fontsize, *delay)

	Mosaiic(filename, *grayscale, *keep, *ascii, *hascolor, *fontsize, *delay)
}

func Mosaiic(filename string, grayscale, keep, ascii, hascolor bool, fontsize, delay int) {
	image := Pixels(filename)
	if image == nil {
		return // File not found.
	}

	h, w := len(image), len(image[0])

	tree := BuildTree(h, w)

	var grayscaled [][]int
	if ascii {
		grayscaled = GrayScale(image)
	}

	var ASCII []byte
	var lines []string
	var colors []color.Color
	var im [][]color.Color

	level := 1
	for i := 0; i < len(tree); i++ {
		if ascii {
			ASCII = Ascii(tree[i], grayscaled)
			lines = make([]string, level)
		}
		if hascolor || !ascii {
			colors = BlockColor(tree[i], image, grayscale)
			im = make([][]color.Color, level)
		} else {
			im = nil
		}

		idx := 0
		for y := 0; y < level*level; y += level {
			if ascii {
				lines[idx] = string(ASCII[y : y+level])
			}
			if hascolor || !ascii {
				im[idx] = colors[y : y+level]
			}
			idx++
		}
		canvas := gg.NewContext(w, h)
		if ascii {
			CreateAsciiImage(canvas, lines, im, tree[i], fontsize, i)
		} else {
			CreateMosaicImage(canvas, im, tree[i], i)
		}
		level <<= 1
	}
}
