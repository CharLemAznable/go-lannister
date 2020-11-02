package sqlite3

import (
    "github.com/CharLemAznable/go-lannister/base"
    "github.com/CharLemAznable/sqlx"
    "github.com/kataras/iris/v12"
)

const (
    queryMerchantById = `
select m.merchant_id      as "MerchantId"
      ,m.merchant_name    as "MerchantName"
      ,m.merchant_code    as "MerchantCode"
  from merchant m
 where m.enabled          = 1
   and m.merchant_id      = :MerchantId
`
    queryMerchantByCode = `
select m.merchant_id      as "MerchantId"
      ,m.merchant_name    as "MerchantName"
      ,m.merchant_code    as "MerchantCode"
  from merchant m
 where m.enabled          = 1
   and m.merchant_code    = :MerchantCode
`
    createMerchant = `
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
    updateMerchant = `
update merchant
   set update_time      = date('now')
      ,merchant_name    = :MerchantName
      ,merchant_code    = :MerchantCode
 where enabled          = 1
   and merchant_id      = :MerchantId
`
    updateAccessorMerchant = `
replace into accessor_merchant
      (accessor_id
      ,merchant_id
      ,enabled)
values(:AccessorId
      ,:MerchantId
      ,1)
`
    queryMerchants = `
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
    queryMerchant = `
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
)

type MerchantManageDao struct {
    db *sqlx.DB
}

func NewMerchantManageDao(db *sqlx.DB) base.MerchantManageDao {
    return &MerchantManageDao{db: db}
}

func (d *MerchantManageDao) QueryMerchantById(merchantId string) (*base.MerchantManage, error) {
    manage := &base.MerchantManage{}
    err := d.db.NamedGet(manage, queryMerchantById,
        iris.Map{"MerchantId": merchantId})
    return manage, err
}

func (d *MerchantManageDao) QueryMerchantByCode(merchantCode string) (*base.MerchantManage, error) {
    manage := &base.MerchantManage{}
    err := d.db.NamedGet(manage, queryMerchantByCode,
        iris.Map{"MerchantCode": merchantCode})
    return manage, err
}

func (d *MerchantManageDao) CreateMerchant(accessorId, merchantId, merchantName, merchantCode string) (int64, error) {
    result, err := d.db.NamedExec(createMerchant, iris.Map{"AccessorId": accessorId,
        "MerchantId": merchantId, "MerchantName": merchantName, "MerchantCode": merchantCode})
    if nil != err {
        return 0, err
    }
    return result.RowsAffected()
}

func (d *MerchantManageDao) UpdateMerchant(accessorId, merchantId, merchantName, merchantCode string) (int64, error) {
    result, err := d.db.NamedExec(updateMerchant, iris.Map{"AccessorId": accessorId,
        "MerchantId": merchantId, "MerchantName": merchantName, "MerchantCode": merchantCode})
    if nil != err {
        return 0, err
    }
    return result.RowsAffected()
}

func (d *MerchantManageDao) UpdateAccessorMerchant(accessorId, merchantId string) (int64, error) {
    result, err := d.db.NamedExec(updateAccessorMerchant,
        iris.Map{"AccessorId": accessorId, "MerchantId": merchantId})
    if nil != err {
        return 0, err
    }
    return result.RowsAffected()
}

func (d *MerchantManageDao) QueryMerchants(accessorId string) ([]*base.MerchantManage, error) {
    merchants := make([]*base.MerchantManage, 0)
    err := d.db.NamedSelect(&merchants, queryMerchants,
        iris.Map{"AccessorId": accessorId})
    return merchants, err
}

func (d *MerchantManageDao) QueryMerchant(accessorId, merchantId string) (*base.MerchantManage, error) {
    merchant := &base.MerchantManage{}
    err := d.db.NamedGet(merchant, queryMerchant, iris.Map{
        "AccessorId": accessorId, "MerchantId": merchantId})
    return merchant, err
}

const (
    queryMerchantVerify = `
select m.merchant_id  as "MerchantId"
  from merchant m
 where m.enabled      = 1
   and m.merchant_id  = :MerchantId
`
    queryAccessorMerchantVerifies = `
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
)

type MerchantVerifyDao struct {
    db *sqlx.DB
}

func NewMerchantVerifyDao(db *sqlx.DB) base.MerchantVerifyDao {
    return &MerchantVerifyDao{db: db}
}

func (d *MerchantVerifyDao) QueryMerchant(merchantId string) (*base.MerchantVerify, error) {
    verify := &base.MerchantVerify{}
    err := d.db.NamedGet(verify, queryMerchantVerify,
        iris.Map{"MerchantId": merchantId})
    return verify, err
}

func (d *MerchantVerifyDao) QueryAccessorMerchants(accessorId, merchantId string) ([]*base.MerchantVerify, error) {
    verifies := make([]*base.MerchantVerify, 0)
    err := d.db.NamedSelect(&verifies, queryAccessorMerchantVerifies,
        iris.Map{"AccessorId": accessorId, "MerchantId": merchantId})
    return verifies, err
}

func init() {
    base.RegisterMerchantManageDao("sqlite3", NewMerchantManageDao)
    base.RegisterMerchantVerifyDao("sqlite3", NewMerchantVerifyDao)
}
