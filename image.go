// Alex Eidt
// Contains some nice utils for image segmentation and processing.

package main

import (
	"fmt"
	"image"
	"image/color"
	"image/jpeg"
	"image/png"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/fogleman/gg"
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font/gofont/gomono"
)

// Draws the Mosaic Image as a PNG image.
func CreateMosaicImage(canvas *gg.Context, colors [][]color.Color, indices []int) {
	top_y := 0.0
	y_count := 0
	for y := 0; y < len(indices); y += 2 {
		top_x := 0.0
		x_count := 0
		for x := 1; x < len(indices); x += 2 {
			R, G, B, _ := colors[y_count][x_count].RGBA()
			x_count++
			canvas.DrawRectangle(top_x, top_y, float64(indices[x])-top_x, float64(indices[y])-top_y)
			canvas.SetRGB(float64(R)/256, float64(G)/256, float64(B)/256)
			canvas.Fill()
			top_x = float64(indices[x])
		}
		top_y = float64(indices[y])
		y_count++
	}
}

// Draws the ASCII matrix as a PNG image.
func CreateAsciiImage(
	canvas *gg.Context,
	lines []string,
	colors [][]color.Color,
	indices []int,
	fontsize int,
) {
	font, err := truetype.Parse(gomono.TTF)
	if err != nil {
		panic("Go Mono Font not found.")
	}
	face := truetype.NewFace(font, &truetype.Options{
		Size: float64(fontsize),
	})
	canvas.SetFontFace(face)
	hascolor := colors != nil
	if !hascolor {
		canvas.SetRGB(0, 0, 0)
	}
	_, ht := canvas.MeasureString(string(lines[0][0]))
	count_y := 0
	row := ht
	for y := 0; y < len(indices); y += 2 {
		for row < float64(indices[y]) {
			count_x := 0
			column := 0.0
			for x := 1; x < len(indices); x += 2 {
				for column < float64(indices[x]) {
					if hascolor {
						R, G, B, _ := colors[count_y][count_x].RGBA()
						canvas.SetRGB(float64(R)/256, float64(G)/256, float64(B)/256)
					}
					canvas.DrawString(string(lines[count_y][count_x]), column, row)
					column += ht
				}
				count_x++
			}
			row += ht
		}
		count_y++
	}
}

// Given a matrix of colors representing color values of an image, segments this
// "image" matrix into blocks of certain sizes determined by the values
// in "indices" and returns the new image matrix as a flattened array of colors.
func BlockColor(indices []int, image [][]color.Color, grayscale bool) []color.Color {
	colors := make([]color.Color, len(indices)*len(indices)/4)
	count := 0
	start_y := 0
	for y := 0; y < len(indices); y += 2 {
		start_x := 0
		for x := 1; x < len(indices); x += 2 {
			area := uint32(0)
			rgb := make([]uint32, 3)
			// Find the average color for every block.
			for _, valy := range image[start_y:indices[y]] {
				for _, valx := range valy[start_x:indices[x]] {
					var R, G, B uint32
					if grayscale {
						R, G, B, _ = color.GrayModel.Convert(valx).RGBA()
					} else {
						R, G, B, _ = valx.RGBA()
					}
					rgb[0] += R / 257
					rgb[1] += G / 257
					rgb[2] += B / 257
					area++
				}
			}
			start_x = indices[x]
			rgb[0] /= area
			rgb[1] /= area
			rgb[2] /= area
			colors[count] = color.NRGBA{uint8(rgb[0]), uint8(rgb[1]), uint8(rgb[2]), 1}
			count++
		}
		start_y = indices[y]
	}
	return colors
}

// Given a matrix of integers representing grayscale values, segments this
// "image" matrix into blocks of certain sizes determined by the values
// in "indices" and returns the new image matrix as a flattened array of ASCII
// characters.
func Ascii(indices []int, grayscaled [][]int) []byte {
	chars := " `.,|'\\/~!_-;:)(\"><?*+7j1ilJyc&vt0$VruoI=wzCnY32LTxs4Zkm5hg6qfU9paOS#eX8D%bdRPGFK@AMQNWHEB"
	ascii := make([]byte, len(indices)*len(indices)/4)
	count := 0
	start_y := 0
	for y := 0; y < len(indices); y += 2 {
		start_x := 0
		for x := 1; x < len(indices); x += 2 {
			// Sum together all pixels in a given "block".
			sum, area := 0, 0
			for _, valy := range grayscaled[start_y:indices[y]] {
				for _, valx := range valy[start_x:indices[x]] {
					sum += valx
					area++
				}
			}
			start_x = indices[x]
			index := len(chars) - (sum/area)*len(chars)/255
			// Normalize ASCII index.
			if index > len(chars)-1 {
				index = len(chars) - 1
			} else if index < 0 {
				index = 0
			}
			ascii[count] = chars[index]
			count++
		}
		start_y = indices[y]
	}
	return ascii
}

// Parses an image into a matrix of integers representing
// the grayscale value (from 0 - 255) of each pixel.
func AsciiChars(image [][]color.Color) [][]int {
	grayscaled := make([][]int, len(image))
	for y := 0; y < len(image); y++ {
		grayscaled[y] = make([]int, len(image[y]))
		for x := 0; x < len(image[y]); x++ {
			R, G, B, _ := image[y][x].RGBA()
			grayscaled[y][x] = (int(R+R+R+B+G+G+G+G) / 257) / 8
		}
	}
	return grayscaled
}

// Given a "filename" of an image, parses it and returns an
// array of colors (dimensions same as original image).
func Pixels(filename string) [][]color.Color {
	// Read image from "filename".
	file, err := os.Open(filename)
	if err != nil {
		fmt.Printf("%s not found.", filename)
		return nil
	}
	defer file.Close()

	var im image.Image
	if strings.HasSuffix(filename, "jpg") {
		im, err = jpeg.Decode(file)
	} else if strings.HasSuffix(filename, "png") {
		im, err = png.Decode(file)
	} else {
		im, _, err = image.Decode(file)
	}

	if err != nil {
		fmt.Printf("%s is an invalid image format. Could not parse.\n", filename)
		return nil
	}
	bounds := im.Bounds().Max

	pixels := make([][]color.Color, bounds.Y)
	// Fill in "pixels" with colors of the image.
	for y := 0; y < bounds.Y; y++ {
		pixels[y] = make([]color.Color, bounds.X)
		for x := 0; x < bounds.X; x++ {
			pixels[y][x] = im.At(x, y)
		}
	}
	return pixels
}

// Copies the given "image" into a new image. New image
// can also be a grayscaled version of the original.
func CopyImage(img [][]color.Color, index int, grayscale bool) {
	w, h := len(img[0]), len(img)
	new_img := image.NewRGBA(
		image.Rectangle{
			image.Point{0, 0},
			image.Point{w, h},
		},
	)
	if grayscale {
		for y := 0; y < h; y++ {
			for x := 0; x < w; x++ {
				new_img.Set(x, y, color.GrayModel.Convert(img[y][x]))
			}
		}
	} else {
		for y := 0; y < h; y++ {
			for x := 0; x < w; x++ {
				new_img.Set(x, y, img[y][x])
			}
		}
	}
	output, err := os.Create(filepath.Join(frames, strconv.Itoa(index)+".png"))
	if err != nil {
		panic(err)
	}
	png.Encode(output, new_img)
}
