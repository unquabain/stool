package main

import (
	"encoding/json"
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"
	"unicode"

	"golang.org/x/exp/constraints"

	yaml "gopkg.in/yaml.v3"
)

type PathChunkType int

const (
	PCTEmpty PathChunkType = iota
	PCTDot
	PCTBrace
	PCTIndex
	PCTMember
	PCTFunction
	PCTStar
)

type Path []rune

func NewPath(s string) *Path {
	r := []rune(s)
	return (*Path)(&r)
}

func (p *Path) String() string {
	return string(*p)
}

func (p *Path) RuneArray() []rune {
	return []rune(*p)
}

func scanForBrace(s []rune) int {
	state := make([]rune, 0)
	var pos int
	for {
		if len(state) == 0 {
			switch s[0] {
			case ']':
				return pos
			case '\'':
				state = append([]rune{'\''}, state...)
			case '"':
				state = append([]rune{'"'}, state...)
			case '(':
				state = append([]rune{')'}, state...)
			case '[':
				state = append([]rune{']'}, state...)
			}
		} else {
			if s[0] == state[0] {
				state = state[1:]
			}
		}
		pos++
		s = s[1:]
		if len(s) == 0 {
			return pos
		}
	}
}

func scanMember(s []rune) (int, bool) {
	allDigits := true
	for i := 0; i < len(s); i++ {
		r := s[i]
		if unicode.IsDigit(r) {
			continue
		}
		if unicode.IsLetter(r) {
			allDigits = false
			continue
		}
		if r == '_' || r == '-' {
			allDigits = false
			continue
		}
		return i, allDigits
	}
	return len(s), allDigits
}

func isAllDigits(s []rune) bool {
	for _, r := range s {
		if !unicode.IsDigit(r) {
			return false
		}
	}
	return true
}

func (p *Path) chunk() (string, PathChunkType) {
	s := p.RuneArray()
	if len(s) == 0 {
		return ``, PCTEmpty
	}
	switch s[0] {
	case '.':
		*p = Path(s[1:])
		return `.`, PCTDot
	case '[':
		s = s[1:]
		clen := scanForBrace(s)
		chunk := string(s[:clen])
		chunk = strings.TrimSpace(chunk)
		*p = Path(s[clen+1:])
		if chunk == `*` {
			return chunk, PCTStar
		}
		if strings.HasPrefix(chunk, `'`) && strings.HasSuffix(chunk, `'`) {
			return strings.Trim(chunk, `'`), PCTMember
		}
		if strings.HasPrefix(chunk, `"`) && strings.HasSuffix(chunk, `"`) {
			return strings.Trim(chunk, `"`), PCTMember
		}
		if isAllDigits([]rune(chunk)) {
			return chunk, PCTIndex
		}
		return chunk, PCTBrace
	default:
		mlen, allDigits := scanMember(s)
		chunk := string(s[:mlen])
		*p = Path(s[mlen:])
		s = p.RuneArray()
		if len(s) != 0 && s[0] == '(' {
			*p = Path([]rune(strings.TrimLeft(string(s), `()`)))
			return chunk, PCTFunction
		}
		if allDigits {
			return chunk, PCTIndex
		}
		return chunk, PCTMember
	}

}

func evalIndex(data []any, index int) []any {
	var part int
	for _, value := range data {
		switch array := value.(type) {
		case []any:
			if len(array) <= index {
				continue
			}
			data[part] = array[index]
			part++
		case map[string]any:
			if len(array) <= index {
				continue
			}
			data[part] = array[fmt.Sprint(index)]
			part++
		case map[any]any:
			if len(array) <= index {
				continue
			}
			data[part] = array[index]
			part++
		}
	}
	return data[:part]
}

func evalMember(data []any, member string) []any {
	var (
		part  int
		value any
		ok    bool
	)
	for _, item := range data {
		var dict map[string]any
		dict, ok = item.(map[string]any)
		if ok {
			value, ok = dict[member]
			if !ok {
				continue
			}
		} else {
			var dict map[any]any
			dict, ok = item.(map[any]any)
			value, ok = dict[member]
			if !ok {
				continue
			}
		}
		data[part] = value
		part++
	}
	return data[:part]
}

