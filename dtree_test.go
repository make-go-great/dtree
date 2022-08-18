package dtree

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDecide(t *testing.T) {
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

	// outcome nodes
	declineOffer, err := NewOutcome("decline")
	require.NoError(t, err)
	acceptOffer, err := NewOutcome("accept")
	require.NoError(t, err)

	// condition nodes
	salary, err := NewCondition("salary >= 50000")
	require.NoError(t, err)
	commutationHour, err := NewCondition("commutation_hour >= 2")
	require.NoError(t, err)
	freeCoffee, err := NewCondition("free_coffee == true")
	require.NoError(t, err)

	// branches
	salary.AddBranch(true, &Node{Condition: commutationHour})
	salary.AddBranch(false, &Node{Outcome: declineOffer})
	commutationHour.AddBranch(true, &Node{Outcome: declineOffer})
	commutationHour.AddBranch(false, &Node{Condition: freeCoffee})
	freeCoffee.AddBranch(true, &Node{Outcome: acceptOffer})
	freeCoffee.AddBranch(false, &Node{Outcome: declineOffer})

	tree := &Tree{Root: &Node{Condition: salary}}
	for _, tc := range []struct {
		in  map[string]interface{}
		out string
	}{
		{
			in: map[string]interface{}{
				"salary":           49999,
				"commutation_hour": 1,
				"free_coffee":      false,
			},
			out: "decline",
		},
		{
			in: map[string]interface{}{
				"salary":           49999,
				"commutation_hour": 1,
				"free_coffee":      true,
			},
			out: "decline",
		},
		{
			in: map[string]interface{}{
				"salary":           49999,
				"commutation_hour": 2,
				"free_coffee":      false,
			},
			out: "decline",
		},
		{
			in: map[string]interface{}{
				"salary":           49999,
				"commutation_hour": 2,
				"free_coffee":      true,
			},
			out: "decline",
		},

		{
			in: map[string]interface{}{
				"salary":           50000,
				"commutation_hour": 1,
				"free_coffee":      false,
			},
			out: "decline",
		},
		{
			in: map[string]interface{}{
				"salary":           50000,
				"commutation_hour": 1,
				"free_coffee":      true,
			},
			out: "accept",
		},
		{
			in: map[string]interface{}{
				"salary":           50000,
				"commutation_hour": 2,
				"free_coffee":      false,
			},
			out: "decline",
		},
		{
			in: map[string]interface{}{
				"salary":           50000,
				"commutation_hour": 2,
				"free_coffee":      true,
			},
			out: "decline",
		},
		{
			in: map[string]interface{}{
				"salary": 49999,
			},
			out: "decline",
		},
		{
			in: map[string]interface{}{
				"salary":           50000,
				"commutation_hour": 2,
			},
			out: "decline",
		},
	} {
		outcome, terr := tree.Decide(tc.in)
		require.NoError(t, terr)
		require.Equal(t, tc.out, outcome)
	}
}

