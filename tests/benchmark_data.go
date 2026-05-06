package tests

// Benchmark data setup functions - shared across all benchmark files.
// newUint256FromInt64 is defined per build tag in uint256_helper{,_holiman}.go.

func createComplexDynamicTuplesData() TestComplexDynamicTuplesCall {
	users := make([]User2, len(testUserData))
	for i, u := range testUserData {
		users[i] = User2{
			Id: newUint256FromInt64(u.Id),
			Profile: UserProfile{
				Name:   u.Name,
				Emails: u.Emails,
				Metadata: UserMetadata2{
					CreatedAt: newUint256FromInt64(u.CreatedAt),
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
						Value:       newUint256FromInt64(999),
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
