package foobar

// Tree -
type Tree struct {
	level int
	data  *TreeItem
}

// TreeItem -
type TreeItem struct {
	val       int
	leftNode  *TreeItem
	rightNode *TreeItem
}

var tree = &Tree{
	level: 3,
	data: &TreeItem{
		val: 21,
	},
}

// SPrintTree -
func SPrintTree(tree *Tree) {
	var treeNodes = make([][]int, tree.level)
	treeIter(treeNodes, tree.data, 0, "right")
}

func treeIter(treeNodes [][]int, t *TreeItem, l int, dir string) {
	if t.leftNode == nil && t.rightNode == nil {
		return
	}

	if dir == "left" {
		treeNodes[l] = append([]int{t.val}, treeNodes[l]...)
		treeIter(treeNodes, t.leftNode, l+1, "right")
		treeIter(treeNodes, t.rightNode, l+1, "right")
	} else {
		treeNodes[l] = append(treeNodes[l], t.val)
		treeIter(treeNodes, t.rightNode, l+1, "left")
		treeIter(treeNodes, t.leftNode, l+1, "left")
	}
}
