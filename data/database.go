// The above code defines a `Database` struct with methods for storing, retrieving, and deleting
// key-value pairs in a Consul key-value store.
package data

import (
	"encoding/json"
	"fmt"

	"github.com/hashicorp/consul/api"
)

type Database struct {
	client *api.Client
}

func NewDatabase() (*Database, error) {
	client, err := api.NewClient(api.DefaultConfig())
	if err != nil {
		return nil, err
	}

	return &Database{client: client}, nil
}

// This `Put` method in the `Database` struct is used to store a key-value pair in the Consul key-value
// store. Here's a breakdown of what it does:
func (db *Database) Put(keyType string, name string, version string, value interface{}) (string, error) {
	kv := db.client.KV()
	// Form the key using the keyType, name, and version
	key := fmt.Sprintf("%s/%s/%s", keyType, name, version)
	jsonValue, err := json.Marshal(value)
	if err != nil {
		return "", err
	}
	p := &api.KVPair{Key: key, Value: jsonValue}
	_, err = kv.Put(p, nil)
	if err != nil {
		return "", err
	}
	return key, nil
}

// The `Get` method in the `Database` struct is used to retrieve a value from the Consul key-value
// store based on the provided key. Here's a breakdown of what it does:
func (db *Database) Get(key string, value interface{}) error {
	kv := db.client.KV()
	pair, _, err := kv.Get(key, nil)
	if err != nil {
		return err
	}
	if pair == nil {
		return nil
	}
	err = json.Unmarshal(pair.Value, value)
	if err != nil {
		return err
	}
	return nil
}

// The `Delete` method in the `Database` struct is used to delete a key-value pair from the Consul
// key-value store based on the provided key. Here's a breakdown of what it does:
func (db *Database) Delete(key string) error {
	kv := db.client.KV()
	_, err := kv.Delete(key, nil)
	if err != nil {
		return err
	}
	return nil
}

// The `List` method in the `Database` struct is used to list all key-value pairs in the Consul
// key-value store that match the provided key prefix. Here's a breakdown of what it does:
func (db *Database) List(keyPrefix string) (map[string]interface{}, error) {
	kv := db.client.KV()
	pairs, _, err := kv.List(keyPrefix, nil)
	if err != nil {
		return nil, err
	}
	result := make(map[string]interface{})
	for _, pair := range pairs {
		var value interface{}
		err := json.Unmarshal(pair.Value, &value)
		if err != nil {
			return nil, err
		}
		result[pair.Key] = value
	}
	return result, nil
}
