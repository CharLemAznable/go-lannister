package app

import (
    "github.com/CharLemAznable/go-lannister/base"
    "github.com/CharLemAznable/gokits"
    "github.com/CharLemAznable/sqlx"
    "github.com/kataras/iris/v12"
    "github.com/kataras/iris/v12/mvc"
    "time"
)

type MerchantManageController struct {
    dao base.MerchantManageDao
}

func (c *MerchantManageController) BeforeActivation(b mvc.BeforeActivation) {
    _ = b.Dependencies().Inject(&c.dao) // does not support func register yet
    b.Handle(PostMapping("/{accessorId}/manage-merchant", "ManageMerchant"))
    b.Handle(GetMapping("/{accessorId}/query-merchants-info", "QueryMerchantsInfo"))
    b.Handle(GetMapping("/{accessorId}/{merchantId}/query-info", "QueryMerchantInfo"))
}

func (c *MerchantManageController) ManageMerchant(ctx iris.Context) {
    accessorId := ctx.Params().Get("accessorId")
    req := &base.MerchantManage{}
    _ = ctx.ReadJSON(req)

    merchant := &base.MerchantManage{}
    create := false
    if "" != req.MerchantId {
        // 根据商户标识查询
        merchant, _ = c.dao.QueryMerchantById(req.MerchantId)
    }
    if "" == merchant.MerchantId && "" != req.MerchantCode {
        // 根据商户编码查询
        merchant, _ = c.dao.QueryMerchantByCode(req.MerchantCode)
    }
    if "" == merchant.MerchantId {
        // 未传商户标识和编码/未查询到商户, 则分配商户标识, 创建新商户
        merchant.MerchantId = gokits.NextId()
        create = true
    }

    // 根据请求信息填充商户信息
    if "" != req.MerchantName {
        merchant.MerchantName = req.MerchantName
    }
    if "" != req.MerchantCode {
        merchant.MerchantCode = req.MerchantCode
    }
    authorizeId := gokits.Condition(gokits.ToBool(
        req.AuthorizeAll), "0", accessorId).(string)

    affected := int64(0)
    var err error
    if create {
        affected, err = c.dao.CreateMerchant(accessorId,
            merchant.MerchantId, merchant.MerchantName, merchant.MerchantCode)
    } else {
        affected, err = c.dao.UpdateMerchant(accessorId,
            merchant.MerchantId, merchant.MerchantName, merchant.MerchantCode)
    }
    if nil != err {
        ctx.Application().Logger().Errorf("ManageMerchant Create/Update Merchant: %s", err.Error())
    }
    _, e := c.dao.UpdateAccessorMerchant(authorizeId, merchant.MerchantId)
    if nil != e {
        ctx.Application().Logger().Errorf("ManageMerchant Update Accessor Merchant: %s", e.Error())
    }
    ctx.Application().Logger().Debugf("Create/Update Merchant by accessor %s merchant: %#v", accessorId, *merchant)

    if 0 >= affected {
        _, _ = ctx.JSON(base.BaseResp{ErrorCode: "MANAGE_MERCHANT_FAILED", ErrorDesc: err.Error()})
    } else {
        _, _ = ctx.JSON(merchant)
    }
}

func (c *MerchantManageController) QueryMerchantsInfo(ctx iris.Context) {
    accessorId := ctx.Params().Get("accessorId")
    merchants, err := c.dao.QueryMerchants(accessorId)
    if nil != err {
        ctx.Application().Logger().Errorf("QueryMerchantsInfo: %s", err.Error())
    }
    ctx.Application().Logger().Debugf("QueryMerchantsInfo by accessor %s merchants: %#v", accessorId, merchants)
    _, _ = ctx.JSON(iris.Map{"merchants": merchants})
}

func (c *MerchantManageController) QueryMerchantInfo(ctx iris.Context) {
    accessorId := ctx.Params().Get("accessorId")
    merchantId := ctx.Params().Get("merchantId")
    merchant, err := c.dao.QueryMerchant(accessorId, merchantId)
    if nil != err {
        ctx.Application().Logger().Errorf("QueryMerchantInfo: %s", err.Error())
    }
    ctx.Application().Logger().Debugf("QueryMerchantInfo by accessor %s merchant: %#v", accessorId, *merchant)
    _, _ = ctx.JSON(merchant)
}

/****************************************************************************************************/

type MerchantVerifyCache struct {
    dao                      base.MerchantVerifyDao
    tableMerchant            *gokits.CacheTable
    lifeSpanMerchant         time.Duration
    tableAccessorMerchant    *gokits.CacheTable
    lifeSpanAccessorMerchant time.Duration
}

