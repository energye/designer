//----------------------------------------
//
// Copyright © yanghy. All Rights Reserved.
//
// Licensed under Apache License Version 2.0, January 2004
//
// https://www.apache.org/licenses/LICENSE-2.0
//
//----------------------------------------

package env

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"github.com/energye/designer/cmd/dflag"
	"github.com/energye/designer/pkg/config"
)

type jsonPathToken struct {
	key      string
	hasKey   bool
	index    int
	hasIndex bool
}

type jsonMatch struct {
	path  string
	value any
}

type orderedObject struct {
	members []orderedMember
}

type orderedMember struct {
	key   string
	value any
}

func RunConfig(args dflag.Args) {
	if err := runConfig(args, os.Stdin, os.Stdout, os.Stderr); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
}

func runConfig(args dflag.Args, in io.Reader, out, errOut io.Writer) error {
	configFile := filepath.Join(config.Path(), "config.json")
	if args.Contains("w") {
		if len(args.Positionals()) > 0 {
			return errors.New("energy env: write accepts only one -w expression")
		}
		return writeConfig(configFile, args.Get("w"), in, out)
	}
	positionals := args.Positionals()
	if len(positionals) > 1 {
		return errors.New("energy env: read accepts only one key or json path")
	}
	return readConfig(configFile, firstArg(positionals), out)
}

func firstArg(args []string) string {
	if len(args) == 0 {
		return ""
	}
	return args[0]
}

func readConfig(configFile, query string, out io.Writer) error {
	root, err := loadJSONFile(configFile)
	if err != nil {
		return err
	}
	query = strings.TrimSpace(query)
	if query == "" {
		return printJSON(out, root)
	}
	if isJSONPath(query) {
		value, ok := getPath(root, query)
		if ok {
			fmt.Fprintf(out, "%s=%s\n", query, formatValue(value))
			return nil
		}
		matches, ok := findIndexedKey(root, query)
		if !ok || len(matches) == 0 {
			fmt.Fprintln(out)
			return nil
		}
		sortMatches(matches)
		for _, match := range matches {
			fmt.Fprintf(out, "%s=%s\n", match.path, formatValue(match.value))
		}
		return nil
	}
	matches := findKey(root, query)
	sortMatches(matches)
	for _, match := range matches {
		fmt.Fprintf(out, "%s=%s\n", match.path, formatValue(match.value))
	}
	return nil
}

func writeConfig(configFile, expr string, in io.Reader, out io.Writer) error {
	expr = strings.TrimSpace(expr)
	if expr == "" {
		return errors.New("energy env: -w requires key=value or path=value")
	}
	name, value, ok := strings.Cut(expr, "=")
	name = strings.TrimSpace(name)
	if !ok || name == "" {
		return errors.New("energy env: -w requires key=value or path=value")
	}
	root, err := loadJSONFile(configFile)
	if err != nil {
		return err
	}
	targetPath := name
	if !isJSONPath(name) {
		matches := findKey(root, name)
		sortMatches(matches)
		if len(matches) == 0 {
			return fmt.Errorf("energy env: key not found: %s", name)
		}
		if len(matches) > 1 {
			targetPath, err = selectMatch(matches, in, out)
			if err != nil {
				return err
			}
		} else {
			targetPath = matches[0].path
		}
	} else if _, ok := getPath(root, name); !ok {
		matches, indexed := findIndexedKey(root, name)
		if !indexed || len(matches) == 0 {
			return fmt.Errorf("energy env: json path not found: %s", name)
		}
		sortMatches(matches)
		if len(matches) > 1 {
			targetPath, err = selectMatch(matches, in, out)
			if err != nil {
				return err
			}
		} else {
			targetPath = matches[0].path
		}
	}
	if err = setPath(root, targetPath, value); err != nil {
		return err
	}
	return saveJSONFile(configFile, root)
}

func loadJSONFile(configFile string) (any, error) {
	data, err := os.ReadFile(configFile)
	if err != nil {
		return nil, err
	}
	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.UseNumber()
	root, err := decodeOrderedValue(decoder)
	if err != nil {
		return nil, err
	}
	return root, nil
}

func saveJSONFile(configFile string, root any) error {
	var data bytes.Buffer
	writePrettyJSON(&data, root, 0)
	data.WriteByte('\n')
	mode := os.FileMode(0644)
	if info, statErr := os.Stat(configFile); statErr == nil {
		mode = info.Mode()
	}
	return os.WriteFile(configFile, data.Bytes(), mode)
}

func printJSON(out io.Writer, value any) error {
	var data bytes.Buffer
	writePrettyJSON(&data, value, 0)
	_, err := fmt.Fprintln(out, data.String())
	return err
}

