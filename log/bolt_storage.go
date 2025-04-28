package log

import (
	"encoding/json"
	"fmt"
	"github.com/boltdb/bolt"
	"strconv"
	"sync"
)

// TODO: Тут надо подумать на самом деле, щас как будто мы просто будем добавлять в журнал, а потом удалять, но можно как-нибудь красиво использовать каналы и соответственно просто rollback записей.

type boltStorage struct {
	unimplementedStorage
	db         *bolt.DB
	mutex      sync.RWMutex
	bucketName []byte
}

func NewBoltStorage(dbPath string, bucketName string) (Storage, error) {
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
		return nil, fmt.Errorf("failed to initialize bolt storage: %s", err.Error())
	}

	return &boltStorage{
		db:         db,
		mutex:      sync.RWMutex{},
		bucketName: []byte(bucketName),
	}, nil
}

func (s *boltStorage) Append(entry *Entry) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	return s.db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(s.bucketName)
		if bucket == nil {
			return fmt.Errorf("no such bucket: %s", s.bucketName)
		}

		valueBytes, err := json.Marshal(entry)
		if err != nil {
			return fmt.Errorf("failed to encode entry: %s", err.Error())
		}
		keyBytes, err := json.Marshal(entry.Index)
		if err != nil {
			return fmt.Errorf("failed to encode entry: %s", err.Error())
		}

		lastIndex, _ := s.LastIndex()
		if lastIndex != entry.Index-1 {
			return fmt.Errorf("failed to append entry with unexpected index %d, last was %d", entry.Index, lastIndex)
		}

		err = bucket.Put(keyBytes, valueBytes)
		if err != nil {
			return fmt.Errorf("failed to append entry: %s", err.Error())
		}
		return nil
	})
}

func (s *boltStorage) Get(index int) (*Entry, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	var entry *Entry
	err := s.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(s.bucketName)
		if bucket == nil {
			return fmt.Errorf("no such bucket: %s", s.bucketName)
		}

		key, err := json.Marshal(index)
		if err != nil {
			return fmt.Errorf("failed to encode index: %s", err.Error())
		}

		valueBytes := bucket.Get(key)
		if err = json.Unmarshal(valueBytes, &entry); err != nil {
			return fmt.Errorf("failed to decode entry: %s", err.Error())
		}
		return nil
	})
	return entry, err
}

func (s *boltStorage) FirstIndex() (int, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	firstIndex := 0

	err := s.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(s.bucketName)
		if bucket == nil {
			return fmt.Errorf("no such bucket: %s", s.bucketName)
		}

		cursor := bucket.Cursor()
		key, _ := cursor.First()
		if key == nil {
			firstIndex = 1
			return nil
		}
		firstIndex, _ = strconv.Atoi(string(key))
		return nil
	})
	return firstIndex, err
}

func (s *boltStorage) LastIndex() (int, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	firstIndex := 0

	err := s.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(s.bucketName)
		if bucket == nil {
			return fmt.Errorf("no such bucket: %s", s.bucketName)
		}

		cursor := bucket.Cursor()
		key, _ := cursor.Last()
		if key == nil {
			firstIndex = 0
			return nil
		}
		firstIndex, _ = strconv.Atoi(string(key))
		return nil
	})
	return firstIndex, err
}

func (s *boltStorage) TermByIndex(index int) (int, error) {
	entry, err := s.Get(index)
	if err != nil {
		return 0, err
	}
	return entry.Term, nil
}

func (s *boltStorage) Truncate(index int) error {
	return nil
}
