package apptest

import (
    "errors"
    "github.com/CharLemAznable/go-lannister/base"
    "github.com/CharLemAznable/sqlx"
    "strings"
)

type ConfigDao struct{}

func NewConfigDao(_ *sqlx.DB) base.ConfigDao {
    return &ConfigDao{}
}

func (d *ConfigDao) InsertMerchantAPIParams(
    accessorId, merchantId, apiName string,
    params map[string]string) error {
    if "error-api" == strings.ToLower(apiName) {
        return errors.New("MockError")
    }
    apiParams := getMerchantApiParams(merchantId, apiName)
    for key, value := range params {
        apiParams[key] = value
    }
    return nil
}

func init() {
    base.RegisterConfigDao("", NewConfigDao)
}
