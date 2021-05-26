package main

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
