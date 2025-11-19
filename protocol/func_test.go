package protocol

import (
	"errors"
	"regexp"
	"strings"
	"testing"
)

func AssertIsProtocolError(t *testing.T, err error, errType, errDetails, errInfo any) {
	t.Helper()
	var e *Error
	if !errors.As(err, &e) {
		t.Fatalf("%T cannot be turned into an *Error", err)
	}

	switch et := errType.(type) {
	case string:
		musteq(t, et, e.Type)
	case *regexp.Regexp:
		if !et.MatchString(e.Type) {
			t.Fatalf("type %q doesn't match", e.Type)
		}
	default:
		t.Fatalf("%T is not a known type", errType)
	}

	switch ed := errDetails.(type) {
	case string:
		musteq(t, strings.ToLower(ed), strings.ToLower(e.Details))
	case *regexp.Regexp:
		if !ed.MatchString(e.Details) {
			t.Fatalf("details %q doesn't match", e.Details)
		}
	default:
		t.Fatalf("%T is not a known type", errDetails)
	}

	switch ed := errInfo.(type) {
	case string:
		musteq(t, ed, e.DevInfo)
	case *regexp.Regexp:
		if !ed.MatchString(e.DevInfo) {
			t.Fatalf("devinfo %q doesn't match", e.DevInfo)
		}
	default:
		t.Fatalf("%T is not a known type", errInfo)
	}
}
