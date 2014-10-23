package sessions

import (
	"time"
)

type Session interface {

	// This returns the Session ID that can/will be stored in the clients cookie.
	// Even if the session is new a session id should be returned.
	ID() (string, error)

	// The Expiry time when the session times out.
	Expiry() time.Time

	// Set when the session timeout.
	SetExpiry(time.Time)

	// Set Session Object
	Set(key, object interface{}) error

	// Get Session Object
	Get(key interface{}) (interface{}, error)

	// Delete Session Object
	// We Shouldn't error if the session object no longer exists.
	Delete(key interface{}) error

	// Keys returns all the keys within the session.
	Keys() ([]interface{}, error)

	// Purge
	// Remove all values assigned with the sesion.
	Purge() error

	// Support Writing to Disc(k?)
	GobDecode([]byte) error

	// If you've got to write it you've got to read it.
	GobEncode() ([]byte, error)
}

type SessionStore interface {

	// Get A Session Based on this ID
	// If no Session can be found a new empty session should be created.
	Get(id string) (Session, error)

	// Save the current session to the store.
	Set(s Session) error

	// Delete the Session
	Delete(id string) error

	// All, list all sessions in the store.
	All() ([]Session, error)
}
