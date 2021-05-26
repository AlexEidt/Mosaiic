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
	fontsize := flag.Int("f", 12, "Font Size for ASCII Graphics.")
	ascii := flag.Bool("ascii", false, "Use ASCII Graphics.")
	hascolor := flag.Bool("color", true, "Include color with ASCII.")
	grayscale := flag.Bool("grayscale", false, "Grayscale the image.")
	keep := flag.Bool("keep", false, "Keep frames used for GIF.")
	delay := flag.Int("d", 100, "GIF Frame Delay in 1/100 of a second. (default 100).")

	flag.Parse()

	args := flag.Args()

	filename := args[0]

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
		if hascolor {
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
			if hascolor {
				im[idx] = colors[y : y+level]
			}
			idx++
		}
		dc := gg.NewContext(w, h)
		if ascii {
			CreateAsciiImage(dc, lines, im, tree[i], fontsize, i)
		} else {
			CreateMosaicImage(dc, im, tree[i], i)
		}
		level <<= 1
	}
}
