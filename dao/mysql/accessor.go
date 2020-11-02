package mysql

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
update accessor a
   set a.update_time        = now()
`
    if "" != req.AccessorName {
        dynamicSql += ",a.accessor_name      = :AccessorName"
    }
    if "" != req.AccessorPubKey {
        dynamicSql += ",a.accessor_pub_key   = :AccessorPubKey"
    }
    if "" != req.PayNotifyUrl {
        dynamicSql += ",a.pay_notify_url     = :PayNotifyUrl"
    }
    if "" != req.RefundNotifyUrl {
        dynamicSql += ",a.refund_notify_url  = :RefundNotifyUrl"
    }
    dynamicSql += `
 where a.accessor_id        = :AccessorId
   and a.enabled            = 1
`
    return dynamicSql, nil
}

func (s *AccessorManageSql) UpdateKeyPairById() string {
    return `
update accessor a
   set a.update_time    = now()
      ,a.pub_key        = :PubKey
      ,a.prv_key        = :PrvKey
      ,a.nonsense       = :Nonsense
 where a.accessor_id    = :AccessorId
   and a.nonsense      != :Nonsense
   and a.enabled        = 1
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
    common.RegisterAccessorManageSql("mysql", &AccessorManageSql{})
    common.RegisterAccessorVerifySql("mysql", &AccessorVerifySql{})

    base.RegisterAccessorManageDao("mysql", common.NewAccessorManageDao)
    base.RegisterAccessorVerifyDao("mysql", common.NewAccessorVerifyDao)
}
