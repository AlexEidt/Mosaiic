package main

import (
	"fmt"
	"image/color"
	"image/jpeg"
	"os"
	"path/filepath"
	"strconv"

	"github.com/fogleman/gg"
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font/gofont/gomono"
)

var count int

func main() {
	//Mosaiic()
	im := Pixels("test.jpg")
	h, w := len(im), len(im[0])
	tree := BuildTree(h, w)
	//grayscaled := GrayScale(im)
	colorMap := make([][][]color.Color, len(tree))
	level := 1
	for i := 0; i < len(tree); i++ {
		colors := BlockColor(tree[i], im)
		//ascii := Ascii(tree[i], grayscaled)
		//lines := make([]string, level)
		im := make([][]color.Color, level)
		idx := 0
		for y := 0; y < level*level; y += level {
			//lines[idx] = string(ascii[y : y+level])
			im[idx] = colors[y : y+level]
			idx++
		}
		colorMap[i] = im
		//CreateAsciiImage(lines, im, w, h)
		//CreateMosaicImage(im, tree[i], w, h)
		level <<= 1
	}
	dc := gg.NewContext(w, h)
	count = 0
	RecursiveMosaic(dc, tree, colorMap, 0)
}

func RecursiveMosaic(dc *gg.Context, tree map[int][]int, colors [][][]color.Color, index int) {
	top_y := 0.0
	y_count := 0
	indices := tree[index]
	for y := 0; y < len(indices); y += 2 {
		top_x := 0.0
		x_count := 0
		for x := 1; x < len(indices); x += 2 {
			R, G, B, _ := colors[index][y_count][x_count].RGBA()
			x_count++
			dc.DrawRectangle(top_x, top_y, float64(indices[x])-top_x, float64(indices[y])-top_y)
			dc.SetRGB(float64(R)/256, float64(G)/256, float64(B)/256)
			dc.Fill()
			top_x = float64(indices[x])
		}
		top_y = float64(indices[y])
		y_count++
	}
	dc.SavePNG(filepath.Join("Data", strconv.Itoa(count)+"color.png"))
	count++
	if index != len(tree)-1 {
		RecursiveMosaic(dc, tree, colors, index+1)
	}
}

func CreateMosaicImage(colors [][]color.Color, indices []int, w, h int) {
	dc := gg.NewContext(w, h)
	top_y := 0.0
	y_count := 0
	for y := 0; y < len(indices); y += 2 {
		top_x := 0.0
		x_count := 0
		for x := 1; x < len(indices); x += 2 {
			R, G, B, _ := colors[y_count][x_count].RGBA()
			x_count++
			dc.DrawRectangle(top_x, top_y, float64(indices[x])-top_x, float64(indices[y])-top_y)
			dc.SetRGB(float64(R)/256, float64(G)/256, float64(B)/256)
			dc.Fill()
			top_x = float64(indices[x])
		}
		top_y = float64(indices[y])
		y_count++
	}
	dc.SavePNG(strconv.Itoa(len(colors)) + "color.png")
}

func CreateAsciiImage(lines []string, colors [][]color.Color, indices []int, w, h int) {
	dc := gg.NewContext(w, h)
	font, err := truetype.Parse(gomono.TTF)
	if err != nil {
		panic("Go Mono Font not found.")
	}
	face := truetype.NewFace(font, &truetype.Options{
		Size: 5,
	})
	dc.SetFontFace(face)
	row, column := 0.0, 0.0
	_, ht := dc.MeasureString(string(lines[0][0]))
	for y, line := range lines {
		column = 0.0
		row += ht
		for x, char := range line {
			R, G, B, _ := colors[y][x].RGBA()
			dc.SetRGB(float64(R)/256, float64(G)/256, float64(B)/256)
			dc.DrawString(string(char), column, row)
			column += ht
		}
	}
	dc.SavePNG(strconv.Itoa(len(lines)) + ".png")
}

