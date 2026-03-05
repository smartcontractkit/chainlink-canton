package common

import (
	"encoding/json"
	"fmt"
	"strings"

	eABI "github.com/ethereum/go-ethereum/accounts/abi"
)

type EntryOutput struct {
	Indexed bool   `json:"indexed,omitempty"`
	Name    string `json:"name,omitempty"`
	Type    string `json:"type,omitempty"`
}

type EntryInput struct {
	Indexed bool   `json:"indexed,omitempty"`
	Name    string `json:"name,omitempty"`
	Type    string `json:"type,omitempty"`
}

type Entry struct {
	Name            string        `json:"name,omitempty"`
	Anonymous       bool          `json:"anonymous,omitempty"`
	Constant        bool          `json:"constant,omitempty"`
	Payable         bool          `json:"payable,omitempty"`
	StateMutability string        `json:"stateMutability,omitempty"`
	Type            string        `json:"type,omitempty"`
	Inputs          []EntryInput  `json:"inputs"`
	Outputs         []EntryOutput `json:"outputs"`
}

type JSONABI struct {
	Entrys []Entry `json:"entrys,omitempty"`
}

func (abi *JSONABI) GetFunctionSignature(name string) (string, error) {
	for _, entry := range abi.Entrys {
		if entry.Name == name {
			var types = make([]string, len(entry.Inputs))
			for i, input := range entry.Inputs {
				types[i] = input.Type
			}
			return fmt.Sprintf("%v(%v)", name, strings.Join(types, ",")), nil
		}
	}
	return "", fmt.Errorf("entry with name %v not found in abi", name)
}

func (abi *JSONABI) GetInputParser(method string) (eABI.Arguments, error) {
	arguments := eABI.Arguments{}
	for _, entry := range abi.Entrys {
		if entry.Name == method {
			for _, out := range entry.Inputs {
				ty, err := eABI.NewType(out.Type, "", nil)
				if err != nil {
					return nil, fmt.Errorf("invalid param %s: %+v", out.Type, err)
				}
				arguments = append(arguments, eABI.Argument{
					Name:    out.Name,
					Type:    ty,
					Indexed: out.Indexed,
				})
			}
			return arguments, nil
		}
	}
	return nil, fmt.Errorf("not found")
}

func (abi *JSONABI) GetOutputParser(method string) (eABI.Arguments, error) {
	arguments := eABI.Arguments{}
	for _, entry := range abi.Entrys {
		if entry.Name == method {
			for _, out := range entry.Outputs {
				ty, err := eABI.NewType(out.Type, "", nil)
				if err != nil {
					return nil, fmt.Errorf("invalid param %s: %+v", out.Type, err)
				}
				arguments = append(arguments, eABI.Argument{
					Name:    out.Name,
					Type:    ty,
					Indexed: out.Indexed,
				})
			}
			return arguments, nil
		}
	}
	return nil, fmt.Errorf("not found")
}

func LoadJSONABI(jsonStr string) (*JSONABI, error) {
	var entries []Entry
	err := json.Unmarshal([]byte(jsonStr), &entries)
	if err != nil {
		return nil, fmt.Errorf("failed to parse ABI JSON: %w", err)
	}
	return &JSONABI{Entrys: entries}, nil
}
