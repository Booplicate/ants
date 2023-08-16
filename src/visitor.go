package main

// Struct keeps the nodes that should be visited
type Visitor map[Node]struct{}

// func (self *Visitor) ToSlice() []Node {
//     // NOTE: 0 length so go doesn't fill the underlaying array
//     // with zero nodes
//     rv := make([]Node, 0, len(*self))
//     for k := range *self {
//         rv = append(rv, k)
//     }
//     return rv
// }

func NewVisitor(nodes ...Node) Visitor {
	v := make(Visitor, len(nodes))
	v.AddNode(nodes...)
	return v
}

func (self *Visitor) AddNode(nodes ...Node) {
	for _, n := range nodes {
		(*self)[n] = struct{}{}
	}
}

func (self *Visitor) RemoveNode(nodes ...Node) {
	for _, n := range nodes {
		delete(*self, n)
	}
}
