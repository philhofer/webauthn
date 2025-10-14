package webauthn

import (
	"encoding/json"
	"testing"

	"github.com/philhofer/webauthn/protocol"
)

func noerr(t *testing.T, err error) {
	t.Helper()
	if err != nil {
		t.Fatal(err)
	}
}

func nonnil[T any](t *testing.T, v *T) {
	if v == nil {
		t.Helper()
		t.Fatal("non-nil")
	}
}

func errequal(t *testing.T, err error, text string) {
	if err.Error() != text {
		t.Helper()
		t.Fatalf("got %q want %q", err.Error(), text)
	}
}

func musteq[T comparable](t *testing.T, want, got T) {
	if want != got {
		t.Helper()
		t.Fatalf("got %+v want +%v", got, want)
	}
}

func TestWithRegistrationRelyingPartyID(t *testing.T) {
	testCases := []struct {
		name         string
		have         *Config
		opts         []RegistrationOption
		expectedID   string
		expectedName string
		err          string
	}{
		{
			name: "OptionDefinedInConfig",
			have: &Config{
				RPID:          "https://example.com",
				RPDisplayName: "Test Display Name",
				RPOrigins:     []string{"https://example.com"},
			},
			opts:         nil,
			expectedID:   "https://example.com",
			expectedName: "Test Display Name",
		},
		{
			name: "OptionDefinedInConfigAndOpts",
			have: &Config{
				RPID:          "https://example.com",
				RPDisplayName: "Test Display Name",
				RPOrigins:     []string{"https://example.com"},
			},
			opts:         []RegistrationOption{WithRegistrationRelyingPartyID("https://a.example.com"), WithRegistrationRelyingPartyName("Test Display Name2")},
			expectedID:   "https://a.example.com",
			expectedName: "Test Display Name2",
		},
		{
			name: "OptionDefinedInConfigWithNoErrAndInOptsWithError",
			have: &Config{
				RPID:          "https://example.com",
				RPDisplayName: "Test Display Name",
				RPOrigins:     []string{"https://example.com"},
			},
			opts: []RegistrationOption{WithRegistrationRelyingPartyID("---::~!!~@#M!@OIK#N!@IOK@@@@@@@@@@"), WithRegistrationRelyingPartyName("Test Display Name2")},
			err:  "error generating credential creation: the relying party id failed to validate as it's not a valid uri with error: parse \"---::~!!~@\": first path segment in URL cannot contain colon",
		},
		{
			name: "OptionDefinedInOpts",
			have: &Config{
				RPOrigins: []string{"https://example.com"},
			},
			opts:         []RegistrationOption{WithRegistrationRelyingPartyID("https://example.com"), WithRegistrationRelyingPartyName("Test Display Name")},
			expectedID:   "https://example.com",
			expectedName: "Test Display Name",
		},
		{
			name: "OptionDisplayNameNotDefined",
			have: &Config{
				RPOrigins: []string{"https://example.com"},
			},
			opts: []RegistrationOption{WithRegistrationRelyingPartyID("https://example.com")},
			err:  "error generating credential creation: the relying party display name must be provided via the configuration or a functional option for a creation",
		},
		{
			name: "OptionIDNotDefined",
			have: &Config{
				RPOrigins: []string{"https://example.com"},
			},
			opts: []RegistrationOption{WithRegistrationRelyingPartyName("Test Display Name")},
			err:  "error generating credential creation: the relying party id must be provided via the configuration or a functional option for a creation",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			w, err := New(tc.have)
			noerr(t, err)

			user := &defaultUser{}

			creation, _, err := w.BeginRegistration(user, tc.opts...)
			if tc.err != "" {
				errequal(t, err, tc.err)
			} else {
				noerr(t, err)
				nonnil(t, creation)
				musteq(t, tc.expectedID, creation.Response.RelyingParty.ID)
				musteq(t, tc.expectedName, creation.Response.RelyingParty.Name)
			}
		})
	}
}

func TestRegistration_FinishRegistrationFailure(t *testing.T) {
	user := &defaultUser{
		id: []byte("123"),
	}

	session := SessionData{
		UserID: []byte("ABC"),
	}

	webauthn := &WebAuthn{}

	credential, err := webauthn.FinishRegistration(user, session, nil)
	if err == nil {
		t.Errorf("FinishRegistration() error = nil, want %v", protocol.ErrBadRequest.Type)
	}

	if credential != nil {
		t.Errorf("FinishRegistration() credential = %v, want nil", credential)
	}
}

func TestEntityEncoding(t *testing.T) {
	testCases := []struct {
		name           string
		b64            bool
		have, expected string
	}{
		{"ShouldEncodeBase64", true, "abc", `{"name":"","displayName":"","id":"YWJj"}`},
		{"ShouldEncodeString", false, "abc", `{"name":"","displayName":"","id":"abc"}`},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			entityUser := protocol.UserEntity{}

			if tc.b64 {
				entityUser.ID = protocol.URLEncodedBase64(tc.have)
			} else {
				entityUser.ID = tc.have
			}

			data, err := json.Marshal(entityUser)
			noerr(t, err)
			musteq(t, tc.expected, string(data))
		})
	}
}
