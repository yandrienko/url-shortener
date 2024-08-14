package remove_test

import (
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
	"url-shortener/internal/http-server/handlers/url/remove"
	"url-shortener/internal/http-server/handlers/url/remove/mocks"
	"url-shortener/internal/lib/logger/handlers/slogdiscard"
)

func TestRemoveHandler(t *testing.T) {
	cases := []struct {
		name      string
		alias     string
		respError string
		mockError error
	}{
		{
			name:  "test with no error",
			alias: "test_alias",
		},
		{
			name:      "missing alias",
			alias:     " ",
			respError: "invalid request",
		},
		// TODO add test cases
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			urlRemoverMock := mocks.NewUrlRemover(t)

			if tc.respError == "" || tc.mockError != nil {
				urlRemoverMock.On("RemoveUrl", tc.alias).
					Return(int64(1), tc.mockError).
					Once()
			}

			handler := remove.New(slogdiscard.NewDiscardLogger(), urlRemoverMock)
			req, err := http.NewRequest(http.MethodDelete, "/"+tc.alias, nil)
			require.NoError(t, err)

			rr := httptest.NewRecorder()
			router := chi.NewRouter()
			router.Delete("/{alias}", handler.ServeHTTP)
			router.ServeHTTP(rr, req)
			require.Equal(t, rr.Code, http.StatusOK)

			body := rr.Body.String()
			var resp remove.Response
			require.NoError(t, json.Unmarshal([]byte(body), &resp))
			require.Equal(t, tc.respError, resp.Error)
		})
	}
}
