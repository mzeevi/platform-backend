package v1_test

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/dana-team/platform-backend/src/routes/mocks"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/dana-team/platform-backend/src/types"
	"github.com/stretchr/testify/assert"
)

const (
	configMapNamespace = testName + "-configmap-ns"
	configMapName      = testName + "-configmap"
	configKey          = "key"
	configValue        = "value"
)

func setupConfigMaps() {
	createTestNamespace(configMapNamespace)
	createTestConfigMap(configMapName+"-1", configMapNamespace)
}

// createTestConfigMap creates a test ConfigMap object.
func createTestConfigMap(name, namespace string) {
	configMap := mocks.GetConfigMap(name, namespace, map[string]string{configKey + "-1": configValue + "-1"})
	_, err := client.CoreV1().ConfigMaps(namespace).Create(context.TODO(), &configMap, metav1.CreateOptions{})
	if err != nil {
		panic(err)
	}
}

func TestGetConfigMap(t *testing.T) {
	type requestParams struct {
		name      string
		namespace string
	}

	type want struct {
		statusCode int
		response   types.ConfigMap
	}

	cases := map[string]struct {
		requestParams requestParams
		want          want
	}{
		"ShouldSucceedGettingConfigMap": {
			requestParams: requestParams{
				namespace: configMapNamespace,
				name:      configMapName + "-1",
			},
			want: want{
				statusCode: http.StatusOK,
				response:   mocks.GetConfigMapType(map[string]string{configKey + "-1": configValue + "-1"}),
			},
		},
		"ShouldFailWithBadRequestInvalidURI": {
			requestParams: requestParams{
				namespace: "",
				name:      configMapName + "-1",
			},
			want: want{
				statusCode: http.StatusBadRequest,
				response:   types.ConfigMap{},
			},
		},
	}

	for name, test := range cases {
		t.Run(name, func(t *testing.T) {
			baseURI := fmt.Sprintf("/v1/namespaces/%s/configmaps/%s", test.requestParams.namespace, test.requestParams.name)
			request, err := http.NewRequest(http.MethodGet, baseURI, nil)
			assert.NoError(t, err)
			writer := httptest.NewRecorder()
			router.ServeHTTP(writer, request)

			assert.Equal(t, test.want.statusCode, writer.Code)

			if writer.Code == http.StatusOK {
				var response types.ConfigMap
				err = json.Unmarshal(writer.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Equal(t, test.want.response, response)
			}
		})
	}
}
