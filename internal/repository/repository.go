package repository

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"time"

	"github.com/google/uuid"
	"go.etcd.io/bbolt"
	bolt "go.etcd.io/bbolt"
)

type Invitation struct {
	Token   uuid.UUID // Invitation token uuid - Key
	Email   string    // User Email
	Created time.Time
}

type Peer struct {
	PeerKey string // Hash Email+Device - Key
	Email   string // User Email
	Device  string // User Device Name
	Created time.Time
}

const bucketInvitations = `Invites`

type Repository struct {
	db *bbolt.DB
}

// NewRepository Creates a new repository to a specified file
func NewRepository(dbFile string) (r *Repository, err error) {
	db, err := bolt.Open(dbFile, 0600, nil)
	if err != nil {
		return nil, err
	}

	return &Repository{
		db: db,
	}, nil
}

// Close will close the repository including the database
func (r *Repository) Close() error {
	return r.db.Close()
}

// AddInvitation adds a new Invitation to the repostitory
func (r *Repository) AddInvitation(i Invitation) error {
	return r.db.Update(func(tx *bolt.Tx) error {
		bt, err := tx.CreateBucketIfNotExists([]byte(bucketInvitations))
		if err != nil {
			return err
		}

		val, err := encodeValue(i)
		if err != nil {
			return err
		}

		return bt.Put([]byte(i.Token.String()), val)
	})
}

// GetInvitation adds a new Invitation to the repostitory
func (r *Repository) GetInvitation(token uuid.UUID) (i *Invitation, err error) {
	return i, r.db.View(func(tx *bolt.Tx) error {
		bt := tx.Bucket([]byte(bucketInvitations))
		if bt == nil {
			return fmt.Errorf("bucket %s not found", bucketInvitations)
		}

		val := bt.Get([]byte(token.String()))
		if val == nil {
			return fmt.Errorf("invitation for token %s not found", token)
		}
		return decodeValue(&i, val)
	})
}

// RevokeInvitation adds a new Invitation to the repostitory
func (r *Repository) RevokeInvitation(token uuid.UUID) error {
	return r.db.Update(func(tx *bolt.Tx) error {
		bt := tx.Bucket([]byte(bucketInvitations))
		if bt == nil {
			return fmt.Errorf("bucket %s not found", bucketInvitations)
		}

		return bt.Delete([]byte(token.String()))
	})
}

func encodeValue(v interface{}) (val []byte, err error) {
	var b bytes.Buffer
	enc := gob.NewEncoder(&b)
	err = enc.Encode(v)
	if err != nil {
		return nil, err
	}
	return b.Bytes(), nil
}

func decodeValue(v interface{}, b []byte) error {
	bf := bytes.NewBuffer(b)
	dec := gob.NewDecoder(bf)
	err := dec.Decode(v)
	return err
}
