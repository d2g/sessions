package boltsessionstore

import (
	"bytes"
	"encoding/gob"
	"github.com/boltdb/bolt"
	"github.com/d2g/sessions"
	"log"
)

type BoltStore struct {
	DB *bolt.DB
}

const (
	bucketname string = "sessions"
)

func (b *BoltStore) Get(id string) (sessions.Session, error) {
	var err error

	s, err := sessions.NewDefaultSession()
	if err != nil {
		return s, err
	}

	if id != "" {
		var bo []byte

		err = b.DB.View(func(tx *bolt.Tx) error {
			bkt := tx.Bucket([]byte(bucketname))
			if bkt == nil {
				return nil
			}

			bo = bkt.Get([]byte(id))
			return nil
		})
		if err != nil {
			return s, err
		}

		dec := gob.NewDecoder(bytes.NewBuffer(bo))
		if err := dec.Decode(&s); err != nil {
			return s, err
		}
	}

	return s, err
}

func (b *BoltStore) Set(s sessions.Session) error {
	sessionid, err := s.ID()
	if err != nil {
		return err
	}

	buf := new(bytes.Buffer)
	enc := gob.NewEncoder(buf)
	if err := enc.Encode(s); err != nil {
		return err
	}

	err = b.DB.Update(func(tx *bolt.Tx) error {
		bkt, err := tx.CreateBucketIfNotExists([]byte(bucketname))
		if err != nil {
			return err
		}

		err = bkt.Put([]byte(sessionid), buf.Bytes())
		return err
	})

	return err
}

func (b *BoltStore) Delete(id string) error {
	err := b.DB.Update(func(tx *bolt.Tx) error {
		bkt, err := tx.CreateBucketIfNotExists([]byte(bucketname))
		if err != nil {
			return err
		}
		err = bkt.Delete([]byte(id))
		return err
	})

	if err != nil {
		log.Println("Error: Deleting Session from Datastore")
	}

	return err
}

func (b *BoltStore) All() ([]sessions.Session, error) {
	s := make([]sessions.Session, 0, 0)

	err := b.DB.View(func(tx *bolt.Tx) error {
		bkt := tx.Bucket([]byte(bucketname))
		c := bkt.Cursor()

		for k, v := c.First(); k != nil; k, v = c.Next() {

			session, err := sessions.NewDefaultSession()
			if err != nil {
				return err
			}

			dec := gob.NewDecoder(bytes.NewBuffer(v))
			if err := dec.Decode(&session); err != nil {
				//Broken Session try and remove it.
				b.Delete(string(k))
				log.Panicf("Warning: Deleting Broken Session \"%s\"\n", string(v))
				continue
			}
			s = append(s, session)
		}

		return nil
	})

	return s, err
}
