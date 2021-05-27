// Alex Eidt
// Accepts user input and creates the Mosaic GIFs.

package main

import (
	"flag"
	"fmt"
	"image/color"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"

	"github.com/fogleman/gg"
)

const (
	frames = "Frames"
	gifs   = "GIFs"
)

func main() {
	// Create output directory "Frames" to store frames.
	frames_dir := filepath.Join(".", frames)
	os.MkdirAll(frames_dir, os.ModePerm)
	// Create output directory "GIFs" to store GIFs
	gifs_dir := filepath.Join(".", gifs)
	os.MkdirAll(gifs_dir, os.ModePerm)

	// Parse command line args.
	fontsize := flag.Int("font", 12, "Font Size for ASCII Graphics.")
	ascii := flag.Bool("ascii", false, "Use ASCII Graphics.")
	hascolor := flag.Bool("color", false, "Include color with ASCII.")
	grayscale := flag.Bool("grayscale", false, "Grayscale the image.")
	keep := flag.Bool("keep", false, "Keep frames used for GIF.")
	background := flag.Bool("b", false, "White background for ASCII if included. Transparent otherwise.")
	fps := flag.Float64("fps", 1.0, "GIF Frames per second.")

	flag.Parse()

	args := flag.Args()

	filename := args[0]
	output := args[1]

	// Create Moasic Frames.
	n := Mosaiic(filename, *grayscale, *ascii, *background, *hascolor, *fontsize)

	// Remove frames if specified.
	if !*keep && n != -1 {
		if !*ascii {
			n++
		}
		for i := 0; i < n; i++ {
			defer os.Remove(filepath.Join(frames, strconv.Itoa(i)+".png"))
		}
	}

	// Call Python Script to create GIF from frames.
	cmd := exec.Command("python", "process.py", output, fmt.Sprintf("%f", *fps))
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr
	cmd.Run()
}

// Creates each Mosaic Frame.
func Mosaiic(
	filename string,
	grayscale,
	ascii,
	background,
	hascolor bool,
	fontsize int,
) int {
	pixels := Pixels(filename)
	if pixels == nil {
		return -1 // File not found.
	}

	h, w := len(pixels), len(pixels[0])

	tree := BuildTree(h, w)

	var chars [][]int
	if ascii {
		chars = AsciiChars(pixels)
	} else {
		CopyImage(pixels, len(tree), grayscale)
	}

	var ASCII []byte
	var lines []string
	var colors []color.Color
	var im [][]color.Color

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
		canvas.SavePNG(filepath.Join(frames, strconv.Itoa(i)+".png"))
		level <<= 1
	}
	return len(tree)
}
