package app

import (
    . "github.com/CharLemAznable/go-lannister/base"
    . "github.com/CharLemAznable/go-lannister/elf"
    "github.com/CharLemAznable/gokits"
    "github.com/kataras/iris/v12"
    "github.com/kataras/iris/v12/mvc"
    "time"
)

type AccessorManageController struct {
    dao AccessorManageDao
}

func (c *AccessorManageController) BeforeActivation(b mvc.BeforeActivation) {
    _ = b.Dependencies().Inject(&c.dao) // does not support func register yet
    b.Handle(GetMapping("/{accessorId}/query-info", "QueryAccessorInfo"))
    b.Handle(PostMapping("/{accessorId}/update-info", "UpdateAccessorInfo"))
    b.Handle(PostMapping("/{accessorId}/reset-info", "ResetKeyPair"))
}

func (c *AccessorManageController) QueryAccessorInfo(ctx iris.Context) {
    accessorId := ctx.Params().Get("accessorId")
    manage, err := c.dao.QueryAccessor(accessorId)
    if err != nil {
        ctx.Application().Logger().Errorf("QueryAccessorInfo: %s", err.Error())
    }
    ctx.Application().Logger().Debugf("Query accessor %s info: %#v", accessorId, *manage)
    _, _ = ctx.JSON(manage)
}

func (c *AccessorManageController) UpdateAccessorInfo(ctx iris.Context) {
    accessorId := ctx.Params().Get("accessorId")
    req := &AccessorManage{}
    _ = ctx.ReadJSON(req)
    update, err := c.dao.UpdateAccessor(accessorId, req)
    if err != nil {
        ctx.Application().Logger().Errorf("UpdateAccessorInfo: %s", err.Error())
    }
    ctx.Application().Logger().Debugf("Update accessor %s info: %#v", accessorId, *req)
    _, _ = ctx.Text(gokits.Condition(1 == update, "SUCCESS", "FAILED").(string))
}

func (c *AccessorManageController) ResetKeyPair(ctx iris.Context) {
    accessorId := ctx.Params().Get("accessorId")
    nonsense := ctx.URLParam("nonsense")
    keyPair, _ := GenerateKeyPairDefault()
    privateKeyString, _ := keyPair.PrivateKeyEncoded()
    publicKeyString, _ := keyPair.PublicKeyEncoded()
    err := c.dao.UpdateKeyPair(accessorId, nonsense, publicKeyString, privateKeyString)
    if err != nil {
        ctx.Application().Logger().Errorf("ResetKeyPair: %s", err.Error())
    }
    ctx.Application().Logger().Debugf("Reset Accessor %s PublicKey: %s", accessorId, publicKeyString)
    manage, _ := c.dao.QueryAccessor(accessorId)
    _, _ = ctx.JSON(iris.Map{"PubKey": manage.PubKey})
}

/****************************************************************************************************/

type AccessorVerifyDaoCache struct {
    dao   AccessorVerifyDao
    table *gokits.CacheTable
}

func NewAccessorVerifyDaoCache(dao AccessorVerifyDao) *AccessorVerifyDaoCache {
    table := gokits.CacheExpireAfterWrite("AccessorVerifyDaoCache.table")
    cache := &AccessorVerifyDaoCache{dao: dao, table: table}
    table.SetDataLoader(cache.accessorVerifyLoader)
    return cache
}

func (c *AccessorVerifyDaoCache) accessorVerifyLoader(accessorId interface{}, _ ...interface{}) (*gokits.CacheItem, error) {
    verify, err := c.dao.QueryAccessorById(accessorId.(string))
    if nil != err {
        return nil, err
    }
    // 缓存1min
    return gokits.NewCacheItem(accessorId, time.Minute, verify), nil
}

func (c *AccessorVerifyDaoCache) queryAccessorById(accessorId string) (*AccessorVerify, error) {
    value, err := c.table.Value(accessorId)
    if nil != err {
        return nil, err
    }
    return value.Data().(*AccessorVerify), nil
}

var (
    accessorIdIllegal = BaseResp{
        ErrorCode: "ACCESSOR_ID_ILLEGAL",
        ErrorDesc: "AccessorId is Illegal",
    }
    nonsenseEmpty = BaseResp{
        ErrorCode: "NONSENSE_EMPTY",
        ErrorDesc: "Nonsense is Empty",
    }
    signatureEmpty = BaseResp{
        ErrorCode: "SIGNATURE_EMPTY",
        ErrorDesc: "Signature is Empty",
    }
    signatureMismatch = BaseResp{
        ErrorCode: "SIGNATURE_MISMATCH",
        ErrorDesc: "Signature mismatch",
    }
)

func AccessorVerifyInterceptor(ctx iris.Context, cache *AccessorVerifyDaoCache) {
    accessorId := ctx.Params().Get("accessorId")
    if "" == accessorId {
        ctx.Next()
        return
    }

    ctx.Application().Logger().Debugf("AccessorVerify accessorId: %s", accessorId)

    verify, err := cache.queryAccessorById(accessorId)
    if nil != err {
        ctx.StopWithJSON(iris.StatusOK, accessorIdIllegal)
        return
    }

    nonsense := ctx.URLParam(NONSENSE)
    if "" == nonsense {
        ctx.StopWithJSON(iris.StatusOK, nonsenseEmpty)
        return
    }

    signature := ctx.URLParam(SIGNATURE)
    if "" == signature {
        ctx.StopWithJSON(iris.StatusOK, signatureEmpty)
        return
    }

    paramMap := map[string]string{}
    _ = ctx.ReadBody(&paramMap)
    urlParamMap := ctx.URLParams()
    for key, value := range urlParamMap {
        paramMap[key] = value
    }

    if err := verify.Verify(paramMap); nil != err {
        ctx.Application().Logger().Debugf("AccessorVerify error: %s", err.Error())
        ctx.StopWithJSON(iris.StatusOK, signatureMismatch)
        return
    }

    ctx.Application().Logger().Debug("AccessorVerify done")
    ctx.Next()
}

/****************************************************************************************************/

func init() {
    RegisterDependency("lannister.AccessorManageDao", GetAccessorManageDao)
    RegisterController("lannister.AccessorManageController", &AccessorManageController{})

    RegisterDependency("lannister.AccessorVerifyDao", GetAccessorVerifyDao)
    RegisterDependency("lannister.AccessorVerifyDaoCache", NewAccessorVerifyDaoCache)
    RegisterMiddleware("lannister.AccessorVerifyInterceptor", AccessorVerifyInterceptor)
}