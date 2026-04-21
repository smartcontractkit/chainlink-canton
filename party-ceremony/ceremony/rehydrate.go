package ceremony

import (
	"bytes"
	"encoding/json"
)

// Rehydrate converts an any value (typically map[string]any after a JSON
// round-trip through the Reporter) back into a concrete type T. It uses
// UseNumber to prevent integer fields from coercing to float64.
func Rehydrate[T any](v any) (T, bool) {
	b, err := json.Marshal(v)
	if err != nil {
		var zero T
		return zero, false
	}
	var out T
	dec := json.NewDecoder(bytes.NewReader(b))
	dec.UseNumber()
	if err := dec.Decode(&out); err != nil {
		var zero T
		return zero, false
	}

	return out, true
}
