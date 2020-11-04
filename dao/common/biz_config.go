package common

import (
    "github.com/CharLemAznable/go-lannister/base"
    "github.com/CharLemAznable/sqlx"
)

type ConfigSql interface {
    InsertMerchantAPIParams(arg interface{}) (string, error)
}

var configSqlRegistry = NewSqlRegistry("ConfigSql")

func RegisterConfigSql(name string, sql ConfigSql) {
    configSqlRegistry.Register(name, sql)
}

func GetConfigSql(db *sqlx.DB) ConfigSql {
    return configSqlRegistry.GetSql(db).(ConfigSql)
}

type ConfigSqlParam struct {
    AccessorId string
    MerchantId string
    ApiName    string
    Params     map[string]string
}

type ConfigDao struct {
    db  *sqlx.DB
    sql ConfigSql
}

func NewConfigDao(db *sqlx.DB) base.ConfigDao {
    return &ConfigDao{db: db, sql: GetConfigSql(db)}
}

func (d *ConfigDao) InsertMerchantAPIParams(
    accessorId, merchantId, apiName string,
    params map[string]string) error {
    sqlParam := &ConfigSqlParam{
        AccessorId: accessorId,
        MerchantId: merchantId,
        ApiName:    apiName,
        Params:     params,
    }
    _, err := d.db.DynamicNamedExec(
        d.sql.InsertMerchantAPIParams, sqlParam)
    return err
}
