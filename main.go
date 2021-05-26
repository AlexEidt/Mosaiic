package main

import (
	"flag"
	"image"
	"image/color"
	"image/draw"
	"image/gif"
	"os"
	"path/filepath"
	"strconv"

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
	background := flag.Bool("b", false, "White background for ASCII if included. Transparent otherwise.")
	delay := flag.Int("delay", 100, "GIF Frame Delay in 1/100 of a second.")
	// background
	// contrast

	flag.Parse()

	args := flag.Args()

	filename := args[0]
	output := args[1]

	Mosaiic(filename, output, *grayscale, *keep, *ascii, *background, *hascolor, *fontsize, *delay)
}

func Mosaiic(
	filename,
	output string,
	grayscale,
	keep,
	ascii,
	background,
	hascolor bool,
	fontsize,
	delay int,
) {
	pixels := Pixels(filename)
	if pixels == nil {
		return // File not found.
	}

	h, w := len(pixels), len(pixels[0])

	tree := BuildTree(h, w)

	var chars [][]int
	if ascii {
		chars = AsciiChars(pixels)
	}

	var ASCII []byte
	var lines []string
	var colors []color.Color
	var im [][]color.Color

	disposal := make([]byte, len(tree))
	// Prevent gif frame stacking/crossfading.
	for i := 0; i < len(tree); i++ {
		disposal[i] = gif.DisposalBackground
	}
	output_GIF := &gif.GIF{Disposal: disposal}

	level := 1
	for i := 0; i < len(tree); i++ {
		if ascii {
			ASCII = Ascii(tree[i], chars)
			lines = make([]string, level)
		}
		if hascolor || !ascii || grayscale {
			colors = BlockColor(tree[i], pixels, grayscale)
			im = make([][]color.Color, level)
		} else {
			im = nil
		}

		idx := 0
		for y := 0; y < level*level; y += level {
			if ascii {
				lines[idx] = string(ASCII[y : y+level])
			}
			if hascolor || !ascii || grayscale {
				im[idx] = colors[y : y+level]
			}
			idx++
		}

		canvas := gg.NewContext(w, h)
		if ascii {
			if background {
				// Set background of ASCII images to be white.
				canvas.SetRGB(1, 1, 1)
				canvas.Clear()
			}
			CreateAsciiImage(canvas, lines, im, tree[i], fontsize)
		} else {
			CreateMosaicImage(canvas, im, tree[i])
		}
		fname := filepath.Join("Data", strconv.Itoa(i)+".png")
		canvas.SavePNG(fname)
		level <<= 1

		f, _ := os.Open(fname)
		if !keep {
			defer os.Remove(fname)
		}
		defer f.Close()

		img, _, _ := image.Decode(f)

		var palette color.Palette
		if hascolor || !ascii || grayscale {
			palette = colors
		} else {
			palette = color.Palette{
				color.Transparent,
				color.White,
				color.Black,
			}
		}

		p := image.NewPaletted(
			img.Bounds(),
			palette,
		)

		//draw.Draw(p, p.Bounds(), img, img.Bounds().Min, draw.Over)
		//draw.Draw(p, img.Bounds(), img, image.Point{}, draw.Over)
		draw.Draw(p, p.Bounds(), img, img.Bounds().Min, draw.Over)
		output_GIF.Image = append(output_GIF.Image, p)
		output_GIF.Delay = append(output_GIF.Delay, delay)
	}

	file, err := os.Create(filepath.Join("Data", output+".gif"))
	if err != nil {
		panic(err)
	}
	gif.EncodeAll(file, output_GIF)
}