func NewMerchantVerifyCache(dao base.MerchantVerifyDao, config *base.Config) *MerchantVerifyCache {
    tableMerchant := gokits.NewCacheExpireAfterWrite("MerchantVerifyCache.tableMerchant")
    tableAccessorMerchant := gokits.NewCacheExpireAfterWrite("MerchantVerifyCache.tableAccessorMerchant")
    cache := &MerchantVerifyCache{dao: dao,
        tableMerchant:            tableMerchant,
        lifeSpanMerchant:         time.Duration(config.MerchantVerifyCacheInMills) * time.Millisecond,
        tableAccessorMerchant:    tableAccessorMerchant,
        lifeSpanAccessorMerchant: time.Duration(config.AccessorMerchantVerifyCacheInMills) * time.Millisecond}
    tableMerchant.SetDataLoader(cache.merchantVerifyLoader)
    tableAccessorMerchant.SetDataLoader(cache.accessorMerchantVerifyLoader)
    return cache
}

func (c *MerchantVerifyCache) merchantVerifyLoader(merchantId interface{}, _ ...interface{}) (*gokits.CacheItem, error) {
    verify, err := c.dao.QueryMerchant(merchantId.(string))
    if nil != err {
        return nil, err
    }
    return gokits.NewCacheItem(merchantId, c.lifeSpanMerchant, verify), nil
}

func (c *MerchantVerifyCache) queryMerchantVerify(merchantId string) (*base.MerchantVerify, error) {
    value, err := c.tableMerchant.Value(merchantId)
    if nil != err {
        return nil, err
    }
    return value.Data().(*base.MerchantVerify), nil
}

type merchantVerifyCacheKey struct {
    accessorId string
    merchantId string
}

func (c *MerchantVerifyCache) accessorMerchantVerifyLoader(key interface{}, _ ...interface{}) (*gokits.CacheItem, error) {
    cacheKey := key.(*merchantVerifyCacheKey)
    verifies, err := c.dao.QueryAccessorMerchants(cacheKey.accessorId, cacheKey.merchantId)
    if nil != err {
        return nil, err
    }
    return gokits.NewCacheItem(cacheKey, c.lifeSpanAccessorMerchant, verifies), nil
}

func (c *MerchantVerifyCache) queryAccessorMerchantVerify(accessorId, merchantId string) ([]*base.MerchantVerify, error) {
    value, err := c.tableAccessorMerchant.Value(&merchantVerifyCacheKey{
        accessorId: accessorId, merchantId: merchantId})
    if nil != err {
        return nil, err
    }
    return value.Data().([]*base.MerchantVerify), nil
}

/****************************************************************************************************/

var (
    merchantIdIllegal = base.BaseResp{
        ErrorCode: "MERCHANT_ID_ILLEGAL",
        ErrorDesc: "MerchantId is Illegal",
    }
    merchantAccessUnauthorized = base.BaseResp{
        ErrorCode: "MERCHANT_ACCESS_UNAUTHORIZED",
        ErrorDesc: "Merchant access Unauthorized",
    }
)

func MerchantVerifyInterceptor(ctx iris.Context, cache *MerchantVerifyCache) {
    accessorId := ctx.Params().Get("accessorId")
    if "" == accessorId {
        ctx.Next()
        return
    }
    merchantId := ctx.Params().Get("merchantId")
    if "" == merchantId {
        ctx.Next()
        return
    }

    ctx.Application().Logger().Debugf("MerchantVerify accessorId: %s, merchantId: %s", accessorId, merchantId)

    _, err := cache.queryMerchantVerify(merchantId)
    if nil != err {
        ctx.StopWithJSON(iris.StatusOK, merchantIdIllegal)
        return
    }

    verifies, err := cache.queryAccessorMerchantVerify(accessorId, merchantId)
    if nil != err || 0 == len(verifies) {
        ctx.StopWithJSON(iris.StatusOK, merchantAccessUnauthorized)
        return
    }

    ctx.Application().Logger().Debug("MerchantVerify done")
    ctx.Next()
}

/****************************************************************************************************/

func init() {
    RegisterDependency(func(config *base.Config, db *sqlx.DB) (
        base.MerchantManageDao, base.MerchantVerifyDao, *MerchantVerifyCache) {
        merchantManageDao := base.GetMerchantManageDao(db)
        merchantVerifyDao := base.GetMerchantVerifyDao(db)
        merchantVerifyCache := NewMerchantVerifyCache(merchantVerifyDao, config)
        return merchantManageDao, merchantVerifyDao, merchantVerifyCache
    })
    RegisterMiddleware(MerchantVerifyInterceptor)
    RegisterController(&MerchantManageController{})
}
