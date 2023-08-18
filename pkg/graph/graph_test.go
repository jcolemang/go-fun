
package graph

import (
    "testing"
    "language/pkg/graph"
)

// TestHelloName calls greetings.Hello with a name, checking
// for a valid return value.
func TestGetMinimalColor(t *testing.T) {
    adjacentColors := []int{1}
    availableColors := []int{0, 1}
    color, _ := graph.GetMinimalColor(adjacentColors, availableColors)
    if color != 0 {
        t.Fatalf("Incorrect!")
    }
}
