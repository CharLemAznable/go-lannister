package sqlite3

import (
    "github.com/CharLemAznable/go-lannister/base"
    "github.com/CharLemAznable/go-lannister/dao/common"
)

type AccessorManageSql struct{}

func (s *AccessorManageSql) QueryAccessorById() string {
    return `
select a.accessor_id        as "AccessorId"
      ,a.accessor_name      as "AccessorName"
      ,a.accessor_pub_key   as "AccessorPubKey"
      ,a.pay_notify_url     as "PayNotifyUrl"
      ,a.refund_notify_url  as "RefundNotifyUrl"
      ,a.pub_key            as "PubKey"
  from accessor a
 where a.accessor_id        = :AccessorId
   and a.enabled            = 1
`
}

func (s *AccessorManageSql) UpdateAccessorInfo(arg interface{}) (string, error) {
    req := arg.(*base.AccessorManage)
    dynamicSql := `
update accessor
   set update_time          = date('now')
`
    if "" != req.AccessorName {
        dynamicSql += ",accessor_name        = :AccessorName"
    }
    if "" != req.AccessorPubKey {
        dynamicSql += ",accessor_pub_key     = :AccessorPubKey"
    }
    if "" != req.PayNotifyUrl {
        dynamicSql += ",pay_notify_url       = :PayNotifyUrl"
    }
    if "" != req.RefundNotifyUrl {
        dynamicSql += ",refund_notify_url    = :RefundNotifyUrl"
    }
    dynamicSql += `
 where accessor_id          = :AccessorId
   and enabled              = 1
`
    return dynamicSql, nil
}

func (s *AccessorManageSql) UpdateKeyPairById() string {
    return `
update accessor
   set update_time      = date('now')
      ,pub_key          = :PubKey
      ,prv_key          = :PrvKey
      ,nonsense         = :Nonsense
 where accessor_id      = :AccessorId
   and nonsense        != :Nonsense
   and enabled          = 1
`
}

type AccessorVerifySql struct{}

func (s *AccessorVerifySql) QueryAccessorVerify() string {
    return `
select a.accessor_id        as "AccessorId"
      ,a.accessor_pub_key   as "AccessorPubKey"
  from accessor a
 where a.accessor_id        = :AccessorId
   and a.enabled            = 1
`
}

func init() {
    common.RegisterAccessorManageSql("sqlite3", &AccessorManageSql{})
    common.RegisterAccessorVerifySql("sqlite3", &AccessorVerifySql{})

    base.RegisterAccessorManageDao("sqlite3", common.NewAccessorManageDao)
    base.RegisterAccessorVerifyDao("sqlite3", common.NewAccessorVerifyDao)
}
