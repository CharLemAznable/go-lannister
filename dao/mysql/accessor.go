package mysql

import (
    "github.com/CharLemAznable/go-lannister/base"
    "github.com/CharLemAznable/sqlx"
    "github.com/kataras/iris/v12"
)

const queryAccessorById = `
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

func updateAccessorInfo(arg interface{}) (string, error) {
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

const updateKeyPairById = `
update accessor a
   set a.update_time    = now()
      ,a.pub_key        = :PubKey
      ,a.prv_key        = :PrvKey
      ,a.nonsense       = :Nonsense
 where a.accessor_id    = :AccessorId
   and a.nonsense      != :Nonsense
   and a.enabled        = 1
`

type AccessorManageDao struct {
    db *sqlx.DB
}

func NewAccessorManageDao(db *sqlx.DB) base.AccessorManageDao {
    return &AccessorManageDao{db: db}
}

func (d *AccessorManageDao) QueryAccessor(accessorId string) (*base.AccessorManage, error) {
    manage := &base.AccessorManage{}
    err := d.db.NamedGet(manage, queryAccessorById,
        iris.Map{"AccessorId": accessorId})
    return manage, err
}

func (d *AccessorManageDao) UpdateAccessor(accessorId string, manage *base.AccessorManage) (int64, error) {
    manage.AccessorId = accessorId
    result, err := d.db.DynamicNamedExec(updateAccessorInfo, manage)
    if nil != err {
        return 0, err
    }
    return result.RowsAffected()
}

func (d *AccessorManageDao) UpdateKeyPair(accessorId, nonsense, pubKey, prvKey string) error {
    _, err := d.db.NamedExec(updateKeyPairById, iris.Map{
        "AccessorId": accessorId, "Nonsense": nonsense, "PubKey": pubKey, "PrvKey": prvKey})
    return err
}

const queryAccessorVerify = `
select a.accessor_id        as "AccessorId"
      ,a.accessor_pub_key   as "AccessorPubKey"
  from accessor a
 where a.accessor_id        = :AccessorId
   and a.enabled            = 1
`

type AccessorVerifyDao struct {
    db *sqlx.DB
}

func NewAccessorVerifyDao(db *sqlx.DB) base.AccessorVerifyDao {
    return &AccessorVerifyDao{db: db}
}

func (d *AccessorVerifyDao) QueryAccessorById(accessorId string) (*base.AccessorVerify, error) {
    verify := &base.AccessorVerify{}
    err := d.db.NamedGet(verify, queryAccessorVerify,
        iris.Map{"AccessorId": accessorId})
    return verify, err
}

func init() {
    base.RegisterAccessorManageDao("mysql", NewAccessorManageDao)
    base.RegisterAccessorVerifyDao("mysql", NewAccessorVerifyDao)
}
