package abi

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

// Regular expressions compiled once at package level
var (
	// Function: function name(type1,type2) [payable|view|pure] [returns(type3,type4)]
	// Match basic function structure, handle parameters and returns manually
	functionRegex = regexp.MustCompile(`^function\s+(\w+)\s*\(.*\)\s*(payable|view|pure)?(?:\s+returns\s*\(.*\))?$`)

	// Event: event name(type1 indexed name1, type2 name2)
	eventRegex = regexp.MustCompile(`^event\s+(\w+)\s*\(([^)]*)\)$`)

	// Constructor: constructor(type1,type2) [payable]
	constructorRegex = regexp.MustCompile(`^constructor\s*\(([^)]*)\)\s*(payable)?$`)

	// Fallback/Receive: fallback() [payable] or receive() [payable]
	fallbackRegex = regexp.MustCompile(`^(fallback|receive)\s*\(\s*\)\s*(payable)?$`)

	// Struct: struct Name { type1 name1; type2 name2; }
	structRegex = regexp.MustCompile(`^struct\s+(\w+)\s*\{\s*([^}]*)\s*\}$`)

	// Parameter with optional indexed and name: type [indexed] [name]
	paramRegex = regexp.MustCompile(`^(\S+)(?:\s+(indexed))?(?:\s+(\w+))?$`)

	// Type without tuple: matches types like uint256, address[], bytes32[4], etc.
	typeWithoutTupleRegex = regexp.MustCompile(`^(\w+)((\[\d*\])+)?$`)
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

		item, err := parseLineWithStructs(line, structs)
		if err != nil {
			return nil, fmt.Errorf("failed to parse line '%s': %w", line, err)
		}
		if item != nil {
			jsonABI = append(jsonABI, item)
		}
	}

	if len(jsonABI) == 0 {
		return nil, fmt.Errorf("no valid ABI items found")
	}

	jsonBytes, err := json.Marshal(jsonABI)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal JSON: %w", err)
	}

	return jsonBytes, nil
}

// isStructSignature checks if a line is a struct definition
func isStructSignature(line string) bool {
	return structRegex.MatchString(line)
}

// parseLineWithStructs parses a single line of human-readable ABI with struct context
func parseLineWithStructs(line string, structs map[string][]map[string]interface{}) (map[string]interface{}, error) {
	// Try to match function
	item, err := parseFunctionWithStructs(line, structs)
	if err != nil {
		return nil, err
	}
	if item != nil {
		return item, nil
	}

	// Try to match event
	item, err = parseEventWithStructs(line, structs)
	if err != nil {
		return nil, err
	}
	if item != nil {
		return item, nil
	}

	// Try to match constructor
	item, err = parseConstructorWithStructs(line, structs)
	if err != nil {
		return nil, err
	}
	if item != nil {
		return item, nil
	}

	// Try to match fallback/receive
	if item := parseFallback(line); item != nil {
		return item, nil
	}

	return nil, fmt.Errorf("unrecognized ABI line format: %s", line)
}

// parseLine parses a single line of human-readable ABI
func parseLine(line string) (map[string]interface{}, error) {
	return parseLineWithStructs(line, nil)
}

