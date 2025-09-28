package utils

import (
	"testing"
)

func TestJWT(t *testing.T) {

	// generate
	at, err := GenerateAccessToken(17)
	if err != nil {
		t.Errorf("GenerateAccessToken err %s", err)
	}
	t.Log(at)

	// parse
	atClaims, err := ParseAccessToken(at)
	if err != nil {
		t.Errorf("ParseAccessToken err %s", err)
	}
	subjectStr, err := atClaims.GetSubject()
	if err != nil {
		t.Errorf("GetSubject err %s", err)
	}

	iss, ok := atClaims["iss"].(string)
	if ok {
		t.Logf("iss: %s", iss)
	}

	t.Log(subjectStr)
}
