package app

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	"study/web/gin/app/routes"
	"study/web/gin/core/server"

	"github.com/stretchr/testify/assert"
)

// 测试 `ApiGetUser` 路由方法
//
// 发送 `GET` 请求, 确认响应结果
func TestGetUser(t *testing.T) {
	// 创建一个请求对象
	req, _ := http.NewRequest(http.MethodGet, "/api/users", nil)

	// 创建一个测试用的 `ResponseRecorder` 对象
	w := httptest.NewRecorder()
	// 启动 http 服务并发送请求
	server.Engine.ServeHTTP(w, req)

	// 确认请求处理正确
	assert.Equal(t, http.StatusOK, w.Code)

	// 结果反序列化
	var resp routes.ResponseData
	json.Unmarshal(w.Body.Bytes(), &resp)

	// 确认返回正确的 code
	assert.Equal(t, routes.OkCode, resp.Code)

	//确认返回的 Payload 为长度为 2 的切片
	payloads := resp.Payload.([]any)
	assert.Len(t, payloads, 2)

	var user routes.User
	MapToStruct.Decode(payloads[0].(map[string]any), &user)
	assert.Equal(t, "001", user.Id)
	assert.Equal(t, "Alvin", user.Name)
	assert.Equal(t, routes.GenderM, user.Gender)
	assert.Equal(t, "1981-03-17", user.Birthday.Format(time.DateOnly))

	MapToStruct.Decode(payloads[1].(map[string]any), &user)
	assert.Equal(t, "002", user.Id)
	assert.Equal(t, "Emma", user.Name)
	assert.Equal(t, routes.GenderF, user.Gender)
	assert.Equal(t, "1985-03-29", user.Birthday.Format(time.DateOnly))
}

// 测试 `ApiGetUser` 路由方法
//
// 发送 `GET` 请求并传递 URL 参数, 确认响应结果
func TestGetUserWithParameters(t *testing.T) {
	// 创建一个请求对象
	req, _ := http.NewRequest(
		http.MethodGet,
		fmt.Sprintf("/api/users?name=%v", url.QueryEscape("Emma")),
		nil,
	)

	// 创建一个测试用的 `ResponseRecorder` 对象
	w := httptest.NewRecorder()
	// 启动 http 服务并发送请求
	server.Engine.ServeHTTP(w, req)

	// 确认请求处理正确
	assert.Equal(t, http.StatusOK, w.Code)

	// 结果反序列化
	var resp routes.ResponseData
	json.Unmarshal(w.Body.Bytes(), &resp)

	// 确认返回正确的 code
	assert.Equal(t, routes.OkCode, resp.Code)

	//确认返回的 Payload 为长度为 2 的切片
	payloads := resp.Payload.([]any)
	assert.Len(t, payloads, 1)

	var user routes.User
	MapToStruct.Decode(payloads[0].(map[string]any), &user)
	assert.Equal(t, "001", user.Id)
	assert.Equal(t, "Emma", user.Name)
	assert.Equal(t, routes.GenderM, user.Gender)
	assert.Equal(t, "1981-03-17", user.Birthday.Format(time.DateOnly))
}

// 测试 `ApiPostUser` 路由方法
//
// 发送 `POST` 请求, 确认响应结果
func TestPostUser(t *testing.T) {
	data, _ := json.Marshal(&routes.UserForm{
		Name:       "Emma",
		Gender:     routes.GenderF,
		BirthYear:  1985,
		BirthMonth: 3,
		BirthDay:   29,
	})

	req, _ := http.NewRequest(http.MethodPost, "/api/users", bytes.NewBuffer(data))

	// 创建一个测试用的 `ResponseRecorder` 对象
	w := httptest.NewRecorder()
	// 启动 http 服务并发送请求
	server.Engine.ServeHTTP(w, req)

	// 确认请求处理正确
	assert.Equal(t, http.StatusOK, w.Code)

	// 结果反序列化
	var resp routes.ResponseData
	json.Unmarshal(w.Body.Bytes(), &resp)

	// 确认返回正确的 code
	assert.Equal(t, routes.OkCode, resp.Code)

	var user routes.User

	// 将 Payload 反序列化为 user 结构体
	MapToStruct.Decode(resp.Payload.(map[string]any), &user)

	// 确认 Payload 符合预期
	assert.Equal(t, "003", user.Id)
	assert.Equal(t, "Emma", user.Name)
	assert.Equal(t, routes.GenderF, user.Gender)
	assert.Equal(t, "1985-03-29", user.Birthday.Format(time.DateOnly))
}

// 测试 `ApiPostUser` 路由方法
//
// 发送 `POST` 请求, 包含错误的请求 body, 确认响应结果中包含的错误信息
func TestPostUserByWrongData(t *testing.T) {
	data, _ := json.Marshal(&routes.UserForm{
		Name:       "Emma",
		Gender:     routes.Gender("X"), // 传递错误的请求参数
		BirthYear:  1985,
		BirthMonth: 3,
		BirthDay:   29,
	})

	req, _ := http.NewRequest(http.MethodPost, "/api/users", bytes.NewBuffer(data))

	// 创建一个测试用的 `ResponseRecorder` 对象
	w := httptest.NewRecorder()
	// 启动 http 服务并发送请求
	server.Engine.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	// 结果反序列化
	var resp routes.ResponseData
	json.Unmarshal(w.Body.Bytes(), &resp)

	// 确认返回 code 400
	assert.Equal(t, routes.InputErrorCode, resp.Code)

	var er routes.ErrorResult

	// 将 Payload 反序列化为 ErrorResult 结构体
	MapToStruct.Decode(resp.Payload.(map[string]any), &er)

	// 确认 Payload 符合预期
	assert.Len(t, er.ErrorFields, 1)
	assert.Equal(t, "gender", er.ErrorFields[0].Name)
	assert.Equal(t, "Gender必须是[F M]中的一个", er.ErrorFields[0].Error)
}

// 测试 `ApiGetUserById` 路由方法
//
// 发送 `GET` 请求, 通过用户 ID 获取 User 对象
func TestApiGetUserById(t *testing.T) {
	req, _ := http.NewRequest(http.MethodGet, "/api/users/002", nil)

	// 创建一个测试用的 `ResponseRecorder` 对象
	w := httptest.NewRecorder()
	// 启动 http 服务并发送请求
	server.Engine.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	// 结果反序列化
	var resp routes.ResponseData
	json.Unmarshal(w.Body.Bytes(), &resp)

	// 确认返回结果正确
	assert.Equal(t, routes.OkCode, resp.Code)

	var user routes.User

	MapToStruct.Decode(resp.Payload, &user)
	assert.Equal(t, "002", user.Id)
	assert.Equal(t, "Emma", user.Name)
	assert.Equal(t, routes.GenderF, user.Gender)
	assert.Equal(t, "1985-03-29", user.Birthday.Format(time.DateOnly))
}
