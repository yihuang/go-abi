# Events Generation API Documentation

This document describes the events generation API provided by the go-abi generator for Ethereum smart contract events.

## Overview

The go-abi generator automatically generates Go code for Ethereum smart contract events, including:

- **Event structs** - Complete event data structures
- **Event data structs** - Separate structs for non-indexed event data
- **Topic encoding/decoding** - Methods for working with indexed parameters
- **Data encoding/decoding** - Methods for working with non-indexed parameters
- **Event signatures** - Pre-calculated event topic constants

## Generated Code Structure

For each event in your ABI, the generator creates:

### 1. Event Topic Constants

```go
// Event signatures
var (
    // Transfer(address,address,uint256)
    TransferTopic = [32]byte{...}
    // ComplexEvent(string,uint256[],address)
    ComplexEventTopic = [32]byte{...}
    // UserCreated((address,string,uint256),address)
    UserCreatedTopic = [32]byte{...}
)
```

These constants contain the Keccak256 hash of the event signature, used as the first topic in Ethereum logs.

### 2. Main Event Struct

```go
// Transfer represents an ABI event
type Transfer struct {
    From  common.Address
    To    common.Address
    Value *big.Int
}
```

Contains all event parameters (both indexed and non-indexed).

### 3. Event Data Struct (for non-indexed parameters)

```go
// TransferData represents the non-indexed data of Transfer event
type TransferData struct {
    Value *big.Int
}
```

Contains only non-indexed parameters that go into the event data section.

### 4. Topic Encoding/Decoding Methods

```go
// EncodeTopics encodes indexed fields of Transfer event to topics
func (e Transfer) EncodeTopics() ([][32]byte, error)

// DecodeTopics decodes indexed fields of Transfer event from topics
func (e *Transfer) DecodeTopics(topics [][32]byte) error
```

### 5. Data Encoding/Decoding Methods

```go
// EncodedSize returns the total encoded size of TransferData
func (t TransferData) EncodedSize() int

// Encode encodes TransferData to ABI bytes
func (t TransferData) Encode() ([]byte, error)

// EncodeTo encodes TransferData to ABI bytes in the provided buffer
func (t TransferData) EncodeTo(buf []byte) (int, error)

// Decode decodes TransferData from ABI bytes
func (t *TransferData) Decode(data []byte) error
```

## Usage Examples

### Basic Event Usage

```go
import (
    "math/big"
    "github.com/ethereum/go-ethereum/common"
)

// Create a Transfer event
transfer := Transfer{
    From:  common.HexToAddress("0x742d35Cc6634C0532925a3b8Dc9F2a5C3B8Dc9F2"),
    To:    common.HexToAddress("0x742d35Cc6634C0532925a3b8Dc9F2a5C3B8Dc9F3"),
    Value: big.NewInt(1000000000000000000), // 1 ETH
}

// Encode topics for indexed fields
topics, err := transfer.EncodeTopics()
if err != nil {
    // handle error
}

// topics[0] = TransferTopic (event signature)
// topics[1] = encoded From address
// topics[2] = encoded To address
```

### Working with Event Data

```go
// Create event data (non-indexed parameters)
data := TransferData{
    Value: big.NewInt(1000000000000000000),
}

// Encode event data
encodedData, err := data.Encode()
if err != nil {
    // handle error
}

// Decode event data
var decodedData TransferData
if err := decodedData.Decode(encodedData); err != nil {
    // handle error
}
```

### Complete Event Processing

```go
// Simulate receiving an Ethereum log
func processLog(topics [][32]byte, data []byte) error {
    // Check if this is a Transfer event
    if topics[0] != TransferTopic {
        return fmt.Errorf("not a Transfer event")
    }

    // Decode indexed fields from topics
    var transfer Transfer
    if err := transfer.DecodeTopics(topics); err != nil {
        return fmt.Errorf("failed to decode topics: %w", err)
    }

    // Decode non-indexed fields from data
    var transferData TransferData
    if err := transferData.Decode(data); err != nil {
        return fmt.Errorf("failed to decode data: %w", err)
    }

    // Now you have the complete event
    fmt.Printf("Transfer from %s to %s value %s\n",
        transfer.From, transfer.To, transferData.Value)

    return nil
}
```

### Complex Event Example

```go
// Complex event with mixed indexed and non-indexed parameters
complexEvent := ComplexEvent{
    Message: "Test message",
    Numbers: []*big.Int{big.NewInt(100), big.NewInt(200), big.NewInt(300)},
    Sender:  common.HexToAddress("0x742d35Cc6634C0532925a3b8Dc9F2a5C3B8Dc9F4"),
}

// Encode topics (only indexed fields)
topics, err := complexEvent.EncodeTopics()
if err != nil {
    // handle error
}

// topics[0] = ComplexEventTopic
// topics[1] = encoded Sender address

// Encode data (non-indexed fields)
data := ComplexEventData{
    Message: complexEvent.Message,
    Numbers: complexEvent.Numbers,
}

encodedData, err := data.Encode()
if err != nil {
    // handle error
}
```

## API Reference

### Event Struct Methods

