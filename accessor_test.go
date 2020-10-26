package lannister

import (
    "github.com/CharLemAznable/go-lannister/dao/dummy"
    . "github.com/CharLemAznable/go-lannister/elf"
    "github.com/CharLemAznable/go-lannister/types"
    "github.com/CharLemAznable/gokits"
    "github.com/kataras/iris/v12/httptest"
    "github.com/stretchr/testify/assert"
    "testing"
)

func TestAccessor(t *testing.T) {
    a := assert.New(t)

    application := Application()
    e := httptest.New(t, application.app)

    responseAccessorIdIllegal := e.GET("/lannister/1000/query-info").
        Expect().Status(httptest.StatusOK).Body()
    resultAccessorIdIllegal := gokits.UnJson(
        responseAccessorIdIllegal.Raw(),
        &types.AccessorManage{}).(*types.AccessorManage)
    a.Equal("ACCESSOR_ID_ILLEGAL", resultAccessorIdIllegal.ErrorCode)
    a.Equal("AccessorId is Illegal", resultAccessorIdIllegal.ErrorDesc)

    responseNonsenseEmpty := e.GET("/lannister/1001/query-info").
        Expect().Status(httptest.StatusOK).Body()
    resultNonsenseEmpty := gokits.UnJson(
        responseNonsenseEmpty.Raw(),
        &types.AccessorManage{}).(*types.AccessorManage)
    a.Equal("NONSENSE_EMPTY", resultNonsenseEmpty.ErrorCode)
    a.Equal("Nonsense is Empty", resultNonsenseEmpty.ErrorDesc)

    responseSignatureEmpty := e.GET("/lannister/1001/query-info").
        WithQuery("nonsense", "query").
        Expect().Status(httptest.StatusOK).Body()
    resultSignatureEmpty := gokits.UnJson(
        responseSignatureEmpty.Raw(),
        &types.AccessorManage{}).(*types.AccessorManage)
    a.Equal("SIGNATURE_EMPTY", resultSignatureEmpty.ErrorCode)
    a.Equal("Signature is Empty", resultSignatureEmpty.ErrorDesc)

    responseSignatureMismatch := e.GET("/lannister/1001/query-info").
        WithQuery("nonsense", "query").
        WithQuery("signature", "signature").
        Expect().Status(httptest.StatusOK).Body()
    resultSignatureMismatch := gokits.UnJson(
        responseSignatureMismatch.Raw(),
        &types.AccessorManage{}).(*types.AccessorManage)
    a.Equal("SIGNATURE_MISMATCH", resultSignatureMismatch.ErrorCode)
    a.Equal("Signature mismatch", resultSignatureMismatch.ErrorDesc)

    signatureQuery, _ := SHA1WithRSA.SignBase64ByKeyString(
        "nonsense=query", dummy.PrivateKeyString)
    responseQuery := e.GET("/lannister/1001/query-info").
        WithQuery("nonsense", "query").
        WithQuery("signature", signatureQuery).
        Expect().Status(httptest.StatusOK).Body()
    resultQuery := gokits.UnJson(responseQuery.Raw(),
        &types.AccessorManage{}).(*types.AccessorManage)
    a.Equal("1001", resultQuery.AccessorId)
    a.Equal("1001", resultQuery.AccessorName)
    a.Equal(dummy.PublicKeyString, resultQuery.AccessorPubKey)
    a.Equal("", resultQuery.PayNotifyUrl)
    a.Equal("", resultQuery.RefundNotifyUrl)
    a.Equal("", resultQuery.PubKey)

    signatureUpdate, _ := SHA1WithRSA.SignBase64ByKeyString(
        "accessorName=test&accessorPubKey="+dummy.PublicKeyString+"&nonsense=update&"+
            "payNotifyUrl=PayNotifyUrl&refundNotifyUrl=RefundNotifyUrl", dummy.PrivateKeyString)
    e.POST("/lannister/1001/update-info").
        WithQuery("nonsense", "update").
        WithQuery("signature", signatureUpdate).
        WithJSON(&types.AccessorManage{AccessorName: "test",
            AccessorPubKey:  dummy.PublicKeyString,
            PayNotifyUrl:    "PayNotifyUrl",
            RefundNotifyUrl: "RefundNotifyUrl",}).
        Expect().Status(httptest.StatusOK).Body().Equal("SUCCESS")
    signatureQuery, _ = SHA1WithRSA.SignBase64ByKeyString(
        "nonsense=query", dummy.PrivateKeyString)
    responseQuery = e.GET("/lannister/1001/query-info").
        WithQuery("nonsense", "query").
        WithQuery("signature", signatureQuery).
        Expect().Status(httptest.StatusOK).Body()
    resultQuery = gokits.UnJson(responseQuery.Raw(),
        &types.AccessorManage{}).(*types.AccessorManage)
    a.Equal("1001", resultQuery.AccessorId)
    a.Equal("test", resultQuery.AccessorName)
    a.Equal(dummy.PublicKeyString, resultQuery.AccessorPubKey)
    a.Equal("", resultQuery.PayNotifyUrl)
    a.Equal("", resultQuery.RefundNotifyUrl)
    a.Equal("", resultQuery.PubKey)

    signatureReset, _ := SHA1WithRSA.SignBase64ByKeyString(
        "nonsense=reset", dummy.PrivateKeyString)
    responseReset := e.POST("/lannister/1001/reset-info").
        WithQuery("nonsense", "reset").
        WithQuery("signature", signatureReset).
        Expect().Status(httptest.StatusOK).Body()
    resultReset := gokits.UnJson(responseReset.Raw(),
        &types.AccessorManage{}).(*types.AccessorManage)
    responseQuery = e.GET("/lannister/1001/query-info").
        WithQuery("nonsense", "query").
        WithQuery("signature", signatureQuery).
        Expect().Status(httptest.StatusOK).Body()
    resultQuery = gokits.UnJson(responseQuery.Raw(),
        &types.AccessorManage{}).(*types.AccessorManage)
    a.Equal("1001", resultQuery.AccessorId)
    a.Equal("test", resultQuery.AccessorName)
    a.Equal(dummy.PublicKeyString, resultQuery.AccessorPubKey)
    a.Equal("", resultQuery.PayNotifyUrl)
    a.Equal("", resultQuery.RefundNotifyUrl)
    a.Equal(resultReset.PubKey, resultQuery.PubKey)
}
