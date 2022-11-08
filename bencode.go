package bencode

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

type Item struct {
	Str  string
	Num  int
	List []Item
	Dict map[string]Item
}

func isNum(in byte) bool {
	return in >= '0' && in <= '9'
}

func ParseToMap(in []byte) (Item, error) {
	rtn, _, err := parseToMapRemaining(in)
	return rtn, err
}
func parseToMapRemaining(in []byte) (Item, []byte, error) {
	out := Item{}
	if len(in) == 0 {
		return out, nil, nil
	}
	var err error
	var remaining []byte
	switch {
	case isNum(in[0]):
		out.Str, remaining, err = parseString(in)
	case in[0] == 'i':
		out.Num, remaining, err = parseInt(in)
	case in[0] == 'l':
		out.List, remaining, err = parseList(in[1:])
	case in[0] == 'd':
		out.Dict, remaining, err = parseDict(in[1:])
	default:
		return out, nil, nil
	}
	if err != nil {
		return out, nil, err
	}
	return out, remaining, nil
}
func parseDict(in []byte) (map[string]Item, []byte, error) {
	rtn := make(map[string]Item)

	for {
		if len(in) == 0 {
			return rtn, nil, fmt.Errorf("missing end of list")
		}
		if in[0] == 'e' { // this is normal exit
			return rtn, in[1:], nil
		}
		key, in1, err := parseToMapRemaining(in)
		if err != nil {
			return rtn, nil, fmt.Errorf("failed to parse dictionary key: %w", err)
		}
		if len(key.Str) == 0 {
			return rtn, nil, fmt.Errorf("dictionary key has zero length")
		}
		in = in1
		var value Item
		value, in1, err = parseToMapRemaining(in)
		if err != nil {
			return rtn, nil, fmt.Errorf("failed to parse dictionary value for key %s: %w", key.Str, err)
		}
		rtn[key.Str] = value
		in = in1
	}

}
func parseList(in []byte) ([]Item, []byte, error) {
	var rtn []Item
	for {
		if len(in) == 0 {
			return rtn, nil, fmt.Errorf("missing end of list")
		}
		if in[0] == 'e' { // this is normal exit
			return rtn, in[1:], nil
		}
		item, in1, err := parseToMapRemaining(in)
		if err != nil {
			return rtn, nil, fmt.Errorf("failed to parse list item: %w", err)
		}
		rtn = append(rtn, item)
		in = in1
	}
}
func parseInt(in []byte) (int, []byte, error) {
	sections := strings.SplitN(string(in[1:]), "e", 2)
	if len(sections) != 2 {
		return 0, nil, errors.New("integer format")
	}
	value, err := strconv.Atoi(sections[0])
	if err != nil {
		return 0, nil, fmt.Errorf("invalid integer: %w", err)
	}
	return value, []byte(sections[1]), nil
}
func parseString(in []byte) (string, []byte, error) {
	sections := strings.SplitN(string(in), ":", 2)
	if len(sections) != 2 {
		return "", nil, errors.New("invalid string format")
	}
	value, err := strconv.Atoi(sections[0])
	if err != nil {
		return "", nil, fmt.Errorf("string prefix number is not valid: %w", err)
	}
	if len(sections[1]) < value {
		return "", nil, errors.New("missing data needed for string")
	}
	return sections[1][:value], []byte(sections[1][value:]), nil
}
