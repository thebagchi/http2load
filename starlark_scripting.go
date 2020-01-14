package main

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"sort"
	"strconv"
	"strings"

	"go.starlark.net/starlark"
	"go.starlark.net/starlarkstruct"
)

// Module Names
const (
	JSONModule = "json.star"
	CSVModule  = "csv.star"
)

// LoadJSON loads json module in starlark script
func LoadJSON() *starlarkstruct.Struct {
	return starlarkstruct.FromStringDict(
		starlarkstruct.Default,
		starlark.StringDict{
			"ToJSON":   starlark.NewBuiltin("json", ToJSON),
			"FromJSON": starlark.NewBuiltin("object", FromJSON),
		},
	)
}

// LoadCSV loads csv module in starlark script
func LoadCSV() *starlarkstruct.Struct {
	return starlarkstruct.FromStringDict(
		starlarkstruct.Default,
		starlark.StringDict{
			"ToCSV":   starlark.NewBuiltin("csv", ToCSV),
			"FromCSV": starlark.NewBuiltin("object", FromCSV),
		},
	)
}

// AsString starlark value to string
func AsString(input starlark.Value) string {
	value, err := strconv.Unquote(input.String())
	if nil != err {
		return ""
	}
	return value
}

// Unmarshal starlark value to interface
func Unmarshal(v starlark.Value) (interface{}, error) {
	switch v.Type() {
	case "NoneType":
		return nil, nil
	case "bool":
		return v.Truth() == starlark.True, nil
	case "int":
		return starlark.AsInt32(v)
	case "float":
		if float, ok := starlark.AsFloat(v); ok {
			return float, nil
		} else {
			return nil, fmt.Errorf("couldn't parse float")
		}
	case "string":
		return strconv.Unquote(v.String())
	case "dict":
		if dict, ok := v.(*starlark.Dict); ok {
			var values = map[string]interface{}{}
			for _, key := range dict.Keys() {
				value, _, err := dict.Get(key)
				if err != nil {
					return nil, err
				}
				temp, err := Unmarshal(value)
				if err != nil {
					return nil, err
				}
				values[AsString(key)] = temp
			}
			return values, nil
		} else {
			return nil, fmt.Errorf("error parsing dict. invalid type: %v", v)
		}
	case "list":
		if list, ok := v.(*starlark.List); ok {
			var element starlark.Value
			var iterator = list.Iterate()
			var value = make([]interface{}, 0)
			for iterator.Next(&element) {
				temp, err := Unmarshal(element)
				if err != nil {
					return nil, err
				}
				value = append(value, temp)
			}
			iterator.Done()
			return value, nil
		} else {
			return nil, fmt.Errorf("error parsing list. invalid type: %v", v)
		}
	case "tuple":
		if tuple, ok := v.(starlark.Tuple); ok {
			var element starlark.Value
			var iterator = tuple.Iterate()
			var value = make([]interface{}, 0)
			for iterator.Next(&element) {
				temp, err := Unmarshal(element)
				if err != nil {
					return nil, err
				}
				value = append(value, temp)
			}
			iterator.Done()
			return value, nil
		} else {
			return nil, fmt.Errorf("error parsing dict. invalid type: %v", v)
		}
	case "set":
		return nil, fmt.Errorf("sets aren't yet supported")
	default:
		return nil, fmt.Errorf("unrecognized starlark type: %s", v.Type())
	}
}

// Marshal interface to starlark value
func Marshal(v interface{}) (starlark.Value, error) {
	switch x := v.(type) {
	case nil:
		return starlark.None, nil
	case bool:
		return starlark.Bool(x), nil
	case string:
		return starlark.String(x), nil
	case int:
		return starlark.MakeInt(x), nil
	case float64:
		return starlark.Float(x), nil
	case []interface{}:
		var elements = make([]starlark.Value, 0)
		for _, value := range x {
			element, err := Marshal(value)
			if err != nil {
				return nil, err
			}
			elements = append(elements, element)
		}
		return starlark.NewList(elements), nil
	case map[string]interface{}:
		dict := &starlark.Dict{}
		for key, value := range x {
			element, err := Marshal(value)
			if err != nil {
				return nil, err
			}
			if err = dict.SetKey(starlark.String(key), element); err != nil {
				return nil, err
			}
		}
		return dict, nil
	default:
		return nil, fmt.Errorf("unknown type %T", v)
	}
}