func findKey(root any, key string) []jsonMatch {
	var matches []jsonMatch
	var walk func(any, string)
	walk = func(node any, path string) {
		switch value := node.(type) {
		case *orderedObject:
			for _, member := range value.members {
				nextPath := joinPath(path, member.key)
				if member.key == key {
					matches = append(matches, jsonMatch{path: nextPath, value: member.value})
				}
				walk(member.value, nextPath)
			}
		case []any:
			for i, item := range value {
				walk(item, fmt.Sprintf("%s[%d]", path, i))
			}
		}
	}
	walk(root, "")
	return matches
}

func findIndexedKey(root any, query string) ([]jsonMatch, bool) {
	baseKey, suffix, ok := splitIndexedKey(query)
	if !ok {
		return nil, false
	}
	baseMatches := findKey(root, baseKey)
	matches := make([]jsonMatch, 0, len(baseMatches))
	for _, match := range baseMatches {
		value := match.value
		matched := true
		for _, token := range suffix {
			next, ok := applyToken(value, token)
			if !ok {
				matched = false
				break
			}
			value = next
		}
		if matched {
			matches = append(matches, jsonMatch{path: match.path + indexedSuffixPath(suffix), value: value})
		}
	}
	return matches, true
}

func splitIndexedKey(query string) (string, []jsonPathToken, bool) {
	if strings.Contains(query, ".") {
		return "", nil, false
	}
	indexStart := strings.Index(query, "[")
	if indexStart <= 0 {
		return "", nil, false
	}
	baseKey := query[:indexStart]
	suffix, err := parsePathPart(query[indexStart:], query)
	if err != nil {
		return "", nil, false
	}
	for _, token := range suffix {
		if !token.hasIndex || token.hasKey {
			return "", nil, false
		}
	}
	return baseKey, suffix, true
}

func indexedSuffixPath(tokens []jsonPathToken) string {
	var result strings.Builder
	for _, token := range tokens {
		result.WriteString("[")
		result.WriteString(strconv.Itoa(token.index))
		result.WriteString("]")
	}
	return result.String()
}

func joinPath(parent, key string) string {
	if parent == "" {
		return key
	}
	return parent + "." + key
}

func sortMatches(matches []jsonMatch) {
	sort.Slice(matches, func(i, j int) bool {
		leftDepth := pathDepth(matches[i].path)
		rightDepth := pathDepth(matches[j].path)
		if leftDepth != rightDepth {
			return leftDepth < rightDepth
		}
		if len(matches[i].path) != len(matches[j].path) {
			return len(matches[i].path) < len(matches[j].path)
		}
		return matches[i].path < matches[j].path
	})
}

func pathDepth(path string) int {
	if path == "" {
		return 0
	}
	depth := 1
	for _, r := range path {
		if r == '.' || r == '[' {
			depth++
		}
	}
	return depth
}

func selectMatch(matches []jsonMatch, in io.Reader, out io.Writer) (string, error) {
	fmt.Fprintln(out, "Multiple keys matched:")
	for i, match := range matches {
		fmt.Fprintf(out, "%d. %s %s %s\n", i+1, match.path, valueType(match.value), formatValue(match.value))
	}
	fmt.Fprint(out, "Select path (0 to cancel): ")
	reader := bufio.NewReader(in)
	line, err := reader.ReadString('\n')
	if err != nil && !errors.Is(err, io.EOF) {
		return "", err
	}
	line = strings.TrimSpace(line)
	if line == "" || line == "0" {
		return "", errors.New("energy env: write canceled")
	}
	idx, err := strconv.Atoi(line)
	if err != nil || idx < 1 || idx > len(matches) {
		return "", errors.New("energy env: invalid selection")
	}
	return matches[idx-1].path, nil
}

func getPath(root any, path string) (any, bool) {
	tokens, err := parsePath(path)
	if err != nil {
		return nil, false
	}
	current := root
	for _, token := range tokens {
		next, ok := applyToken(current, token)
		if !ok {
			return nil, false
		}
		current = next
	}
	return current, true
}

