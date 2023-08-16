## Solution to Travelling Salesman problem using Ant colony algorithm

### Idea
1. Spawn N ants (workers, goroutines)
2. Wait for an ant to finish walking though the graph
3. Compare ant path with the best known path
4. - if the path is better, increase priorities of each node connection so they are more likely to be picked again
   - if the path is the same, increase the priorities again (but perhaps by a smaller amount) to compensate decaying
   - if the path is worse, decay priorities of each node connection so they are less likely to be picked again
5. Each ant can walk up to K times before the program exits

### Stack
- Go 1.20/1.21

### Notes for future improvements
1. Perhaps it could be smarter to run ants in groups and change the priorities once the group finishes walking, right now the first ants can influence the ants that are currently walking (it might be more realistic, but not useful for the problem)
2. Normalise points distance. The algorithm parameters heavily depend on the distance and units used (as well as the number of ants and number of nodes). By normalising distance to `(0.0, 1.0)` (exclusive), the algorithm could require less configuration
3. Perhaps instead of decaying the nodes in the "bad" route, all nodes could decay some every interval
4. The chance to pick a node depends on the weight of the connection between the current and potential node. The weight depends on the distance and connection priority ("pheromone" level). Priority is hard to configurate, it can be affected by
    - priority increase amount
    - priority refresh amount
    - priority decay amount
    - max priority value
    - (less so, but still) min priority value
    - priority factor (high factor means priority strongly affects weight)
5. Some visualisation would be nice:
    - current best path
    - nodes priorities
6. Easier configuration
    - env vars/program arguments
    - read graph from a file

