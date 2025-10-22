package abi

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

// Regular expressions compiled once at package level to avoid recompilation
var (
	structLineRegex      = regexp.MustCompile(`^struct\s+\w+\s*\{\s*.*\s*\}$`)
	structDefinitionRegex = regexp.MustCompile(`^struct\s+(\w+)\s*\{\s*(.*)\s*\}$`)
	eventRegex           = regexp.MustCompile(`^event\s+(\w+)\s*\(([^)]*)\)$`)
	constructorRegex     = regexp.MustCompile(`^constructor\s*\(([^)]*)\)\s*(\w*)$`)
	fallbackRegex        = regexp.MustCompile(`^(fallback|receive)\s*\(\s*\)\s*(\w*)$`)
)

// ParseHumanReadableABI parses human-readable ABI definitions and converts them to JSON ABI format
func ParseHumanReadableABI(humanABI []string) ([]byte, error) {
	// First pass: extract and parse all struct definitions
	structs, err := parseStructs(humanABI)
	if err != nil {
		return nil, fmt.Errorf("failed to parse structs: %w", err)
	}

	// Second pass: parse all non-struct signatures with struct context
	var jsonABI []map[string]interface{}
	for _, line := range humanABI {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "//") {
			continue
		}

		// Skip struct definitions - they're only used for type resolution
		if isStructSignature(line) {
			continue
		}

		item, err := parseHumanReadableLineWithStructs(line, structs)
		if err != nil {
			return nil, fmt.Errorf("failed to parse line '%s': %w", line, err)
		}
		if item != nil {
			jsonABI = append(jsonABI, item)
		}
	}

	// Convert to JSON bytes
	jsonBytes, err := json.Marshal(jsonABI)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal JSON: %w", err)
	}

	return jsonBytes, nil
}

// parseHumanReadableLine parses a single line of human-readable ABI
func parseHumanReadableLine(line string) (map[string]interface{}, error) {
	return parseHumanReadableLineWithStructs(line, nil)
}

// parseHumanReadableLineWithStructs parses a single line of human-readable ABI with struct context
func parseHumanReadableLineWithStructs(line string, structs map[string][]map[string]interface{}) (map[string]interface{}, error) {
	// Try to match function
	if item, err := parseFunctionWithStructs(line, structs); err == nil && item != nil {
		return item, nil
	}

	// Try to match event
	if item, err := parseEventWithStructs(line, structs); err == nil && item != nil {
		return item, nil
	}

	// Try to match constructor
	if item, err := parseConstructorWithStructs(line, structs); err == nil && item != nil {
		return item, nil
	}

	// Try to match fallback/receive
	if item, err := parseFallback(line); err == nil && item != nil {
		return item, nil
	}

	return nil, fmt.Errorf("unrecognized ABI line format")
}

// isStructSignature checks if a line is a struct definition
func isStructSignature(line string) bool {
	return structLineRegex.MatchString(line)
}

// parseStructs parses struct definitions from a list of lines
func parseStructs(lines []string) (map[string][]map[string]interface{}, error) {
	structs := make(map[string][]map[string]interface{})

	// First pass: create shallow structs (without resolving nested struct references)
	shallowStructs := make(map[string][]map[string]interface{})
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "//") {
			continue
		}

		matches := structDefinitionRegex.FindStringSubmatch(line)
		if matches == nil {
			continue
		}

		name := matches[1]
		properties := matches[2]

		// Parse properties (split by semicolon)
		propList := strings.Split(properties, ";")
		components := []map[string]interface{}{}

		for _, prop := range propList {
			prop = strings.TrimSpace(prop)
			if prop == "" {
				continue
			}

			// Parse each property as a parameter
			parts := strings.Fields(prop)
			if len(parts) < 1 {
				continue
			}

			paramType := parts[0]
			paramName := ""
			if len(parts) > 1 {
				paramName = parts[1]
			}

			// For struct parsing, we don't validate types yet
			component := map[string]interface{}{
				"name": paramName,
				"type": paramType,
			}
			components = append(components, component)
		}

		if len(components) == 0 {
			return nil, fmt.Errorf("invalid struct signature (no properties): %s", line)
		}

		shallowStructs[name] = components
	}

	// Second pass: resolve nested struct references
	for name, parameters := range shallowStructs {
		resolved, err := resolveStructComponents(parameters, shallowStructs, make(map[string]bool))
		if err != nil {
			return nil, err
		}
		structs[name] = resolved
	}

	return structs, nil
}

