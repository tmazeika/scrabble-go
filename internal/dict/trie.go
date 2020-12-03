package dict

type Node struct {
	edges  map[Letter]*Node
	accept bool
}

func NewNode() *Node {
	return &Node{
		edges: make(map[Letter]*Node),
	}
}

func (n *Node) Edges() map[Letter]*Node {
	if n == nil {
		return nil
	}
	return n.edges
}

func (n *Node) Accept() bool {
	return n != nil && n.accept
}

func (n *Node) Search(word Word) *Node {
	if len(word) == 0 {
		return n
	}
	next, ok := n.edges[word.Head()]
	if !ok {
		return nil
	}
	return next.Search(word.Tail())
}

func (n *Node) Insert(word Word) {
	if len(word) == 0 {
		n.accept = true
	} else if next, ok := n.edges[word.Head()]; ok {
		next.Insert(word.Tail())
	} else {
		next := NewNode()
		n.edges[word.Head()] = next
		next.Insert(word.Tail())
	}
}
