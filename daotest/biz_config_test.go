package daotest

import (
    "github.com/CharLemAznable/go-lannister/app"
    "github.com/CharLemAznable/go-lannister/base"
    "github.com/CharLemAznable/gokits"
    "github.com/CharLemAznable/sqlx"
    "github.com/kataras/iris/v12/httptest"
    "github.com/stretchr/testify/assert"
    "testing"
)

func TestConfig(t *testing.T) {
    a := assert.New(t)

    for _, testConfig := range TestConfigSet {
        application := app.Application(func(config *base.Config) {
            config.DriverName = testConfig["DriverName"]
            config.DataSourceName = testConfig["DataSourceName"]
        })
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
        responseConfig = e.POST("/lannister/1001/1001/TEST-API/config").
            WithQuery("nonsense", "config").
            WithQuery("signature", signatureConfig).
            WithJSON(map[string]string{
                "param-name-0": "param-value-0",
                "param-name-1": "param-value-1",
            }).
            Expect().Status(httptest.StatusOK).Body()
        a.Equal("SUCCESS", responseConfig.Raw())

        db, _ := sqlx.Open(testConfig["DriverName"], testConfig["DataSourceName"])
        merchantApiParams := make([]struct {
            ParamName  string `db:"ParamName"`
            ParamValue string `db:"ParamValue"`
        }, 0)
        _ = db.Select(&merchantApiParams, `
select mp.param_name  as "ParamName"
      ,mp.param_value as "ParamValue"
  from merchant_api_params  mp
 where mp.merchant_id       = '1001'
   and mp.api_name          = 'test-api'
`)
        a.Equal(2, len(merchantApiParams))
        a.Equal("param-name-0", merchantApiParams[0].ParamName)
        a.Equal("param-value-0", merchantApiParams[0].ParamValue)
        a.Equal("param-name-1", merchantApiParams[1].ParamName)
        a.Equal("param-value-1", merchantApiParams[1].ParamValue)
    }
}