func evalStar(data []any) []any {
	out := make([]any, 0)
	for _, item := range data {
		switch v := item.(type) {
		case []any:
			out = append(out, v...)
		case map[string]any:
			for _, vi := range v {
				out = append(out, vi)
			}
		default:
			continue
		}
	}
	return out
}

func evalFuncLen(data []any) []any {
	for idx, item := range data {
		switch v := item.(type) {
		case []any:
			data[idx] = len(v)
		case map[string]any:
			data[idx] = len(v)
		case string:
			data[idx] = len(v)
		default:
			data[idx] = 1
		}
	}
	return data
}

func evalFuncJSON(data []any) ([]any, error) {
	for idx, item := range data {
		bytes, err := json.Marshal(item)
		if err != nil {
			return nil, fmt.Errorf(`could not marshal item %d as JSON: %w`, idx, err)
		}
		data[idx] = string(bytes)
	}
	return data, nil
}

func evalFuncJSONPretty(data []any) ([]any, error) {
	for idx, item := range data {
		bytes, err := json.MarshalIndent(item, ``, `    `)
		if err != nil {
			return nil, fmt.Errorf(`could not marshal item %d as JSON-Pretty: %w`, idx, err)
		}
		data[idx] = string(bytes)
	}
	return data, nil
}

func evalFuncYAML(data []any) ([]any, error) {
	for idx, item := range data {
		bytes, err := yaml.Marshal(item)
		if err != nil {
			return nil, fmt.Errorf(`could not marshal item %d as YAML: %w`, idx, err)
		}
		data[idx] = string(bytes)
	}
	return data, nil
}

func evalFuncKeys(data []any) []any {
	var part int
	for _, item := range data {
		switch dict := item.(type) {
		case map[string]any:
			keys := make([]any, 0, len(dict))
			for key := range dict {
				keys = append(keys, key)
			}
			data[part] = keys
			part++
		case map[any]any:
			keys := make([]any, 0, len(dict))
			for key := range dict {
				key = append(keys, fmt.Sprint(key))
			}
			data[part] = keys
			part++
		}
	}
	return data[:part]
}

func evalFuncFlatten(data []any) []any {
	var out = make([]any, 0)
	for _, item := range data {
		switch value := item.(type) {
		case []any:
			out = append(out, value...)
		case map[string]any:
			for _, v := range value {
				out = append(out, v)
			}
		case map[any]any:
			for _, v := range value {
				out = append(out, v)
			}
		default:
			out = append(out, value)
		}
	}
	return out
}

func evalFuncResults(data []any) []any {
	return []any{data}
}

func evalFuncJSONEval(data []any) ([]any, error) {
	var part int
	for rnum, item := range data {
		switch result := item.(type) {
		case string:
			var value any
			err := json.Unmarshal([]byte(result), &value)
			if err != nil {
				log.Print(result)
				return nil, fmt.Errorf(`could not unmarshal result %d as JSON: %w`, rnum, err)
			}
			data[part] = value
			part++
		}
	}
	return data[:part], nil
}

func evalFuncYAMLEval(data []any) ([]any, error) {
	var part int
	for rnum, item := range data {
		switch result := item.(type) {
		case string:
			var value any
			err := yaml.Unmarshal([]byte(result), &value)
			if err != nil {
				return nil, fmt.Errorf(`could not unmarshal result %d as YAML: %w`, rnum, err)
			}
			data[part] = value
			part++
		}
	}
	return data[:part], nil
}

func evalFunction(data []any, function string) ([]any, error) {
	switch function {
	case `len`, `length`:
		return evalFuncLen(data), nil
	case `json`, `js`:
		return evalFuncJSON(data)
	case `jsonpretty`, `jspretty`, `jpretty`:
		return evalFuncJSONPretty(data)
	case `jsoneval`, `jeval`:
		return evalFuncJSONEval(data)
	case `yaml`, `yml`:
		return evalFuncYAML(data)
	case `yamleval`, `yeval`:
		return evalFuncYAMLEval(data)
	case `keys`:
		return evalFuncKeys(data), nil
	case `flatten`, `flat`:
		return evalFuncFlatten(data), nil
	case `results`:
		return evalFuncResults(data), nil
	default:
		return nil, fmt.Errorf(`unknown function`)
	}
}

