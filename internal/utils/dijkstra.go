package utils

import (
	"container/heap"
	"math"
	"ngevent/internal/model"
)

type Item struct {
	Node     string
	Priority float64
}

type PriorityQueue []*Item

func (pq PriorityQueue) Len() int {
	return len(pq)
}

func (pq PriorityQueue) Less(i, j int) bool {
	return pq[i].Priority < pq[j].Priority
}

func (pq PriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
}

func (pq *PriorityQueue) Push(x interface{}) {
	*pq = append(*pq, x.(*Item))
}

func (pq *PriorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	item := old[n-1]
	*pq = old[:n-1]

	return item
}

func Dijkstra(graph model.Graph, start string) (map[string]float64, map[string]string) {
	dist := make(map[string]float64)
	prev := make(map[string]string)

	for node := range graph {
		dist[node] = math.Inf(1)
	}
	dist[start] = 0

	pq := &PriorityQueue{}
	heap.Init(pq)
	heap.Push(pq, &Item{Node: start, Priority: 0})

	for pq.Len() > 0 {
		current := heap.Pop(pq).(*Item)

		for _, edge := range graph[current.Node] {
			newDist := dist[current.Node] + edge.Weight

			if newDist < dist[edge.To] {
				dist[edge.To] = newDist
				prev[edge.To] = current.Node

				heap.Push(pq, &Item{
					Node:     edge.To,
					Priority: newDist,
				})
			}
		}
	}

	return dist, prev
}

func BuildPath(prev map[string]string, target string) []string {
	var path []string

	for target != "" {
		path = append([]string{target}, path...)
		target = prev[target]
	}

	return path
}
