package rpc

import (
	"distributed-kv-store/store"
	"errors"
	"log"
)

// Args represents the arguments for RPC calls.
type Args struct {
	Key   string
	Value string
}

// Reply represents the reply for RPC calls.
type Reply struct {
	Value string
}

// KVStore represents the RPC server for the key-value store.
type KVStore struct {
	store *store.Store
}

// NewKVStore creates and returns a new KVStore RPC server instance.
func NewKVStore(s *store.Store) *KVStore {
	return &KVStore{
		store: s,
	}
}

// Get retrieves a value from the store.
func (k *KVStore) Get(args *Args, reply *Reply) error {
	log.Printf("RPC Get request for key: %s\n", args.Key)
	val, err := k.store.Get(args.Key)
	if err != nil {
		return errors.New("key not found")
	}
	reply.Value = val
	return nil
}

// Put stores a key-value pair in the store.
func (k *KVStore) Put(args *Args, reply *Reply) error {
	log.Printf("RPC Put request for key: %s, value: %s\n", args.Key, args.Value)
	err := k.store.Put(args.Key, args.Value)
	if err != nil {
		return err
	}
	reply.Value = "OK" // Indicate success
	return nil
}

