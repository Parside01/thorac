package storage

import (
	"encoding/json"
	"fmt"
	"github.com/boltdb/bolt"
	"sync"
)

type boltStateStorage struct {
	unimplementedStateStorage
	db         *bolt.DB
	mutex      sync.RWMutex
	bucketName []byte

	keyCurrentTerm []byte
	keyVotedFor    []byte
}

func NewBoltStateStorage(dbPath string, bucketName string) (StateStorage, error) {
	db, err := bolt.Open(dbPath, 0600, nil)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	err = db.Update(func(tx *bolt.Tx) error {
		_, err = tx.CreateBucketIfNotExists([]byte(bucketName))
		if err != nil {
			return fmt.Errorf("failed to create bucket: %s", err.Error())
		}
		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to initialize bolt StateStorage: %s", err.Error())
	}

	return &boltStateStorage{
		db:         db,
		mutex:      sync.RWMutex{},
		bucketName: []byte(bucketName),
		// TODO: Перенести в конфиг.
		keyCurrentTerm: []byte("current_term"),
		keyVotedFor:    []byte("voted_for"),
	}, nil
}

func (s *boltStateStorage) GetTerm() (int, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	term := 0
	err := s.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(s.bucketName)
		if bucket == nil {
			return fmt.Errorf("no such bucket: %s", s.bucketName)
		}

		currentTermBytes := bucket.Get(s.keyCurrentTerm)
		if currentTermBytes == nil {
			term = 0
			return nil
		}

		if err := json.Unmarshal(currentTermBytes, &term); err != nil {
			return fmt.Errorf("failed to decode currentTerm: %s", err.Error())
		}
		return nil
	})
	return term, err
}

func (s *boltStateStorage) SetTerm(term int) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	return s.db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(s.bucketName)
		if bucket == nil {
			return fmt.Errorf("no such bucket: %s", s.bucketName)
		}

		termBytes, err := json.Marshal(term)
		if err != nil {
			return fmt.Errorf("failed to encode term: %s", err.Error())
		}

		if err := bucket.Put(s.keyCurrentTerm, termBytes); err != nil {
			return fmt.Errorf("failed to set current term: %s", err.Error())
		}
		return nil
	})
}

func (s *boltStateStorage) GetVotedFor() (int, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	votedFor := 0
	err := s.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(s.bucketName)
		if bucket == nil {
			return fmt.Errorf("no such bucket: %s", s.bucketName)
		}

		votedForBytes := bucket.Get(s.keyVotedFor)
		if votedForBytes == nil {
			votedFor = 0
			return nil
		}

		if err := json.Unmarshal(votedForBytes, &votedFor); err != nil {
			return fmt.Errorf("failed to decode voted for: %s", err.Error())
		}
		return nil
	})
	return votedFor, err
}

func (s *boltStateStorage) SetVotedFor(term int) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	return s.db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(s.bucketName)
		if bucket == nil {
			return fmt.Errorf("no such bucket: %s", s.bucketName)
		}

		votedForBytes, err := json.Marshal(term)
		if err != nil {
			return fmt.Errorf("failed to encode voted for: %s", err.Error())
		}

		if err := bucket.Put(s.keyCurrentTerm, votedForBytes); err != nil {
			return fmt.Errorf("failed to set voted for: %s", err.Error())
		}
		return nil
	})
}
