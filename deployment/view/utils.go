package view

import (
	"strings"

	apiv2 "github.com/digital-asset/dazl-client/v8/go/api/com/daml/ledger/api/v2"
)

func NormalizeNumeric0Val(s string) string {
	if strings.HasSuffix(s, ".") {
		return strings.TrimSuffix(s, ".")
	}
	return s
}

func FieldVal(rec *apiv2.Record, label string) *apiv2.Value {
	for _, f := range rec.GetFields() {
		if f.GetLabel() == label {
			return f.GetValue()
		}
	}
	return nil
}

func PartyVal(v *apiv2.Value) string {
	if v == nil {
		return ""
	}
	if x, ok := v.Sum.(*apiv2.Value_Party); ok {
		return x.Party
	}
	return ""
}

func TextVal(v *apiv2.Value) string {
	if v == nil {
		return ""
	}
	if x, ok := v.Sum.(*apiv2.Value_Text); ok {
		return x.Text
	}
	return ""
}

func NumericVal(v *apiv2.Value) string {
	if v == nil {
		return ""
	}
	if x, ok := v.Sum.(*apiv2.Value_Numeric); ok {
		return x.Numeric
	} // keep "...."
	return ""
}

func BoolVal(v *apiv2.Value) bool {
	if v == nil {
		return false
	}
	if x, ok := v.Sum.(*apiv2.Value_Bool); ok {
		return x.Bool
	}
	return false
}

func TextListVal(v *apiv2.Value) []string {
	if v == nil {
		return nil
	}
	x, ok := v.Sum.(*apiv2.Value_List)
	if !ok {
		return nil
	}
	out := make([]string, 0, len(x.List.GetElements()))
	for _, e := range x.List.GetElements() {
		out = append(out, TextVal(e))
	}
	return out
}

func GenMapEntriesVal(v *apiv2.Value) []*apiv2.GenMap_Entry {
	if v == nil {
		return nil
	}
	if x, ok := v.Sum.(*apiv2.Value_GenMap); ok {
		return x.GenMap.GetEntries()
	}
	return nil
}

func RecordVal(v *apiv2.Value) *apiv2.Record {
	if v == nil {
		return nil
	}
	if x, ok := v.Sum.(*apiv2.Value_Record); ok {
		return x.Record
	}
	return nil
}
