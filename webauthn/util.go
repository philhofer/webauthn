package webauthn

import (
	"bytes"
	"slices"
)

func isByteArrayInSlice(needle []byte, haystack ...[]byte) (valid bool) {
	return slices.ContainsFunc(haystack, func(b []byte) bool { return bytes.Equal(b, needle) })
}

func isCredentialsAllowedMatchingOwned(allowedCredentialIDs [][]byte, credentials []Credential) (valid bool) {
	var credential Credential

allowed:
	for _, allowedCredentialID := range allowedCredentialIDs {
		for _, credential = range credentials {
			if bytes.Equal(credential.ID, allowedCredentialID) {
				continue allowed
			}
		}

		return false
	}

	return true
}
