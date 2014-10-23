package sessions_test

import (
	"github.com/d2g/sessions"
	"testing"
	"time"
)

func TestDefaultSession(t *testing.T) {
	s, err := sessions.NewDefaultSession()
	if err != nil {
		t.Fatalf("Error: creating new session:%s\n", err.Error())
	}

	id, err := s.ID()
	if err != nil {
		t.Fatalf("Error: getting session id:%s\n", err.Error())
	}

	if id == "" {
		t.Fatalf("Error: Blank session ID generated\n")
	}

	if !s.Expiry().Equal(time.Time{}) {
		t.Fatalf("Error: Default time not as expected\n")
	}

	tn := time.Now()
	s.SetExpiry(tn)
	if !s.Expiry().Equal(tn) {
		t.Fatalf("Error: expiry got %v wanted %v\n", s.Expiry(), tn)
	}

	err = s.Set("Key", "Value")
	if err != nil {
		t.Fatalf("Error: setting value to session:%s\n", err.Error())
	}

	uv, err := s.Get("Key")
	if err != nil {
		t.Fatalf("Error: getting value from session:%s\n", err.Error())
	}

	v, ok := uv.(string)
	if !ok {
		t.Fatalf("Error: Unable to turn result back to string\n")
	}

	if v != "Value" {
		t.Fatalf("Error: expected \"Value\" received \"%s\"\n", v)
	}

	keys, err := s.Keys()
	if err != nil {
		t.Fatalf("Error: getting keys from session:%s\n", err.Error())
	}

	if len(keys) != 1 {
		t.Fatalf("Error: expected 1 key got :%s\n", len(keys))
	}

	k, ok := keys[0].(string)
	if !ok {
		t.Fatalf("Error: Unable to turn key back to string\n")
	}

	if k != "Key" {
		t.Fatalf("Error: expected \"Key\" received \"%s\"\n", k)
	}

	err = s.Delete("Key")
	if err != nil {
		t.Fatalf("Error: deleting key:\"%s\"\n", err.Error())
	}

	keys, err = s.Keys()
	if err != nil {
		t.Fatalf("Error: getting keys (after delete) from session:%s\n", err.Error())
	}

	if len(keys) != 0 {
		t.Fatalf("Error: expected 0 key got :%s\n", len(keys))
	}

	err = s.Set("Key", "Value")
	if err != nil {
		t.Fatalf("Error: setting value to session:%s\n", err.Error())
	}

	keys, err = s.Keys()
	if err != nil {
		t.Fatalf("Error: getting keys (prior to purge) from session:%s\n", err.Error())
	}

	if len(keys) != 1 {
		t.Fatalf("Error: (prior to purge) expected 1 key got :%s\n", len(keys))
	}

	err = s.Purge()
	if err != nil {
		t.Fatalf("Error: purging from session:%s\n", err.Error())
	}

	keys, err = s.Keys()
	if err != nil {
		t.Fatalf("Error: getting keys (after purge) from session:%s\n", err.Error())
	}

	if len(keys) != 0 {
		t.Fatalf("Error: (after purge) expected 0 key got :%s\n", len(keys))
	}

}

func TestDefaultSessionEncoding(t *testing.T) {
	s, err := sessions.NewDefaultSession()
	if err != nil {
		t.Fatalf("Error: creating new session:%s\n", err.Error())
	}

	err = s.Set("Key", "Value")
	if err != nil {
		t.Fatalf("Error: setting value to session:%s\n", err.Error())
	}

	b, err := s.GobEncode()
	if err != nil {
		t.Fatalf("Error: encoding session:%s\n", err.Error())
	}

	err = s.Purge()
	if err != nil {
		t.Fatalf("Error: purging session:%s\n", err.Error())
	}

	keys, err := s.Keys()
	if err != nil {
		t.Fatalf("Error: getting keys from session:%s\n", err.Error())
	}

	if len(keys) != 0 {
		t.Fatalf("Error: expected 0 key got :%s\n", len(keys))
	}

	err = s.GobDecode(b)
	if err != nil {
		t.Fatalf("Error: decoding session:%s\n", err.Error())
	}

	uv, err := s.Get("Key")
	if err != nil {
		t.Fatalf("Error: getting value from session:%s\n", err.Error())
	}

	v, ok := uv.(string)
	if !ok {
		t.Fatalf("Error: Unable to turn result back to string\n")
	}

	if v != "Value" {
		t.Fatalf("Error: expected \"Value\" received \"%s\"\n", v)
	}
}
