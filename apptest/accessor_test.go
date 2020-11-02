package apptest

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

    application := app.Application(func(config *app.Config) {
        config.AccessorVerifyCacheInMills = 1
        config.MerchantVerifyCacheInMills = 1
        config.AccessorMerchantVerifyCacheInMills = 1
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
    e.POST("/lannister/1001/update-info").
        WithQuery("nonsense", "update").
        WithQuery("signature", signatureUpdate).
        WithJSON(&base.AccessorManage{AccessorName: "test",
            AccessorPubKey:  PublicKeyString,
            PayNotifyUrl:    "PayNotifyUrl",
            RefundNotifyUrl: "RefundNotifyUrl",}).
        Expect().Status(httptest.StatusOK).Body().Equal("SUCCESS")
    signatureQuery, _ = gokits.SHA1WithRSA.SignBase64ByRSAKeyString(
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

    signatureReset, _ := gokits.SHA1WithRSA.SignBase64ByRSAKeyString(
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

func TestAccessorError(t *testing.T) {
    a := assert.New(t)

    application := app.Application(func(config *app.Config) {
        config.AccessorVerifyCacheInMills = 1
        config.MerchantVerifyCacheInMills = 1
        config.AccessorMerchantVerifyCacheInMills = 1
    })
    e := httptest.New(t, application.App())

    responseAccessorIdIllegal := e.GET("/lannister/1000/query-info").
        Expect().Status(httptest.StatusOK).Body()
    resultAccessorIdIllegal := gokits.UnJson(
        responseAccessorIdIllegal.Raw(),
        &base.AccessorManage{}).(*base.AccessorManage)
    a.Equal("ACCESSOR_ID_ILLEGAL", resultAccessorIdIllegal.ErrorCode)
    a.Equal("AccessorId is Illegal", resultAccessorIdIllegal.ErrorDesc)

    responseNonsenseEmpty := e.GET("/lannister/1001/query-info").
        Expect().Status(httptest.StatusOK).Body()
    resultNonsenseEmpty := gokits.UnJson(
        responseNonsenseEmpty.Raw(),
        &base.AccessorManage{}).(*base.AccessorManage)
    a.Equal("NONSENSE_EMPTY", resultNonsenseEmpty.ErrorCode)
    a.Equal("Nonsense is Empty", resultNonsenseEmpty.ErrorDesc)

    responseSignatureEmpty := e.GET("/lannister/1001/query-info").
        WithQuery("nonsense", "query").
        Expect().Status(httptest.StatusOK).Body()
    resultSignatureEmpty := gokits.UnJson(
        responseSignatureEmpty.Raw(),
        &base.AccessorManage{}).(*base.AccessorManage)
    a.Equal("SIGNATURE_EMPTY", resultSignatureEmpty.ErrorCode)
    a.Equal("Signature is Empty", resultSignatureEmpty.ErrorDesc)

    responseSignatureMismatch := e.GET("/lannister/1001/query-info").
        WithQuery("nonsense", "query").
        WithQuery("signature", "signature").
        Expect().Status(httptest.StatusOK).Body()
    resultSignatureMismatch := gokits.UnJson(
        responseSignatureMismatch.Raw(),
        &base.AccessorManage{}).(*base.AccessorManage)
    a.Equal("SIGNATURE_MISMATCH", resultSignatureMismatch.ErrorCode)
    a.Equal("Signature mismatch", resultSignatureMismatch.ErrorDesc)

    signatureQuery, _ := gokits.SHA1WithRSA.SignBase64ByRSAKeyString(
        "nonsense=query", PrivateKeyString)
    responseQuery := e.GET("/lannister/1002/query-info").
        WithQuery("nonsense", "query").
        WithQuery("signature", signatureQuery).
        Expect().Status(httptest.StatusOK).Body()
    resultQuery := gokits.UnJson(responseQuery.Raw(),
        &base.AccessorManage{}).(*base.AccessorManage)
    a.Equal("", resultQuery.AccessorId)
    a.Equal("", resultQuery.AccessorName)
    a.Equal("", resultQuery.AccessorPubKey)
    a.Equal("", resultQuery.PayNotifyUrl)
    a.Equal("", resultQuery.RefundNotifyUrl)
    a.Equal("", resultQuery.PubKey)

    signatureUpdate, _ := gokits.SHA1WithRSA.SignBase64ByRSAKeyString(
        "accessorName=test&accessorPubKey="+PublicKeyString+"&nonsense=update&"+
            "payNotifyUrl=PayNotifyUrl&refundNotifyUrl=RefundNotifyUrl", PrivateKeyString)
    e.POST("/lannister/1002/update-info").
        WithQuery("nonsense", "update").
        WithQuery("signature", signatureUpdate).
        WithJSON(&base.AccessorManage{AccessorName: "test",
            AccessorPubKey:  PublicKeyString,
            PayNotifyUrl:    "PayNotifyUrl",
            RefundNotifyUrl: "RefundNotifyUrl",}).
        Expect().Status(httptest.StatusOK).Body().Equal("FAILED")

    signatureReset, _ := gokits.SHA1WithRSA.SignBase64ByRSAKeyString(
        "nonsense=reset", PrivateKeyString)
    responseReset := e.POST("/lannister/1002/reset-info").
        WithQuery("nonsense", "reset").
        WithQuery("signature", signatureReset).
        Expect().Status(httptest.StatusOK).Body()
    resultReset := gokits.UnJson(responseReset.Raw(),
        &base.AccessorManage{}).(*base.AccessorManage)
    a.Equal("", resultReset.PubKey)
}
