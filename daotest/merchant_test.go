package daotest

import (
    "github.com/CharLemAznable/go-lannister/app"
    "github.com/CharLemAznable/go-lannister/base"
    "github.com/CharLemAznable/gokits"
    "github.com/kataras/iris/v12/httptest"
    "github.com/stretchr/testify/assert"
    "testing"
)

func TestMerchant(t *testing.T) {
    a := assert.New(t)

    for _, testConfig := range TestConfigSet {
        application := app.Application(func(config *app.Config) {
            config.DriverName = testConfig["DriverName"]
            config.DataSourceName = testConfig["DataSourceName"]
        })
        e := httptest.New(t, application.App())

        signatureQueryAll, _ := gokits.SHA1WithRSA.SignBase64ByRSAKeyString(
            "nonsense=queryAll", PrivateKeyString)
        responseQueryAll := e.GET("/lannister/1001/query-merchants-info").
            WithQuery("nonsense", "queryAll").
            WithQuery("signature", signatureQueryAll).
            Expect().Status(httptest.StatusOK).Body()
        resultQueryAll := gokits.UnJson(responseQueryAll.Raw(),
            &map[string][]*base.MerchantManage{}).(*map[string][]*base.MerchantManage)
        a.Equal(1, len((*resultQueryAll)["merchants"]))
        resultQueryFirst := (*resultQueryAll)["merchants"][0]
        a.Equal("1001", resultQueryFirst.MerchantId)
        a.Equal("1001", resultQueryFirst.MerchantName)
        a.Equal("m1001", resultQueryFirst.MerchantCode)

        signatureCreate, _ := gokits.SHA1WithRSA.SignBase64ByRSAKeyString(
            "merchantCode=mm1001&merchantId=1001&merchantName=createById&nonsense=create", PrivateKeyString)
        responseCreate := e.POST("/lannister/1001/manage-merchant").
            WithQuery("nonsense", "create").
            WithQuery("signature", signatureCreate).
            WithJSON(&base.MerchantManage{
                MerchantId:   "1001",
                MerchantName: "createById",
                MerchantCode: "mm1001",}).
            Expect().Status(httptest.StatusOK).Body()
        resultCreate := gokits.UnJson(responseCreate.Raw(),
            &map[string]string{}).(*map[string]string)
        a.Equal("Create/Update Success", (*resultCreate)["message"])
        a.Equal("1001", (*resultCreate)["merchantId"])

        signatureQuery, _ := gokits.SHA1WithRSA.SignBase64ByRSAKeyString(
            "nonsense=query", PrivateKeyString)
        responseQuery := e.GET("/lannister/1001/1001/query-info").
            WithQuery("nonsense", "query").
            WithQuery("signature", signatureQuery).
            Expect().Status(httptest.StatusOK).Body()
        resultQuery := gokits.UnJson(responseQuery.Raw(),
            &base.MerchantManage{}).(*base.MerchantManage)
        a.Equal("1001", resultQuery.MerchantId)
        a.Equal("createById", resultQuery.MerchantName)
        a.Equal("mm1001", resultQuery.MerchantCode)

        signatureCreate, _ = gokits.SHA1WithRSA.SignBase64ByRSAKeyString(
            "merchantCode=mm1001&merchantName=createByCode&nonsense=create", PrivateKeyString)
        responseCreate = e.POST("/lannister/1001/manage-merchant").
            WithQuery("nonsense", "create").
            WithQuery("signature", signatureCreate).
            WithJSON(&base.MerchantManage{
                MerchantCode: "mm1001",
                MerchantName: "createByCode",}).
            Expect().Status(httptest.StatusOK).Body()
        resultCreate = gokits.UnJson(responseCreate.Raw(),
            &map[string]string{}).(*map[string]string)
        a.Equal("Create/Update Success", (*resultCreate)["message"])
        a.Equal("1001", (*resultCreate)["merchantId"])

        signatureQuery, _ = gokits.SHA1WithRSA.SignBase64ByRSAKeyString(
            "nonsense=query", PrivateKeyString)
        responseQuery = e.GET("/lannister/1001/1001/query-info").
            WithQuery("nonsense", "query").
            WithQuery("signature", signatureQuery).
            Expect().Status(httptest.StatusOK).Body()
        resultQuery = gokits.UnJson(responseQuery.Raw(),
            &base.MerchantManage{}).(*base.MerchantManage)
        a.Equal("1001", resultQuery.MerchantId)
        a.Equal("createByCode", resultQuery.MerchantName)
        a.Equal("mm1001", resultQuery.MerchantCode)

        signatureCreate, _ = gokits.SHA1WithRSA.SignBase64ByRSAKeyString(
            "authorizeAll=true&merchantCode=m2001&merchantId=2001&merchantName=create&nonsense=create", PrivateKeyString)
        responseCreate = e.POST("/lannister/1001/manage-merchant").
            WithQuery("nonsense", "create").
            WithQuery("signature", signatureCreate).
            WithJSON(&base.MerchantManage{
                MerchantId:   "2001",
                MerchantName: "create",
                MerchantCode: "m2001",
                AuthorizeAll: "true"}).
            Expect().Status(httptest.StatusOK).Body()
        resultCreate = gokits.UnJson(responseCreate.Raw(),
            &map[string]string{}).(*map[string]string)
        a.Equal("Create/Update Success", (*resultCreate)["message"])
        newMerchantId := (*resultCreate)["merchantId"]

        signatureQuery, _ = gokits.SHA1WithRSA.SignBase64ByRSAKeyString(
            "nonsense=query", PrivateKeyString)
        responseQuery = e.GET("/lannister/1001/" + newMerchantId + "/query-info").
            WithQuery("nonsense", "query").
            WithQuery("signature", signatureQuery).
            Expect().Status(httptest.StatusOK).Body()
        resultQuery = gokits.UnJson(responseQuery.Raw(),
            &base.MerchantManage{}).(*base.MerchantManage)
        a.Equal(newMerchantId, resultQuery.MerchantId)
        a.Equal("create", resultQuery.MerchantName)
        a.Equal("m2001", resultQuery.MerchantCode)

        signatureQuery, _ = gokits.SHA1WithRSA.SignBase64ByRSAKeyString(
            "nonsense=query", PrivateKeyString)
        responseQuery = e.GET("/lannister/1002/" + newMerchantId + "/query-info").
            WithQuery("nonsense", "query").
            WithQuery("signature", signatureQuery).
            Expect().Status(httptest.StatusOK).Body()
        resultQuery = gokits.UnJson(responseQuery.Raw(),
            &base.MerchantManage{}).(*base.MerchantManage)
        a.Equal(newMerchantId, resultQuery.MerchantId)
        a.Equal("create", resultQuery.MerchantName)
        a.Equal("m2001", resultQuery.MerchantCode)
    }
}