// parseFunctionWithStructs parses a function definition with struct context
func parseFunctionWithStructs(line string, structs map[string][]map[string]interface{}) (map[string]interface{}, error) {
	matches := functionRegex.FindStringSubmatch(line)
	if matches == nil {
		return nil, nil
	}

	name := matches[1]
	inputsStr := ""
	outputsStr := ""

	// Manually extract parameters section
	openParen := strings.Index(line, "(")
	if openParen != -1 {
		// Find the matching closing parenthesis for parameters
		parenCount := 1
		for i := openParen + 1; i < len(line); i++ {
			if line[i] == '(' {
				parenCount++
			} else if line[i] == ')' {
				parenCount--
				if parenCount == 0 {
					inputsStr = line[openParen+1 : i]
					break
				}
			}
		}
	}

	// Manually extract returns section if it exists
	returnsIndex := -1
	if strings.Contains(line, "returns") {
		returnsIndex = strings.Index(line, "returns")
		if returnsIndex != -1 {
			// Find the opening parenthesis after "returns"
			openParen := strings.Index(line[returnsIndex:], "(")
			if openParen != -1 {
				start := returnsIndex + openParen + 1
				// Find the matching closing parenthesis
				parenCount := 1
				for i := start; i < len(line); i++ {
					if line[i] == '(' {
						parenCount++
					} else if line[i] == ')' {
						parenCount--
						if parenCount == 0 {
							outputsStr = line[start:i]
							break
						}
					}
				}
			}
		}
	}

	// Extract state mutability manually - look for payable/view/pure between parameters and returns
	stateMutability := "nonpayable"
	if returnsIndex != -1 {
		// Look for state mutability between the end of parameters and "returns"
		endOfParams := openParen + len(inputsStr) + 2 // position after closing parenthesis of parameters
		if endOfParams < returnsIndex {
			between := strings.TrimSpace(line[endOfParams:returnsIndex])
			if between == "payable" {
				stateMutability = "payable"
			} else if between == "view" {
				stateMutability = "view"
			} else if between == "pure" {
				stateMutability = "pure"
			}
		}
	} else {
		// No returns clause, look for state mutability after parameters
		endOfParams := openParen + len(inputsStr) + 2 // position after closing parenthesis of parameters
		if endOfParams < len(line) {
			remaining := strings.TrimSpace(line[endOfParams:])
			if remaining == "payable" {
				stateMutability = "payable"
			} else if remaining == "view" {
				stateMutability = "view"
			} else if remaining == "pure" {
				stateMutability = "pure"
			}
		}
	}

	inputs, err := parseParametersWithStructs(inputsStr, false, structs)
	if err != nil {
		return nil, err
	}

	outputs := []map[string]interface{}{}
	if outputsStr != "" {
		outputs, err = parseParametersWithStructs(outputsStr, false, structs)
		if err != nil {
			return nil, err
		}
	}

	return map[string]interface{}{
		"type":            "function",
		"name":            name,
		"inputs":          inputs,
		"outputs":         outputs,
		"stateMutability": stateMutability,
	}, nil
}

// parseFunction parses a function definition
func parseFunction(line string) (map[string]interface{}, error) {
	return parseFunctionWithStructs(line, nil)
}

// parseEventWithStructs parses an event definition with struct context
func parseEventWithStructs(line string, structs map[string][]map[string]interface{}) (map[string]interface{}, error) {
	matches := eventRegex.FindStringSubmatch(line)
	if matches == nil {
		return nil, nil
	}

	name := matches[1]
	inputsStr := matches[2]

	inputs, err := parseParametersWithStructs(inputsStr, true, structs)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"type":      "event",
		"name":      name,
		"inputs":    inputs,
		"anonymous": false,
	}, nil
}

// parseEvent parses an event definition
func parseEvent(line string) (map[string]interface{}, error) {
	return parseEventWithStructs(line, nil)
}

// parseConstructorWithStructs parses a constructor definition with struct context
func parseConstructorWithStructs(line string, structs map[string][]map[string]interface{}) (map[string]interface{}, error) {
	matches := constructorRegex.FindStringSubmatch(line)
	if matches == nil {
		return nil, nil
	}

	inputsStr := matches[1]
	stateMutability := matches[2]

	if stateMutability == "" {
		stateMutability = "nonpayable"
	}

	inputs, err := parseParametersWithStructs(inputsStr, false, structs)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"type":            "constructor",
		"inputs":          inputs,
		"stateMutability": stateMutability,
	}, nil
}

// parseConstructor parses a constructor definition
func parseConstructor(line string) (map[string]interface{}, error) {
	return parseConstructorWithStructs(line, nil)
}