func setPath(root any, path, rawValue string) error {
	tokens, err := parsePath(path)
	if err != nil {
		return err
	}
	if len(tokens) == 0 {
		return errors.New("energy env: empty json path")
	}
	parent := root
	for _, token := range tokens[:len(tokens)-1] {
		next, ok := applyToken(parent, token)
		if !ok {
			return fmt.Errorf("energy env: json path not found: %s", path)
		}
		parent = next
	}
	last := tokens[len(tokens)-1]
	switch node := parent.(type) {
	case *orderedObject:
		if !last.hasKey || last.hasIndex {
			return fmt.Errorf("energy env: json path not found: %s", path)
		}
		member := node.find(last.key)
		if member == nil {
			return fmt.Errorf("energy env: json path not found: %s", path)
		}
		oldValue := member.value
		newValue, err := convertValue(oldValue, rawValue, path)
		if err != nil {
			return err
		}
		member.value = newValue
		return nil
	case []any:
		if !last.hasIndex || last.index < 0 || last.index >= len(node) {
			return fmt.Errorf("energy env: json path not found: %s", path)
		}
		oldValue := node[last.index]
		newValue, err := convertValue(oldValue, rawValue, path)
		if err != nil {
			return err
		}
		node[last.index] = newValue
		return nil
	default:
		return fmt.Errorf("energy env: json path not found: %s", path)
	}
}

func (m *orderedObject) find(key string) *orderedMember {
	if m == nil {
		return nil
	}
	for i := range m.members {
		if m.members[i].key == key {
			return &m.members[i]
		}
	}
	return nil
}

func applyToken(node any, token jsonPathToken) (any, bool) {
	if token.hasKey {
		object, ok := node.(*orderedObject)
		if !ok {
			return nil, false
		}
		member := object.find(token.key)
		if member == nil {
			return nil, false
		}
		return member.value, true
	}
	if token.hasIndex {
		array, ok := node.([]any)
		if !ok || token.index < 0 || token.index >= len(array) {
			return nil, false
		}
		return array[token.index], true
	}
	return nil, false
}

func parsePath(path string) ([]jsonPathToken, error) {
	path = strings.TrimSpace(path)
	if path == "" {
		return nil, errors.New("energy env: empty json path")
	}
	parts := strings.Split(path, ".")
	tokens := make([]jsonPathToken, 0, len(parts))
	for _, part := range parts {
		if part == "" {
			return nil, fmt.Errorf("energy env: invalid json path: %s", path)
		}
		parsed, err := parsePathPart(part, path)
		if err != nil {
			return nil, err
		}
		tokens = append(tokens, parsed...)
	}
	return tokens, nil
}

func parsePathPart(part, fullPath string) ([]jsonPathToken, error) {
	var tokens []jsonPathToken
	nameEnd := strings.Index(part, "[")
	if nameEnd == -1 {
		return []jsonPathToken{{key: part, hasKey: true}}, nil
	}
	if nameEnd > 0 {
		tokens = append(tokens, jsonPathToken{key: part[:nameEnd], hasKey: true})
	}
	rest := part[nameEnd:]
	for rest != "" {
		if !strings.HasPrefix(rest, "[") {
			return nil, fmt.Errorf("energy env: invalid json path: %s", fullPath)
		}
		end := strings.Index(rest, "]")
		if end <= 1 {
			return nil, fmt.Errorf("energy env: invalid json path: %s", fullPath)
		}
		idx, err := strconv.Atoi(rest[1:end])
		if err != nil || idx < 0 {
			return nil, fmt.Errorf("energy env: invalid json path: %s", fullPath)
		}
		tokens = append(tokens, jsonPathToken{index: idx, hasIndex: true})
		rest = rest[end+1:]
	}
	if len(tokens) == 0 {
		return nil, fmt.Errorf("energy env: invalid json path: %s", fullPath)
	}
	return tokens, nil
}

func isJSONPath(name string) bool {
	return strings.ContainsAny(name, ".[")
}

func convertValue(oldValue any, rawValue, path string) (any, error) {
	switch oldValue.(type) {
	case string:
		return rawValue, nil
	case bool:
		value := strings.ToLower(strings.TrimSpace(rawValue))
		if value == "true" {
			return true, nil
		}
		if value == "false" {
			return false, nil
		}
		return nil, fmt.Errorf("energy env: %s requires bool value", path)
	case json.Number:
		return convertNumber(oldValue.(json.Number), rawValue, path)
	case *orderedObject, []any:
		return nil, fmt.Errorf("energy env: %s is %s and cannot be written as a scalar value", path, valueType(oldValue))
	default:
		return nil, fmt.Errorf("energy env: %s has unsupported type %s", path, valueType(oldValue))
	}
}

func convertNumber(oldValue json.Number, rawValue, path string) (json.Number, error) {
	value := strings.TrimSpace(rawValue)
	if value == "" {
		return "", fmt.Errorf("energy env: %s requires number value", path)
	}
	if isIntegerNumber(oldValue.String()) {
		num, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return "", fmt.Errorf("energy env: %s requires integer value", path)
		}
		return json.Number(strconv.FormatInt(num, 10)), nil
	}
	num, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return "", fmt.Errorf("energy env: %s requires number value", path)
	}
	return json.Number(strconv.FormatFloat(num, 'f', -1, 64)), nil
}

