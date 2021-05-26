// Alex Eidt
// Processes image dimensions into maps containing information on how to
// divide up the image (by coordinates) by increasing powers of two.

package main

// For a given image with a width "w" and a height "h",
// builds two trees mapping levels of division to indices
// along both dimensions. Both trees values are zipped up
// such that values for the height appear on even indices
// and values for the width appear on odd indices.
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

// Given a tree (as described in the function comment for the "Build"
// function), this function removes values with corrupted lists
// and also updates each value to be a running sum of the dimension.
//
// For a 1080 by 720 image, this is what the first layers of the
// map would look like when we are looking at the width (1080) given
// the tree build using "Build".
//
// 0 -> [1080]
// 1 -> [540, 1080]
// 2 -> [270, 540, 810, 1080]
// 3 -> [135, 270, 405, 540, 675, 810, 945, 1080]
// 4 -> ...
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

// Builds a tree mapping level of division to a list of coordinates
// in the image.
//
// For a 1080 by 720 image, this is what the first layers of the
// map would look like when we are looking at the width (1080)
// as the "val".
//
// 0 -> [1080]
// 1 -> [540, 540]
// 2 -> [270, 270, 270, 270]
// 3 -> [135, 135, 135, 135, 135, 135, 135, 135]
// 4 -> ...
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