// resolveStructComponents recursively resolves struct references in parameter components
func resolveStructComponents(parameters []map[string]interface{}, structs map[string][]map[string]interface{}, ancestors map[string]bool) ([]map[string]interface{}, error) {
	components := []map[string]interface{}{}

	for _, param := range parameters {
		paramType := param["type"].(string)

		// If already a tuple, keep it as-is
		if strings.HasPrefix(paramType, "tuple") {
			components = append(components, param)
			continue
		}

		// Check if this is a struct reference
		if nestedStruct, exists := structs[paramType]; exists {
			// Detect circular references
			if ancestors[paramType] {
				return nil, fmt.Errorf("circular reference detected: %s", paramType)
			}

			// Recursively resolve nested structs
			newAncestors := make(map[string]bool)
			for k, v := range ancestors {
				newAncestors[k] = v
			}
			newAncestors[paramType] = true

			resolvedComponents, err := resolveStructComponents(nestedStruct, structs, newAncestors)
			if err != nil {
				return nil, err
			}

			// Create tuple type with components and internalType
			tupleParam := map[string]interface{}{
				"name":         param["name"],
				"type":         "tuple",
				"internalType": "struct " + paramType,
				"components":   resolvedComponents,
			}
			components = append(components, tupleParam)
		} else {
			// Not a struct, validate it's a valid Solidity type
			if _, err := normalizeType(paramType); err != nil {
				return nil, fmt.Errorf("unknown type: %s", paramType)
			}
			components = append(components, param)
		}
	}

	return components, nil
}

// parseFunction parses a function definition
func parseFunction(line string) (map[string]interface{}, error) {
	return parseFunctionWithStructs(line, nil)
}

// parseFunctionWithStructs parses a function definition with struct context
func parseFunctionWithStructs(line string, structs map[string][]map[string]interface{}) (map[string]interface{}, error) {
	// Check if this is a function line
	if !strings.HasPrefix(line, "function ") {
		return nil, nil
	}

	// Extract function name
	line = strings.TrimPrefix(line, "function ")
	nameEnd := strings.Index(line, "(")
	if nameEnd == -1 {
		return nil, nil
	}
	name := strings.TrimSpace(line[:nameEnd])
	line = line[nameEnd:]

	// Find the matching closing parenthesis for inputs
	bracketCount := 0
	inputsEnd := -1
	for i, ch := range line {
		if ch == '(' {
			bracketCount++
		} else if ch == ')' {
			bracketCount--
			if bracketCount == 0 {
				inputsEnd = i
				break
			}
		}
	}
	if inputsEnd == -1 {
		return nil, nil
	}

	inputsStr := line[1:inputsEnd]
	line = line[inputsEnd+1:]

	// Parse state mutability and returns
	stateMutability := "nonpayable"
	var outputsStr string

	// Check for state mutability
	line = strings.TrimSpace(line)
	if strings.HasPrefix(line, "view ") || strings.HasPrefix(line, "pure ") || strings.HasPrefix(line, "payable ") || line == "view" || line == "pure" || line == "payable" {
		parts := strings.Fields(line)
		stateMutability = parts[0]
		line = strings.TrimSpace(strings.TrimPrefix(line, parts[0]))
	}

	// Check for returns
	if strings.HasPrefix(line, "returns (") {
		line = strings.TrimPrefix(line, "returns ")
		// Find the matching closing parenthesis for returns
		bracketCount = 0
		returnsEnd := -1
		for i, ch := range line {
			if ch == '(' {
				bracketCount++
			} else if ch == ')' {
				bracketCount--
				if bracketCount == 0 {
					returnsEnd = i
					break
				}
			}
		}
		if returnsEnd == -1 {
			return nil, nil
		}
		outputsStr = line[1:returnsEnd] // Skip the opening '('
	}

	// Parse inputs
	inputs, err := parseParametersWithStructs(inputsStr, structs)
	if err != nil {
		return nil, fmt.Errorf("failed to parse inputs for function %s: %w", name, err)
	}

	// Parse outputs
	var outputs []map[string]interface{}
	if outputsStr != "" {
		outputs, err = parseParametersWithStructs(outputsStr, structs)
		if err != nil {
			return nil, fmt.Errorf("failed to parse outputs for function %s: %w", name, err)
		}
	} else {
		outputs = []map[string]interface{}{}
	}

	return map[string]interface{}{
		"type":            "function",
		"name":            name,
		"inputs":          inputs,
		"outputs":         outputs,
		"stateMutability": stateMutability,
	}, nil
}