func typedCompare[T constraints.Ordered](l, r T, comparison string) bool {
	switch strings.TrimSpace(comparison) {
	case `<`:
		return l < r
	case `<=`:
		return l <= r
	case `>`:
		return l > r
	case `>=`:
		return l >= r
	case `==`:
		return l == r
	case `!=`:
		return l != r
	default:
		return false
	}
}

func compare(lval any, rval string, comparison string) bool {
	switch lv := lval.(type) {
	case int:
		rv, err := strconv.Atoi(rval)
		if err != nil {
			return false
		}
		return typedCompare(lv, rv, comparison)
	case string:
		rval, err := strconv.Unquote(rval)
		if err != nil {
			return false
		}
		return typedCompare(lv, rval, comparison)
	case float64:
		rv, err := strconv.ParseFloat(rval, 64)
		if err != nil {
			return false
		}
		return typedCompare(lv, rv, comparison)
	default:
		return false
	}
}

func evalBrace(data []any, expression string) ([]any, error) {
	matches := regexp.MustCompile(`^(.*?)(<=?|>=?|[!=]=)(.*)$`).FindStringSubmatch(expression)
	if matches == nil {
		return nil, fmt.Errorf(`don't know how to interpret %q`, expression)
	}
	lpath := strings.TrimSpace(matches[1])
	comparison := strings.TrimSpace(matches[2])
	rval := strings.TrimSpace(matches[3])
	out := make([]any, 0)
	for idx, item := range data {
		switch subitems := item.(type) {
		case []any:
		arrayloop:
			for _, subitem := range subitems {
				lmatches, err := Evaluate(subitem, lpath)
				if err != nil {
					return nil, fmt.Errorf(`could not evaluate lpath %q for array item %d: %w`, lpath, idx, err)
				}
				if len(lmatches) == 0 {
					continue
				}
				for _, lval := range lmatches {
					if compare(lval, rval, comparison) {
						out = append(out, subitem)
						continue arrayloop
					}
				}
			}
		case map[string]any:
		stringloop:
			for _, subitem := range subitems {
				lmatches, err := Evaluate(subitem, lpath)
				if err != nil {
					return nil, fmt.Errorf(`could not evaluate lpath %q for dict item %d: %w`, lpath, idx, err)
				}
				if len(lmatches) == 0 {
					continue
				}
				for _, lval := range lmatches {
					if compare(lval, rval, comparison) {
						out = append(out, subitem)
						continue stringloop
					}
				}
			}
		case map[any]any:
		anyloop:
			for _, subitem := range subitems {
				lmatches, err := Evaluate(subitem, lpath)
				if err != nil {
					return nil, fmt.Errorf(`could not evaluate lpath %q for map item %d: %w`, lpath, idx, err)
				}
				if len(lmatches) == 0 {
					continue
				}
				for _, lval := range lmatches {
					if compare(lval, rval, comparison) {
						out = append(out, subitem)
						continue anyloop
					}
				}
			}
		}
	}
	return out, nil
}

func Evaluate(data any, path string) ([]any, error) {
	p := NewPath(path)
	results := []any{data}
	for {
		if len(results) == 0 {
			return results, nil
		}
		chunk, chunkType := p.chunk()
		switch chunkType {
		case PCTEmpty:
			return results, nil
		case PCTDot:
			continue
		case PCTIndex:
			idx, err := strconv.Atoi(chunk)
			if err != nil {
				return nil, fmt.Errorf(`incorrectly interpreted %q as an index: %w`, chunk, err)
			}
			results = evalIndex(results, idx)
		case PCTMember:
			results = evalMember(results, chunk)
		case PCTStar:
			results = evalStar(results)
		case PCTBrace:
			var err error
			results, err = evalBrace(results, chunk)
			if err != nil {
				return nil, fmt.Errorf(`unable to evaluate expression in brace %q: %w`, chunk, err)
			}
		case PCTFunction:
			var err error
			results, err = evalFunction(results, chunk)
			if err != nil {
				return nil, fmt.Errorf(`unable to evaluate function %q: %w`, chunk, err)
			}
		default:
			return nil, fmt.Errorf(`unknown chunk type %q pulled from path %q`, chunk, p.String())
		}
	}
}
