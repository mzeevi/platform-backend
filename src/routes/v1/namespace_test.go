package v1_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/dana-team/platform-backend/src/routes/mocks"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/dana-team/platform-backend/src/types"
	"github.com/stretchr/testify/assert"
)

const (
	nsName = testName + "-namespace"
)

func setupNamespaces() {
	createTestNamespace(nsName + "-1")
	createTestNamespace(nsName + "-2")

}

func TestGetNamespaces(t *testing.T) {
	type want struct {
		statusCode int
		response   types.NamespaceList
	}

	cases := map[string]struct {
		want want
	}{
		"ShouldSucceedGettingNamespaces": {
			want: want{
				statusCode: http.StatusOK,
				response:   types.NamespaceList{Count: 2, Namespaces: []types.Namespace{{Name: nsName + "-1"}, {Name: nsName + "-2"}}},
			},
		},
	}

	for name, test := range cases {
		t.Run(name, func(t *testing.T) {
			baseURI := "/v1/namespaces/"
			request, err := http.NewRequest(http.MethodGet, baseURI, nil)
			assert.NoError(t, err)
			writer := httptest.NewRecorder()
			router.ServeHTTP(writer, request)
			assert.Equal(t, test.want.statusCode, writer.Code)

			if writer.Code == http.StatusOK {
				var response types.NamespaceList
				err = json.Unmarshal(writer.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.GreaterOrEqual(t, response.Count, test.want.response.Count)
				for _, ns := range test.want.response.Namespaces {
					assert.Contains(t, response.Namespaces, ns)
				}
			}
		})
	}
}

func TestGetNamespace(t *testing.T) {
	type requestParams struct {
		name string
	}

	type want struct {
		statusCode int
		response   types.Namespace
	}

	cases := map[string]struct {
		requestParams requestParams
		want          want
	}{
		"ShouldSucceedGettingNamespace": {
			requestParams: requestParams{
				name: nsName + "-1",
			},
			want: want{
				statusCode: http.StatusOK,
				response:   mocks.GetNamespaceType(nsName + "-1"),
			},
		},
	}

	for name, test := range cases {
		t.Run(name, func(t *testing.T) {
			baseURI := fmt.Sprintf("/v1/namespaces/%s", test.requestParams.name)
			request, err := http.NewRequest(http.MethodGet, baseURI, nil)
			assert.NoError(t, err)

			writer := httptest.NewRecorder()
			router.ServeHTTP(writer, request)

			assert.Equal(t, test.want.statusCode, writer.Code)

			if writer.Code == http.StatusOK {
				var response types.Namespace
				err = json.Unmarshal(writer.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Equal(t, test.want.response, response)
			}
		})
	}
}

func TestCreateNamespace(t *testing.T) {
	type bodyParams struct {
		ns types.Namespace
	}

	type want struct {
		statusCode int
		response   types.Namespace
	}

	cases := map[string]struct {
		bodyParams bodyParams
		want       want
	}{
		"ShouldSucceedCreatingNamespace": {
			bodyParams: bodyParams{
				ns: mocks.GetNamespaceType(nsName + "-3"),
			},
			want: want{
				statusCode: http.StatusOK,
				response:   mocks.GetNamespaceType(nsName + "-3"),
			},
		},
		"ShouldFailWithBadRequestBody": {
			bodyParams: bodyParams{
				ns: types.Namespace{},
			},
			want: want{
				statusCode: http.StatusBadRequest,
				response:   types.Namespace{},
			},
		},
	}

	for name, test := range cases {
		t.Run(name, func(t *testing.T) {
			body, err := json.Marshal(test.bodyParams.ns)
			assert.NoError(t, err)

			baseURI := fmt.Sprintf("/v1/namespaces/")
			request, err := http.NewRequest(http.MethodPost, baseURI, bytes.NewBuffer(body))
			assert.NoError(t, err)
			request.Header.Set(contentType, applicationJson)

			writer := httptest.NewRecorder()
			router.ServeHTTP(writer, request)

			assert.Equal(t, test.want.statusCode, writer.Code)

			if writer.Code == http.StatusOK {
				var response types.Namespace
				err = json.Unmarshal(writer.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Equal(t, test.want.response, response)
			}
		})
	}
}

func TestDeleteNamespace(t *testing.T) {
	type requestParams struct {
		name string
	}

	type want struct {
		statusCode int
		response   map[string]string
	}

	cases := map[string]struct {
		requestParams requestParams
		want          want
	}{
		"ShouldSucceedDeletingNamespace": {
			requestParams: requestParams{
				name: cappName + "-4",
			},
			want: want{
				statusCode: http.StatusOK,
				response:   map[string]string{"message": fmt.Sprintf("Deleted namespace successfully %q", cappName+"-4")},
			},
		},
	}

	createTestNamespace(cappName + "-4")
	for name, test := range cases {
		t.Run(name, func(t *testing.T) {
			baseURI := fmt.Sprintf("/v1/namespaces/%s", test.requestParams.name)
			request, err := http.NewRequest(http.MethodDelete, baseURI, nil)
			assert.NoError(t, err)
			request.Header.Set(contentType, applicationJson)

			writer := httptest.NewRecorder()
			router.ServeHTTP(writer, request)

			assert.Equal(t, test.want.statusCode, writer.Code)

			if writer.Code == http.StatusOK {
				var response map[string]string
				err = json.Unmarshal(writer.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Equal(t, test.want.response, response)
			}
		})
	}
}
