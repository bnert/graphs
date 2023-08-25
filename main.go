package main

type Edge[V comparable] struct {
	Value V
}

func EmptyEdges[T comparable]() map[string]*Edge[T] {
	return make(map[string]*Edge[T], 0)
}

type Node[E comparable] struct {
	id    string
	connected map[string]*Edge[E]
}

func (n *Node[E]) Id() string {
	return n.id
}

func(n *Node[E]) AssociateFx(other *Node[E], edge *Edge[E]) {
	if other == nil {
		return
	}

	e := edge
	if e == nil {
		e = new(Edge[E])
	}
	n.connected[other.Id()] = edge
}

func(n *Node[E]) Connected() map[string]*Edge[E] {
	return n.connected
}

func NewNode[E comparable](id string) *Node[E] {
	return &Node[E]{id, EmptyEdges[E]()}
}


type Graph[E comparable] struct {
	nodes map[string]*Node[E]
}

func (g *Graph[E]) Nodes() map[string]*Node[E] {
	return g.nodes
}

func (g *Graph[E]) Associate(nodeLeft, nodeRight *Node[E], edge *Edge[E]) {
	if nodeLeft == nil && nodeRight == nil {
		return
	}

	if nodeLeft == nil {
		g.nodes[nodeRight.Id()] = nodeRight
	}

	if nodeRight == nil {
		g.nodes[nodeLeft.Id()] = nodeLeft
	}

	if n, ok := g.nodes[nodeLeft.Id()]; ok {
		n.AssociateFx(nodeRight, edge)
	} else {
		nodeLeft.AssociateFx(nodeRight, edge)
		g.nodes[nodeLeft.Id()] = nodeLeft
	}

	if _, ok := g.nodes[nodeRight.Id()]; !ok {
		g.nodes[nodeRight.Id()] = nodeRight
	}
}

func (g *Graph[E]) Node(id string) *Node[E] {
	n, ok := g.nodes[id]
	if !ok {
		return nil
	}
	return n
}

func (g *Graph[E]) AddNode(n *Node[E]) {
	g.nodes[n.Id()] = n
}

func (g *Graph[E]) NodeEdge(left, right string) (*Node[E], *Edge[E], *Node[E]) {
	l := g.Node(left)
	if l == nil {
		return nil, nil, nil
	}

	edge, ok := l.Connected()[right]
	if !ok {
		return l, nil, nil
	}

	return l, edge, g.Node(right)
}

func (g *Graph[E]) Path(fromNode, toNode string) []string {
	// I.e. find **a** subgraph, which represents "parent"/"child"
	// relationships for nodes searched (which is the preludeGraph)
	//
	// the prelude var holds the nodes which have been searched in the bfs
	// run
	prelude, preludeGraph := Bfs(g, fromNode, toNode)

	// Some omtimization around some cases
	switch len(prelude) {
	case 0:
		return []string{}
	case 1:
		return []string{prelude[0].Id()}
	case 2:
		return []string{prelude[0].Id(), prelude[1].Id()}
	default:
		break;
	}

	path, _ := Bfs(preludeGraph, toNode, fromNode)
	p := make([]string, len(path))
	for i := 0; i < len(path); i++ {
		// need to "reverse"
		p[len(path) - i - 1] = path[i].Id()
	}

	return p
}

func (g *Graph[E]) Weights(fromNode, toNode string) []E {
	window := 2
	path := g.Path(fromNode, toNode)
	if len(path) == 0 {
		return make([]E, 0)
	}

	result := make([]E, len(path) - 1 )
	for i := 0; i <= len(path) - window; i++ {
		w := path[i:i + window]
		edge, ok := g.Node(w[0]).Connected()[w[1]]
		if ok {
			result[i] = edge.Value
		}
	}

	return result
}

type Connection[E comparable] struct {
	Left  *Node[E]
	Edge  *Edge[E]
	Right *Node[E]
}

func CreateDag[E comparable](connections ...Connection[E]) *Graph[E] {
	g := &Graph[E]{
		nodes: make(map[string]*Node[E]),
	}

	for _, c := range connections {
		g.Associate(c.Left, c.Right, c.Edge)
	}

	return g
}

func CreateDg[E comparable](inverter func(E) E, connections ...Connection[E]) *Graph[E] {
	g := &Graph[E]{
		nodes: make(map[string]*Node[E]),
	}

	for _, c := range connections {
		g.Associate(c.Left, c.Right, c.Edge)
		g.Associate(c.Right, c.Left, &Edge[E]{inverter(c.Edge.Value)})
	}

	return g
}

