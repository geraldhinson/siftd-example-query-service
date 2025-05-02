package main

import (
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
)

type QueryParam struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

type Output struct {
	Query           string       `json:"query"`
	QueryParameters []QueryParam `json:"queryParameters"`
}

// Infer SQL type -> API type
func inferType(value string) string {
	val := strings.ToLower(value)

	// Case-insensitive cast match
	typeCastRegex := regexp.MustCompile(`::\s*([a-zA-Z0-9_ \[\]]+)`)
	if matches := typeCastRegex.FindStringSubmatch(val); len(matches) > 1 {
		switch strings.TrimSpace(matches[1]) {
		case "boolean":
			return "BOOLEAN"
		case "smallint", "int2":
			return "SHORT"
		case "int", "integer", "int4":
			return "INTEGER"
		case "bigint", "int8":
			return "LONG"
		case "float", "float4", "real":
			return "FLOAT"
		case "double", "double precision", "float8":
			return "DOUBLE"
		case "varchar", "character varying", "text", "char":
			return "STRING"
		case "uuid":
			return "GUID"
		case "date":
			return "DATE"
		case "timestamp", "timestamp without time zone", "timestamp with time zone", "timestamptz":
			return "TIMESTAMP"
		case "json", "jsonb":
			return "JSON"
		case "varchar[]", "text[]", "character varying[]":
			return "ARRAY_VARCHAR"
		case "int[]", "integer[]":
			return "ARRAY_INTEGER"
		case "date[]":
			return "ARRAY_DATE"
		default:
			return "STRING"
		}
	}

	// Literal fallback
	val = strings.Trim(val, "'\"")
	if val == "true" || val == "false" {
		return "BOOLEAN"
	}
	if _, err := strconv.ParseInt(val, 10, 16); err == nil {
		return "SHORT"
	}
	if _, err := strconv.ParseInt(val, 10, 32); err == nil {
		return "INTEGER"
	}
	if _, err := strconv.ParseInt(val, 10, 64); err == nil {
		return "LONG"
	}
	if _, err := strconv.ParseFloat(val, 32); err == nil {
		return "FLOAT"
	}
	if _, err := strconv.ParseFloat(val, 64); err == nil {
		return "DOUBLE"
	}
	if strings.HasPrefix(val, "[") && strings.HasSuffix(val, "]") {
		return "ARRAY_VARCHAR"
	}

	return "STRING"
}

func ConvertSQL(input string) (*Output, error) {
	cteRegex := regexp.MustCompile(`(?is)WITH\s+vars\s+AS\s*\(\s*SELECT\s+(.*?)\)\s*`)
	cteMatch := cteRegex.FindStringSubmatch(input)
	if cteMatch == nil {
		return nil, fmt.Errorf("no CTE 'vars' block found")
	}

	varBlock := cteMatch[1]
	paramLines := strings.Split(varBlock, ",")

	var params []QueryParam
	for _, line := range paramLines {
		line = strings.TrimSpace(line)
		parts := regexp.MustCompile(`(?i)\s+AS\s+`).Split(line, 2)
		if len(parts) != 2 {
			continue
		}
		val := strings.TrimSpace(parts[0])
		name := strings.Trim(parts[1], `" '`)
		params = append(params, QueryParam{Name: name, Type: inferType(val)})
	}

	// Remove the CTE block and replace vars.{x} with {x}
	withoutCTE := cteRegex.ReplaceAllString(input, "")
	withoutVars := regexp.MustCompile(`vars\.(\w+)`).ReplaceAllString(withoutCTE, `{$1}`)
	withoutVars = strings.ReplaceAll(withoutVars, ", vars", "")

	cleanedQuery := strings.Join(strings.Fields(withoutVars), " ") // collapses all whitespace
	output := Output{
		Query:           cleanedQuery,
		QueryParameters: params,
	}

	return &output, nil
}

func ToJSON(output *Output) ([]byte, error) {
	return json.MarshalIndent(output, "", "  ")
}

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: convert_sql <input_file.sql> <output_file.json>")
		return
	}

	inputPath := os.Args[1]
	outputPath := os.Args[2]

	data, err := os.ReadFile(inputPath)
	if err != nil {
		fmt.Println("Error reading input:", err)
		os.Exit(1)
	}

	result, err := ConvertSQL(string(data))
	if err != nil {
		fmt.Println("Error converting SQL:", err)
		os.Exit(1)
	}

	jsonBytes, err := ToJSON(result)
	if err != nil {
		fmt.Println("Error converting to JSON:", err)
		os.Exit(1)
	}

	if err := os.WriteFile(outputPath, jsonBytes, 0644); err != nil {
		fmt.Println("Error writing file:", err)
		os.Exit(1)
	}

	fmt.Println("✅ Conversion complete! Output written to:", outputPath)
}
