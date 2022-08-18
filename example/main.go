package main

import (
	"fmt"
	"os"

	"github.com/make-go-great/dtree"
)

func main() {
	// The following json represents tree:
	//              salary >= 50000
	//             /              \
	//           yes              no
	//           /                  \
	//  commutation_hour >= 2h    decline
	//          /            \
	//         no             yes
	//        /                \
	//  free_coffee == true  decline
	//     /          \
	//   yes           no
	//   /              \
	// accept         decline

	bytes, err := os.ReadFile("example/example_tree.json")
	if err != nil {
		panic(err)
	}

	tree, err := dtree.NewTreeFromJson(bytes)
	if err != nil {
		panic(err)
	}

	params := map[string]interface{}{
		"salary":           100000,
		"commutation_hour": 2,
		"free_coffee":      true,
	}
	outcome, err := tree.Decide(params)
	if err != nil {
		panic(err)
	}
	fmt.Println(outcome == "decline")
}
