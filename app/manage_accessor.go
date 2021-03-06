package app

import (
    "github.com/CharLemAznable/go-lannister/base"
    "github.com/CharLemAznable/gokits"
    "github.com/CharLemAznable/sqlx"
    "github.com/kataras/iris/v12"
    "github.com/kataras/iris/v12/mvc"
    "time"
)

type AccessorManageController struct {
    dao base.AccessorManageDao
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
    req := &base.AccessorManage{}
    _ = ctx.ReadJSON(req)
    _, err := c.dao.UpdateAccessor(accessorId, req)
    if err != nil {
        ctx.Application().Logger().Errorf("UpdateAccessorInfo: %s", err.Error())
    }
    ctx.Application().Logger().Debugf("Update accessor %s info: %#v", accessorId, *req)

    manage, _ := c.dao.QueryAccessor(accessorId)
    _, _ = ctx.JSON(manage)
}

func (c *AccessorManageController) ResetKeyPair(ctx iris.Context) {
    accessorId := ctx.Params().Get("accessorId")
    nonsense := ctx.URLParam("nonsense")
    keyPair, _ := gokits.GenerateRSAKeyPairDefault()
    privateKeyString, _ := keyPair.RSAPrivateKeyEncoded()
    publicKeyString, _ := keyPair.RSAPublicKeyEncoded()
    _, err := c.dao.UpdateKeyPair(accessorId, nonsense, publicKeyString, privateKeyString)
    if err != nil {
        ctx.Application().Logger().Errorf("ResetKeyPair: %s", err.Error())
    }
    ctx.Application().Logger().Debugf("Reset Accessor %s PublicKey: %s", accessorId, publicKeyString)

    manage, _ := c.dao.QueryAccessor(accessorId)
    _, _ = ctx.JSON(manage)
}

/****************************************************************************************************/

type AccessorVerifyCache struct {
    dao      base.AccessorVerifyDao
    table    *gokits.CacheTable
    lifeSpan time.Duration
}

func NewAccessorVerifyCache(dao base.AccessorVerifyDao, config *base.Config) *AccessorVerifyCache {
    table := gokits.NewCacheExpireAfterWrite("AccessorVerifyCache.table")
    cache := &AccessorVerifyCache{dao: dao, table: table,
        lifeSpan: time.Duration(config.AccessorVerifyCacheInMills) * time.Millisecond}
    table.SetDataLoader(cache.accessorVerifyLoader)
    return cache
}

func (c *AccessorVerifyCache) accessorVerifyLoader(accessorId interface{}, _ ...interface{}) (*gokits.CacheItem, error) {
    verify, err := c.dao.QueryAccessor(accessorId.(string))
    if nil != err {
        return nil, err
    }
    return gokits.NewCacheItem(accessorId, c.lifeSpan, verify), nil
}

func (c *AccessorVerifyCache) queryAccessorVerify(accessorId string) (*base.AccessorVerify, error) {
    value, err := c.table.Value(accessorId)
    if nil != err {
        return nil, err
    }
    return value.Data().(*base.AccessorVerify), nil
}

/****************************************************************************************************/

var (
    accessorIdIllegal = base.BaseResp{
        ErrorCode: "ACCESSOR_ID_ILLEGAL",
        ErrorDesc: "AccessorId is Illegal",
    }
    nonsenseEmpty = base.BaseResp{
        ErrorCode: "NONSENSE_EMPTY",
        ErrorDesc: "Nonsense is Empty",
    }
    signatureEmpty = base.BaseResp{
        ErrorCode: "SIGNATURE_EMPTY",
        ErrorDesc: "Signature is Empty",
    }
    signatureMismatch = base.BaseResp{
        ErrorCode: "SIGNATURE_MISMATCH",
        ErrorDesc: "Signature mismatch",
    }
)

func AccessorVerifyInterceptor(ctx iris.Context, cache *AccessorVerifyCache) {
    accessorId := ctx.Params().Get("accessorId")
    if "" == accessorId {
        ctx.Next()
        return
    }

    ctx.Application().Logger().Debugf("AccessorVerify accessorId: %s", accessorId)

    verify, err := cache.queryAccessorVerify(accessorId)
    if nil != err {
        ctx.StopWithJSON(iris.StatusOK, accessorIdIllegal)
        return
    }

    nonsense := ctx.URLParam(base.NONSENSE)
    if "" == nonsense {
        ctx.StopWithJSON(iris.StatusOK, nonsenseEmpty)
        return
    }

    signature := ctx.URLParam(base.SIGNATURE)
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
    RegisterDependency(func(config *base.Config, db *sqlx.DB) (
        base.AccessorManageDao, base.AccessorVerifyDao, *AccessorVerifyCache) {
        accessorManageDao := base.GetAccessorManageDao(db)
        accessorVerifyDao := base.GetAccessorVerifyDao(db)
        accessorVerifyCache := NewAccessorVerifyCache(accessorVerifyDao, config)
        return accessorManageDao, accessorVerifyDao, accessorVerifyCache
    })
    RegisterMiddleware(AccessorVerifyInterceptor)
    RegisterController(&AccessorManageController{})
}
