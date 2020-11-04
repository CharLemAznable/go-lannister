package sqlite3

import (
    "github.com/CharLemAznable/go-lannister/base"
    "github.com/CharLemAznable/go-lannister/dao/common"
)

type ConfigSql struct{}

func (s *ConfigSql) InsertMerchantAPIParams(arg interface{}) (string, error) {
    sqlParam := arg.(*common.ConfigSqlParam)
    dynamicSql := `
replace into merchant_api_params
      (merchant_id
      ,api_name
      ,param_name
      ,param_value
      ,update_time
      ,update_accessor)
values
`
    for key, value := range sqlParam.Params {
        dynamicSql += `
      (:MerchantId
      ,lower(:ApiName)
      ,'` + key + `'
      ,'` + value + `'
      ,date('now')
      ,:AccessorId),`
    }
    return dynamicSql[:(len(dynamicSql) - 1)], nil
}

func init() {
    common.RegisterConfigSql("sqlite3", &ConfigSql{})

    base.RegisterConfigDao("sqlite3", common.NewConfigDao)
}
