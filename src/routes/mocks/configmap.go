package mocks

import (
	"github.com/dana-team/platform-backend/src/types"
	"github.com/dana-team/platform-backend/src/utils"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// GetConfigMap returns a mock ConfigMap object.
func GetConfigMap(name, namespace string, data map[string]string) corev1.ConfigMap {
	configMap := corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
		Data: data,
	}

	return configMap
}

// GetConfigMapType returns a mock ConfigMap type object.
func GetConfigMapType(data map[string]string) types.ConfigMap {
	configMap := types.ConfigMap{
		Data: utils.ConvertMapToKeyValue(data),
	}

	return configMap
}
