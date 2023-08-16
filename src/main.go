package main

import (
	"fmt"
	"math"
)

// Max runs each worker does
const MAX_RUNS int = 100
const MAX_WORKERS int = 1000

// The higher the factor, the more the distance
// affects path choice
const DISTANCE_FACTOR float64 = 2.0

// The higher the factor, the more node priority
// affects path choice
const PRIORITY_FACTOR float64 = 2.0

// Priority increase when found new best route
const PRIORITY_INC float64 = 1.0

// Priority increase when found a route as good as the current best
const PRIORITY_REFRESH float64 = 0.75

// Priority decay when found a route worse than the current best
const PRIORITY_DECAY float64 = -0.001

// Calculates weight based on the distance and priority
func GetWeight(dist float64, priority float64) float64 {
	return math.Pow(1.0/dist, DISTANCE_FACTOR) * math.Pow(priority, PRIORITY_FACTOR)
}

// Weighted choice for node selection based on the distance between nodes
func SelectNextNode(from Node, nodes Visitor, manager *PriorityManager) Node {
	choices := make([]Choice[Node], len(nodes))
	for n := range nodes {
		if n == from {
			panic("The from Node is in the slice of available Nodes")
		}
		dist := from.GetSqrDistanceTo(n)
		priority := manager.GetPriority(from, n)
		weight := GetWeight(dist, priority)
		// fmt.Printf("%v->%v: %v\n", from, p, weight)
		c := Choice[Node]{Weight: weight, Item: n}
		choices = append(choices, c)
	}
	res := WeightedChoice(choices...)
	// fmt.Printf("Picked %v, weight was %v\n", res.Item, res.Weight)
	return res.Item
}

// Calculates path length from the nodes forming the path
func GetPathLength(nodes Path) float64 {
	if len(nodes) < 2 {
		return 0.0
	}

	length := 0.0
	current := nodes[0]
	for _, next := range nodes[1:] {
		length += current.GetDistanceTo(next)
		current = next
	}

	return length
}

func main() {
	base := Node{0.0, 0.0}
	all_nodes := []Node{
		base,
		{2.0, 1.0},
		{3.0, 2.5},
		{2.6, 4.0},
		{5.0, 5.0},
		{-1.0, -4.0},

		// {-7.5, -21.5},
		// {4.0, 4.0},
		// {13.3, 7.5},
		// {0.0, 7.7},
		// {7.7, 7.7},
		// {8.0, 8.0},
	}

	results_ch := make(chan RunResult, MAX_WORKERS)

	total_runs := 0
	total_best_runs := 0
	best_length := math.MaxFloat64
	best_path := make(Path, 0, len(all_nodes))

	manager := NewPriorityManager()

	launch := func() {
		start := RandomChoice[Node](all_nodes)
		visit_next := NewVisitor(all_nodes...)
		visit_next.RemoveNode(start)
		go LaunchWorker(start, visit_next, manager, results_ch)
	}

	for i := 0; i < MAX_WORKERS; i++ {
		launch()
	}

	for res := range results_ch {
		l := GetPathLength(res.Path)
		var prio_change float64

		if l < best_length {
			total_best_runs += 1
			best_length = l
			best_path = res.Path
			prio_change = PRIORITY_INC
		} else if l == best_length {
			prio_change = PRIORITY_REFRESH
		} else {
			prio_change = PRIORITY_DECAY
		}
		manager.UpdatePriorityForPath(res.Path, prio_change)

		total_runs += 1

		if total_runs >= MAX_RUNS*MAX_WORKERS {
			break
		}
		launch()
	}

	fmt.Printf("Total runs: %v\n", total_runs)
	fmt.Printf("Total runs with best paths: %v\n", total_best_runs)
	fmt.Printf("Best path %v, length %v\n", best_path, best_length)
	fmt.Println("Priorities:")
	tmp := make(map[NodeCon]struct{})
	for k := range manager.inner {
		// NOTE: avoid printing duplicates
		if _, ok := tmp[k]; !ok {
			tmp[NodeCon{k.To, k.From}] = struct{}{}
			fmt.Printf("%v->%v: %v\n", k.From, k.To, manager.inner[k])
		}
	}
}

type RunResult struct {
	Path Path
}

func LaunchWorker(
	start Node,
	visit_next Visitor,
	manager *PriorityManager,
	ch chan RunResult,
) {
	path := make(Path, 0, len(visit_next)+2)
	path = append(path, start)
	current := start
	// fmt.Printf("Started at %v\n", current)

	for len(visit_next) > 0 {
		next := SelectNextNode(current, visit_next, manager)
		path = append(path, next)
		visit_next.RemoveNode(next)
		current = next
		// fmt.Printf("Went to %v\n", next)
	}

	path = append(path, start)
	// fmt.Printf("Returned to %v\n", start)
	// path_length := GetPathLength(path)
	// fmt.Printf("Path: %.2f\n", path_length)
	ch <- RunResult{path}
}
