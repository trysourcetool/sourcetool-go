package conv

import (
	"encoding/json"

	"github.com/samber/lo"
)

func SafeValue[T comparable](in *T) T {
	if in == nil {
		return lo.Empty[T]()
	}
	return *in
}

func NilValue[T comparable](in T) *T {
	if lo.IsEmpty(in) {
		return nil
	}
	return &in
}

func StringValue(in any) string {
	switch v := in.(type) {
	case string:
		return v
	default:
		res, err := json.Marshal(in)
		if err != nil {
			return ""
		}
		return string(res)
	}
}
