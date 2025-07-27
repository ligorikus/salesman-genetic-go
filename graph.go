package main

import "math/rand"

type Graph struct {
	graph map[int]map[int]int
}

func NewGraph(countOfCities int, maxDistance int) Graph {
	var graph Graph
	graph.graph = make(map[int]map[int]int)
	for i := 0; i < countOfCities; i++ {
		graph.graph[i] = make(map[int]int)
		for j := 0; j < countOfCities; j++ {
			if i == j {
				graph.graph[i][j] = 0
				continue
			}

			if i > j {
				graph.graph[i][j] = graph.graph[j][i]
			} else {
				distance := rand.Intn(maxDistance) + 1
				graph.graph[i][j] = distance
			}

		}
	}

	return graph
}
