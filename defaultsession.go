package sessions

import (
	"bytes"
	"crypto/rand"
	"encoding/base32"
	"encoding/gob"
	"io"
	"time"
)

type defaultSession struct {
	// The Session ID (The key for the record in the datastore)
	id string

	// Expiary Date
	// This is when the data stored in this session expires and therefore should be GC
	expires time.Time

	// Stuff in the session
	values map[interface{}]interface{}
}

func NewDefaultSession() (*defaultSession, error) {
	session := new(defaultSession)
	session.values = make(map[interface{}]interface{})
	_, err := session.ID()
	return session, err
}

func (t *defaultSession) ID() (string, error) {
	if t.id == "" {
		k := make([]byte, 32)
		if _, err := io.ReadFull(rand.Reader, k); err != nil {
			return "", err
		}
		t.id = base32.StdEncoding.EncodeToString(k)
	}
	return t.id, nil
}

func (t *defaultSession) Expiry() time.Time {
	return t.expires
}

func (t *defaultSession) SetExpiry(i time.Time) {
	t.expires = i
}

func (t *defaultSession) Set(key, object interface{}) error {
	t.values[key] = object
	return nil
}

func (t *defaultSession) Get(key interface{}) (interface{}, error) {
	return t.values[key], nil
}

func (t *defaultSession) Delete(key interface{}) error {
	delete(t.values, key)
	return nil
}

func (t *defaultSession) Keys() ([]interface{}, error) {
	var keys []interface{}
	for key := range t.values {
		keys = append(keys, key)
	}
	return keys, nil
}

func (t *defaultSession) Purge() error {
	for key := range t.values {
		delete(t.values, key)
	}
	return nil
}

func (t *defaultSession) GobEncode() ([]byte, error) {
	encoded := struct {
		ID      string
		Expires time.Time
		Values  map[interface{}]interface{}
	}{
		t.id,
		t.expires,
		t.values,
	}

	//If the ID hasn't be encoded
	if encoded.ID == "" {
		var err error
		encoded.ID, err = t.ID()
		if err != nil {
			return []byte{}, err
		}
	}

	buf := new(bytes.Buffer)
	enc := gob.NewEncoder(buf)
	if err := enc.Encode(encoded); err != nil {
		return []byte{}, err
	}

	return buf.Bytes(), nil
}

func (t *defaultSession) GobDecode(data []byte) error {
	decoded := struct {
		ID      string
		Expires time.Time
		Values  map[interface{}]interface{}
	}{}

	dec := gob.NewDecoder(bytes.NewBuffer(data))
	if err := dec.Decode(&decoded); err != nil {
		return err
	}

	t.id = decoded.ID
	t.expires = decoded.Expires
	t.values = decoded.Values
	return nil
}
