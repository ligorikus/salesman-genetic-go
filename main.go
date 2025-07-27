package main

import (
	"context"
	"time"
)

const MaxDistance = 99
const WorkSeconds = 1
const CountOfCities = 10
const CountOfChains = 10
const Survivors = 5
const MutationRate = 0.5

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), WorkSeconds*time.Second)
	defer cancel()

	graph := NewGraph(CountOfCities, MaxDistance)
	population := NewPopulation(graph, CountOfChains)
	evConfig := &EvolutionConfig{
		survivorRate: Survivors,
		mutationRate: MutationRate,
	}
	population.evolution(ctx, evConfig)
}
