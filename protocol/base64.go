package protocol

import (
	"bytes"
	"encoding/base64"
	"fmt"
)

// URLEncodedBase64 represents a byte slice holding URL-encoded base64 data.
// When fields of this type are unmarshalled from JSON, the data is base64
// decoded into a byte slice.
type URLEncodedBase64 []byte

func (e URLEncodedBase64) String() string {
	return base64.RawURLEncoding.EncodeToString(e)
}

// UnmarshalJSON base64 decodes a URL-encoded value, storing the result in the
// provided byte slice.
func (e *URLEncodedBase64) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		*e = nil
		return nil
	}
	// Trim the leading and trailing quotes from raw JSON data (the whole value part).
	if len(data) < 2 || data[0] != '"' || data[len(data)-1] != '"' {
		return fmt.Errorf("%q cannot be unmarshaled as base64-encoded data", data)
	}
	data = data[1 : len(data)-1]

	// FIXME(philhofer): is this correct?
	// I don't understand why the upstream library
	// is trimming padding characters.
	//
	// Trim the trailing equal characters.
	data = bytes.TrimRight(data, "=")

	out := make([]byte, base64.RawURLEncoding.DecodedLen(len(data)))
	n, err := base64.RawURLEncoding.Decode(out, data)
	if err != nil {
		return err
	}
	*e = out[:n]
	return nil
}

// MarshalJSON base64 encodes a non URL-encoded value, storing the result in the
// provided byte slice.
func (e URLEncodedBase64) MarshalJSON() ([]byte, error) {
	if e == nil {
		return []byte("null"), nil
	}
	dst := make([]byte, 0, base64.RawURLEncoding.EncodedLen(len(e))+2)
	dst = append(dst, '"')
	dst = base64.RawURLEncoding.AppendEncode(dst, e)
	dst = append(dst, '"')
	return dst, nil
}
