package sqlite3

import (
    "github.com/CharLemAznable/go-lannister/base"
    "github.com/CharLemAznable/sqlx"
    "github.com/kataras/iris/v12"
)

const queryAccessorByIdSql = `
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

func updateAccessorInfoSql(arg interface{}) (string, error) {
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

const updateKeyPairByIdSql = `
update accessor
   set update_time      = date('now')
      ,pub_key          = :PubKey
      ,prv_key          = :PrvKey
      ,nonsense         = :Nonsense
 where accessor_id      = :AccessorId
   and nonsense        != :Nonsense
   and enabled          = 1
`

type AccessorManageDao struct {
    db *sqlx.DB
}

func NewAccessorManageDao(db *sqlx.DB) base.AccessorManageDao {
    return &AccessorManageDao{db: db}
}

func (d *AccessorManageDao) QueryAccessor(accessorId string) (*base.AccessorManage, error) {
    manage := &base.AccessorManage{}
    err := d.db.NamedGet(manage, queryAccessorByIdSql,
        iris.Map{"AccessorId": accessorId})
    return manage, err
}

func (d *AccessorManageDao) UpdateAccessor(accessorId string, manage *base.AccessorManage) (int64, error) {
    manage.AccessorId = accessorId
    result, err := d.db.DynamicNamedExec(updateAccessorInfoSql, manage)
    if nil != err {
        return 0, err
    }
    return result.RowsAffected()
}

func (d *AccessorManageDao) UpdateKeyPair(accessorId, nonsense, pubKey, prvKey string) error {
    _, err := d.db.NamedExec(updateKeyPairByIdSql, iris.Map{
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
    base.RegisterAccessorManageDao("sqlite3", NewAccessorManageDao)
    base.RegisterAccessorVerifyDao("sqlite3", NewAccessorVerifyDao)
}
