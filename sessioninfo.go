package sessions

import (
	"log"
	"net/http"
	"time"
)

type SessionInfo struct {
	Cookie struct {
		Name   string
		Path   string
		Domain string
	}
	Timeout time.Duration
	Store   SessionStore
	Cache   RequestSessions
}

// Get the Session Id From the current Request.
func (t *SessionInfo) GetSessionID(request *http.Request) (string, error) {

	cookie, err := request.Cookie(t.Cookie.Name)
	if err != nil {
		if err == http.ErrNoCookie {
			return "", nil
		} else {
			return "", err
		}
	}

	return cookie.Value, nil
}

// Get The session, try from cache then fallback to store.
func (t *SessionInfo) GetSession(request *http.Request) (Session, error) {
	cachedSession := t.Cache.Get(request)

	if cachedSession.Request != nil {
		return cachedSession.Session, nil
	}

	sessionid, err := t.GetSessionID(request)
	if err != nil {
		log.Printf("Debug: Error Getting Session ID For Request: %s\n", err.Error())
		return nil, err
	}

	session, err := t.Store.Get(sessionid)
	if err != nil {
		log.Printf("Debug: Error Session For Session ID \"%s\" because: %s\n", sessionid, err.Error())
		return nil, err
	}

	t.Cache.Set(RequestSession{
		Request: request,
		Session: session,
	})

	return session, nil
}

// Try and Set the session id in the browsers cookie.
func (t *SessionInfo) SetSessionCookie(response http.ResponseWriter, s Session) error {
	//Cookies.....
	cookie := &http.Cookie{
		Name: t.Cookie.Name,
	}

	if t.Cookie.Path != "" {
		cookie.Path = t.Cookie.Path
	} else {
		cookie.Path = "/"
	}

	if t.Cookie.Domain != "" {
		cookie.Domain = t.Cookie.Domain
	}

	sessionkeys, err := s.Keys()
	if err != nil {
		return err
	}

	//If The Session Is Empty
	if len(sessionkeys) <= 0 || s.Expiry().Before(time.Now()) {
		//Expire the Cookie.
		cookie.Value = ""
		cookie.Expires = time.Now()
		cookie.MaxAge = -1

	} else {

		//Expire the Cookie.
		cookie.Value, err = s.ID()
		if err != nil {
			return err
		}
		cookie.Expires = s.Expiry()

		if int64(^uint(0)>>1) < int64(s.Expiry().Sub(time.Now()).Seconds()) {
			cookie.MaxAge = int(^uint(0) >> 1)
		} else {
			cookie.MaxAge = int(s.Expiry().Sub(time.Now()).Seconds())
		}

	}

	//Send the cookie back
	http.SetCookie(response, cookie)
	return nil
}

// Set Session in the Request Cache.
func (t *SessionInfo) SetSession(request *http.Request, s Session) {
	r := RequestSession{
		Request: request,
		Session: s,
	}
	t.Cache.Set(r)
}

// Persist session to underlying store.
func (t *SessionInfo) PersistSession(request *http.Request) error {
	session, err := t.GetSession(request)
	if err != nil {
		return err
	}

	return t.Store.Set(session)
}

// Clear the Request From the Active Cache. This should always be done when getsession has been called.
func (t *SessionInfo) ClearCache(request *http.Request) {
	t.Cache.Delete(request)
}

// Wrapper function to allow http.handler chaining.
func (t *SessionInfo) GetHandler(c http.Handler) http.Handler {
	wrapper := sessionInfoHandler{
		c,
		t,
	}

	return &wrapper
}

func (t *SessionInfo) SaveSession(w http.ResponseWriter, r *http.Request) {
	//Did we use a session.
	cache := t.Cache.Get(r)
	if cache.Request != nil {
		//Yep we did.

		//Do we have anything in the session?
		keys, err := cache.Session.Keys()
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}

		if len(keys) > 0 {
			//Increase the session expiry.
			cache.Session.SetExpiry(time.Now().Add(t.Timeout))

			//Save the updated session to cache.
			t.SetSession(r, cache.Session)

			//Store the session to disk.
			err := t.PersistSession(r)
			if err != nil {
				http.Error(w, err.Error(), 500)
				return
			}

			//Save the session
			//This is needed when start the session for the first time.
			err = t.SetSessionCookie(w, cache.Session)
			if err != nil {
				http.Error(w, err.Error(), 500)
				return
			}
		}
	}
}

type sessionInfoHandler struct {
	http.Handler
	*SessionInfo
}

func (t sessionInfoHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	//Cookies have to be updated before we write back to the client.
	//This means this requests is the update from the last session.
	session, err := t.GetSession(r)
	if err == nil {
		//Tell the client to update it's cookie.
		err = t.SetSessionCookie(w, session)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
	}

	//Call the inner servehttp.
	t.Handler.ServeHTTP(w, r)

	//Save the session
	t.SessionInfo.SaveSession(w, r)

	//Always clear our cache to free up resource.
	t.SessionInfo.ClearCache(r)
}
