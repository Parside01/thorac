package storage

import (
	"encoding/json"
	"fmt"
	"github.com/boltdb/bolt"
	"strconv"
	"sync"
	"thorac/core"
)

// TODO: Тут надо подумать на самом деле, щас как будто мы просто будем добавлять в журнал, а потом удалять, но можно как-нибудь красиво использовать каналы и соответственно просто rollback записей.

type boltLogStorage struct {
	unimplementedLogStorage
	db         *bolt.DB
	mutex      sync.RWMutex
	bucketName []byte
}

func NewBoltLogStorage(dbPath string, bucketName string) (core.LogStorage, error) {
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
		return nil, fmt.Errorf("failed to initialize bolt LogStorage: %s", err.Error())
	}

	return &boltLogStorage{
		db:         db,
		mutex:      sync.RWMutex{},
		bucketName: []byte(bucketName),
	}, nil
}

func (s *boltLogStorage) Append(entry *core.LogEntry) error {
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

func (s *boltLogStorage) Get(index int) (*core.LogEntry, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	var LogEntry *core.LogEntry
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
		if err = json.Unmarshal(valueBytes, &LogEntry); err != nil {
			return fmt.Errorf("failed to decode LogEntry: %s", err.Error())
		}
		return nil
	})
	return LogEntry, err
}

func (s *boltLogStorage) FirstIndex() (int, error) {
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

func (s *boltLogStorage) LastIndex() (int, error) {
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

func (s *boltLogStorage) TermByIndex(index int) (int, error) {
	LogEntry, err := s.Get(index)
	if err != nil {
		return 0, err
	}
	return LogEntry.Term, nil
}

func (s *boltLogStorage) Truncate(index int) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	return s.db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(s.bucketName)
		if bucket == nil {
			return fmt.Errorf("no such bucket: %s", s.bucketName)
		}

		cursor := bucket.Cursor()

		minIndex, err := json.Marshal(index)
		if err != nil {
			return fmt.Errorf("failed to encode index: %s", err.Error())
		}

		for key, _ := cursor.Seek(minIndex); key != nil; key, _ = cursor.Next() {
			err = bucket.Delete(key)
			if err != nil {
				return fmt.Errorf("failed to delete key: %s", err.Error())
			}
		}
		return nil
	})
}

func (s *boltLogStorage) Close() error {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	return s.db.Close()
}
