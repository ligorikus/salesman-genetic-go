package main

import (
	"context"
	"fmt"
	"time"
)

const MaxDistance = 99

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	graph := NewGraph(100, MaxDistance)
	population := NewPopulation(graph, 10)
	fmt.Println(population.chains[0])
	evConfig := &EvolutionConfig{
		survivorRate: 5,
		mutationRate: 0.1,
	}
	population.evolution(ctx, evConfig)
}
