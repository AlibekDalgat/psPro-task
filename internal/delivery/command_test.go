package delivery

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"psPro-task/internal/models"
	"psPro-task/internal/service"
	mock_service "psPro-task/internal/service/mocks"
	"strconv"
	"testing"
)

func TestDelivery_createCommand(t *testing.T) {
	type mockBehavior func(s *mock_service.MockCommand, command models.Command)

	testTable := []struct {
		name           string
		inputBody      string
		inputCommand   models.Command
		mockBehavior   mockBehavior
		expStatusCode  int
		expRequestBody string
	}{
		{
			name:      "OK",
			inputBody: `{"script": "echo test"}`,
			inputCommand: models.Command{
				Script: "echo test",
			},
			mockBehavior: func(s *mock_service.MockCommand, command models.Command) {
				s.EXPECT().CreateCommand(command).Return(1, nil)
				s.EXPECT().ExecuteCommand(gomock.Any(), gomock.Any()).DoAndReturn(func(_ int, _ string) {
				}).AnyTimes()
			},
			expStatusCode:  http.StatusOK,
			expRequestBody: `{"id":1}`,
		},
		{
			name:           "Wrong json",
			inputBody:      `{"ssss": "fail"}`,
			inputCommand:   models.Command{},
			mockBehavior:   func(s *mock_service.MockCommand, command models.Command) {},
			expStatusCode:  http.StatusBadRequest,
			expRequestBody: `{"message":"неверное содержание json"}`,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			command := mock_service.NewMockCommand(c)
			testCase.mockBehavior(command, testCase.inputCommand)

			services := &service.Service{Command: command}
			handler := NewHandler(services)

			r := gin.New()
			r.POST("/api", handler.createCommand)

			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/api", bytes.NewBufferString(testCase.inputBody))

			r.ServeHTTP(w, req)
			assert.Equal(t, w.Code, testCase.expStatusCode)
			assert.Equal(t, w.Body.String(), testCase.expRequestBody)

		})
	}
}

func TestDelivery_listCommands(t *testing.T) {
	testTable := []struct {
		name             string
		mockBehavior     func(r *mock_service.MockCommand)
		expectedHTTPCode int
		expectedResponse string
	}{
		{
			name: "OK",
			mockBehavior: func(command *mock_service.MockCommand) {
				stdout := "test"
				commands := []models.Command{
					{
						Id:     1,
						Stdout: &stdout,
						Stderr: nil,
						Script: "echo test",
					},
				}
				command.EXPECT().GetAllCommands().Return(commands, nil)
			},
			expectedHTTPCode: http.StatusOK,
			expectedResponse: `[
                {
                    "id": 1,
					"echo test",
					"created_at": "2024-05-13T20:11:55.611813Z"
                }
            ]`,
		},
		{
			name: "Fail",
			mockBehavior: func(command *mock_service.MockCommand) {
				command.EXPECT().GetAllCommands().Return(nil, errors.New("some error"))
			},
			expectedHTTPCode: http.StatusInternalServerError,
			expectedResponse: `{"message":"Ошибка во время работы сервера"}`,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()

			commandMock := mock_service.NewMockCommand(mockCtrl)
			testCase.mockBehavior(commandMock)

			services := &service.Service{Command: commandMock}
			handler := NewHandler(services)

			r := gin.New()
			r.GET("/api", handler.listCommands)

			w := httptest.NewRecorder()
			req := httptest.NewRequest("GET", "/api", nil)

			r.ServeHTTP(w, req)

			assert.Equal(t, testCase.expectedHTTPCode, w.Code)
			assert.Equal(t, testCase.expectedResponse, w.Body.String())
		})
	}
}

func TestDelivery_stopCommand(t *testing.T) {
	testTable := []struct {
		name             string
		input            string
		mockBehavior     func(r *mock_service.MockCommand, input int)
		expectedHTTPCode int
		expectedResponse string
	}{
		{
			name:  "OK",
			input: "1",
			mockBehavior: func(command *mock_service.MockCommand, input int) {
				command.EXPECT().StopCommand(input).Return(nil)
			},
			expectedHTTPCode: http.StatusOK,
			expectedResponse: `{"status":"Сигнал отправлен"}`,
		},
		{
			name:             "Wrong id parametr",
			input:            "abc",
			mockBehavior:     func(command *mock_service.MockCommand, input int) {},
			expectedHTTPCode: http.StatusBadRequest,
			expectedResponse: `{"message":"Неверный параметр id"}`,
		},
		{
			name:  "Stop stopped command",
			input: "1",
			mockBehavior: func(command *mock_service.MockCommand, input int) {
				command.EXPECT().StopCommand(input).Return(errors.New("Команда 1 уже остановлена"))
			},
			expectedHTTPCode: http.StatusInternalServerError,
			expectedResponse: `{"message":"Команда 1 уже остановлена"}`,
		},
		{
			name:  "Non-exist running command",
			input: "2",
			mockBehavior: func(command *mock_service.MockCommand, input int) {
				command.EXPECT().StopCommand(input).Return(errors.New("Команды с id 2 не найдено"))
			},
			expectedHTTPCode: http.StatusInternalServerError,
			expectedResponse: `{"message":"Команды с id 2 не найдено"}`,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			commandMock := mock_service.NewMockCommand(c)

			if id, err := strconv.Atoi(testCase.input); err == nil {
				testCase.mockBehavior(commandMock, id)
			}

			services := &service.Service{Command: commandMock}
			handler := NewHandler(services)

			r := gin.New()
			r.PATCH("/api/stop/:id", handler.stopCommand)

			w := httptest.NewRecorder()
			req := httptest.NewRequest("PATCH", fmt.Sprintf("/api/stop/%s", testCase.input), nil)

			r.ServeHTTP(w, req)
			assert.Equal(t, w.Code, testCase.expectedHTTPCode)
			assert.Equal(t, w.Body.String(), testCase.expectedResponse)
		})
	}
}

