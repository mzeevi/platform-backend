package controllers

import (
	"github.com/dana-team/platform-backend/src/types"
)

// convertKeyValueToMap converts a slice of KeyValue pairs to a map
// with string keys and values.
func convertKeyValueToMap(kvList []types.KeyValue) map[string]string {
	values := make(map[string]string)
	for _, kv := range kvList {
		values[kv.Key] = kv.Value
	}
	return values
}

// convertMapToKeyValue converts a map with string keys and values
// to a slice of KeyValue pairs.
func convertMapToKeyValue(values map[string]string) []types.KeyValue {
	var kvList []types.KeyValue
	for k, v := range values {
		kvList = append(kvList, types.KeyValue{Key: k, Value: v})
	}
	return kvList
}

// convertKeyValueToByteMap converts a slice of KeyValue pairs
// to a map with string keys and byte slice values.
func convertKeyValueToByteMap(kvList []types.KeyValue) map[string][]byte {
	data := map[string][]byte{}
	for _, kv := range kvList {
		data[kv.Key] = []byte(kv.Value)
	}
	return data
}
