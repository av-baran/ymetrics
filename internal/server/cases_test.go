package server

import "net/http"

type handlersTests struct {
	name          string
	request       string
	method        string
	body          string
	expectedCode  int
	expectedBody  string
	expectedZBody string
	wantErr       bool
}

var (
	handlersTestCases = []handlersTests{
		{
			name:         "GET params",
			request:      "/update/gauge/g1/1.0",
			body:         "",
			method:       http.MethodGet,
			expectedCode: http.StatusMethodNotAllowed,
			expectedBody: "",
			wantErr:      true,
		},
		{
			name:         "PUT params",
			request:      "/update/gauge/g1/1.0",
			body:         "",
			method:       http.MethodPut,
			expectedCode: http.StatusMethodNotAllowed,
			expectedBody: "",
			wantErr:      true,
		},
		{
			name:         "DELETE params",
			request:      "/update/gauge/g1/1.0",
			body:         "",
			method:       http.MethodDelete,
			expectedCode: http.StatusMethodNotAllowed,
			expectedBody: "",
			wantErr:      true,
		},
		{
			name:         "POST - params update bad type",
			request:      "/update/unknowntype/g1/1.0",
			body:         "",
			method:       http.MethodPost,
			expectedCode: http.StatusNotImplemented,
			expectedBody: "",
			wantErr:      true,
		},
		{
			name:         "POST - params update gauge with bad value",
			request:      "/update/gauge/g1/fuuuu",
			body:         "",
			method:       http.MethodPost,
			expectedCode: http.StatusBadRequest,
			expectedBody: "",
			wantErr:      true,
		},
		{
			name:         "POST - params update counter with bad value",
			request:      "/update/counter/c1/fuuuu",
			body:         "",
			method:       http.MethodPost,
			expectedCode: http.StatusBadRequest,
			expectedBody: "",
			wantErr:      true,
		},
		{
			name:         "POST - params update gauge ok",
			request:      "/update/gauge/g1/1.01",
			body:         "",
			method:       http.MethodPost,
			expectedCode: http.StatusOK,
			expectedBody: "",
			wantErr:      false,
		},
		{
			name:         "POST - params update counter ok",
			request:      "/update/counter/c1/1",
			body:         "",
			method:       http.MethodPost,
			expectedCode: http.StatusOK,
			expectedBody: "",
			wantErr:      false,
		},
		{
			name:          "GET - params get gauge ok",
			request:       "/value/gauge/g1",
			body:          "",
			method:        http.MethodGet,
			expectedCode:  http.StatusOK,
			expectedBody:  "1.01",
			expectedZBody: "1.01",
			wantErr:       false,
		},
		{
			name:          "GET - params get counter ok",
			request:       "/value/counter/c1",
			body:          "",
			method:        http.MethodGet,
			expectedCode:  http.StatusOK,
			expectedBody:  "1",
			expectedZBody: "1",
			wantErr:       false,
		},
		{
			name:         "GET - params get unknown gauge",
			request:      "/value/gauge/unknownname",
			body:         "",
			method:       http.MethodGet,
			expectedCode: http.StatusNotFound,
			expectedBody: "",
			wantErr:      true,
		},
		{
			name:         "GET - params get unknown counter",
			request:      "/value/counter/unknownname",
			body:         "",
			method:       http.MethodGet,
			expectedCode: http.StatusNotFound,
			expectedBody: "",
			wantErr:      true,
		},
	}

	jsonHandlersTestCases = []handlersTests{
		{
			name:          "json counter POST - OK",
			request:       "/update/",
			method:        http.MethodPost,
			body:          `{"id":"counter1","type":"counter","delta":5}`,
			expectedCode:  http.StatusOK,
			expectedBody:  `{"id":"counter1","type":"counter","delta":5}`,
			expectedZBody: `{"id":"counter1","type":"counter","delta":5}`,
			wantErr:       false,
		},
		{
			name:          "json second counter POST - OK",
			request:       "/update/",
			method:        http.MethodPost,
			body:          `{"id":"counter1","type":"counter","delta":5}`,
			expectedCode:  http.StatusOK,
			expectedBody:  `{"id":"counter1","type":"counter","delta":15}`,
			expectedZBody: `{"id":"counter1","type":"counter","delta":10}`,
			wantErr:       false,
		},
		{
			name:          "json gauge POST - OK",
			request:       "/update/",
			method:        http.MethodPost,
			body:          `{"id":"gauge1","type":"gauge","value":5}`,
			expectedCode:  http.StatusOK,
			expectedBody:  `{"id":"gauge1","type":"gauge","value":5}`,
			expectedZBody: `{"id":"gauge1","type":"gauge","value":5}`,
			wantErr:       false,
		},
		{
			name:          "json gauge test change POST - OK",
			request:       "/update/",
			method:        http.MethodPost,
			body:          `{"id":"gauge1","type":"gauge","value":0.1}`,
			expectedCode:  http.StatusOK,
			expectedBody:  `{"id":"gauge1","type":"gauge","value":0.1}`,
			expectedZBody: `{"id":"gauge1","type":"gauge","value":0.1}`,
			wantErr:       false,
		},
		{
			name:         "json gauge POST - unknown type",
			request:      "/update/",
			method:       http.MethodPost,
			body:         `{"id":"unknown_name1","type":"unknown","value":5}`,
			expectedCode: http.StatusNotImplemented,
			expectedBody: "",
			wantErr:      true,
		},
	}
)
