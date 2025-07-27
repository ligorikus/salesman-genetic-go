package main

import (
	"math/rand"
)

type Chain []int

func NewChain(size int) Chain {
	chain := make([]int, size)
	for i := range chain {
		chain[i] = i
	}
	rand.Shuffle(len(chain), func(i, j int) { chain[i], chain[j] = chain[j], chain[i] })
	return chain
}

func (chain Chain) chainWeight(graph Graph) int {
	weight := 0
	for i := 0; i < len(chain)-1; i++ {
		weight += graph.graph[chain[i]][chain[i+1]]
	}
	return weight
}

func (chain Chain) mutation(startPoint int, endPoint int) Chain {
	result := make(Chain, len(chain))
	copy(result, chain)
	mutationSlice := result[startPoint:endPoint]
	rand.Shuffle(len(mutationSlice), func(i, j int) {
		mutationSlice[i], mutationSlice[j] = mutationSlice[j], mutationSlice[i]
	})
	return result
}

func (chain Chain) crossbreeding(partner Chain, startPoint int, endPoint int) Chain {
	child := make(Chain, len(chain))
	for i := 0; i < len(child); i++ {
		child[i] = -1
	}

	partnerGenotype := partner[startPoint:endPoint]
	copy(child[startPoint:endPoint], partnerGenotype)

	childMap := make(map[int]int)
	for _, i := range child {
		childMap[i] = i
	}
	childIndex := endPoint
	parentIndex := endPoint
	maxIndex := len(chain)

	incrementIndex := func(currentIndex int, maxIndex int) int {
		currentIndex++
		return currentIndex % maxIndex
	}

	for {
		_, ok := childMap[chain[parentIndex]]
		if !ok {
			child[childIndex] = chain[parentIndex]
			childMap[chain[parentIndex]] = chain[parentIndex]
			childIndex = incrementIndex(childIndex, maxIndex)
		}
		parentIndex = incrementIndex(parentIndex, maxIndex)
		if parentIndex == endPoint {
			break
		}
	}
	return child
}

func (chain *Chain) mutateByRate(mutationRate float64) (Chain, bool) {
	rateMultiplier := getMultiplierByFloat(mutationRate)
	if tryLuck(mutationRate, rateMultiplier) {
		mutationRangeMin, mutationRangeMax := randomRange(len(*chain))
		return chain.mutation(mutationRangeMin, mutationRangeMax), true
	}
	return nil, false
}
