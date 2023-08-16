package main

import (
	"fmt"
	"math"
)

// Struct represents a 2D point
// TODO: perhaps it worth caching results of the methods in this struct
// Like distance and weight are expensive to calculate
type Node struct {
	X float64
	Y float64
}

// Stringify interface implimentation
func (self Node) String() string {
	return fmt.Sprintf("(%v, %v)", self.X, self.Y)
}

// Returns squared distance from this node to another
func (self Node) GetSqrDistanceTo(other Node) float64 {
	return math.Pow(other.X-self.X, 2) + math.Pow(other.Y-self.Y, 2)
}

// Returns distance from this node to another
func (self Node) GetDistanceTo(other Node) float64 {
	return math.Sqrt(self.GetSqrDistanceTo(other))
}

// // Returns weight for the other node when going from this to it
// // (depends on the distance)
// func (self Node) GetWeightToOther(other Node) float64 {
//     dist := self.GetSqrDistanceTo(other)
//     weight := math.Pow(1.0/dist, DISTANCE_FACTOR)
//     return weight
//     // raw_weight := math.Pow(1.0/dist, DISTANCE_FACTOR) * WEIGHT_MULTI * DISTANCE_FACTOR
//     // // NOTE: ceil ensures that we will get at least 1 and never 0
//     // return int(math.Ceil(raw_weight))
// }

type Path []Node

// Represents a connection between two nodes
type NodeCon struct {
	From Node
	To   Node
}
