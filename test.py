    

def build(tree, val, i):
    if i not in tree:
        tree[i] = []
    half = val // 2
    other_half = half + (val % 2)
    if half + other_half > 2:
        tree[i].append(half)
        build(tree, half, i + 1)

        tree[i].append(other_half)
        build(tree, other_half, i + 1)


def build_tree(val):
    tree = {0: [val]}

    build(tree, val, 1)

    expected = 1 << (len(tree) - 1)
    for i in range(len(tree) - 1, -1, -1):

        if len(tree[i]) == expected:
            break
        else:
            del tree[i]

        expected >>= 1

    return tree


def main():
    w, h = 1080, 720
    tree_w = build_tree(w)
    tree_h = build_tree(h)

    for i in range(min(len(tree_w), len(tree_h))):
        print(i, len(tree_w[i]), tree_w[i])
        print(i, len(tree_h[i]), tree_h[i])
        print()


if __name__ == '__main__':
    main()