func Mosaiic() {
	im := Pixels("test.jpg")

	h, w := len(im), len(im[0])

	tree := BuildTree(h, w)

	grayscaled := GrayScale(im)

	level := 1
	for i := 0; i < len(tree); i++ {
		colors := BlockColor(tree[i], im)
		ascii := Ascii(tree[i], grayscaled)
		lines := make([]string, level)
		im := make([][]color.Color, level)
		idx := 0
		for y := 0; y < level*level; y += level {
			lines[idx] = string(ascii[y : y+level])
			im[idx] = colors[y : y+level]
			idx++
		}
		//CreateAsciiImage(lines, im, w, h)
		CreateMosaicImage(im, tree[i], w, h)
		level <<= 1
	}
}

func BlockColor(indices []int, image [][]color.Color) []color.Color {
	colors := make([]color.Color, len(indices)*len(indices)/4)
	count := 0
	start_y := 0
	for y := 0; y < len(indices); y += 2 {
		start_x := 0
		for x := 1; x < len(indices); x += 2 {
			area := uint32(0)
			rgb := make([]uint32, 3)
			for _, valy := range image[start_y:indices[y]] {
				for _, valx := range valy[start_x:indices[x]] {
					R, G, B, _ := valx.RGBA()
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

func Ascii(indices []int, grayscaled [][]int) []byte {
	//chars := " `.,|^'\\/~!_-;:)(\"><¬?*+7j1ilJyc&vt0$VruoI=wzCnY32LTxs4Zkm5hg6qfU9paOS#£eX8D%bdRPGFK@AMQNWHEB"
	chars := " `.,|'\\/~!_-;:)(\"><?*+7j1ilJyc&vt0$VruoI=wzCnY32LTxs4Zkm5hg6qfU9paOS#eX8D%bdRPGFK@AMQNWHEB"
	ascii := make([]byte, len(indices)*len(indices)/4)
	count := 0
	start_y := 0
	for y := 0; y < len(indices); y += 2 {
		start_x := 0
		for x := 1; x < len(indices); x += 2 {
			sum, area := 0, 0
			for _, valy := range grayscaled[start_y:indices[y]] {
				for _, valx := range valy[start_x:indices[x]] {
					sum += valx
					area++
				}
			}
			start_x = indices[x]
			index := len(chars) - (sum/area)*len(chars)/255
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

func GrayScale(image [][]color.Color) [][]int {
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

func Pixels(filename string) [][]color.Color {
	// Read image from "filename".
	file, err := os.Open(filename)
	if err != nil {
		fmt.Printf("%s not found.", filename)
		return nil
	}
	defer file.Close()

	im, err := jpeg.Decode(file)
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

func BuildTree(w int, h int) map[int][]int {
	tree_w := map[int][]int{0: {w}}
	tree_h := map[int][]int{0: {h}}

	Build(tree_w, w, 1)
	Build(tree_h, h, 1)

	ProcessTree(tree_w)
	ProcessTree(tree_h)

	size := len(tree_w)
	if size > len(tree_h) {
		size = len(tree_h)
	}

	tree := make(map[int][]int)

	for i := 0; i < size; i++ {
		tree[i] = make([]int, len(tree_h[i])*2)
		new_idx := 0
		for index := 0; index < len(tree_h[i]); index++ {
			tree[i][new_idx] = tree_w[i][index]
			new_idx++
			tree[i][new_idx] = tree_h[i][index]
			new_idx++
		}
	}
	return tree
}

func ProcessTree(tree map[int][]int) {
	size := len(tree) - 1
	expected := 1 << size
	for i := size; i >= 0; i-- {
		if len(tree[i]) == expected {
			for index := 1; index < len(tree[i]); index++ {
				tree[i][index] += tree[i][index-1]
			}
		} else {
			delete(tree, i)
		}
		expected >>= 1
	}
}

func Build(tree map[int][]int, val int, i int) {
	if _, ok := tree[i]; !ok {
		tree[i] = make([]int, 0)
	}
	half := val / 2
	other_half := half + (val % 2)
	if half+other_half > 2 {
		tree[i] = append(tree[i], half)
		Build(tree, half, i+1)

		tree[i] = append(tree[i], other_half)
		Build(tree, other_half, i+1)
	}
}
