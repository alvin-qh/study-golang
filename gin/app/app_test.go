package app

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"study-gin/app/routes"
	"study-gin/core/server"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetHelloWithoutParameters(t *testing.T) {
	req, _ := http.NewRequest(http.MethodGet, "/api/hello", nil)

	w := httptest.NewRecorder()
	server.Engine.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "{\"name\":\"Alvin\",\"text\":\"Hello Alvin\"}", w.Body.String())
}

func TestPostHello(t *testing.T) {
	data, _ := json.Marshal(&routes.Agreeing{
		Name:   "Emma",
		Gender: "F",
		Age:    32,
	})

	req, _ := http.NewRequest(http.MethodPost, "/api/hello", bytes.NewBuffer(data))

	w := httptest.NewRecorder()
	server.Engine.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var answer routes.Answer
	json.Unmarshal(w.Body.Bytes(), &answer)
	assert.True(t, answer.Ok)
	assert.Equal(t, "你好, Emma小姐", answer.Answer)
}

func TestPostHelloWithWrongData(t *testing.T) {
	data, _ := json.Marshal(&routes.Agreeing{
		Name:   "Emma",
		Gender: "X",
		Age:    32,
	})

	req, _ := http.NewRequest(http.MethodPost, "/api/hello", bytes.NewBuffer(data))

	w := httptest.NewRecorder()
	server.Engine.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var er routes.ErrorResult
	json.Unmarshal(w.Body.Bytes(), &er)

	assert.Equal(t, "input_error", er.Code)
	assert.Len(t, er.Errors, 1)
	assert.Equal(t, "Gender必须是[F M]中的一个", er.Errors["gender"])
}
