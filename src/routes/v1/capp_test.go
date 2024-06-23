package v1_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/dana-team/platform-backend/src/utils"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/dana-team/platform-backend/src/types"
	"github.com/stretchr/testify/assert"
)

const (
	cappNamespace = testName + "-capp-ns"
	cappName      = testName + "-capp"
)

func setupCapps() {
	createTestNamespace(cappNamespace)
	createTestCapp(cappName+"-1", cappNamespace, map[string]string{labelKey + "-1": labelValue + "-1"}, nil)
	createTestCapp(cappName+"-2", cappNamespace, map[string]string{labelKey + "-2": labelValue + "-2"}, nil)
}

// createTestCapp creates a test Capp object.
func createTestCapp(name, namespace string, labels, annotations map[string]string) {
	capp := utils.GetStubCapp(name, namespace, labels, annotations)
	err := dynClient.Create(context.TODO(), &capp)
	if err != nil {
		panic(err)
	}
}

func TestGetCapps(t *testing.T) {
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
		response   types.CappList
	}

	cases := map[string]struct {
		requestParams requestParams
		want          want
	}{
		"ShouldSucceedGettingCapps": {
			requestParams: requestParams{
				namespace: cappNamespace,
			},
			want: want{
				statusCode: http.StatusOK,
				response:   types.CappList{Count: 2, Capps: []string{cappName + "-1", cappName + "-2"}},
			},
		},
		"ShouldFailWithBadRequestInvalidURI": {
			requestParams: requestParams{
				namespace: "",
			},
			want: want{
				statusCode: http.StatusBadRequest,
				response:   types.CappList{},
			},
		},
		"ShouldSucceedGettingCappsWithLabelSelector": {
			requestParams: requestParams{
				namespace: cappNamespace,
				labelSelector: selector{
					keys:   []string{labelKey + "-1"},
					values: []string{labelValue + "-1"},
				},
			},
			want: want{
				statusCode: http.StatusOK,
				response:   types.CappList{Count: 1, Capps: []string{cappName + "-1"}},
			},
		},
		"ShouldFailGettingCappsWithInvalidLabelSelector": {
			requestParams: requestParams{
				namespace: cappNamespace,
				labelSelector: selector{
					keys:   []string{labelKey + "-1"},
					values: []string{labelValue + " 1"},
				},
			},
			want: want{
				statusCode: http.StatusBadRequest,
				response:   types.CappList{},
			},
		},
		"ShouldSucceedGettingNoCappsWithLabelSelector": {
			requestParams: requestParams{
				namespace: cappNamespace,
				labelSelector: selector{
					keys:   []string{labelKey + "-3"},
					values: []string{labelValue + "-3"},
				},
			},
			want: want{
				statusCode: http.StatusOK,
				response:   types.CappList{Count: 0, Capps: nil},
			},
		},
	}

	for name, test := range cases {
		t.Run(name, func(t *testing.T) {
			params := url.Values{}

			for i, key := range test.requestParams.labelSelector.keys {
				params.Add(labelSelectorKey, key+"="+test.requestParams.labelSelector.values[i])
			}

			baseURI := fmt.Sprintf("/v1/namespaces/%s/capps/", test.requestParams.namespace)

			request, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s?%s", baseURI, params.Encode()), nil)
			assert.NoError(t, err)
			writer := httptest.NewRecorder()
			router.ServeHTTP(writer, request)

			assert.Equal(t, test.want.statusCode, writer.Code)

			if writer.Code == http.StatusOK {
				var response types.CappList
				err = json.Unmarshal(writer.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Equal(t, test.want.response, response)
			}
		})
	}
}

