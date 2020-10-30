package app

import (
    . "github.com/CharLemAznable/go-lannister/base"
    . "github.com/CharLemAznable/go-lannister/elf"
    "github.com/CharLemAznable/gokits"
    "github.com/kataras/iris/v12"
    "github.com/kataras/iris/v12/mvc"
    "time"
)

const (
    messageKey    = "message"
    merchantIdKey = "merchantId"
)

type MerchantManageController struct {
    dao MerchantManageDao
}

func (c *MerchantManageController) BeforeActivation(b mvc.BeforeActivation) {
    _ = b.Dependencies().Inject(&c.dao) // does not support func register yet
    b.Handle(PostMapping("/{accessorId}/manage-merchant", "ManageMerchant"))
    b.Handle(GetMapping("/{accessorId}/query-merchants-info", "QueryMerchantsInfo"))
    b.Handle(GetMapping("/{accessorId}/{merchantId}/query-info", "QueryMerchantInfo"))
}

func (c *MerchantManageController) ManageMerchant(ctx iris.Context) {
    accessorId := ctx.Params().Get("accessorId")
    req := &MerchantManage{}
    _ = ctx.ReadJSON(req)

    merchant := &MerchantManage{}
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
    if "" != req.MerchantCode {
        merchant.MerchantCode = req.MerchantCode
    }
    if "" != req.MerchantName {
        merchant.MerchantName = req.MerchantName
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
    _, err = c.dao.UpdateAccessorMerchant(authorizeId, merchant.MerchantId)
    if nil != err {
        ctx.Application().Logger().Errorf("ManageMerchant Update Accessor Merchant: %s", err.Error())
    }
    ctx.Application().Logger().Debugf("Create/Update Merchant by accessor %s merchant: %#v", accessorId, *merchant)

    if 0 >= affected {
        _, _ = ctx.JSON(iris.Map{messageKey: "Create/Update Failed"})
    } else {
        _, _ = ctx.JSON(iris.Map{messageKey: "Create/Update Success", merchantIdKey: merchant.MerchantId})
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

type MerchantVerifyDaoCache struct {
    dao                   MerchantVerifyDao
    tableMerchant         *gokits.CacheTable
    tableAccessorMerchant *gokits.CacheTable
}

func NewMerchantVerifyDaoCache(dao MerchantVerifyDao) *MerchantVerifyDaoCache {
    tableMerchant := gokits.CacheExpireAfterWrite("MerchantVerifyDaoCache.tableMerchant")
    tableAccessorMerchant := gokits.CacheExpireAfterWrite("MerchantVerifyDaoCache.tableAccessorMerchant")
    cache := &MerchantVerifyDaoCache{dao: dao,
        tableMerchant:         tableMerchant,
        tableAccessorMerchant: tableAccessorMerchant}
    tableMerchant.SetDataLoader(cache.merchantVerifyLoader)
    tableAccessorMerchant.SetDataLoader(cache.accessorMerchantVerifyLoader)
    return cache
}

func (c *MerchantVerifyDaoCache) merchantVerifyLoader(merchantId interface{}, _ ...interface{}) (*gokits.CacheItem, error) {
    verify, err := c.dao.QueryMerchant(merchantId.(string))
    if nil != err {
        return nil, err
    }
    // 缓存1min
    return gokits.NewCacheItem(merchantId, time.Minute, verify), nil
}

func (c *MerchantVerifyDaoCache) queryMerchantById(merchantId string) (*MerchantVerify, error) {
    value, err := c.tableMerchant.Value(merchantId)
    if nil != err {
        return nil, err
    }
    return value.Data().(*MerchantVerify), nil
}

type merchantVerifyDaoCacheKey struct {
    accessorId string
    merchantId string
}

func (c *MerchantVerifyDaoCache) accessorMerchantVerifyLoader(key interface{}, _ ...interface{}) (*gokits.CacheItem, error) {
    cacheKey := key.(*merchantVerifyDaoCacheKey)
    verifies, err := c.dao.QueryAccessorMerchants(cacheKey.accessorId, cacheKey.merchantId)
    if nil != err {
        return nil, err
    }
    // 缓存1min
    return gokits.NewCacheItem(cacheKey, time.Minute, verifies), nil
}

func (c *MerchantVerifyDaoCache) queryAccessorMerchantById(accessorId, merchantId string) ([]*MerchantVerify, error) {
    value, err := c.tableAccessorMerchant.Value(&merchantVerifyDaoCacheKey{
        accessorId: accessorId, merchantId: merchantId})
    if nil != err {
        return nil, err
    }
    return value.Data().([]*MerchantVerify), nil
}

var (
    merchantIdIllegal = BaseResp{
        ErrorCode: "MERCHANT_ID_ILLEGAL",
        ErrorDesc: "MerchantId is Illegal",
    }
    merchantAccessUnauthorized = BaseResp{
        ErrorCode: "MERCHANT_ACCESS_UNAUTHORIZED",
        ErrorDesc: "Merchant access Unauthorized",
    }
)

func MerchantVerifyInterceptor(ctx iris.Context, cache *MerchantVerifyDaoCache) {
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

    _, err := cache.queryMerchantById(merchantId)
    if nil != err {
        ctx.StopWithJSON(iris.StatusOK, merchantIdIllegal)
        return
    }

    verifies, err := cache.queryAccessorMerchantById(accessorId, merchantId)
    if nil != err || 0 == len(verifies) {
        ctx.StopWithJSON(iris.StatusOK, merchantAccessUnauthorized)
        return
    }

    ctx.Application().Logger().Debug("MerchantVerify done")
    ctx.Next()
}

/****************************************************************************************************/

func init() {
    RegisterDependency("lannister.MerchantManageDao", GetMerchantManageDao)
    RegisterController("lannister.MerchantManageController", &MerchantManageController{})

    RegisterDependency("lannister.MerchantVerifyDao", GetMerchantVerifyDao)
    RegisterDependency("lannister.MerchantVerifyDaoCache", NewMerchantVerifyDaoCache)
    RegisterMiddleware("lannister.MerchantVerifyInterceptor", MerchantVerifyInterceptor)
}
