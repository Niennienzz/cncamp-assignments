package api

import (
	"io/ioutil"
	"net/http"
	"os"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
)

func TestBasicRoutes(t *testing.T) {
	testcases := []struct {
		description string

		route  string
		method string

		expectedError bool
		expectedCode  int
		expectedBody  string
	}{
		{
			description:   "health check",
			route:         "/healthz",
			method:        fiber.MethodGet,
			expectedError: false,
			expectedCode:  200,
			expectedBody:  "OK",
		},
		{
			description:   "non existing route",
			route:         "/i-dont-exist",
			method:        fiber.MethodGet,
			expectedError: false,
			expectedCode:  404,
			expectedBody:  "Cannot GET /i-dont-exist",
		},
	}

	os.Setenv("SQLITE_FILE", "sqlite_test.db")
	srv := New().(*apiImpl)

	for _, tc := range testcases {
		req, _ := http.NewRequest(
			tc.method,
			tc.route,
			nil,
		)

		res, err := srv.fiberApp.Test(req, -1)
		assert.Equalf(t, tc.expectedError, err != nil, tc.description)
		if tc.expectedError {
			continue
		}
		assert.Equalf(t, tc.expectedCode, res.StatusCode, tc.description)

		body, err := ioutil.ReadAll(res.Body)
		assert.Nilf(t, err, tc.description)
		assert.Equalf(t, tc.expectedBody, string(body), tc.description)
	}
}
