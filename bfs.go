package main

func Bfs[E comparable] (g *Graph[E], start, end string) ([]*Node[E], *Graph[E]) {
	queue := make([]*Node[E], 0)
	visited := make(map[string]int, 0)
	subg := CreateDag[E]()
	subg.AddNode(NewNode[E](start))

	startNode := g.Node(start)
	if startNode == nil {
		return make([]*Node[E], 0), subg
	}

	endNode := g.Node(end)
	if endNode == nil {
		return make([]*Node[E], 0), subg
	}

	queue = append(queue, g.Node(start))

	for i := 0; len(queue) != 0; i++ {
		head := queue[0] // i.e. peek(queue)
		queue = queue[1:] // i.e. pop(queue)
		visited[head.Id()] = i

		for nodeId := range head.Connected() {
			_, inVisited := visited[nodeId]
			if !inVisited {
				parent, edge, child := g.NodeEdge(head.Id(), nodeId)

				queue = append(queue, child)
				subg.Associate(
					NewNode[E](child.Id()),
					NewNode[E](parent.Id()),
					edge,
				)
			}
		}

		if head.Id() == end {
			break
		}
	}


	r := make([]*Node[E], len(visited))
	for nodeId, order := range visited {
		r[order] = g.Node(nodeId)
	}

	return r, subg
}

