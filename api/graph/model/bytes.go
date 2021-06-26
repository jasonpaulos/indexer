package model

import (
	"encoding/base64"
	"fmt"
	"io"

	"github.com/99designs/gqlgen/graphql"
)

func MarshalBytes(b []byte) graphql.Marshaler {
	return graphql.WriterFunc(func(w io.Writer) {
		io.WriteString(w, fmt.Sprintf("\"%s\"", base64.StdEncoding.EncodeToString(b)))
	})
}

func UnmarshalBytes(v interface{}) ([]byte, error) {
	if v == nil {
		return []byte{}, nil
	}

	switch v := v.(type) {
	case string:
		return base64.StdEncoding.DecodeString(v)
	case []byte:
		return v, nil
	default:
		return nil, fmt.Errorf("%T is not a byte slice", v)
	}
}
