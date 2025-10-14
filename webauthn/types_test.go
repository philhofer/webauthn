package webauthn

import (
	"testing"

	"github.com/philhofer/webauthn/protocol"
)

type defaultUser struct {
	id          []byte
	credentials []Credential
}

var _ User = (*defaultUser)(nil)

func (user *defaultUser) WebAuthnID() []byte {
	return user.id
}

func (user *defaultUser) WebAuthnName() string {
	return "newUser"
}

func (user *defaultUser) WebAuthnDisplayName() string {
	return "New User"
}

func (user *defaultUser) WebAuthnCredentials() []Credential {
	return user.credentials
}

func TestNew(t *testing.T) {
	testCases := []struct {
		description string
		config      *Config
		err         string
	}{
		{
			"ShouldPassMinimalConfig",
			&Config{
				RPID:      "https://example.com/",
				RPOrigins: []string{"https://example.com"},
			},
			"",
		},
		{
			"ShouldFailBadRPID",
			&Config{
				RPID:      "%%&&",
				RPOrigins: []string{"https://example.com"},
			},
			`error occurred validating the configuration: field 'RPID' is not a valid URI: parse "%%&&": invalid URL escape "%%&"`,
		},
		{
			"ShouldFailNoRPOrigins",
			&Config{
				RPID: "https://example.com/",
			},
			"error occurred validating the configuration: must provide at least one value to the 'RPOrigins' field",
		},
		{
			"ShouldFailBadTopOrigins",
			&Config{
				RPID:                        "https://example.com/",
				RPOrigins:                   []string{"https://example.com"},
				RPTopOriginVerificationMode: protocol.TopOriginExplicitVerificationMode,
			},
			"error occurred validating the configuration: must provide at least one value to the 'RPTopOrigins' field when 'RPTopOriginVerificationMode' field is set to protocol.TopOriginExplicitVerificationMode",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			w, err := New(tc.config)

			if tc.err == "" {
				nonnil(t, w)
				noerr(t, err)
				noerr(t, tc.config.validate())
			} else {
				errequal(t, err, tc.err)
				if tc.config.validate() == nil {
					t.Fatal("no error?")
				}
			}
		})
	}
}
