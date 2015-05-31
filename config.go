package ssh

import (
	"bytes"
	"encoding/gob"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/opt"
)

// Config is a struct encapsulating the config.
type Config struct {
	Path string
}

// ConfigData is the raw data that's going to be stored in the db.
type ConfigData map[string]interface{}

// ConfigItem is a map to store in the config.
type ConfigItem interface {
	ToMap() ConfigData
	FromMap(ConfigData)
}

// NewConfig creates a new config given a path.
func NewConfig(path string) *Config {
	config := &Config{
		Path: path,
	}

	return config
}

// Put puts a value given a field.
func (config *Config) Put(key string, item ConfigItem) error {
	data, err := encodeData(item.ToMap())
	if err != nil {
		return err
	}

	db, err := config.openDB()
	if err != nil {
		return err
	}
	defer db.Close()

	err = db.Put([]byte(key), data, &opt.WriteOptions{Sync: true})

	return err
}

// Get gets a value from the database given a key.
func (config *Config) Get(key string, item ConfigItem) error {
	db, err := config.openDB()
	if err != nil {
		return err
	}
	defer db.Close()

	data, err := db.Get([]byte(key), nil)
	if err != nil {
		return err
	}

	value, err := decodeData(data)
	if err != nil {
		return err
	}

	item.FromMap(value)
	return nil
}

// List returns all config items
func (config *Config) List() []ConfigData {
	items := []ConfigData{}

	db, err := config.openDB()
	if err != nil {
		return items
	}
	defer db.Close()

	iter := db.NewIterator(nil, nil)
	for iter.Next() {
		value, err := decodeData(iter.Value())
		if err != nil {
		}
		items = append(items, value)
	}

	iter.Release()

	return items
}

// HasKeys returns true if the database has keys.
func (config *Config) HasKeys() bool {
	db, err := config.openDB()
	if err != nil {
		return false
	}
	defer db.Close()
	iter := db.NewIterator(nil, nil)

	hasKeys := iter.Next()
	iter.Release()

	return hasKeys
}

func encodeData(data map[string]interface{}) ([]byte, error) {
	var buffer bytes.Buffer
	encoder := gob.NewEncoder(&buffer)
	err := encoder.Encode(data)

	if err != nil {
		return nil, err
	}

	return buffer.Bytes(), nil
}

func decodeData(data []byte) (map[string]interface{}, error) {
	buffer := bytes.NewBuffer(data)
	var item map[string]interface{}

	encoder := gob.NewDecoder(buffer)
	err := encoder.Decode(&item)

	if err != nil {
		return nil, err
	}

	return item, nil
}

func (config *Config) openDB() (*leveldb.DB, error) {
	db, err := leveldb.OpenFile(config.Path, nil)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func init() {
}
