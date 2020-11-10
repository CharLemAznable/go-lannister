package apptest

import (
    "github.com/CharLemAznable/go-lannister/app"
    "github.com/CharLemAznable/go-lannister/base"
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
        &base.BaseResp{}).(*base.BaseResp)
    a.Equal("CONFIG_PARAMS_ILLEGAL", resultConfig.ErrorCode)
    a.Equal("Config Params is Illegal", resultConfig.ErrorDesc)

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
    resultConfig = gokits.UnJson(responseConfig.Raw(),
        &base.BaseResp{}).(*base.BaseResp)
    a.Equal("CONFIG_FAILED", resultConfig.ErrorCode)
    a.Equal("MockError", resultConfig.ErrorDesc)

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
    resultConfig = gokits.UnJson(responseConfig.Raw(),
        &base.BaseResp{}).(*base.BaseResp)
    a.Equal("", resultConfig.ErrorCode)
    a.Equal("", resultConfig.ErrorDesc)

    a.Equal("param-value-0", merchantApiParams["1001"]["test-api"]["param-name-0"])
    a.Equal("param-value-1", merchantApiParams["1001"]["test-api"]["param-name-1"])
}
