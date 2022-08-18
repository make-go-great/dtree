package dtree

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDecide(t *testing.T) {
	//              salary >= 50000
	//             /              \
	//           yes              no
	//           /                  \
	//  commutation_hour >= 2h  decline
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
	salary.Branches[true] = commutationHour
	salary.Branches[false] = declineOffer
	commutationHour.Branches[true] = declineOffer
	commutationHour.Branches[false] = freeCoffee
	freeCoffee.Branches[true] = acceptOffer
	freeCoffee.Branches[false] = declineOffer

	tree := &Tree{Root: salary}
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
