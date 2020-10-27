package apptest_test

import (
    "github.com/CharLemAznable/go-lannister/app"
    "github.com/CharLemAznable/go-lannister/apptest"
    . "github.com/CharLemAznable/go-lannister/base"
    . "github.com/CharLemAznable/go-lannister/elf"
    "github.com/CharLemAznable/gokits"
    "github.com/kataras/iris/v12/httptest"
    "github.com/stretchr/testify/assert"
    "testing"
)

func TestMerchant(t *testing.T) {
    a := assert.New(t)

    application := app.Application()
    e := httptest.New(t, application.App())

    signatureQueryAll, _ := SHA1WithRSA.SignBase64ByKeyString(
        "nonsense=queryAll", apptest.PrivateKeyString)
    responseQueryAll := e.GET("/lannister/1001/query-merchants-info").
        WithQuery("nonsense", "queryAll").
        WithQuery("signature", signatureQueryAll).
        Expect().Status(httptest.StatusOK).Body()
    resultQueryAll := gokits.UnJson(responseQueryAll.Raw(),
        &map[string][]*MerchantManage{}).(*map[string][]*MerchantManage)
    a.Equal(1, len((*resultQueryAll)["merchants"]))
    resultQueryFirst := (*resultQueryAll)["merchants"][0]
    a.Equal("1001", resultQueryFirst.MerchantId)
    a.Equal("", resultQueryFirst.MerchantName)
    a.Equal("m1001", resultQueryFirst.MerchantCode)

    signatureCreate, _ := SHA1WithRSA.SignBase64ByKeyString(
        "merchantCode=mm1001&merchantId=1001&merchantName=createById&nonsense=create", apptest.PrivateKeyString)
    responseCreate := e.POST("/lannister/1001/manage-merchant").
        WithQuery("nonsense", "create").
        WithQuery("signature", signatureCreate).
        WithJSON(&MerchantManage{
            MerchantId:   "1001",
            MerchantName: "createById",
            MerchantCode: "mm1001",}).
        Expect().Status(httptest.StatusOK).Body()
    resultCreate := gokits.UnJson(responseCreate.Raw(),
        &map[string]string{}).(*map[string]string)
    a.Equal("Create/Update Success", (*resultCreate)["message"])
    a.Equal("1001", (*resultCreate)["merchantId"])

    signatureQuery, _ := SHA1WithRSA.SignBase64ByKeyString(
        "nonsense=query", apptest.PrivateKeyString)
    responseQuery := e.GET("/lannister/1001/1001/query-info").
        WithQuery("nonsense", "query").
        WithQuery("signature", signatureQuery).
        Expect().Status(httptest.StatusOK).Body()
    resultQuery := gokits.UnJson(responseQuery.Raw(),
        &MerchantManage{}).(*MerchantManage)
    a.Equal("1001", resultQuery.MerchantId)
    a.Equal("createById", resultQuery.MerchantName)
    a.Equal("mm1001", resultQuery.MerchantCode)

    signatureCreate, _ = SHA1WithRSA.SignBase64ByKeyString(
        "merchantCode=mm1001&merchantName=createByCode&nonsense=create", apptest.PrivateKeyString)
    responseCreate = e.POST("/lannister/1001/manage-merchant").
        WithQuery("nonsense", "create").
        WithQuery("signature", signatureCreate).
        WithJSON(&MerchantManage{
            MerchantCode: "mm1001",
            MerchantName: "createByCode",}).
        Expect().Status(httptest.StatusOK).Body()
    resultCreate = gokits.UnJson(responseCreate.Raw(),
        &map[string]string{}).(*map[string]string)
    a.Equal("Create/Update Success", (*resultCreate)["message"])
    a.Equal("1001", (*resultCreate)["merchantId"])

    signatureQuery, _ = SHA1WithRSA.SignBase64ByKeyString(
        "nonsense=query", apptest.PrivateKeyString)
    responseQuery = e.GET("/lannister/1001/1001/query-info").
        WithQuery("nonsense", "query").
        WithQuery("signature", signatureQuery).
        Expect().Status(httptest.StatusOK).Body()
    resultQuery = gokits.UnJson(responseQuery.Raw(),
        &MerchantManage{}).(*MerchantManage)
    a.Equal("1001", resultQuery.MerchantId)
    a.Equal("createByCode", resultQuery.MerchantName)
    a.Equal("mm1001", resultQuery.MerchantCode)
}
