package main

import (
	"context"
	"fmt"
	"math/rand/v2"
	"sort"
	"sync"
)

type Population struct {
	mutex  sync.Mutex
	graph  Graph
	chains map[int]Chain
}

type EvolutionConfig struct {
	survivorRate int
	mutationRate float64
}

type ChainWeight struct {
	index  int
	weight int
}

func NewPopulation(graph Graph, countOfChains int) Population {
	var population Population
	population.graph = graph

	population.chains = make(map[int]Chain)
	for index := range countOfChains {
		population.chains[index] = NewChain(len(graph.graph))
	}
	return population
}

func (population *Population) evolution(ctx context.Context, config *EvolutionConfig) {
	i := 1
	for {
		ch := make(chan struct{})
		go func() {
			defer close(ch)

			mutateCh := population.mutate(ctx, config.mutationRate)
			survivorsCh := population.filterChainsByRate(population.calculateChainWeights(ctx), config.survivorRate)
			survivors := <-survivorsCh
			crossbreedingCh := population.crossbreed(ctx, survivors)

			newPopulationChains := make([]Chain, 0)
			newPopulationChains = append(newPopulationChains, <-mutateCh...)
			newPopulationChains = append(newPopulationChains, survivors...)
			newPopulationChains = append(newPopulationChains, <-crossbreedingCh...)

			newPopulationChainsMap := make(map[int]Chain)
			for i := range len(newPopulationChains) {
				newPopulationChainsMap[i] = newPopulationChains[i]
			}
			population.mutex.Lock()
			population.chains = newPopulationChainsMap
			population.mutex.Unlock()
		}()

		select {
		case <-ch:
			best := population.sortWeights(<-population.calculateChainWeights(ctx))[0]
			fmt.Println("iteration: ", i, " Best distance - ", best.weight)
		case <-ctx.Done():
			doneCtx := context.Background()
			best := population.sortWeights(<-population.calculateChainWeights(doneCtx))[0]
			fmt.Println("Finally: Best distance - ", best.weight, "; Chain -", population.chains[best.index])
			return
		}
		i++
	}
}

func (population *Population) crossbreed(ctx context.Context, chains []Chain) <-chan []Chain {
	outputCh := make(chan []Chain)
	go func() {
		defer close(outputCh)

		ch := make(chan Chain)
		wg := &sync.WaitGroup{}
		wg.Add(len(chains))

		resultCrossbreeding := make([]Chain, 0)
		defer func() {
			outputCh <- resultCrossbreeding
		}()
		for _, chain := range chains {
			chain := chain
			go func() {
				defer wg.Done()
				population.mutex.Lock()
				index := rand.IntN(len(population.chains) - 1)
				parent := population.chains[index]
				population.mutex.Unlock()
				crossbreedRangeMin, crossbreedRangeMax := randomRange(len(chain))
				select {
				case ch <- chain.crossbreeding(parent, crossbreedRangeMin, crossbreedRangeMax):
				case <-ctx.Done():
					return
				}

			}()
		}
		go func() {
			wg.Wait()
			close(ch)
		}()

		select {
		case v, ok := <-ch:
			if !ok {
				return
			}
			resultCrossbreeding = append(resultCrossbreeding, v)
		case <-ctx.Done():
			return
		}
	}()
	return outputCh
}

func (population *Population) filterChainsByRate(inputCh <-chan []ChainWeight, survivorRate int) <-chan []Chain {
	outputCh := make(chan []Chain)
	go func() {
		defer close(outputCh)
		weights := <-inputCh

		survivors := make([]Chain, 0)
		if len(weights) == 0 {
			outputCh <- survivors
			return
		}

		weights = population.sortWeights(weights)

		for i := range survivorRate {
			survivors = append(survivors, population.chains[weights[i].index])
		}
		outputCh <- survivors
	}()
	return outputCh
}

func (population *Population) sortWeights(weights []ChainWeight) []ChainWeight {
	sort.Slice(weights, func(i, j int) bool {
		return weights[i].weight < weights[j].weight
	})
	return weights
}

func (population *Population) calculateChainWeights(ctx context.Context) <-chan []ChainWeight {
	outputCh := make(chan []ChainWeight)
	go func() {
		defer close(outputCh)

		ch := make(chan ChainWeight)
		go population.countChainsWeight(ctx, ch, population.chains)
		weights := make([]ChainWeight, 0)

		defer func() {
			outputCh <- weights
		}()

		for {
			select {
			case w, ok := <-ch:
				if !ok {
					return
				}
				weights = append(weights, w)
			case <-ctx.Done():
				return
			}
		}
	}()

	return outputCh
}

func (population *Population) countChainsWeight(ctx context.Context, ch chan<- ChainWeight, chains map[int]Chain) {
	wg := &sync.WaitGroup{}

	for index, chain := range chains {
		chain := chain
		index := index

		wg.Add(1)
		go func() {
			defer wg.Done()
			select {
			case ch <- ChainWeight{
				index:  index,
				weight: chain.chainWeight(population.graph),
			}:
			case <-ctx.Done():
				return
			}

		}()
	}
	go func() {
		wg.Wait()
		close(ch)
	}()
}

func (population *Population) mutate(ctx context.Context, mutationRate float64) <-chan []Chain {
	outputCh := make(chan []Chain)

	go func() {
		mutated := make([]Chain, 0)

		defer close(outputCh)
		defer func() {
			outputCh <- mutated
		}()

		ch := make(chan Chain)
		wg := &sync.WaitGroup{}
		wg.Add(len(population.chains))

		for _, chain := range population.chains {
			chain := chain
			go func() {
				defer wg.Done()

				c, ok := chain.mutateByRate(mutationRate)
				if !ok {
					return
				}
				select {
				case ch <- c:
				case <-ctx.Done():
					return
				}
			}()
		}

		go func() {
			wg.Wait()
			close(ch)
		}()

		for {
			select {
			case v, ok := <-ch:
				if !ok {
					return
				}
				mutated = append(mutated, v)
			case <-ctx.Done():
				return
			}
		}
	}()

	return outputCh
}

func randomRange(max int) (int, int) {
	resultMin := rand.IntN(max - 2)
	resultMax := rand.IntN(max-(resultMin+2)) + resultMin + 2

	return resultMin, resultMax
}
