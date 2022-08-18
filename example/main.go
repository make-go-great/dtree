package main

import (
	"fmt"

	"dtree"
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

	jsonBytes := `
{
  "root": {
    "condition": {
      "predicate": "salary >= 50000",
      "branches": [
        {
          "value": true,
          "next_node": {
            "condition": {
              "predicate": "commutation_hour >= 2",
              "branches": [
                {
                  "value": true,
                  "next_node": {
                    "outcome": {
                      "value": "decline"
                    }
                  }
                },
                {
                  "value": false,
                  "next_node": {
                    "condition": {
                      "predicate": "free_coffee == true",
                      "branches": [
                        {
                          "value": true,
                          "next_node": {
                            "outcome": {
                              "value": "accept"
                            }
                          }
                        },
                        {
                          "value": false,
                          "next_node": {
                            "outcome": {
                              "value": "decline"
                            }
                          }
                        }
                      ]
                    }
                  }
                }
              ]
            }
          }
        },
        {
          "value": false,
          "next_node": {
            "outcome": {
              "value": "decline"
            }
          }
        }
      ]
    }
  }
}`

	tree, err := dtree.NewTreeFromJson([]byte(jsonBytes))
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
	fmt.Println(outcome)
}
