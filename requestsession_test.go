package sessions_test

import (
	"github.com/d2g/sessions"
	"net/http"
	"testing"
)

func TestRequestSessions(t *testing.T) {

	r := http.Request{}
	s, err := sessions.NewDefaultSession()
	if err != nil {
		t.Fatalf("Error: creating new session:%s\n", err.Error())
	}

	id, err := s.ID()
	if err != nil {
		t.Fatalf("Error: creating getting session id:%s\n", err.Error())
	}

	rs := sessions.RequestSession{
		Request: &r,
		Session: s,
	}

	rss := sessions.RequestSessions{}
	rss.Set(rs)

	rs1 := rss.Get(&r)
	id1, err := rs1.Session.ID()
	if err != nil {
		t.Fatalf("Error: creating getting session id:%s\n", err.Error())
	}

	if id != id1 {
		t.Fatalf("Error: Ids don't match expected %s got %s\n", id, id1)
	}

	//Set the request session taht already exists (update).
	rss.Set(rs)

	rss.Delete(&r)

	if rss.Get(&r).Request != nil || rss.Get(&r).Request == &r {
		t.Fatalf("Error: session should have been removed from cache.\n")
	}

	//Delete a session thats already been deleted.
	rss.Delete(&r)

}