func TestGetCapp(t *testing.T) {
	type requestParams struct {
		name      string
		namespace string
	}

	type want struct {
		statusCode int
		response   types.Capp
	}

	cases := map[string]struct {
		requestParams requestParams
		want          want
	}{
		"ShouldSucceedGettingCapp": {
			requestParams: requestParams{
				namespace: cappNamespace,
				name:      cappName + "-1",
			},
			want: want{
				statusCode: http.StatusOK,
				response:   utils.GetStubCappType(cappName+"-1", cappNamespace, map[string]string{labelKey + "-1": labelValue + "-1"}, nil),
			},
		},
		"ShouldFailWithBadRequestInvalidURI": {
			requestParams: requestParams{
				namespace: "",
				name:      cappRevisionName + "-1",
			},
			want: want{
				statusCode: http.StatusBadRequest,
				response:   types.Capp{},
			},
		},
	}

	for name, test := range cases {
		t.Run(name, func(t *testing.T) {
			baseURI := fmt.Sprintf("/v1/namespaces/%s/capps/%s", test.requestParams.namespace, test.requestParams.name)
			request, err := http.NewRequest(http.MethodGet, baseURI, nil)
			assert.NoError(t, err)

			writer := httptest.NewRecorder()
			router.ServeHTTP(writer, request)

			assert.Equal(t, test.want.statusCode, writer.Code)

			if writer.Code == http.StatusOK {
				var response types.Capp
				err = json.Unmarshal(writer.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Equal(t, test.want.response, response)
			}
		})
	}
}

func TestCreateCapp(t *testing.T) {
	type bodyParams struct {
		capp types.Capp
	}

	type want struct {
		statusCode int
		response   types.Capp
	}

	cases := map[string]struct {
		bodyParams bodyParams
		want       want
	}{
		"ShouldSucceedCreatingCapp": {
			bodyParams: bodyParams{
				capp: utils.GetStubCappType(cappName+"-4", cappNamespace, map[string]string{labelKey + "-4": labelValue + "-4"}, nil),
			},
			want: want{
				statusCode: http.StatusOK,
				response:   utils.GetStubCappType(cappName+"-4", cappNamespace, map[string]string{labelKey + "-4": labelValue + "-4"}, nil),
			},
		},
		"ShouldFailWithBadRequestInvalidURI": {
			bodyParams: bodyParams{
				capp: utils.GetStubCappType(cappName+"-4", "", map[string]string{labelKey + "-4": labelValue + "-4"}, nil),
			},
			want: want{
				statusCode: http.StatusBadRequest,
				response:   types.Capp{},
			},
		},
		"ShouldFailWithBadRequestBody": {
			bodyParams: bodyParams{
				capp: utils.GetStubCappType("", cappNamespace, map[string]string{labelKey + "-4": labelValue + "-4"}, nil),
			},
			want: want{
				statusCode: http.StatusBadRequest,
				response:   types.Capp{},
			},
		},
	}

	for name, test := range cases {
		t.Run(name, func(t *testing.T) {
			body, err := json.Marshal(test.bodyParams.capp)
			assert.NoError(t, err)

			baseURI := fmt.Sprintf("/v1/namespaces/%s/capps/", test.bodyParams.capp.Metadata.Namespace)
			request, err := http.NewRequest(http.MethodPost, baseURI, bytes.NewBuffer(body))
			assert.NoError(t, err)
			request.Header.Set(contentType, applicationJson)

			writer := httptest.NewRecorder()
			router.ServeHTTP(writer, request)

			assert.Equal(t, test.want.statusCode, writer.Code)

			if writer.Code == http.StatusOK {
				var response types.Capp
				err = json.Unmarshal(writer.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Equal(t, test.want.response, response)
			}
		})
	}
}

func TestUpdateCapp(t *testing.T) {
	type bodyParams struct {
		capp types.Capp
	}

	type want struct {
		statusCode int
		response   types.Capp
	}

	cases := map[string]struct {
		bodyParams bodyParams
		want       want
	}{
		"ShouldSucceedUpdatingCapp": {
			bodyParams: bodyParams{
				capp: utils.GetStubCappType(cappName+"-5", cappNamespace, map[string]string{labelKey + "-5-updated": labelValue + "-5-updated"}, nil),
			},
			want: want{
				statusCode: http.StatusOK,
				response:   utils.GetStubCappType(cappName+"-5", cappNamespace, map[string]string{labelKey + "-5-updated": labelValue + "-5-updated"}, nil),
			},
		},
		"ShouldFailWithBadRequestInvalidURI": {
			bodyParams: bodyParams{
				capp: utils.GetStubCappType(cappName+"-5", "", map[string]string{labelKey + "-5": labelValue + "-5"}, nil),
			},
			want: want{
				statusCode: http.StatusBadRequest,
				response:   types.Capp{},
			},
		},
	}

	createTestCapp(cappName+"-5", cappNamespace, map[string]string{labelKey + "-5": labelValue + "-5"}, nil)
	for name, test := range cases {
		t.Run(name, func(t *testing.T) {
			body, err := json.Marshal(test.bodyParams.capp)
			assert.NoError(t, err)

			baseURI := fmt.Sprintf("/v1/namespaces/%s/capps/%s", test.bodyParams.capp.Metadata.Namespace, test.bodyParams.capp.Metadata.Name)
			request, err := http.NewRequest(http.MethodPut, baseURI, bytes.NewBuffer(body))
			assert.NoError(t, err)
			request.Header.Set(contentType, applicationJson)

			writer := httptest.NewRecorder()
			router.ServeHTTP(writer, request)

			assert.Equal(t, test.want.statusCode, writer.Code)

			if writer.Code == http.StatusOK {
				var response types.Capp
				err = json.Unmarshal(writer.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Equal(t, test.want.response, response)
			}
		})
	}
}

func TestDeleteCapp(t *testing.T) {
	type requestParams struct {
		name      string
		namespace string
	}

	type want struct {
		statusCode int
		response   types.CappError
	}

	cases := map[string]struct {
		requestParams requestParams
		want          want
	}{
		"ShouldSucceedDeletingCapp": {
			requestParams: requestParams{
				name:      cappName + "-6",
				namespace: cappNamespace,
			},
			want: want{
				statusCode: http.StatusOK,
				response:   types.CappError{Message: fmt.Sprintf("Deleted capp %q in namespace %q successfully", cappName+"-6", cappNamespace)},
			},
		},
		"ShouldFailWithBadRequestInvalidURI": {
			requestParams: requestParams{
				name:      cappName + "-6",
				namespace: "",
			},
			want: want{
				statusCode: http.StatusBadRequest,
				response:   types.CappError{Message: ""},
			},
		},
	}

	createTestCapp(cappName+"-6", cappNamespace, map[string]string{labelKey + "-6": labelValue + "-6"}, nil)
	for name, test := range cases {
		t.Run(name, func(t *testing.T) {
			baseURI := fmt.Sprintf("/v1/namespaces/%s/capps/%s", test.requestParams.namespace, test.requestParams.name)
			request, err := http.NewRequest(http.MethodDelete, baseURI, nil)
			assert.NoError(t, err)
			request.Header.Set(contentType, applicationJson)

			writer := httptest.NewRecorder()
			router.ServeHTTP(writer, request)

			assert.Equal(t, test.want.statusCode, writer.Code)

			if writer.Code == http.StatusOK {
				var response types.CappError
				err = json.Unmarshal(writer.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Equal(t, test.want.response, response)
			}
		})
	}
}
