package base

import (
    "github.com/CharLemAznable/sqlx"
)

type ConfigDao interface {
    InsertMerchantAPIParams(
        accessorId, merchantId, apiName string,
        params map[string]string) error
}

var configDaoRegistry = NewDaoRegistry("ConfigDao")

func RegisterConfigDao(name string, builder func(*sqlx.DB) ConfigDao) {
    configDaoRegistry.Register(name, builder)
}

func GetConfigDao(db *sqlx.DB) ConfigDao {
    return configDaoRegistry.GetDao(db).(ConfigDao)
}
