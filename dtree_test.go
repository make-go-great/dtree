package dtree

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTree(t *testing.T) {
	// 			  salary >= 50000
	// 			   /             \
	// 			 yes              no
	// 		     /                  \
	//   commutation_hour >= 2h   decline
	// 		  /              \
	// 		 no              yes
	// 		/                  \
	// free_coffee == true    decline
	// 	   /    \
	// 	 yes    no
	// 	 /        \
	// accept   decline

	// outcome nodes
	declineOffer, err := NewOutcome("decline")
	if err != nil {
		panic(err)
	}
	acceptOffer, err := NewOutcome("accept")
	if err != nil {
		panic(err)
	}

	// condition nodes
	salary, err := NewCondition("salary >= 50000")
	if err != nil {
		panic(err)
	}
	commutationHour, err := NewCondition("commutation_hour >= 2")
	if err != nil {
		panic(err)
	}
	freeCoffee, err := NewCondition("free_coffee == true")
	if err != nil {
		panic(err)
	}

	// branches
	salary.Branches[true] = NewConditionNode(commutationHour)
	salary.Branches[false] = NewOutcomeNode(declineOffer)
	commutationHour.Branches[true] = NewOutcomeNode(declineOffer)
	commutationHour.Branches[false] = NewConditionNode(freeCoffee)
	freeCoffee.Branches[true] = NewOutcomeNode(acceptOffer)
	freeCoffee.Branches[false] = NewOutcomeNode(declineOffer)

	tree := &Tree{Root: NewConditionNode(salary)}
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
		outcome, err := tree.Decide(tc.in)
		require.NoError(t, err)
		require.Equal(t, tc.out, outcome)
	}
}