func TestDelivery_startCommand(t *testing.T) {
	testTable := []struct {
		name             string
		input            string
		mockBehavior     func(r *mock_service.MockCommand, input int)
		expectedHTTPCode int
		expectedResponse string
	}{
		{
			name:  "OK",
			input: "1",
			mockBehavior: func(command *mock_service.MockCommand, input int) {
				command.EXPECT().StartCommand(input).Return(nil)
			},
			expectedHTTPCode: http.StatusOK,
			expectedResponse: `{"status":"Сигнал отправлен"}`,
		},
		{
			name:             "Wrong id parametr",
			input:            "abc",
			mockBehavior:     func(command *mock_service.MockCommand, input int) {},
			expectedHTTPCode: http.StatusBadRequest,
			expectedResponse: `{"message":"Неверный параметр id"}`,
		},
		{
			name:  "Start running command",
			input: "1",
			mockBehavior: func(command *mock_service.MockCommand, input int) {
				command.EXPECT().StartCommand(input).Return(errors.New("Команда 1 уже запущена"))
			},
			expectedHTTPCode: http.StatusInternalServerError,
			expectedResponse: `{"message":"Команда 1 уже запущена"}`,
		},
		{
			name:  "Non-exist running command",
			input: "2",
			mockBehavior: func(command *mock_service.MockCommand, input int) {
				command.EXPECT().StartCommand(input).Return(errors.New("Команды с id 2 не найдено"))
			},
			expectedHTTPCode: http.StatusInternalServerError,
			expectedResponse: `{"message":"Команды с id 2 не найдено"}`,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			commandMock := mock_service.NewMockCommand(c)

			if id, err := strconv.Atoi(testCase.input); err == nil {
				testCase.mockBehavior(commandMock, id)
			}

			services := &service.Service{Command: commandMock}
			handler := NewHandler(services)

			r := gin.New()
			r.PATCH("/api/start/:id", handler.startCommand)

			w := httptest.NewRecorder()
			req := httptest.NewRequest("PATCH", fmt.Sprintf("/api/start/%s", testCase.input), nil)

			r.ServeHTTP(w, req)
			assert.Equal(t, w.Code, testCase.expectedHTTPCode)
			assert.Equal(t, w.Body.String(), testCase.expectedResponse)
		})
	}
}

func TestDelivery_killCommand(t *testing.T) {
	testTable := []struct {
		name             string
		input            string
		mockBehavior     func(r *mock_service.MockCommand, input int)
		expectedHTTPCode int
		expectedResponse string
	}{
		{
			name:  "OK",
			input: "1",
			mockBehavior: func(command *mock_service.MockCommand, input int) {
				command.EXPECT().KillCommand(input).Return(nil)
			},
			expectedHTTPCode: http.StatusOK,
			expectedResponse: `{"status":"Сигнал отправлен"}`,
		},
		{
			name:             "Wrong id parametr",
			input:            "abc",
			mockBehavior:     func(command *mock_service.MockCommand, input int) {},
			expectedHTTPCode: http.StatusBadRequest,
			expectedResponse: `{"message":"Неверный параметр id"}`,
		},
		{
			name:  "Non-exist running command",
			input: "2",
			mockBehavior: func(command *mock_service.MockCommand, input int) {
				command.EXPECT().KillCommand(input).Return(errors.New("Команды с id 2 не найдено"))
			},
			expectedHTTPCode: http.StatusInternalServerError,
			expectedResponse: `{"message":"Команды с id 2 не найдено"}`,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			commandMock := mock_service.NewMockCommand(c)

			if id, err := strconv.Atoi(testCase.input); err == nil {
				testCase.mockBehavior(commandMock, id)
			}

			services := &service.Service{Command: commandMock}
			handler := NewHandler(services)

			r := gin.New()
			r.PATCH("/api/kill/:id", handler.killCommand)

			w := httptest.NewRecorder()
			req := httptest.NewRequest("PATCH", fmt.Sprintf("/api/kill/%s", testCase.input), nil)

			r.ServeHTTP(w, req)
			assert.Equal(t, w.Code, testCase.expectedHTTPCode)
			assert.Equal(t, w.Body.String(), testCase.expectedResponse)
		})
	}
}
