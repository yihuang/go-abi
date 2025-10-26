# Events Generation Quick Reference

## Common Event Patterns

### ERC-20 Transfer Event

```go
// Generated code for Transfer(address,address,uint256)
type Transfer struct {
    From  common.Address
    To    common.Address
    Value *big.Int
}

type TransferData struct {
    Value *big.Int
}

var TransferTopic = [32]byte{...}

// Usage
transfer := Transfer{
    From:  common.HexToAddress("0x123..."),
    To:    common.HexToAddress("0x456..."),
    Value: big.NewInt(1000),
}

// Encode topics (From and To are indexed)
topics, _ := transfer.EncodeTopics()
// topics[0] = TransferTopic
// topics[1] = encoded From
// topics[2] = encoded To

// Encode data (Value is non-indexed)
data := TransferData{Value: transfer.Value}
encodedData, _ := data.Encode()
```

### Complex Event with Mixed Parameters

```go
// Generated for ComplexEvent(string indexed message, uint256[] numbers, address indexed sender)
type ComplexEvent struct {
    Message string
    Numbers []*big.Int
    Sender  common.Address
}

type ComplexEventData struct {
    Message string
    Numbers []*big.Int
}

// Usage
complexEvent := ComplexEvent{
    Message: "Hello World",
    Numbers: []*big.Int{big.NewInt(1), big.NewInt(2)},
    Sender:  common.HexToAddress("0x789..."),
}

// Only Message and Sender are indexed
topics, _ := complexEvent.EncodeTopics()
// topics[0] = ComplexEventTopic
// topics[1] = encoded Message (dynamic type in topic)
// topics[2] = encoded Sender

// Numbers goes in data
data := ComplexEventData{
    Message: complexEvent.Message,
    Numbers: complexEvent.Numbers,
}
encodedData, _ := data.Encode()
```

## Common Operations

### Creating Events from Logs

```go
func processLog(topics [][32]byte, data []byte) {
    switch topics[0] {
    case TransferTopic:
        processTransfer(topics, data)
    case ComplexEventTopic:
        processComplexEvent(topics, data)
    default:
        log.Printf("Unknown event: %x", topics[0])
    }
}

func processTransfer(topics [][32]byte, data []byte) {
    var transfer Transfer
    if err := transfer.DecodeTopics(topics); err != nil {
        log.Printf("Failed to decode transfer topics: %v", err)
        return
    }

    var transferData TransferData
    if err := transferData.Decode(data); err != nil {
        log.Printf("Failed to decode transfer data: %v", err)
        return
    }

    fmt.Printf("Transfer: %s -> %s: %s\n",
        transfer.From, transfer.To, transferData.Value)
}
```

### Event Filtering

```go
// Filter for specific sender
targetSender := common.HexToAddress("0x123...")

func shouldProcessTransfer(topics [][32]byte) bool {
    if len(topics) < 3 || topics[0] != TransferTopic {
        return false
    }

    // Check if From matches target
    var transfer Transfer
    if err := transfer.DecodeTopics(topics); err != nil {
        return false
    }

    return transfer.From == targetSender
}
```

### Batch Event Processing

```go
func processEvents(logs []types.Log) {
    for _, log := range logs {
        switch log.Topics[0] {
        case TransferTopic:
            go processTransferEvent(log)
        case ComplexEventTopic:
            go processComplexEvent(log)
        }
    }
}

func processTransferEvent(log types.Log) {
    var transfer Transfer
    if err := transfer.DecodeTopics(log.Topics); err != nil {
        return
    }

    var data TransferData
    if err := data.Decode(log.Data); err != nil {
        return
    }

    // Process transfer
}
```

## Performance Tips

### Pre-allocate Buffers

```go
// For high-throughput applications
func encodeTransferDataFast(data TransferData) ([]byte, error) {
    size := data.EncodedSize()
    buf := make([]byte, size)
    _, err := data.EncodeTo(buf)
    return buf, err
}
```

### Reuse Event Structs

```go
// Avoid allocations by reusing structs
var transferPool = sync.Pool{
    New: func() interface{} { return &Transfer{} },
}

func getTransfer() *Transfer {
    return transferPool.Get().(*Transfer)
}

func putTransfer(t *Transfer) {
    // Reset fields if needed
    transferPool.Put(t)
}
```

## Error Handling Patterns

### Graceful Error Handling

```go
func safeDecodeTransfer(topics [][32]byte, data []byte) (*Transfer, *TransferData, error) {
    if len(topics) < 3 {
        return nil, nil, fmt.Errorf("insufficient topics: got %d, need 3", len(topics))
    }

    if topics[0] != TransferTopic {
        return nil, nil, fmt.Errorf("not a transfer event")
    }

    var transfer Transfer
    if err := transfer.DecodeTopics(topics); err != nil {
        return nil, nil, fmt.Errorf("topic decoding failed: %w", err)
    }

    var transferData TransferData
    if err := transferData.Decode(data); err != nil {
        return nil, nil, fmt.Errorf("data decoding failed: %w", err)
    }

    return &transfer, &transferData, nil
}
```

### Validation Helpers

```go
func validateTransfer(transfer *Transfer, data *TransferData) error {
    if transfer.From == (common.Address{}) {
        return fmt.Errorf("invalid from address")
    }
    if transfer.To == (common.Address{}) {
        return fmt.Errorf("invalid to address")
    }
    if data.Value.Sign() < 0 {
        return fmt.Errorf("negative transfer value")
    }
    return nil
}
```

## Testing Patterns

### Unit Test Examples

```go
func TestTransferEvent(t *testing.T) {
    from := common.HexToAddress("0x123...")
    to := common.HexToAddress("0x456...")
    value := big.NewInt(1000)

    transfer := Transfer{From: from, To: to, Value: value}

    // Test topic encoding
    topics, err := transfer.EncodeTopics()
    require.NoError(t, err)
    require.Len(t, topics, 3)
    require.Equal(t, TransferTopic, topics[0])

    // Test topic decoding
    var decoded Transfer
    require.NoError(t, decoded.DecodeTopics(topics))
    require.Equal(t, from, decoded.From)
    require.Equal(t, to, decoded.To)

    // Test data encoding
    data := TransferData{Value: value}
    encoded, err := data.Encode()
    require.NoError(t, err)

    // Test data decoding
    var decodedData TransferData
    require.NoError(t, decodedData.Decode(encoded))
    require.Equal(t, 0, value.Cmp(decodedData.Value))
}
```

## Common Gotchas

### 1. Dynamic Types in Topics

Dynamic types (strings, bytes, arrays) in topics are hashed:

```go
// For event: Event(string indexed message, ...)
// The string is Keccak256 hashed before being put in the topic
```

### 2. Address Padding

Addresses in topics are padded to 32 bytes:

```go
// Address: 0x123... (20 bytes)
// In topic: 000000000000000000000000123... (32 bytes)
```

### 3. Big.Int Zero Values

```go
// Use big.NewInt(0) instead of nil for zero values
data := TransferData{Value: big.NewInt(0)}  // Correct
data := TransferData{Value: nil}           // May cause issues
```

### 4. Empty Arrays

```go
// Empty arrays are valid
data := ComplexEventData{
    Message: "test",
    Numbers: []*big.Int{},  // Empty array
}
```

This quick reference covers the most common patterns and operations for working with generated event code.