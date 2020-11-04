package app

import (
    "github.com/CharLemAznable/go-lannister/base"
    "github.com/CharLemAznable/sqlx"
    "github.com/kataras/iris/v12"
    "github.com/kataras/iris/v12/mvc"
)

var (
    configParamsIllegal = base.BaseResp{
        ErrorCode: "CONFIG_PARAMS_ILLEGAL",
        ErrorDesc: "Config Params is Illegal",
    }
)

type ConfigController struct {
    dao base.ConfigDao
}

func (c *ConfigController) BeforeActivation(b mvc.BeforeActivation) {
    _ = b.Dependencies().Inject(&c.dao) // does not support func register yet
    b.Handle(PostMapping("/{accessorId}/{merchantId}/{apiName}/config", "Config"))
}

func (c *ConfigController) Config(ctx iris.Context) {
    accessorId := ctx.Params().Get("accessorId")
    merchantId := ctx.Params().Get("merchantId")
    apiName := ctx.Params().Get("apiName")
    params := make(map[string]string)
    _ = ctx.ReadJSON(&params)
    if 0 == len(params) {
        ctx.StopWithJSON(iris.StatusOK, configParamsIllegal)
        return
    }

    ctx.Application().Logger().Debugf("Configuring accessorId %s merchantId %s "+
        "apiName %s configuration: %#v", accessorId, merchantId, apiName, params)
    err := c.dao.InsertMerchantAPIParams(accessorId, merchantId, apiName, params)
    if err != nil {
        ctx.Application().Logger().Errorf("Configuring Error: %s", err.Error())
        _, _ = ctx.Text("FAILED")
    } else {
        _, _ = ctx.Text("SUCCESS")
    }
}

/****************************************************************************************************/

func init() {
    RegisterDependency(func(_ *base.Config, db *sqlx.DB) base.ConfigDao {
        return base.GetConfigDao(db)
    })
    RegisterController(&ConfigController{})
}
