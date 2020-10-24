package lannister

import (
    "github.com/kataras/iris/v12/httptest"
    "testing"
)

func TestPing(t *testing.T) {
    application := Application()
    e := httptest.New(t, application.app)

    e.GET("/").Expect().Status(httptest.StatusOK).Body().Equal("")

    e.GET("/lannister").Expect().Status(httptest.StatusOK).
        Body().Equal("Lannisters Always pay their debts.")
    e.GET("/lannister/welcome").Expect().Status(httptest.StatusOK).
        Body().Equal("Lannisters Always pay their debts.")
}
