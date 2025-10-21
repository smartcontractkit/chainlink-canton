// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package token

import (
	"github.com/smartcontractkit/chainlink-canton-internal/bindings/da/time"
	"github.com/smartcontractkit/chainlink-canton-internal/codegen/daml"
	stdtime "time"
)

type ChoiceExecutionMetadata struct {
	Meta Metadata
}

type ExtraArgs struct {
	Context ChoiceContext
	Meta    Metadata
}

type Metadata struct {
	Values map[string]string
}

type ChoiceContext struct {
	Values map[string]AnyValue
}

type AnyContractView struct {
}

// Variant AnyValue
// Types that are valid to be assigned to AnyValue:
//
//	AV_Bool
//	AV_ContractId
//	AV_Date
//	AV_Decimal
//	AV_Int
//	AV_List
//	AV_Map
//	AV_Party
//	AV_RelTime
//	AV_Text
//	AV_Time
type AnyValue interface {
	_isAnyValue()
}

type AV_Bool bool

func (v AV_Bool) _isAnyValue() {}

type AV_ContractId daml.ContractId

func (v AV_ContractId) _isAnyValue() {}

type AV_Date stdtime.Time

func (v AV_Date) _isAnyValue() {}

type AV_Decimal string

func (v AV_Decimal) _isAnyValue() {}

type AV_Int int64

func (v AV_Int) _isAnyValue() {}

type AV_List []AnyValue

func (v AV_List) _isAnyValue() {}

type AV_Map map[string]AnyValue

func (v AV_Map) _isAnyValue() {}

type AV_Party daml.Party

func (v AV_Party) _isAnyValue() {}

type AV_RelTime time.RelTime

func (v AV_RelTime) _isAnyValue() {}

type AV_Text string

func (v AV_Text) _isAnyValue() {}

type AV_Time stdtime.Time

func (v AV_Time) _isAnyValue() {}