// parseEvent parses an event definition
func parseEvent(line string) (map[string]interface{}, error) {
	return parseEventWithStructs(line, nil)
}

// parseEventWithStructs parses an event definition with struct context
func parseEventWithStructs(line string, structs map[string][]map[string]interface{}) (map[string]interface{}, error) {
	// Match event with optional indexed parameters
	// Examples:
	// "event Transfer(address from, address to, uint256 value)"
	// "event Transfer(address indexed from, address indexed to, uint256 value)"
	matches := eventRegex.FindStringSubmatch(line)
	if matches == nil {
		return nil, nil
	}

	name := matches[1]
	inputsStr := matches[2]

	// Parse inputs with indexed support
	inputs, err := parseEventParametersWithStructs(inputsStr, structs)
	if err != nil {
		return nil, fmt.Errorf("failed to parse event inputs: %w", err)
	}

	return map[string]interface{}{
		"type":      "event",
		"name":      name,
		"inputs":    inputs,
		"anonymous": false,
	}, nil
}

// parseConstructor parses a constructor definition
func parseConstructor(line string) (map[string]interface{}, error) {
	return parseConstructorWithStructs(line, nil)
}

// parseConstructorWithStructs parses a constructor definition with struct context
func parseConstructorWithStructs(line string, structs map[string][]map[string]interface{}) (map[string]interface{}, error) {
	// Match constructor
	// Examples:
	// "constructor(address owner, uint256 initialSupply)"
	// "constructor(address owner, uint256 initialSupply) payable"
	matches := constructorRegex.FindStringSubmatch(line)
	if matches == nil {
		return nil, nil
	}

	inputsStr := matches[1]
	stateMutability := matches[2]

	// Parse inputs
	inputs, err := parseParametersWithStructs(inputsStr, structs)
	if err != nil {
		return nil, fmt.Errorf("failed to parse constructor inputs: %w", err)
	}

	// Determine state mutability
	if stateMutability == "" {
		stateMutability = "nonpayable"
	} else if stateMutability != "payable" {
		return nil, fmt.Errorf("invalid state mutability for constructor: %s", stateMutability)
	}

	return map[string]interface{}{
		"type":            "constructor",
		"inputs":          inputs,
		"stateMutability": stateMutability,
	}, nil
}

// parseFallback parses fallback and receive function definitions
func parseFallback(line string) (map[string]interface{}, error) {
	// Match fallback function
	// Examples:
	// "fallback()"
	// "fallback() payable"
	// "receive() payable"
	matches := fallbackRegex.FindStringSubmatch(line)
	if matches == nil {
		return nil, nil
	}

	funcType := matches[1]
	stateMutability := matches[2]

	// Determine state mutability
	if stateMutability == "" {
		stateMutability = "nonpayable"
	} else if stateMutability != "payable" {
		return nil, fmt.Errorf("invalid state mutability for %s: %s", funcType, stateMutability)
	}

	return map[string]interface{}{
		"type":            funcType,
		"stateMutability": stateMutability,
	}, nil
}

