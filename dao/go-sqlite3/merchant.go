package sqlite3

import (
    "github.com/CharLemAznable/go-lannister/base"
    "github.com/CharLemAznable/go-lannister/dao/common"
)

type MerchantManageSql struct{}

func (s *MerchantManageSql) QueryMerchantById() string {
    return `
select m.merchant_id      as "MerchantId"
      ,m.merchant_name    as "MerchantName"
      ,m.merchant_code    as "MerchantCode"
  from merchant m
 where m.enabled          = 1
   and m.merchant_id      = :MerchantId
`
}

func (s *MerchantManageSql) QueryMerchantByCode() string {
    return `
select m.merchant_id      as "MerchantId"
      ,m.merchant_name    as "MerchantName"
      ,m.merchant_code    as "MerchantCode"
  from merchant m
 where m.enabled          = 1
   and m.merchant_code    = :MerchantCode
`
}

func (s *MerchantManageSql) CreateMerchant() string {
    return `
replace into merchant
      (merchant_id
      ,merchant_name
      ,merchant_code
      ,enabled
      ,update_accessor)
values(:MerchantId
      ,:MerchantName
      ,:MerchantCode
      ,1
      ,:AccessorId)
`
}

func (s *MerchantManageSql) UpdateMerchant() string {
    return `
update merchant
   set update_time      = date('now')
      ,merchant_name    = :MerchantName
      ,merchant_code    = :MerchantCode
 where enabled          = 1
   and merchant_id      = :MerchantId
`
}

func (s *MerchantManageSql) UpdateAccessorMerchant() string {
    return `
replace into accessor_merchant
      (accessor_id
      ,merchant_id
      ,enabled)
values(:AccessorId
      ,:MerchantId
      ,1)
`
}

func (s *MerchantManageSql) QueryMerchant() string {
    return `
select distinct
       m.merchant_id    as "MerchantId"
      ,m.merchant_name  as "MerchantName"
      ,m.merchant_code  as "MerchantCode"
  from merchant m
      ,accessor a
      ,accessor_merchant r
 where m.merchant_id    = :MerchantId
   and m.enabled        = 1
   and a.accessor_id    = :AccessorId
   and a.enabled        = 1
   and r.merchant_id    = m.merchant_id
   and(r.accessor_id    = 0
    or r.accessor_id    = a.accessor_id)
   and r.enabled        = 1
`
}

func (s *MerchantManageSql) QueryMerchants() string {
    return `
select distinct
       m.merchant_id    as "MerchantId"
      ,m.merchant_name  as "MerchantName"
      ,m.merchant_code  as "MerchantCode"
  from merchant m
      ,accessor a
      ,accessor_merchant r
 where m.enabled        = 1
   and a.accessor_id    = :AccessorId
   and a.enabled        = 1
   and r.merchant_id    = m.merchant_id
   and(r.accessor_id    = 0
    or r.accessor_id    = a.accessor_id)
   and r.enabled        = 1
`
}

type MerchantVerifySql struct{}

func (s *MerchantVerifySql) QueryAccessorMerchantVerifies() string {
    return `
select distinct
       r.accessor_id    as "AccessorId"
      ,r.merchant_id    as "MerchantId"
  from merchant m
      ,accessor a
      ,accessor_merchant r
 where m.merchant_id    = :MerchantId
   and m.enabled        = 1
   and a.accessor_id    = :AccessorId
   and a.enabled        = 1
   and r.merchant_id    = m.merchant_id
   and(r.accessor_id    = 0
    or r.accessor_id    = a.accessor_id)
   and r.enabled        = 1
`
}

func (s *MerchantVerifySql) QueryMerchantVerify() string {
    return `
select m.merchant_id  as "MerchantId"
  from merchant m
 where m.enabled      = 1
   and m.merchant_id  = :MerchantId
`
}

func init() {
    common.RegisterMerchantManageSql("sqlite3", &MerchantManageSql{})
    common.RegisterMerchantVerifySql("sqlite3", &MerchantVerifySql{})

    base.RegisterMerchantManageDao("sqlite3", common.NewMerchantManageDao)
    base.RegisterMerchantVerifyDao("sqlite3", common.NewMerchantVerifyDao)
}