func isIntegerNumber(value string) bool {
	return !strings.ContainsAny(value, ".eE")
}

func formatValue(value any) string {
	switch typed := value.(type) {
	case string:
		return typed
	case json.Number:
		return typed.String()
	case bool:
		return strconv.FormatBool(typed)
	case nil:
		return "null"
	default:
		var data bytes.Buffer
		writeCompactJSON(&data, typed)
		return data.String()
	}
}

func valueType(value any) string {
	switch value.(type) {
	case string:
		return "string"
	case json.Number:
		return "number"
	case bool:
		return "bool"
	case []any:
		return "array"
	case *orderedObject:
		return "object"
	case nil:
		return "null"
	default:
		return fmt.Sprintf("%T", value)
	}
}

func decodeOrderedValue(decoder *json.Decoder) (any, error) {
	token, err := decoder.Token()
	if err != nil {
		return nil, err
	}
	switch typed := token.(type) {
	case json.Delim:
		switch typed {
		case '{':
			object := &orderedObject{}
			for decoder.More() {
				keyToken, err := decoder.Token()
				if err != nil {
					return nil, err
				}
				key, ok := keyToken.(string)
				if !ok {
					return nil, errors.New("energy env: invalid json object key")
				}
				value, err := decodeOrderedValue(decoder)
				if err != nil {
					return nil, err
				}
				object.members = append(object.members, orderedMember{key: key, value: value})
			}
			if _, err = decoder.Token(); err != nil {
				return nil, err
			}
			return object, nil
		case '[':
			var array []any
			for decoder.More() {
				value, err := decodeOrderedValue(decoder)
				if err != nil {
					return nil, err
				}
				array = append(array, value)
			}
			if _, err = decoder.Token(); err != nil {
				return nil, err
			}
			return array, nil
		default:
			return nil, errors.New("energy env: invalid json delimiter")
		}
	default:
		return typed, nil
	}
}

func writePrettyJSON(out *bytes.Buffer, value any, depth int) {
	switch typed := value.(type) {
	case *orderedObject:
		if typed == nil || len(typed.members) == 0 {
			out.WriteString("{}")
			return
		}
		out.WriteString("{\n")
		for i, member := range typed.members {
			writeIndent(out, depth+1)
			writeJSONString(out, member.key)
			out.WriteString(": ")
			writePrettyJSON(out, member.value, depth+1)
			if i < len(typed.members)-1 {
				out.WriteByte(',')
			}
			out.WriteByte('\n')
		}
		writeIndent(out, depth)
		out.WriteByte('}')
	case []any:
		if len(typed) == 0 {
			out.WriteString("[]")
			return
		}
		out.WriteString("[\n")
		for i, item := range typed {
			writeIndent(out, depth+1)
			writePrettyJSON(out, item, depth+1)
			if i < len(typed)-1 {
				out.WriteByte(',')
			}
			out.WriteByte('\n')
		}
		writeIndent(out, depth)
		out.WriteByte(']')
	default:
		writeScalarJSON(out, typed)
	}
}

func writeCompactJSON(out *bytes.Buffer, value any) {
	switch typed := value.(type) {
	case *orderedObject:
		out.WriteByte('{')
		if typed != nil {
			for i, member := range typed.members {
				if i > 0 {
					out.WriteByte(',')
				}
				writeJSONString(out, member.key)
				out.WriteByte(':')
				writeCompactJSON(out, member.value)
			}
		}
		out.WriteByte('}')
	case []any:
		out.WriteByte('[')
		for i, item := range typed {
			if i > 0 {
				out.WriteByte(',')
			}
			writeCompactJSON(out, item)
		}
		out.WriteByte(']')
	default:
		writeScalarJSON(out, typed)
	}
}

func writeScalarJSON(out *bytes.Buffer, value any) {
	switch typed := value.(type) {
	case string:
		writeJSONString(out, typed)
	case json.Number:
		out.WriteString(typed.String())
	case bool:
		out.WriteString(strconv.FormatBool(typed))
	case nil:
		out.WriteString("null")
	default:
		writeJSONString(out, fmt.Sprint(typed))
	}
}

func writeJSONString(out *bytes.Buffer, value string) {
	data, err := json.Marshal(value)
	if err != nil {
		out.WriteString(`""`)
		return
	}
	out.Write(data)
}

func writeIndent(out *bytes.Buffer, depth int) {
	for i := 0; i < depth; i++ {
		out.WriteByte('\t')
	}
}
