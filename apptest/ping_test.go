package apptest

import (
    "github.com/CharLemAznable/go-lannister/app"
    "github.com/kataras/iris/v12/httptest"
    "testing"
)

func TestPing(t *testing.T) {
    application := app.Application(func(config *app.Config) {
        config.LogLevel = "debug"
    })
    e := httptest.New(t, application.App())

    e.GET("/").Expect().Status(httptest.StatusOK).Body().Equal("")

    e.GET("/lannister").Expect().Status(httptest.StatusOK).
        Body().Equal("Lannisters Always pay their debts.")
    e.GET("/lannister/welcome").Expect().Status(httptest.StatusOK).
        Body().Equal("Lannisters Always pay their debts.")
}
