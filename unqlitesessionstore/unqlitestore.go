package unqlitesessionstore

import (
	"bytes"
	"encoding/gob"
	"github.com/d2g/sessions"
	"github.com/d2g/unqlitego"
	"log"
)

type unqliteStore struct {
	collection *unqlitego.Database
}

func New(filename string) (*unqliteStore, error) {
	store := new(unqliteStore)

	var err error
	store.collection, err = unqlitego.NewDatabase(filename)
	if err != nil {
		return nil, err
	}

	return store, nil
}

func (t *unqliteStore) Get(id string) (sessions.Session, error) {
	var err error

	s, err := sessions.NewDefaultSession()
	if err != nil {
		return s, err
	}

	if id != "" {
		byteobject, err := t.collection.Fetch([]byte(id))
		if err != nil {
			if err == unqlitego.UnQLiteError(-6) || err == unqlitego.UnQLiteError(-3) {
				//Not Found is not an error in my world...
				return s, nil
			}

			return s, err
		}

		dec := gob.NewDecoder(bytes.NewBuffer(byteobject))
		if err := dec.Decode(&s); err != nil {
			return s, err
		}

	}

	return s, err
}

func (t *unqliteStore) Set(s sessions.Session) error {
	sessionid, err := s.ID()
	if err != nil {
		return err
	}

	buf := new(bytes.Buffer)
	enc := gob.NewEncoder(buf)
	if err := enc.Encode(s); err != nil {
		return err
	}

	err = t.collection.Store([]byte(sessionid), buf.Bytes())
	if err != nil {
		log.Println("Error: " + err.Error())
		return err
	}

	return nil
}

func (t *unqliteStore) Delete(id string) error {
	err := t.collection.DeleteObject(id)

	if err != nil {
		log.Println("Error: Deleting Session from Datastore")
	}

	return err
}

func (t *unqliteStore) All() ([]sessions.Session, error) {

	s := make([]sessions.Session, 0, 0)

	cursor, err := t.collection.NewCursor()
	defer cursor.Close()

	if err != nil {
		return s, err
	}

	err = cursor.First()
	if err != nil {
		//You Get -28 When There are no records.
		if err == unqlitego.UnQLiteError(-28) {
			return s, nil
		} else {
			return s, err
		}
	}

	for {
		if !cursor.IsValid() {
			break
		}

		session, err := sessions.NewDefaultSession()
		if err != nil {
			return s, err
		}

		value, err := cursor.Value()

		if err != nil {
			log.Println("Error: Cursor Get Value Error:" + err.Error())
		} else {

			dec := gob.NewDecoder(bytes.NewBuffer(value))
			if err := dec.Decode(&session); err != nil {
				key, err := cursor.Key()
				if err != nil {
					log.Println("Error: Cursor Get Key Error:" + err.Error())
					cursor.Delete()
					continue
				} else {
					log.Println("Error: Invalid Session in Datastore:" + string(key))
					cursor.Delete()
					continue
				}
			}

			s = append(s, session)
		}

		err = cursor.Next()
		if err != nil {
			break
		}
	}

	err = cursor.Close()
	if err != nil {
		log.Println("Error Closing Sursor:" + err.Error())
	}
	return s, err
}
