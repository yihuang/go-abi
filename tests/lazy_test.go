//go:build !uint256

package tests

import (
	"testing"

	"github.com/test-go/testify/require"
)

//go:generate go run ../cmd -var LazyTestABI -output lazy.abi.go -lazy

var LazyTestABI = []string{
	"struct Profile { string name; uint64 age; string[] tags }",
	"struct Record { uint32 id; bytes data; bool active }",
	"struct Container { Profile profile; Record[] records }",
	"function getProfile(Profile profile) returns (bool)",
	"function getRecords(Record[] records) returns (uint256)",
	"function getContainer(Container container) returns (bool)",
}

func TestLazyViewBasic(t *testing.T) {
	profile := &Profile{
		Name: "Alice",
		Age:  30,
		Tags: []string{"developer", "blockchain"},
	}

	encoded, err := profile.Encode()
	require.NoError(t, err)

	// Decode as view
	view, bytesRead, err := DecodeProfileView(encoded)
	require.NoError(t, err)
	require.Equal(t, len(encoded), bytesRead)

	// Access individual fields
	name, err := view.Name()
	require.NoError(t, err)
	require.Equal(t, "Alice", name)

	age, err := view.Age()
	require.NoError(t, err)
	require.Equal(t, uint64(30), age)

	// Check raw bytes
	require.Equal(t, encoded, view.Raw())
}

func TestLazyViewSliceAccess(t *testing.T) {
	profile := &Profile{
		Name: "Bob",
		Age:  25,
		Tags: []string{"go", "rust", "solidity"},
	}

	encoded, err := profile.Encode()
	require.NoError(t, err)

	view, _, err := DecodeProfileView(encoded)
	require.NoError(t, err)

	// Tags() returns []string directly (stdlib type)
	tags, err := view.Tags()
	require.NoError(t, err)
	require.Equal(t, 3, len(tags))
	require.Equal(t, "go", tags[0])
	require.Equal(t, "rust", tags[1])
	require.Equal(t, "solidity", tags[2])
}

func TestLazyViewMaterialize(t *testing.T) {
	profile := &Profile{
		Name: "Charlie",
		Age:  35,
		Tags: []string{"ethereum"},
	}

	encoded, err := profile.Encode()
	require.NoError(t, err)

	view, _, err := DecodeProfileView(encoded)
	require.NoError(t, err)

	materialized, err := view.Materialize()
	require.NoError(t, err)
	require.Equal(t, profile.Name, materialized.Name)
	require.Equal(t, profile.Age, materialized.Age)
	require.Equal(t, profile.Tags, materialized.Tags)
}

func TestLazyViewSliceMaterialize(t *testing.T) {
	profile := &Profile{
		Name: "Dave",
		Age:  40,
		Tags: []string{"web3", "defi", "nft"},
	}

	encoded, err := profile.Encode()
	require.NoError(t, err)

	view, _, err := DecodeProfileView(encoded)
	require.NoError(t, err)

	// Tags() returns []string directly (stdlib type, no Materialize needed)
	tags, err := view.Tags()
	require.NoError(t, err)
	require.Equal(t, profile.Tags, tags)
}

func TestLazyViewInvalidData(t *testing.T) {
	// Too short
	_, _, err := DecodeProfileView([]byte{0x00, 0x01})
	require.Error(t, err)

	// Invalid offset - create data with wrong offset
	invalidData := make([]byte, 128)
	// First 32 bytes: offset for name (should be 96 for 3-field tuple)
	// Set it to something wrong
	invalidData[31] = 0x10 // offset = 16, which is wrong
	_, _, err = DecodeProfileView(invalidData)
	require.Error(t, err)
}

func TestLazyViewWithDynamicField(t *testing.T) {
	// Record has dynamic field (bytes data)
	record := &Record{
		Id:     42,
		Data:   []byte{1, 2, 3, 4, 5},
		Active: true,
	}

	encoded, err := record.Encode()
	require.NoError(t, err)

	view, bytesRead, err := DecodeRecordView(encoded)
	require.NoError(t, err)
	require.Equal(t, len(encoded), bytesRead)

	id, err := view.Id()
	require.NoError(t, err)
	require.Equal(t, uint32(42), id)

	active, err := view.Active()
	require.NoError(t, err)
	require.Equal(t, true, active)

	data, err := view.Data()
	require.NoError(t, err)
	require.Equal(t, []byte{1, 2, 3, 4, 5}, data)
}

func TestLazyViewNestedTuple(t *testing.T) {
	container := &Container{
		Profile: Profile{
			Name: "Eve",
			Age:  28,
			Tags: []string{"security", "audit"},
		},
		Records: []Record{
			{Id: 1, Data: []byte{0xaa}, Active: true},
			{Id: 2, Data: []byte{0xbb, 0xcc}, Active: false},
		},
	}

	encoded, err := container.Encode()
	require.NoError(t, err)

	view, bytesRead, err := DecodeContainerView(encoded)
	require.NoError(t, err)
	require.Equal(t, len(encoded), bytesRead)

	// Access nested profile view
	profileView, err := view.Profile()
	require.NoError(t, err)

	name, err := profileView.Name()
	require.NoError(t, err)
	require.Equal(t, "Eve", name)

	age, err := profileView.Age()
	require.NoError(t, err)
	require.Equal(t, uint64(28), age)

	// Access records slice view
	recordsView, err := view.Records()
	require.NoError(t, err)
	require.Equal(t, 2, recordsView.Len())

	record0, err := recordsView.Get(0)
	require.NoError(t, err)

	record0Id, err := record0.Id()
	require.NoError(t, err)
	require.Equal(t, uint32(1), record0Id)

	record0Active, err := record0.Active()
	require.NoError(t, err)
	require.Equal(t, true, record0Active)
}

func TestLazyViewEmptySlice(t *testing.T) {
	profile := &Profile{
		Name: "Frank",
		Age:  50,
		Tags: []string{}, // Empty slice
	}

	encoded, err := profile.Encode()
	require.NoError(t, err)

	view, _, err := DecodeProfileView(encoded)
	require.NoError(t, err)

	// Tags() returns []string directly (stdlib type)
	tags, err := view.Tags()
	require.NoError(t, err)
	require.Equal(t, 0, len(tags))
}

func TestLazyViewRoundTrip(t *testing.T) {
	// Test that we can create a view, materialize it, re-encode, and get the same bytes
	original := &Profile{
		Name: "Grace",
		Age:  33,
		Tags: []string{"layer2", "optimism", "arbitrum"},
	}

	encoded1, err := original.Encode()
	require.NoError(t, err)

	view, _, err := DecodeProfileView(encoded1)
	require.NoError(t, err)

	materialized, err := view.Materialize()
	require.NoError(t, err)

	encoded2, err := materialized.Encode()
	require.NoError(t, err)

	require.Equal(t, encoded1, encoded2)
}

func TestLazyViewMethodCall(t *testing.T) {
	// Test that view works with method call structs
	profile := Profile{
		Name: "Henry",
		Age:  45,
		Tags: []string{"testing"},
	}

	call := &GetProfileCall{
		Profile: profile,
	}

	encoded, err := call.Encode()
	require.NoError(t, err)

	// Decode the call struct as a view
	view, _, err := DecodeGetProfileCallView(encoded)
	require.NoError(t, err)

	// Access the nested profile
	profileView, err := view.Profile()
	require.NoError(t, err)

	name, err := profileView.Name()
	require.NoError(t, err)
	require.Equal(t, "Henry", name)
}