// CSVString creates csv string from interface
func CSVString(v interface{}) string {
	values := make([]string, 0)
	switch x := v.(type) {
	case nil:
		return "nil"
	case int:
		return strconv.FormatInt(int64(x), 10)
	case float64:
		return strconv.FormatFloat(x, 'f', -1, 64)
	case bool:
		return strconv.FormatBool(x)
	case string:
		return x
	case []interface{}:
		for _, value := range x {
			values = append(values, CSVString(value))
		}
	case map[string]interface{}:
		for key, value := range x {
			values = append(values, key, CSVString(value))
		}
	default:
		return ""
	}
	buffer := &bytes.Buffer{}
	writer := csv.NewWriter(buffer)
	_ = writer.Write(values)
	writer.Flush()
	return strings.TrimSpace(buffer.String())
}

// ToCSV creates csv string from starlark value
func ToCSV(
	thread *starlark.Thread,
	builtin *starlark.Builtin,
	args starlark.Tuple,
	kwargs []starlark.Tuple,
) (starlark.Value, error) {
	var value starlark.Value
	if err := starlark.UnpackArgs("csv", args, kwargs, "value", &value); err != nil {
		return starlark.None, err
	}
	native, err := Unmarshal(value)
	if nil != err {
		return starlark.None, err
	}
	return starlark.String(CSVString(native)), nil
}

// ToJSON creates json string from starlark value
func ToJSON(
	thread *starlark.Thread,
	builtin *starlark.Builtin,
	args starlark.Tuple,
	kwargs []starlark.Tuple,
) (starlark.Value, error) {
	var value starlark.Value
	if err := starlark.UnpackArgs("json", args, kwargs, "value", &value); err != nil {
		return starlark.None, err
	}
	native, err := Unmarshal(value)
	if nil != err {
		return starlark.None, err
	}
	data, err := json.Marshal(native)
	if nil != err {
		return starlark.None, err
	}
	return starlark.String(data), nil
}

// FromCSV creates starlark value from csv string
func FromCSV(
	thread *starlark.Thread,
	builtin *starlark.Builtin,
	args starlark.Tuple,
	kwargs []starlark.Tuple,
) (starlark.Value, error) {
	return starlark.None, nil
}

// FromJSON creates starlark value from json string
func FromJSON(
	thread *starlark.Thread,
	builtin *starlark.Builtin,
	args starlark.Tuple,
	kwargs []starlark.Tuple,
) (starlark.Value, error) {
	var content starlark.String
	err := starlark.UnpackArgs("add", args, kwargs, "content", &content)
	if nil != err {
		return starlark.None, err
	}
	var value interface{}
	err = json.Unmarshal([]byte(AsString(content)), &value)
	if nil != err {
		return starlark.None, err
	}
	return Marshal(value)
}

// LoadScript loads a starlark script from file
func LoadScript(filename string) (int, string) {
	repeat := 0
	requests := ""
	thread := &starlark.Thread{
		Load: loader,
	}
	arguments := starlark.StringDict{}
	response, err := starlark.ExecFile(thread, filename, nil, arguments)
	if nil != err {
		fmt.Println("Error: ", err)
	} else {
		var names []string
		for name := range response {
			names = append(names, name)
		}
		sort.Strings(names)
		for _, name := range names {
			v := response[name]
			if strings.Compare(name, "repeat") == 0 {
				value, err := strconv.Atoi(v.String())
				if nil == err {
					repeat = value
				}
			}
			if strings.Compare(name, "requests") == 0 {
				requests = AsString(v)
			}
		}
	}
	return repeat, requests
}

func loader(thread *starlark.Thread, module string) (starlark.StringDict, error) {
	switch module {
	case JSONModule:
		return starlark.StringDict{
			"json": LoadJSON(),
		}, nil
	case CSVModule:
		return starlark.StringDict{
			"csv": LoadCSV(),
		}, nil
	}
	return nil, fmt.Errorf("invalid module")
}