// parseFallback parses fallback and receive function definitions
func parseFallback(line string) map[string]interface{} {
	matches := fallbackRegex.FindStringSubmatch(line)
	if matches == nil {
		return nil
	}

	funcType := matches[1]
	stateMutability := matches[2]

	if stateMutability == "" {
		stateMutability = "nonpayable"
	}

	return map[string]interface{}{
		"type":            funcType,
		"stateMutability": stateMutability,
	}
}

// parseParametersWithStructs parses a comma-separated list of parameters with struct context
func parseParametersWithStructs(paramsStr string, isEvent bool, structs map[string][]map[string]interface{}) ([]map[string]interface{}, error) {
	if strings.TrimSpace(paramsStr) == "" {
		return []map[string]interface{}{}, nil
	}

	// Parse parameters with proper nested parentheses handling
	params, err := splitByCommaOutsideParentheses(paramsStr)
	if err != nil {
		return nil, err
	}

	result := make([]map[string]interface{}, 0, len(params))

	for _, param := range params {
		param = strings.TrimSpace(param)
		if param == "" {
			continue
		}

		// Parse parameter components
		paramMap, err := parseParameterWithStructs(param, isEvent, structs)
		if err != nil {
			return nil, err
		}

		result = append(result, paramMap)
	}

	return result, nil
}

// parseParameterWithStructs parses a single parameter string with struct context
func parseParameterWithStructs(paramStr string, isEvent bool, structs map[string][]map[string]interface{}) (map[string]interface{}, error) {
	// For tuple types, we need special handling
	// Look for opening parenthesis and find matching closing parenthesis
	if strings.HasPrefix(paramStr, "(") {
		// Find the matching closing parenthesis
		parenCount := 0
		for _, ch := range paramStr {
			if ch == '(' {
				parenCount++
			} else if ch == ')' {
				parenCount--
				if parenCount == 0 {
					// Found matching closing parenthesis at position i
					return parseTupleParameterWithStructs(paramStr, isEvent, structs)
				}
			}
		}
	}

	// For regular types, use regex parsing
	matches := paramRegex.FindStringSubmatch(paramStr)
	if matches == nil {
		return nil, fmt.Errorf("invalid parameter format: %s", paramStr)
	}

	typeStr := matches[1]
	indexed := matches[2] == "indexed"
	name := matches[3]

	matches = typeWithoutTupleRegex.FindStringSubmatch(typeStr)
	if matches == nil {
		return nil, fmt.Errorf("invalid type format: %s", typeStr)
	}
	baseType := matches[1]
	arrayPart := matches[2]

	// Check if this is a struct reference
	if structs != nil {
		if structComponents, exists := structs[baseType]; exists {
			// Create tuple array type with components
			result := map[string]interface{}{
				"name":         name,
				"type":         "tuple" + arrayPart,
				"internalType": "struct " + baseType + arrayPart,
				"components":   structComponents,
			}
			if isEvent {
				result["indexed"] = indexed
			}
			return result, nil
		}
	}

	// Validate and normalize type
	var err error
	baseType, err = normalizeType(baseType)
	if err != nil {
		return nil, err
	}

	paramMap := map[string]interface{}{
		"name": name,
		"type": baseType + arrayPart,
	}

	if isEvent {
		paramMap["indexed"] = indexed
	}

	return paramMap, nil
}