// parseParameters parses a comma-separated list of parameters
func parseParameters(paramsStr string) ([]map[string]interface{}, error) {
	return parseParametersWithStructs(paramsStr, nil)
}

// parseParametersWithStructs parses a comma-separated list of parameters with struct context
func parseParametersWithStructs(paramsStr string, structs map[string][]map[string]interface{}) ([]map[string]interface{}, error) {
	if strings.TrimSpace(paramsStr) == "" {
		return []map[string]interface{}{}, nil
	}

	params := strings.Split(paramsStr, ",")
	result := make([]map[string]interface{}, 0, len(params))

	for _, param := range params {
		param = strings.TrimSpace(param)
		if param == "" {
			continue
		}

		// Parse parameter: "type name" or just "type"
		parts := strings.Fields(param)
		if len(parts) == 0 {
			continue
		}

		paramType := parts[0]
		paramName := ""
		if len(parts) > 1 {
			paramName = parts[1]
		}

		// Check if this is a struct reference (handle arrays too)
		if structs != nil {
			// Check for array types with structs
			if strings.HasSuffix(paramType, "[]") {
				baseType := paramType[:len(paramType)-2]
				if structComponents, exists := structs[baseType]; exists {
					// Create tuple array type with components
					result = append(result, map[string]interface{}{
						"name":         paramName,
						"type":         "tuple[]",
						"internalType": "struct " + baseType + "[]",
						"components":   structComponents,
					})
					continue
				}
			}
			// Check for fixed array types with structs
			if idx := strings.Index(paramType, "["); idx != -1 && strings.HasSuffix(paramType, "]") {
				baseType := paramType[:idx]
				if structComponents, exists := structs[baseType]; exists {
					// Create tuple fixed array type with components
					result = append(result, map[string]interface{}{
						"name":         paramName,
						"type":         "tuple" + paramType[idx:],
						"internalType": "struct " + baseType + paramType[idx:],
						"components":   structComponents,
					})
					continue
				}
			}
			// Check for regular struct reference
			if structComponents, exists := structs[paramType]; exists {
				// Create tuple type with components
				result = append(result, map[string]interface{}{
					"name":         paramName,
					"type":         "tuple",
					"internalType": "struct " + paramType,
					"components":   structComponents,
				})
				continue
			}
		}

		// Validate and normalize type
		normalizedType, err := normalizeType(paramType)
		if err != nil {
			return nil, fmt.Errorf("invalid type '%s': %w", paramType, err)
		}

		result = append(result, map[string]interface{}{
			"name": paramName,
			"type": normalizedType,
		})
	}

	return result, nil
}

// parseEventParameters parses event parameters with indexed support
func parseEventParameters(paramsStr string) ([]map[string]interface{}, error) {
	return parseEventParametersWithStructs(paramsStr, nil)
}

