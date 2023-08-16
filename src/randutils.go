/*
This is inspired and based on github.com/jmcvetta/randutil
which itself is based on http://eli.thegreenplace.net/2010/01/22/weighted-random-generation-in-python/
This is updated to be a bit more convenient and uses generics
*/
package main

import (
	"math/rand"
)

// A Choice contains a generic item and a weight controlling the frequency with
// which it will be selected.
type Choice[T any] struct {
	Weight float64
	Item   T
}

/*
WeightedChoice used weighted random selection to return one of the supplied
choices.  Weights of 0 are never selected.  All other weight values are
relative.  E.g. if you have two choices both weighted 3, they will be
returned equally often; and each will be returned 3 times as often as a
choice weighted 1.
*/
func WeightedChoice[T any](choices ...Choice[T]) Choice[T] {
	sum := 0.0
	for _, c := range choices {
		w := c.Weight
		if w < 0.0 {
			panic("Got negative weight for a choice")
		}
		sum += w
	}
	target_weight := rand.Float64() * sum
	for _, c := range choices {
		target_weight -= c.Weight
		if target_weight <= 0.0 {
			return c
		}
	}
	panic("Unreachable code")
}

/*
Classic random choice from a list
NOTE: Assumes the slice is not empty
*/
func RandomChoice[T any](s []T) T {
	if len(s) == 0 {
		panic("Empty slice in RandomChoice")
	}
	idx := rand.Intn(len(s))
	return s[idx]
}
