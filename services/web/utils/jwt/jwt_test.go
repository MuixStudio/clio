package jwt

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

	// GetSubject
	subjectStr, err := atClaims.GetSubject()
	if err != nil {
		t.Errorf("GetSubject err %s", err)
	}
	t.Logf("subjectStr %s", subjectStr)
	//
	//// GetExpirationTime
	//expirationTime, err := atClaims.GetExpirationTime()
	//if err != nil {
	//	t.Errorf("GetExpirationTime err %s", err)
	//}
	//t.Logf("expirationTime %v", expirationTime)
	//
	//if expirationTime.Time.Before(time.Now()) {
	//	t.Log("expired")
	//}
	//
	//iss, ok := atClaims["iss"].(string)
	//if ok {
	//	t.Logf("iss: %s", iss)
	//}

	userId, ok := atClaims["user_id"].(uint32)
	if ok {
		t.Logf("user_id: %d", userId)
	}

	//t.Log(subjectStr)
}
