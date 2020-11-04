package apptest

import (
    "github.com/CharLemAznable/go-lannister/app"
    "github.com/CharLemAznable/gokits"
    "github.com/kataras/iris/v12/httptest"
    "github.com/stretchr/testify/assert"
    "testing"
)

func TestConfig(t *testing.T) {
    a := assert.New(t)

    application := app.Application()
    e := httptest.New(t, application.App())

    signatureConfig, _ := gokits.SHA1WithRSA.SignBase64ByRSAKeyString(
        "nonsense=config", PrivateKeyString)
    responseConfig := e.POST("/lannister/1001/1001/TEST-API/config").
        WithQuery("nonsense", "config").
        WithQuery("signature", signatureConfig).
        Expect().Status(httptest.StatusOK).Body()
    resultConfig := gokits.UnJson(responseConfig.Raw(),
        &map[string]string{}).(*map[string]string)
    a.Equal("CONFIG_PARAMS_ILLEGAL", (*resultConfig)["errorCode"])
    a.Equal("Config Params is Illegal", (*resultConfig)["errorDesc"])

    signatureConfig, _ = gokits.SHA1WithRSA.SignBase64ByRSAKeyString(
        "nonsense=config&param-name-0=param-value-0&param-name-1=param-value-1", PrivateKeyString)
    responseConfig = e.POST("/lannister/1001/1001/ERROR-API/config").
        WithQuery("nonsense", "config").
        WithQuery("signature", signatureConfig).
        WithJSON(map[string]string{
            "param-name-0": "param-value-0",
            "param-name-1": "param-value-1",
        }).
        Expect().Status(httptest.StatusOK).Body()
    a.Equal("FAILED", responseConfig.Raw())

    signatureConfig, _ = gokits.SHA1WithRSA.SignBase64ByRSAKeyString(
        "nonsense=config&param-name-0=param-value-0&param-name-1=param-value-1", PrivateKeyString)
    responseConfig = e.POST("/lannister/1001/1001/TEST-API/config").
        WithQuery("nonsense", "config").
        WithQuery("signature", signatureConfig).
        WithJSON(map[string]string{
            "param-name-0": "param-value-0",
            "param-name-1": "param-value-1",
        }).
        Expect().Status(httptest.StatusOK).Body()
    a.Equal("SUCCESS", responseConfig.Raw())

    a.Equal("param-value-0", merchantApiParams["1001"]["test-api"]["param-name-0"])
    a.Equal("param-value-1", merchantApiParams["1001"]["test-api"]["param-name-1"])
}
