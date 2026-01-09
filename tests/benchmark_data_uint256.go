//go:build uint256

package tests

import (
	"github.com/holiman/uint256"
)

func newUint256FromInt64(v int64) *uint256.Int {
	return uint256.NewInt(uint64(v))
}

// Benchmark data setup functions - shared across all benchmark files
func createComplexDynamicTuplesData() TestComplexDynamicTuplesCall {
	users := make([]User2, len(testUserData))
	for i, u := range testUserData {
		users[i] = User2{
			Id: uint256.NewInt(uint64(u.Id)),
			Profile: UserProfile{
				Name:   u.Name,
				Emails: u.Emails,
				Metadata: UserMetadata2{
					CreatedAt: uint256.NewInt(uint64(u.CreatedAt)),
					Tags:      u.Tags,
				},
			},
		}
	}
	return TestComplexDynamicTuplesCall{Users: users}
}

func createNestedDynamicArraysData() TestNestedDynamicArraysCall {
	return TestNestedDynamicArraysCall{
		Matrix:        createTestMatrix(newUint256FromInt64),
		AddressMatrix: testAddressMatrix,
	}
}

func createDeeplyNestedData() TestDeeplyNestedCall {
	return TestDeeplyNestedCall{
		Data: Level1{
			Level1: Level2{
				Level2: Level3{
					Level3: Level4{
						Value:       uint256.NewInt(999),
						Description: "Deeply nested value",
					},
				},
			},
		},
	}
}

func createFixedArraysData() TestFixedArraysCall {
	return TestFixedArraysCall{
		Addresses: testAddresses5,
		Uints:     createTestUints3(newUint256FromInt64),
		Bytes32s:  testBytes32s2,
	}
}