// parseEventParametersWithStructs parses event parameters with indexed support and struct context
func parseEventParametersWithStructs(paramsStr string, structs map[string][]map[string]interface{}) ([]map[string]interface{}, error) {
	if strings.TrimSpace(paramsStr) == "" {
		return []map[string]interface{}{}, nil
	}

	params := strings.Split(paramsStr, ",")
	result := make([]map[string]interface{}, 0, len(params))

	for _, param := range params {
		param = strings.TrimSpace(param)
		if param == "" {
			continue
		}

		// Parse parameter: "type [indexed] name" or "type name"
		parts := strings.Fields(param)
		if len(parts) == 0 {
			continue
		}

		indexed := false
		paramType := parts[0]
		paramName := ""

		// Check for indexed keyword - it can appear after the type
		if len(parts) > 1 && parts[1] == "indexed" {
			indexed = true
			if len(parts) > 2 {
				paramName = parts[2]
			}
		} else if len(parts) > 1 {
			// No indexed keyword, second part is the name
			paramName = parts[1]
		}

		// Check if this is a struct reference (handle arrays too)
		if structs != nil {
			// Check for array types with structs
			if strings.HasSuffix(paramType, "[]") {
				baseType := paramType[:len(paramType)-2]
				if structComponents, exists := structs[baseType]; exists {
					// Create tuple array type with components
					result = append(result, map[string]interface{}{
						"name":         paramName,
						"type":         "tuple[]",
						"internalType": "struct " + baseType + "[]",
						"components":   structComponents,
						"indexed":      indexed,
					})
					continue
				}
			}
			// Check for fixed array types with structs
			if idx := strings.Index(paramType, "["); idx != -1 && strings.HasSuffix(paramType, "]") {
				baseType := paramType[:idx]
				if structComponents, exists := structs[baseType]; exists {
					// Create tuple fixed array type with components
					result = append(result, map[string]interface{}{
						"name":         paramName,
						"type":         "tuple" + paramType[idx:],
						"internalType": "struct " + baseType + paramType[idx:],
						"components":   structComponents,
						"indexed":      indexed,
					})
					continue
				}
			}
			// Check for regular struct reference
			if structComponents, exists := structs[paramType]; exists {
				// Create tuple type with components
				result = append(result, map[string]interface{}{
					"name":         paramName,
					"type":         "tuple",
					"internalType": "struct " + paramType,
					"components":   structComponents,
					"indexed":      indexed,
				})
				continue
			}
		}

		// Validate and normalize type
		normalizedType, err := normalizeType(paramType)
		if err != nil {
			return nil, fmt.Errorf("invalid type '%s': %w", paramType, err)
		}

		result = append(result, map[string]interface{}{
			"name":    paramName,
			"type":    normalizedType,
			"indexed": indexed,
		})
	}

	return result, nil
}

// normalizeType validates and normalizes Solidity type names
func normalizeType(typeStr string) (string, error) {
	// Handle arrays first (they have higher priority)
	if strings.HasSuffix(typeStr, "[]") {
		elemType := typeStr[:len(typeStr)-2]
		normalizedElem, err := normalizeType(elemType)
		if err != nil {
			return "", err
		}
		return normalizedElem + "[]", nil
	}

	// Handle fixed arrays
	if idx := strings.Index(typeStr, "["); idx != -1 && strings.HasSuffix(typeStr, "]") {
		elemType := typeStr[:idx]
		sizeStr := typeStr[idx+1 : len(typeStr)-1]

		normalizedElem, err := normalizeType(elemType)
		if err != nil {
			return "", err
		}

		if _, err := strconv.Atoi(sizeStr); err != nil {
			return "", fmt.Errorf("invalid array size '%s'", sizeStr)
		}

		return normalizedElem + "[" + sizeStr + "]", nil
	}

	// Handle basic types
	basicTypes := map[string]string{
		"address": "address",
		"bool":    "bool",
		"string":  "string",
		"bytes":   "bytes",
	}

	if normalized, exists := basicTypes[typeStr]; exists {
		return normalized, nil
	}

	// Handle fixed bytes (bytes1 to bytes32)
	if strings.HasPrefix(typeStr, "bytes") {
		if len(typeStr) > 5 {
			sizeStr := typeStr[5:]
			if size, err := strconv.Atoi(sizeStr); err == nil && size >= 1 && size <= 32 {
				return typeStr, nil
			}
		}
		return "", fmt.Errorf("invalid bytes type: %s", typeStr)
	}

	// Handle integers (u)int8 to (u)int256
	if strings.HasPrefix(typeStr, "uint") || strings.HasPrefix(typeStr, "int") {
		// Extract size
		prefix := ""
		if strings.HasPrefix(typeStr, "uint") {
			prefix = "uint"
		} else {
			prefix = "int"
		}

		if len(typeStr) > len(prefix) {
			sizeStr := typeStr[len(prefix):]
			if size, err := strconv.Atoi(sizeStr); err == nil && size >= 8 && size <= 256 && size%8 == 0 {
				return typeStr, nil
			}
		}
		return "", fmt.Errorf("invalid integer type: %s", typeStr)
	}

	// Handle tuples (we'll treat any unrecognized type as a tuple for now)
	// In a real implementation, you might want more sophisticated tuple detection
	return typeStr, nil
}
