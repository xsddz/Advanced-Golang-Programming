package foobar

//   2          2
// 1  3       3   1

type RTree struct {
	val       int
	leftNode  *RTree
	rightNode *RTree
}

func Convert(t *RTree) {
	if t == nil {
		return
	}

	t.leftNode, t.rightNode = t.rightNode, t.leftNode
	Convert(t.leftNode)
	Convert(t.rightNode)
}
