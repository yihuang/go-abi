package abi

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParseHumanReadableABI(t *testing.T) {
	tests := []struct {
		name     string
		input    []string
		expected string
		hasError bool
	}{
		{
			name:  "simple function",
			input: []string{"function transfer(address to, uint256 amount)"},
			expected: `[
				{
					"type": "function",
					"name": "transfer",
					"inputs": [
						{"name": "to", "type": "address"},
						{"name": "amount", "type": "uint256"}
					],
					"outputs": [],
					"stateMutability": "nonpayable"
				}
			]`,
		},
		{
			name:  "function with view and returns",
			input: []string{"function balanceOf(address account) view returns (uint256)"},
			expected: `[
				{
					"type": "function",
					"name": "balanceOf",
					"inputs": [
						{"name": "account", "type": "address"}
					],
					"outputs": [
						{"name": "", "type": "uint256"}
					],
					"stateMutability": "view"
				}
			]`,
		},
		{
			name:  "function with payable",
			input: []string{"function deposit() payable"},
			expected: `[
				{
					"type": "function",
					"name": "deposit",
					"inputs": [],
					"outputs": [],
					"stateMutability": "payable"
				}
			]`,
		},
		{
			name:  "event with indexed parameters",
			input: []string{"event Transfer(address indexed from, address indexed to, uint256 value)"},
			expected: `[
				{
					"type": "event",
					"name": "Transfer",
					"inputs": [
						{"name": "from", "type": "address", "indexed": true},
						{"name": "to", "type": "address", "indexed": true},
						{"name": "value", "type": "uint256", "indexed": false}
					],
					"anonymous": false
				}
			]`,
		},
		{
			name:  "constructor",
			input: []string{"constructor(address owner, uint256 initialSupply)"},
			expected: `[
				{
					"type": "constructor",
					"inputs": [
						{"name": "owner", "type": "address"},
						{"name": "initialSupply", "type": "uint256"}
					],
					"stateMutability": "nonpayable"
				}
			]`,
		},
		{
			name:  "constructor payable",
			input: []string{"constructor(address owner) payable"},
			expected: `[
				{
					"type": "constructor",
					"inputs": [
						{"name": "owner", "type": "address"}
					],
					"stateMutability": "payable"
				}
			]`,
		},
		{
			name:  "fallback function",
			input: []string{"fallback()"},
			expected: `[
				{
					"type": "fallback",
					"stateMutability": "nonpayable"
				}
			]`,
		},
		{
			name:  "receive function",
			input: []string{"receive() payable"},
			expected: `[
				{
					"type": "receive",
					"stateMutability": "payable"
				}
			]`,
		},
		{
			name: "multiple functions",
			input: []string{
				"function transfer(address to, uint256 amount)",
				"function balanceOf(address account) view returns (uint256)",
			},
			expected: `[
				{
					"type": "function",
					"name": "transfer",
					"inputs": [
						{"name": "to", "type": "address"},
						{"name": "amount", "type": "uint256"}
					],
					"outputs": [],
					"stateMutability": "nonpayable"
				},
				{
					"type": "function",
					"name": "balanceOf",
					"inputs": [
						{"name": "account", "type": "address"}
					],
					"outputs": [
						{"name": "", "type": "uint256"}
					],
					"stateMutability": "view"
				}
			]`,
		},
		{
			name:  "function with arrays",
			input: []string{"function batchTransfer(address[] recipients, uint256[] amounts)"},
			expected: `[
				{
					"type": "function",
					"name": "batchTransfer",
					"inputs": [
						{"name": "recipients", "type": "address[]"},
						{"name": "amounts", "type": "uint256[]"}
					],
					"outputs": [],
					"stateMutability": "nonpayable"
				}
			]`,
		},
		{
			name:  "function with fixed arrays",
			input: []string{"function getBalances(address[10] accounts) view returns (uint256[10])"},
			expected: `[
				{
					"type": "function",
					"name": "getBalances",
					"inputs": [
						{"name": "accounts", "type": "address[10]"}
					],
					"outputs": [
						{"name": "", "type": "uint256[10]"}
					],
					"stateMutability": "view"
				}
			]`,
		},
		{
			name:  "function with bytes types",
			input: []string{"function setData(bytes32 key, bytes value)"},
			expected: `[
				{
					"type": "function",
					"name": "setData",
					"inputs": [
						{"name": "key", "type": "bytes32"},
						{"name": "value", "type": "bytes"}
					],
					"outputs": [],
					"stateMutability": "nonpayable"
				}
			]`,
		},
		{
			name:  "function with small integers",
			input: []string{"function smallIntegers(uint8 u8, uint16 u16, uint32 u32, uint64 u64, int8 i8, int16 i16, int32 i32, int64 i64)"},
			expected: `[
				{
					"type": "function",
					"name": "smallIntegers",
					"inputs": [
						{"name": "u8", "type": "uint8"},
						{"name": "u16", "type": "uint16"},
						{"name": "u32", "type": "uint32"},
						{"name": "u64", "type": "uint64"},
						{"name": "i8", "type": "int8"},
						{"name": "i16", "type": "int16"},
						{"name": "i32", "type": "int32"},
						{"name": "i64", "type": "int64"}
					],
					"outputs": [],
					"stateMutability": "nonpayable"
				}
			]`,
		},
		{
			name: "comments and empty lines",
			input: []string{
				"// This is a comment",
				"",
				"function transfer(address to, uint256 amount)",
				"",
				"// Another comment",
				"function balanceOf(address account) view returns (uint256)",
			},
			expected: `[
				{
					"type": "function",
					"name": "transfer",
					"inputs": [
						{"name": "to", "type": "address"},
						{"name": "amount", "type": "uint256"}
					],
					"outputs": [],
					"stateMutability": "nonpayable"
				},
				{
					"type": "function",
					"name": "balanceOf",
					"inputs": [
						{"name": "account", "type": "address"}
					],
					"outputs": [
						{"name": "", "type": "uint256"}
					],
					"stateMutability": "view"
				}
			]`,
		},
		{
			name: "function with struct parameter",
			input: []string{
				"struct User { string name; uint256 balance; }",
				"function updateUser(User user)",
			},
			expected: `[
				{
					"type": "function",
					"name": "updateUser",
					"inputs": [
						{
							"name": "user",
							"type": "tuple",
							"internalType": "struct User",
							"components": [
								{"name": "name", "type": "string"},
								{"name": "balance", "type": "uint256"}
							]
						}
					],
					"outputs": [],
					"stateMutability": "nonpayable"
				}
			]`,
		},
		{
			name: "function with nested structs",
			input: []string{
				"struct Address { string street; string city; }",
				"struct User { string name; Address addr; uint256 balance; }",
				"function createUser(User user)",
			},
			expected: `[
				{
					"type": "function",
					"name": "createUser",
					"inputs": [
						{
							"name": "user",
							"type": "tuple",
							"internalType": "struct User",
							"components": [
								{"name": "name", "type": "string"},
								{
									"name": "addr",
									"type": "tuple",
									"internalType": "struct Address",
									"components": [
										{"name": "street", "type": "string"},
										{"name": "city", "type": "string"}
									]
								},
								{"name": "balance", "type": "uint256"}
							]
						}
					],
					"outputs": [],
					"stateMutability": "nonpayable"
				}
			]`,
		},
		{
			name: "event with struct parameter",
			input: []string{
				"struct TransferData { address from; address to; uint256 amount; }",
				"event Transfer(TransferData data)",
			},
			expected: `[
				{
					"type": "event",
					"name": "Transfer",
					"inputs": [
						{
							"name": "data",
							"type": "tuple",
							"internalType": "struct TransferData",
							"components": [
								{"name": "from", "type": "address"},
								{"name": "to", "type": "address"},
								{"name": "amount", "type": "uint256"}
							],
							"indexed": false
						}
					],
					"anonymous": false
				}
			]`,
		},
		{
			name: "function with struct array",
			input: []string{
				"struct User { string name; uint256 balance; }",
				"function batchUpdate(User[] users)",
			},
			expected: `[
				{
					"type": "function",
					"name": "batchUpdate",
					"inputs": [
						{
							"name": "users",
							"type": "tuple[]",
							"internalType": "struct User[]",
							"components": [
								{"name": "name", "type": "string"},
								{"name": "balance", "type": "uint256"}
							]
						}
					],
					"outputs": [],
					"stateMutability": "nonpayable"
				}
			]`,
		},
		{
			name: "function with nested dynamic arrays",
			input: []string{
				"function processNestedArrays(uint256[][] matrix, address[][2][] deepArray)",
			},
			expected: `[
				{
					"type": "function",
					"name": "processNestedArrays",
					"inputs": [
						{"name": "matrix", "type": "uint256[][]"},
						{"name": "deepArray", "type": "address[][2][]"}
					],
					"outputs": [],
					"stateMutability": "nonpayable"
				}
			]`,
		},
		{
			name: "function with mixed fixed and dynamic arrays",
			input: []string{
				"function processMixedArrays(uint256[5] fixedArray, address[] dynamicArray, bytes32[3][] fixedDynamicArray)",
			},
			expected: `[
				{
					"type": "function",
					"name": "processMixedArrays",
					"inputs": [
						{"name": "fixedArray", "type": "uint256[5]"},
						{"name": "dynamicArray", "type": "address[]"},
						{"name": "fixedDynamicArray", "type": "bytes32[3][]"}
					],
					"outputs": [],
					"stateMutability": "nonpayable"
				}
			]`,
		},
		{
			name: "function with struct arrays and nested arrays",
			input: []string{
				"struct DataPoint { uint256 value; string label; }",
				"function processData(DataPoint[][] dataMatrix, DataPoint[5][] fixedDataArray)",
			},
			expected: `[
				{
					"type": "function",
					"name": "processData",
					"inputs": [
						{
							"name": "dataMatrix",
							"type": "tuple[][]",
							"internalType": "struct DataPoint[][]",
							"components": [
								{"name": "value", "type": "uint256"},
								{"name": "label", "type": "string"}
							]
						},
						{
							"name": "fixedDataArray",
							"type": "tuple[5][]",
							"internalType": "struct DataPoint[5][]",
							"components": [
								{"name": "value", "type": "uint256"},
								{"name": "label", "type": "string"}
							]
						}
					],
					"outputs": [],
					"stateMutability": "nonpayable"
				}
			]`,
		},
		{
			name: "function with deeply nested mixed arrays",
			input: []string{
				"function deepNestedArrays(uint256[][] complexArray, address[][] mixedArray)",
			},
			expected: `[
				{
					"type": "function",
					"name": "deepNestedArrays",
					"inputs": [
						{"name": "complexArray", "type": "uint256[][]"},
						{"name": "mixedArray", "type": "address[][]"}
					],
					"outputs": [],
					"stateMutability": "nonpayable"
				}
			]`,
		},
		{
			name: "int and uint without explicit sizes normalize to 256 bits",
			input: []string{
				"function testIntUint(int value1, uint value2)",
				"function testArrays(int[] values1, uint[10] values2)",
				"function testMixed(int value1, uint value2, int8 value3, uint256 value4)",
			},
			expected: `[
				{
					"type": "function",
					"name": "testIntUint",
					"inputs": [
						{"name": "value1", "type": "int256"},
						{"name": "value2", "type": "uint256"}
					],
					"outputs": [],
					"stateMutability": "nonpayable"
				},
				{
					"type": "function",
					"name": "testArrays",
					"inputs": [
						{"name": "values1", "type": "int256[]"},
						{"name": "values2", "type": "uint256[10]"}
					],
					"outputs": [],
					"stateMutability": "nonpayable"
				},
				{
					"type": "function",
					"name": "testMixed",
					"inputs": [
						{"name": "value1", "type": "int256"},
						{"name": "value2", "type": "uint256"},
						{"name": "value3", "type": "int8"},
						{"name": "value4", "type": "uint256"}
					],
					"outputs": [],
					"stateMutability": "nonpayable"
				}
			]`,
		},
		{
			name: "nested tuple in return",
			input: []string{
				"function communityPool() view returns ((string denom, uint256 amount)[] coins)",
			},
			expected: `[
				{
					"type": "function",
					"name": "communityPool",
					"inputs": [],
					"outputs": [
						{
							"name": "coins",
							"type": "tuple[]",
							"components": [
								{"name": "denom", "type": "string"},
								{"name": "amount", "type": "uint256"}
							]
						}
					],
					"stateMutability": "view"
				}
			]`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := ParseHumanReadableABI(tt.input)
			if tt.hasError {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)

			// Parse both expected and actual as JSON for comparison
			var expectedJSON interface{}
			err = json.Unmarshal([]byte(tt.expected), &expectedJSON)
			require.NoError(t, err)

			var actualJSON interface{}
			err = json.Unmarshal(result, &actualJSON)
			require.NoError(t, err)

			require.Equal(t, expectedJSON, actualJSON)
		})
	}
}

func TestParseHumanReadableABI_Errors(t *testing.T) {
	tests := []struct {
		name  string
		input []string
	}{
		{
			name:  "invalid function format",
			input: []string{"function invalid format"},
		},
		{
			name:  "invalid type",
			input: []string{"function test(uint257 invalid) returns (bool)"},
		},
		{
			name:  "invalid array size",
			input: []string{"function test(uint256[invalid] arr) returns (bool)"},
		},
		{
			name:  "unrecognized line",
			input: []string{"invalid line format"},
		},
		{
			name:  "unprocessed parentheses",
			input: []string{"function communityPool() view returns (tuple(string denom, uint256 amount)[] coins)"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := ParseHumanReadableABI(tt.input)
			require.Error(t, err)
		})
	}
}