func TestMarshal(t *testing.T) {
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

	// outcome nodes
	declineOffer, err := NewOutcome("decline")
	require.NoError(t, err)
	acceptOffer, err := NewOutcome("accept")
	require.NoError(t, err)

	// condition nodes
	salary, err := NewCondition("salary >= 50000")
	require.NoError(t, err)
	commutationHour, err := NewCondition("commutation_hour >= 2")
	require.NoError(t, err)
	freeCoffee, err := NewCondition("free_coffee == true")
	require.NoError(t, err)

	// branches
	salary.AddBranch(true, &Node{Condition: commutationHour})
	salary.AddBranch(false, &Node{Outcome: declineOffer})
	commutationHour.AddBranch(true, &Node{Outcome: declineOffer})
	commutationHour.AddBranch(false, &Node{Condition: freeCoffee})
	freeCoffee.AddBranch(true, &Node{Outcome: acceptOffer})
	freeCoffee.AddBranch(false, &Node{Outcome: declineOffer})

	tree := &Tree{Root: &Node{Condition: salary}}
	jsonBytes, err := json.Marshal(tree)
	require.NoError(t, err)

	newTree, err := NewTreeFromJson(jsonBytes)
	require.NoError(t, err)

	salaryNode := newTree.Root
	require.NotNil(t, salaryNode)
	require.NotNil(t, salaryNode.Condition)
	require.Nil(t, salaryNode.Outcome)
	require.Equal(t, "salary >= 50000", salaryNode.Condition.Predicate)
	require.NotNil(t, salaryNode.Condition.evaluablePredicate)
	require.Len(t, salaryNode.Condition.Branches, 2)
	require.Equal(t, true, salaryNode.Condition.Branches[0].Value)
	require.Equal(t, false, salaryNode.Condition.Branches[1].Value)
	require.Len(t, salaryNode.Condition.branchMap, 2)
	require.Contains(t, salaryNode.Condition.branchMap, true)
	require.Contains(t, salaryNode.Condition.branchMap, false)
	require.Equal(t, salaryNode.Condition.Branches[0].NextNode, salaryNode.Condition.branchMap[true])
	require.Equal(t, salaryNode.Condition.Branches[1].NextNode, salaryNode.Condition.branchMap[false])

	salaryFalseNode := newTree.Root.Condition.branchMap[false]
	require.Nil(t, salaryFalseNode.Condition)
	require.NotNil(t, salaryFalseNode.Outcome)
	require.Equal(t, "decline", salaryFalseNode.Outcome.Value)

	commutationHourNode := newTree.Root.Condition.branchMap[true]
	require.Equal(t, "commutation_hour >= 2", commutationHourNode.Condition.Predicate)
	require.NotNil(t, commutationHourNode.Condition.evaluablePredicate)
	require.Len(t, commutationHourNode.Condition.Branches, 2)
	require.Equal(t, true, commutationHourNode.Condition.Branches[0].Value)
	require.Equal(t, false, commutationHourNode.Condition.Branches[1].Value)
	require.Len(t, commutationHourNode.Condition.branchMap, 2)
	require.Contains(t, commutationHourNode.Condition.branchMap, true)
	require.Contains(t, commutationHourNode.Condition.branchMap, false)
	require.Equal(t, commutationHourNode.Condition.Branches[0].NextNode, commutationHourNode.Condition.branchMap[true])
	require.Equal(t, commutationHourNode.Condition.Branches[1].NextNode, commutationHourNode.Condition.branchMap[false])

	commutationTrueNode := commutationHourNode.Condition.branchMap[true]
	require.Nil(t, commutationTrueNode.Condition)
	require.NotNil(t, commutationTrueNode.Outcome)
	require.Equal(t, "decline", commutationTrueNode.Outcome.Value)

	freeCoffeeNode := commutationHourNode.Condition.branchMap[false]
	require.Equal(t, "free_coffee == true", freeCoffeeNode.Condition.Predicate)
	require.NotNil(t, freeCoffeeNode.Condition.evaluablePredicate)
	require.Len(t, freeCoffeeNode.Condition.Branches, 2)
	require.Equal(t, true, freeCoffeeNode.Condition.Branches[0].Value)
	require.Equal(t, false, freeCoffeeNode.Condition.Branches[1].Value)
	require.Len(t, freeCoffeeNode.Condition.branchMap, 2)
	require.Contains(t, freeCoffeeNode.Condition.branchMap, true)
	require.Contains(t, freeCoffeeNode.Condition.branchMap, false)
	require.Equal(t, freeCoffeeNode.Condition.Branches[0].NextNode, freeCoffeeNode.Condition.branchMap[true])
	require.Equal(t, freeCoffeeNode.Condition.Branches[1].NextNode, freeCoffeeNode.Condition.branchMap[false])

	freeCoffeeFalseNode := freeCoffeeNode.Condition.branchMap[false]
	require.Nil(t, freeCoffeeFalseNode.Condition)
	require.NotNil(t, freeCoffeeFalseNode.Outcome)
	require.Equal(t, "decline", freeCoffeeFalseNode.Outcome.Value)

	freeCoffeeTrueNode := freeCoffeeNode.Condition.branchMap[true]
	require.Nil(t, freeCoffeeTrueNode.Condition)
	require.NotNil(t, freeCoffeeTrueNode.Outcome)
	require.Equal(t, "accept", freeCoffeeTrueNode.Outcome.Value)

	for _, tc := range []struct {
		in  map[string]interface{}
		out string
	}{
		{
			in: map[string]interface{}{
				"salary":           49999,
				"commutation_hour": 1,
				"free_coffee":      false,
			},
			out: "decline",
		},
		{
			in: map[string]interface{}{
				"salary":           49999,
				"commutation_hour": 1,
				"free_coffee":      true,
			},
			out: "decline",
		},
		{
			in: map[string]interface{}{
				"salary":           49999,
				"commutation_hour": 2,
				"free_coffee":      false,
			},
			out: "decline",
		},
		{
			in: map[string]interface{}{
				"salary":           49999,
				"commutation_hour": 2,
				"free_coffee":      true,
			},
			out: "decline",
		},

		{
			in: map[string]interface{}{
				"salary":           50000,
				"commutation_hour": 1,
				"free_coffee":      false,
			},
			out: "decline",
		},
		{
			in: map[string]interface{}{
				"salary":           50000,
				"commutation_hour": 1,
				"free_coffee":      true,
			},
			out: "accept",
		},
		{
			in: map[string]interface{}{
				"salary":           50000,
				"commutation_hour": 2,
				"free_coffee":      false,
			},
			out: "decline",
		},
		{
			in: map[string]interface{}{
				"salary":           50000,
				"commutation_hour": 2,
				"free_coffee":      true,
			},
			out: "decline",
		},
		{
			in: map[string]interface{}{
				"salary": 49999,
			},
			out: "decline",
		},
		{
			in: map[string]interface{}{
				"salary":           50000,
				"commutation_hour": 2,
			},
			out: "decline",
		},
	} {
		outcome, terr := newTree.Decide(tc.in)
		require.NoError(t, terr)
		require.Equal(t, tc.out, outcome)
	}
}
