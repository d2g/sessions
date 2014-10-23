package sessions

import (
	"net/http"
)

type RequestSession struct {
	Request *http.Request
	Session Session
}

type RequestSessions struct {
	RequestSessions []RequestSession
}

func (t *RequestSessions) Get(request *http.Request) RequestSession {
	for _, v := range t.RequestSessions {
		if v.Request == request {
			return v
		}
	}

	return RequestSession{}
}

func (t *RequestSessions) Set(s RequestSession) {
	for i := range t.RequestSessions {
		if t.RequestSessions[i].Request == s.Request {
			t.RequestSessions[i].Session = s.Session
			return
		}
	}
	t.RequestSessions = append(t.RequestSessions, s)
}

func (t *RequestSessions) Delete(request *http.Request) {
	for i := range t.RequestSessions {
		if t.RequestSessions[i].Request == request {
			t.RequestSessions[i] = t.RequestSessions[len(t.RequestSessions)-1]
			t.RequestSessions = t.RequestSessions[:len(t.RequestSessions)-1]
			return
		}
	}

	return
}
