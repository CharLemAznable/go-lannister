package daotest

import (
    "github.com/CharLemAznable/go-lannister/app"
    "github.com/CharLemAznable/go-lannister/base"
    "github.com/CharLemAznable/gokits"
    "github.com/kataras/iris/v12/httptest"
    "github.com/stretchr/testify/assert"
    "testing"
)

func TestAccessor(t *testing.T) {
    a := assert.New(t)

    for _, testConfig := range TestConfigSet {
        application := app.Application(func(config *base.Config) {
            config.DriverName = testConfig["DriverName"]
            config.DataSourceName = testConfig["DataSourceName"]
        })
        e := httptest.New(t, application.App())

        signatureQuery, _ := gokits.SHA1WithRSA.SignBase64ByRSAKeyString(
            "nonsense=query", PrivateKeyString)
        responseQuery := e.GET("/lannister/1001/query-info").
            WithQuery("nonsense", "query").
            WithQuery("signature", signatureQuery).
            Expect().Status(httptest.StatusOK).Body()
        resultQuery := gokits.UnJson(responseQuery.Raw(),
            &base.AccessorManage{}).(*base.AccessorManage)
        a.Equal("1001", resultQuery.AccessorId)
        a.Equal("1001", resultQuery.AccessorName)
        a.Equal(PublicKeyString, resultQuery.AccessorPubKey)
        a.Equal("", resultQuery.PayNotifyUrl)
        a.Equal("", resultQuery.RefundNotifyUrl)
        a.Equal("", resultQuery.PubKey)

        signatureUpdate, _ := gokits.SHA1WithRSA.SignBase64ByRSAKeyString(
            "accessorName=test&accessorPubKey="+PublicKeyString+"&nonsense=update&"+
                "payNotifyUrl=PayNotifyUrl&refundNotifyUrl=RefundNotifyUrl", PrivateKeyString)
        responseUpdate := e.POST("/lannister/1001/update-info").
            WithQuery("nonsense", "update").
            WithQuery("signature", signatureUpdate).
            WithJSON(&base.AccessorManage{AccessorName: "test",
                AccessorPubKey:  PublicKeyString,
                PayNotifyUrl:    "PayNotifyUrl",
                RefundNotifyUrl: "RefundNotifyUrl",}).
            Expect().Status(httptest.StatusOK).Body()
        resultUpdate := gokits.UnJson(responseUpdate.Raw(),
            &base.AccessorManage{}).(*base.AccessorManage)
        a.Equal("1001", resultUpdate.AccessorId)
        a.Equal("test", resultUpdate.AccessorName)
        a.Equal(PublicKeyString, resultUpdate.AccessorPubKey)
        a.Equal("PayNotifyUrl", resultUpdate.PayNotifyUrl)
        a.Equal("RefundNotifyUrl", resultUpdate.RefundNotifyUrl)
        a.Equal("", resultUpdate.PubKey)

        signatureReset, _ := gokits.SHA1WithRSA.SignBase64ByRSAKeyString(
            "nonsense=reset", PrivateKeyString)
        responseReset := e.POST("/lannister/1001/reset-info").
            WithQuery("nonsense", "reset").
            WithQuery("signature", signatureReset).
            Expect().Status(httptest.StatusOK).Body()
        resultReset := gokits.UnJson(responseReset.Raw(),
            &base.AccessorManage{}).(*base.AccessorManage)
        a.Equal("1001", resultReset.AccessorId)
        a.Equal("test", resultReset.AccessorName)
        a.Equal(PublicKeyString, resultReset.AccessorPubKey)
        a.Equal("PayNotifyUrl", resultReset.PayNotifyUrl)
        a.Equal("RefundNotifyUrl", resultReset.RefundNotifyUrl)
        a.NotEqual("", resultReset.PubKey)
    }
}