#### `EncodeTopics() ([][32]byte, error)`

Encodes all indexed fields of the event into Ethereum log topics.

- **Returns**: Array of 32-byte topics where:
  - `topics[0]` is always the event signature
  - Subsequent topics contain encoded indexed parameters
- **Error**: Returns error if encoding fails

#### `DecodeTopics(topics [][32]byte) error`

Decodes indexed fields from Ethereum log topics into the event struct.

- **Parameters**: `topics` - Array of 32-byte topics from Ethereum log
- **Error**: Returns error if topics array is insufficient or decoding fails

### Event Data Struct Methods

#### `EncodedSize() int`

Calculates the total encoded size in bytes required for the data.

- **Returns**: Total size in bytes including dynamic type overhead

#### `Encode() ([]byte, error)`

Encodes the event data struct to ABI-encoded bytes.

- **Returns**: ABI-encoded byte slice
- **Error**: Returns error if encoding fails

#### `EncodeTo(buf []byte) (int, error)`

Encodes the event data struct to the provided buffer.

- **Parameters**: `buf` - Pre-allocated buffer to write encoded data
- **Returns**: Number of bytes written
- **Error**: Returns error if buffer is too small or encoding fails

#### `Decode(data []byte) error`

Decodes ABI-encoded bytes into the event data struct.

- **Parameters**: `data` - ABI-encoded byte slice
- **Error**: Returns error if data is insufficient or decoding fails

## Indexed vs Non-Indexed Parameters

### Indexed Parameters
- Stored in Ethereum log topics
- Limited to 32 bytes each
- Can be efficiently filtered by Ethereum clients
- Generated methods: `EncodeTopics()`, `DecodeTopics()`

### Non-Indexed Parameters
- Stored in event data section
- Can be of any size (strings, arrays, etc.)
- More expensive to filter
- Generated methods: `Encode()`, `Decode()`, etc.

## Supported Data Types

The events generation supports all standard ABI types:

- **Basic types**: `uint8-256`, `int8-256`, `bool`, `address`
- **Fixed types**: `bytes1-32`, fixed arrays
- **Dynamic types**: `string`, `bytes`, dynamic arrays
- **Complex types**: Tuples, nested structs

## Performance Considerations

### Static Size Constants

Each event data struct includes a static size constant:

```go
const TransferDataStaticSize = 32
```

This allows efficient buffer pre-allocation and size calculations.

### Buffer Reuse

Use `EncodeTo()` when you can pre-allocate buffers to avoid memory allocations:

```go
// Pre-allocate buffer
size := data.EncodedSize()
buf := make([]byte, size)

// Encode directly to buffer
written, err := data.EncodeTo(buf)
```

### Native Go Types

The generator uses native Go types (`uint8`, `uint16`, etc.) for common integer sizes to avoid `big.Int` allocations where possible.

## Error Handling

All encoding/decoding methods return errors for:

- Insufficient data/buffer size
- Invalid parameter values
- Type conversion failures
- Memory allocation failures

Always check returned errors:

```go
topics, err := event.EncodeTopics()
if err != nil {
    return fmt.Errorf("failed to encode topics: %w", err)
}
```

## Best Practices

### 1. Validate Topics Before Decoding

```go
func decodeTransfer(topics [][32]byte, data []byte) (*Transfer, *TransferData, error) {
    if len(topics) < 3 {
        return nil, nil, fmt.Errorf("insufficient topics for Transfer event")
    }

    if topics[0] != TransferTopic {
        return nil, nil, fmt.Errorf("not a Transfer event")
    }

    // ... rest of decoding
}
```

### 2. Use Type-Safe Event Creation

```go
func createTransfer(from, to common.Address, value *big.Int) (*Transfer, error) {
    if value.Sign() < 0 {
        return nil, fmt.Errorf("transfer value cannot be negative")
    }

    return &Transfer{
        From:  from,
        To:    to,
        Value: value,
    }, nil
}
```

### 3. Handle Large Arrays Efficiently

For events with large arrays, consider streaming or chunking:

```go
// For very large arrays, process in chunks
func processLargeEvent(data []byte) error {
    var eventData ComplexEventData
    if err := eventData.Decode(data); err != nil {
        return err
    }

    // Process numbers in chunks to avoid memory issues
    chunkSize := 1000
    for i := 0; i < len(eventData.Numbers); i += chunkSize {
        end := i + chunkSize
        if end > len(eventData.Numbers) {
            end = len(eventData.Numbers)
        }
        processChunk(eventData.Numbers[i:end])
    }

    return nil
}
```

## Testing

The generated code includes comprehensive tests:

- Event topic encoding/decoding round-trip
- Event data encoding/decoding round-trip
- Mixed indexed/non-indexed parameter handling
- Complex nested structure support

Run tests with:

```bash
go test ./tests -v -run "TestEvent"
```

## Compatibility

- **Go version**: 1.18+
- **Ethereum ABI**: Full compatibility with Ethereum ABI specification
- **go-ethereum**: Compatible with `github.com/ethereum/go-ethereum` types
- **Performance**: Optimized for minimal allocations and maximum throughput