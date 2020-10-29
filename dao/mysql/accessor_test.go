package mysql

import (
    "github.com/CharLemAznable/go-lannister/app"
    "github.com/CharLemAznable/go-lannister/base"
    . "github.com/CharLemAznable/go-lannister/elf"
    "github.com/CharLemAznable/gokits"
    "github.com/kataras/iris/v12/httptest"
    "github.com/stretchr/testify/assert"
    "testing"
)

func TestAccessor(t *testing.T) {
    a := assert.New(t)

    application := app.Application(func(config *app.Config) {
        config.DriverName = "mysql"
        config.DataSourceName = mysqlDataSourceName
    })
    e := httptest.New(t, application.App())

    signatureQuery, _ := SHA1WithRSA.SignBase64ByKeyString(
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

    signatureUpdate, _ := SHA1WithRSA.SignBase64ByKeyString(
        "accessorName=test&accessorPubKey="+PublicKeyString+"&nonsense=update&"+
            "payNotifyUrl=PayNotifyUrl&refundNotifyUrl=RefundNotifyUrl", PrivateKeyString)
    e.POST("/lannister/1001/update-info").
        WithQuery("nonsense", "update").
        WithQuery("signature", signatureUpdate).
        WithJSON(&base.AccessorManage{AccessorName: "test",
            AccessorPubKey:  PublicKeyString,
            PayNotifyUrl:    "PayNotifyUrl",
            RefundNotifyUrl: "RefundNotifyUrl",}).
        Expect().Status(httptest.StatusOK).Body().Equal("SUCCESS")
    signatureQuery, _ = SHA1WithRSA.SignBase64ByKeyString(
        "nonsense=query", PrivateKeyString)
    responseQuery = e.GET("/lannister/1001/query-info").
        WithQuery("nonsense", "query").
        WithQuery("signature", signatureQuery).
        Expect().Status(httptest.StatusOK).Body()
    resultQuery = gokits.UnJson(responseQuery.Raw(),
        &base.AccessorManage{}).(*base.AccessorManage)
    a.Equal("1001", resultQuery.AccessorId)
    a.Equal("test", resultQuery.AccessorName)
    a.Equal(PublicKeyString, resultQuery.AccessorPubKey)
    a.Equal("PayNotifyUrl", resultQuery.PayNotifyUrl)
    a.Equal("RefundNotifyUrl", resultQuery.RefundNotifyUrl)
    a.Equal("", resultQuery.PubKey)

    signatureReset, _ := SHA1WithRSA.SignBase64ByKeyString(
        "nonsense=reset", PrivateKeyString)
    responseReset := e.POST("/lannister/1001/reset-info").
        WithQuery("nonsense", "reset").
        WithQuery("signature", signatureReset).
        Expect().Status(httptest.StatusOK).Body()
    resultReset := gokits.UnJson(responseReset.Raw(),
        &base.AccessorManage{}).(*base.AccessorManage)
    responseQuery = e.GET("/lannister/1001/query-info").
        WithQuery("nonsense", "query").
        WithQuery("signature", signatureQuery).
        Expect().Status(httptest.StatusOK).Body()
    resultQuery = gokits.UnJson(responseQuery.Raw(),
        &base.AccessorManage{}).(*base.AccessorManage)
    a.Equal("1001", resultQuery.AccessorId)
    a.Equal("test", resultQuery.AccessorName)
    a.Equal(PublicKeyString, resultQuery.AccessorPubKey)
    a.Equal("PayNotifyUrl", resultQuery.PayNotifyUrl)
    a.Equal("RefundNotifyUrl", resultQuery.RefundNotifyUrl)
    a.Equal(resultReset.PubKey, resultQuery.PubKey)
}
