package sessions_test

import (
	"errors"
	"github.com/d2g/sessions"
	"net/http"
	"testing"
)

const TESTSESSIONID = "MNJ3PZ34RYBOBNEKSWILFMAXTBDMPVYEZLMVVR5TXVBTIQHHJISA===="

type MockStore struct{}

func (t *MockStore) Get(id string) (sessions.Session, error) {
	if id == "ERROR" {
		return nil, errors.New("Mock Get Error")
	} else {
		s, err := sessions.NewDefaultSession()
		if err != nil {
			return nil, err
		}

		if id == TESTSESSIONID {
			err = s.GobDecode([]byte{44, 255, 129, 3, 1, 2, 255, 130, 0, 1, 3, 1, 2, 73, 68, 1, 12, 0, 1, 7, 69, 120, 112, 105, 114, 101, 115, 1, 255, 132, 0, 1, 6, 86, 97, 108, 117, 101, 115, 1, 255, 134, 0, 0, 0, 16, 255, 131, 5, 1, 1, 4, 84, 105, 109, 101, 1, 255, 132, 0, 0, 0, 45, 255, 133, 4, 1, 1, 29, 109, 97, 112, 91, 105, 110, 116, 101, 114, 102, 97, 99, 101, 32, 123, 125, 93, 105, 110, 116, 101, 114, 102, 97, 99, 101, 32, 123, 125, 1, 255, 134, 0, 1, 16, 1, 16, 0, 0, 93, 255, 130, 1, 56, 77, 78, 74, 51, 80, 90, 51, 52, 82, 89, 66, 79, 66, 78, 69, 75, 83, 87, 73, 76, 70, 77, 65, 88, 84, 66, 68, 77, 80, 86, 89, 69, 90, 76, 77, 86, 86, 82, 53, 84, 88, 86, 66, 84, 73, 81, 72, 72, 74, 73, 83, 65, 61, 61, 61, 61, 2, 1, 6, 115, 116, 114, 105, 110, 103, 12, 5, 0, 3, 75, 101, 121, 6, 115, 116, 114, 105, 110, 103, 12, 7, 0, 5, 86, 97, 108, 117, 101, 0})
			if err != nil {
				return nil, err
			}
		}

		return s, nil
	}
}

func (t *MockStore) Set(s sessions.Session) error {
	return nil
}

func (t *MockStore) Delete(id string) error {
	return nil
}

func (t *MockStore) All() ([]sessions.Session, error) {
	return []sessions.Session{}, nil
}

func TestSessionInfo(t *testing.T) {
	si := sessions.SessionInfo{}
	si.Cookie.Name = "SESSIONID"
	si.Store = &MockStore{}
	si.Cache = sessions.RequestSessions{make([]sessions.RequestSession, 0)}

	r, err := http.NewRequest("GET", "http://example.com", nil)
	if err != nil {
		t.Fatalf("Error: creating dummy request:%s\n", err.Error())
	}

	id, err := si.GetSessionID(r)
	if err != nil {
		t.Fatalf("Error: getting session id:%s\n", err.Error())
	}

	if id != "" {
		t.Fatalf("Error: didn't expect session id got %s\n", id)
	}

	c := http.Cookie{
		Name:  si.Cookie.Name,
		Value: TESTSESSIONID,
	}
	r.AddCookie(&c)

	id, err = si.GetSessionID(r)
	if err != nil {
		t.Fatalf("Error: getting session id:%s\n", err.Error())
	}

	if id != c.Value {
		t.Fatalf("Error: expected session id %s got %s\n", id, c.Value)
	}

	s, err := si.GetSession(r)
	if err != nil {
		t.Fatalf("Error: getting session %s\n", err.Error())
	}

	id, err = s.ID()
	if err != nil {
		t.Fatalf("Error: getting session id for returned session:%s\n", err.Error())
	}

	if id != TESTSESSIONID {
		t.Fatalf("Error: wrong session expected %s got %s\n", TESTSESSIONID, id)
	}

	// The session we just got from the store should now be cached.
	rs := si.Cache.Get(r)
	id, err = rs.Session.ID()
	if err != nil {
		t.Fatalf("Error: getting session id for cached session:%s\n", err.Error())
	}

	if id != TESTSESSIONID {
		t.Fatalf("Error: wrong session expected %s got %s\n", TESTSESSIONID, id)
	}

	si.ClearCache(r)

	// Get the session from the cache but this time it should be nil
	rs = si.Cache.Get(r)
	if rs.Request != nil {
		t.Fatalf("Error: Cache should be clear.\n")
	}

	si.SetSession(r, s)

	// The session we just set to the store should now be cached.
	rs = si.Cache.Get(r)
	id, err = rs.Session.ID()
	if err != nil {
		t.Fatalf("Error: getting session id for cached session:%s\n", err.Error())
	}

	if id != TESTSESSIONID {
		t.Fatalf("Error: wrong session expected %s got %s\n", TESTSESSIONID, id)
	}

	//Persist the session to the store.
	err = si.PersistSession(r)
	if err != nil {
		t.Fatalf("Error: persisting session to store:%s\n", err.Error())
	}

}
