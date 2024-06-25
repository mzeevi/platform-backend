package v1_test

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/dana-team/platform-backend/src/routes/mocks"
	"github.com/dana-team/platform-backend/src/types"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

const (
	cappRevisionNamespace = testName + "-capp-revision-ns"
	cappRevisionName      = testName + "-capp-revision"
)

func setupCappRevisions() {
	createTestNamespace(cappRevisionNamespace)
	createTestCappRevision(cappRevisionName+"-1", cappRevisionNamespace, map[string]string{labelKey + "-1": labelValue + "-1"}, nil)
	createTestCappRevision(cappRevisionName+"-2", cappRevisionNamespace, map[string]string{labelKey + "-2": labelValue + "-2"}, nil)
}

// createTestCappRevision creates a test CappRevision object.
func createTestCappRevision(name, namespace string, labels, annotations map[string]string) {
	cappRevision := mocks.GetCappRevision(name, namespace, labels, annotations)
	err := dynClient.Create(context.TODO(), &cappRevision)
	if err != nil {
		panic(err)
	}
}

func TestGetCappRevisions(t *testing.T) {
	type selector struct {
		keys   []string
		values []string
	}

	type requestParams struct {
		namespace     string
		labelSelector selector
	}

	type want struct {
		statusCode int
		response   types.CappRevisionList
	}

	cases := map[string]struct {
		requestParams requestParams
		want          want
	}{
		"ShouldSucceedGettingCappRevisions": {
			requestParams: requestParams{
				namespace: cappRevisionNamespace,
			},
			want: want{
				statusCode: http.StatusOK,
				response:   types.CappRevisionList{Count: 2, CappRevisions: []string{cappRevisionName + "-1", cappRevisionName + "-2"}},
			},
		},
		"ShouldFailWithBadRequestInvalidURI": {
			requestParams: requestParams{
				namespace: "",
			},
			want: want{
				statusCode: http.StatusBadRequest,
				response:   types.CappRevisionList{},
			},
		},
		"ShouldSucceedGettingCappRevisionsWithLabelSelector": {
			requestParams: requestParams{
				namespace: cappRevisionNamespace,
				labelSelector: selector{
					keys:   []string{labelKey + "-1"},
					values: []string{labelValue + "-1"},
				},
			},
			want: want{
				statusCode: http.StatusOK,
				response:   types.CappRevisionList{Count: 1, CappRevisions: []string{cappRevisionName + "-1"}},
			},
		},
		"ShouldFailGettingCappRevisionsWithInvalidLabelSelector": {
			requestParams: requestParams{
				namespace: cappRevisionNamespace,
				labelSelector: selector{
					keys:   []string{labelKey + "-1"},
					values: []string{labelValue + " 1"},
				},
			},
			want: want{
				statusCode: http.StatusBadRequest,
				response:   types.CappRevisionList{},
			},
		},
		"ShouldSucceedGettingNoCappRevisionsWithLabelSelector": {
			requestParams: requestParams{
				namespace: cappRevisionNamespace,
				labelSelector: selector{
					keys:   []string{labelKey + "-3"},
					values: []string{labelValue + "-3"},
				},
			},
			want: want{
				statusCode: http.StatusOK,
				response:   types.CappRevisionList{Count: 0, CappRevisions: nil},
			},
		},
	}

	for name, test := range cases {
		t.Run(name, func(t *testing.T) {
			params := url.Values{}

			for i, key := range test.requestParams.labelSelector.keys {
				params.Add(labelSelectorKey, key+"="+test.requestParams.labelSelector.values[i])
			}

			baseURI := fmt.Sprintf("/v1/namespaces/%s/capprevisions/", test.requestParams.namespace)
			request, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s?%s", baseURI, params.Encode()), nil)
			assert.NoError(t, err)
			writer := httptest.NewRecorder()
			router.ServeHTTP(writer, request)

			assert.Equal(t, test.want.statusCode, writer.Code)

			if writer.Code == http.StatusOK {
				var response types.CappRevisionList
				err = json.Unmarshal(writer.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Equal(t, test.want.response, response)
			}
		})
	}
}

func TestGetCappRevision(t *testing.T) {
	type requestParams struct {
		name      string
		namespace string
	}

	type want struct {
		statusCode int
		response   types.CappRevision
	}

	cases := map[string]struct {
		requestParams requestParams
		want          want
	}{
		"ShouldSucceedGettingCappRevision": {
			requestParams: requestParams{
				namespace: cappRevisionNamespace,
				name:      cappRevisionName + "-1",
			},
			want: want{
				statusCode: http.StatusOK,
				response: mocks.GetCappRevisionType(cappRevisionName+"-1", cappRevisionNamespace,
					map[string]string{labelKey + "-1": labelValue + "-1"}, nil),
			},
		},
		"ShouldFailWithBadRequestInvalidURI": {
			requestParams: requestParams{
				namespace: "",
				name:      cappRevisionName + "-1",
			},
			want: want{
				statusCode: http.StatusBadRequest,
				response:   types.CappRevision{},
			},
		},
	}

	for name, test := range cases {
		t.Run(name, func(t *testing.T) {
			baseURI := fmt.Sprintf("/v1/namespaces/%s/capprevisions/%s", test.requestParams.namespace, test.requestParams.name)
			request, err := http.NewRequest(http.MethodGet, baseURI, nil)
			assert.NoError(t, err)
			writer := httptest.NewRecorder()
			router.ServeHTTP(writer, request)

			assert.Equal(t, test.want.statusCode, writer.Code)

			if writer.Code == http.StatusOK {
				var response types.CappRevision
				err = json.Unmarshal(writer.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Equal(t, test.want.response, response)
			}
		})
	}
}