// parseTupleParameterWithStructs parses a tuple parameter with struct context
func parseTupleParameterWithStructs(paramStr string, isEvent bool, structs map[string][]map[string]interface{}) (map[string]interface{}, error) {
	// Find the matching closing parenthesis for the tuple content
	parenCount := 0
	tupleEnd := -1
	for i, ch := range paramStr {
		if ch == '(' {
			parenCount++
		} else if ch == ')' {
			parenCount--
			if parenCount == 0 {
				tupleEnd = i
				break
			}
		}
	}

	if tupleEnd == -1 {
		return nil, fmt.Errorf("unbalanced parentheses in tuple: %s", paramStr)
	}

	// Extract the content inside the tuple parentheses
	content := strings.TrimSpace(paramStr[1:tupleEnd])

	// Parse the tuple components
	components, err := parseParametersWithStructs(content, false, structs)
	if err != nil {
		return nil, err
	}

	// Extract name and array info from the part after the tuple
	name := ""
	isArray := false
	isFixedArray := false
	arraySize := ""

	if tupleEnd+1 < len(paramStr) {
		remaining := strings.TrimSpace(paramStr[tupleEnd+1:])
		if remaining != "" {
			// Check for array types
			if strings.HasPrefix(remaining, "[]") {
				isArray = true
				// Update the name to remove the array brackets
				name = strings.TrimSpace(strings.TrimPrefix(remaining, "[]"))
			} else if bracketIdx := strings.Index(remaining, "["); bracketIdx != -1 && strings.HasSuffix(remaining, "]") {
				isFixedArray = true
				arraySize = remaining[bracketIdx+1 : len(remaining)-1]
				// Update the name to remove the fixed array brackets
				name = strings.TrimSpace(remaining[:bracketIdx])
			} else {
				// No array, just a name
				name = remaining
			}
		}
	}

	paramMap := map[string]interface{}{
		"name":       name,
		"type":       "tuple",
		"components": components,
	}

	if isArray {
		paramMap["type"] = "tuple[]"
	} else if isFixedArray {
		paramMap["type"] = "tuple[" + arraySize + "]"
	}

	// Only add indexed field for events
	// For functions, don't include the indexed field at all

	return paramMap, nil
}

// splitByCommaOutsideParentheses splits a string by commas that are not inside parentheses
func splitByCommaOutsideParentheses(s string) ([]string, error) {
	var parts []string
	var current strings.Builder
	parenCount := 0

	for _, ch := range s {
		if ch == '(' {
			parenCount++
			current.WriteRune(ch)
		} else if ch == ')' {
			parenCount--
			current.WriteRune(ch)
		} else if ch == ',' && parenCount == 0 {
			// Only split on commas that are not inside parentheses
			part := strings.TrimSpace(current.String())
			if part != "" {
				parts = append(parts, part)
			}
			current.Reset()
		} else {
			current.WriteRune(ch)
		}
	}

	// Add the last part
	part := strings.TrimSpace(current.String())
	if part != "" {
		parts = append(parts, part)
	}

	// Validate that all parentheses are balanced
	if parenCount != 0 {
		return nil, fmt.Errorf("unbalanced parentheses in parameter string: %s", s)
	}

	return parts, nil
}

// normalizeType validates and normalizes Solidity type names
func normalizeType(typeStr string) (string, error) {
	// Handle arrays first
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

		if len(typeStr) == len(prefix) {
			// No size specified - normalize to 256 bits (Solidity default)
			return prefix + "256", nil
		}

		if len(typeStr) > len(prefix) {
			sizeStr := typeStr[len(prefix):]
			if size, err := strconv.Atoi(sizeStr); err == nil && size >= 8 && size <= 256 {
				return typeStr, nil
			}
		}
		return "", fmt.Errorf("invalid integer type: %s", typeStr)
	}

	// Handle tuple types (already handled in parseParameter)
	if strings.HasPrefix(typeStr, "(") && strings.HasSuffix(typeStr, ")") {
		return typeStr, nil
	}

	// For now, treat any unrecognized type as a potential struct reference
	// This allows the parsing to continue and the struct resolution can happen later
	return typeStr, nil
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

		matches := structRegex.FindStringSubmatch(line)
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

			var err error
			paramType, err = normalizeType(paramType)
			if err != nil {
				return nil, fmt.Errorf("invalid type in struct %s: %s", name, paramType)
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

		matches := typeWithoutTupleRegex.FindStringSubmatch(paramType)
		if matches == nil {
			return nil, fmt.Errorf("invalid type format in struct: %s", paramType)
		}

		baseType := matches[1]
		arrayPart := matches[2]

		// Check if this is a struct reference
		if nestedStruct, exists := structs[baseType]; exists {
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
				"type":         "tuple" + arrayPart,
				"internalType": "struct " + baseType + arrayPart,
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
