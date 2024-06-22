package controllers

import (
	"fmt"
	"github.com/dana-team/platform-backend/src/types"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	"regexp"
	"strings"
)

// checkLabelSelectorFormat checks if the input string is in the format "a=b,c=d,e=f".
func checkLabelSelectorFormat(labelSelector string) bool {
	pattern := `^([^\s=,]+)=([^\s=,]+)(,([^\s=,]+)=([^\s=,]+))*$`
	re := regexp.MustCompile(pattern)
	return re.MatchString(labelSelector)
}

// parseLabelSelector parses the label selector into a map of key-value pairs.
func parseLabelSelector(labelSelector string) (map[string]string, error) {
	if len(labelSelector) == 0 {
		return nil, nil
	}

	if !checkLabelSelectorFormat(labelSelector) {
		return nil, k8serrors.NewBadRequest(fmt.Sprintf("format of labelSelector %q is invalid, must be of fomat 'a=b,c=d,e=f'", labelSelector))
	}

	pairs := strings.Split(labelSelector, ",")
	labels := make(map[string]string, len(pairs))
	for _, pair := range pairs {
		kv := strings.Split(pair, "=")
		if len(kv) != 2 {
			return nil, k8serrors.NewBadRequest(fmt.Sprintf("format of key-value pair %q is invalid, must of format 'a=b'", pair))
		}
		labels[kv[0]] = kv[1]
	}

	return labels, nil
}

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
