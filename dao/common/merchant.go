package common

import (
    "github.com/CharLemAznable/go-lannister/base"
    "github.com/CharLemAznable/sqlx"
    "github.com/kataras/iris/v12"
)

type MerchantManageSql interface {
    QueryMerchantById() string
    QueryMerchantByCode() string
    CreateMerchant() string
    UpdateMerchant() string
    UpdateAccessorMerchant() string
    QueryMerchants() string
    QueryMerchant() string
}

var merchantManageSqlRegistry = NewSqlRegistry("MerchantManageSql")

func RegisterMerchantManageSql(name string, sql MerchantManageSql) {
    merchantManageSqlRegistry.Register(name, sql)
}

func GetMerchantManageSql(db *sqlx.DB) MerchantManageSql {
    return merchantManageSqlRegistry.GetSql(db).(MerchantManageSql)
}

type MerchantManageDao struct {
    db  *sqlx.DB
    sql MerchantManageSql
}

func NewMerchantManageDao(db *sqlx.DB) base.MerchantManageDao {
    return &MerchantManageDao{db: db, sql: GetMerchantManageSql(db)}
}

func (d *MerchantManageDao) QueryMerchantById(merchantId string) (*base.MerchantManage, error) {
    manage := &base.MerchantManage{}
    err := d.db.NamedGet(manage, d.sql.QueryMerchantById(),
        iris.Map{"MerchantId": merchantId})
    return manage, err
}

func (d *MerchantManageDao) QueryMerchantByCode(merchantCode string) (*base.MerchantManage, error) {
    manage := &base.MerchantManage{}
    err := d.db.NamedGet(manage, d.sql.QueryMerchantByCode(),
        iris.Map{"MerchantCode": merchantCode})
    return manage, err
}

func (d *MerchantManageDao) CreateMerchant(accessorId, merchantId, merchantName, merchantCode string) (int64, error) {
    result, err := d.db.NamedExec(d.sql.CreateMerchant(), iris.Map{"AccessorId": accessorId,
        "MerchantId": merchantId, "MerchantName": merchantName, "MerchantCode": merchantCode})
    if nil != err {
        return 0, err
    }
    return result.RowsAffected()
}

func (d *MerchantManageDao) UpdateMerchant(accessorId, merchantId, merchantName, merchantCode string) (int64, error) {
    result, err := d.db.NamedExec(d.sql.UpdateMerchant(), iris.Map{"AccessorId": accessorId,
        "MerchantId": merchantId, "MerchantName": merchantName, "MerchantCode": merchantCode})
    if nil != err {
        return 0, err
    }
    return result.RowsAffected()
}

func (d *MerchantManageDao) UpdateAccessorMerchant(accessorId, merchantId string) (int64, error) {
    result, err := d.db.NamedExec(d.sql.UpdateAccessorMerchant(),
        iris.Map{"AccessorId": accessorId, "MerchantId": merchantId})
    if nil != err {
        return 0, err
    }
    return result.RowsAffected()
}

func (d *MerchantManageDao) QueryMerchants(accessorId string) ([]*base.MerchantManage, error) {
    merchants := make([]*base.MerchantManage, 0)
    err := d.db.NamedSelect(&merchants, d.sql.QueryMerchants(),
        iris.Map{"AccessorId": accessorId})
    return merchants, err
}

func (d *MerchantManageDao) QueryMerchant(accessorId, merchantId string) (*base.MerchantManage, error) {
    merchant := &base.MerchantManage{}
    err := d.db.NamedGet(merchant, d.sql.QueryMerchant(), iris.Map{
        "AccessorId": accessorId, "MerchantId": merchantId})
    return merchant, err
}

type MerchantVerifySql interface {
    QueryMerchantVerify() string
    QueryAccessorMerchantVerifies() string
}

var merchantVerifySqlRegistry = NewSqlRegistry("MerchantVerifySql")

func RegisterMerchantVerifySql(name string, sql MerchantVerifySql) {
    merchantVerifySqlRegistry.Register(name, sql)
}

func GetMerchantVerifySql(db *sqlx.DB) MerchantVerifySql {
    return merchantVerifySqlRegistry.GetSql(db).(MerchantVerifySql)
}

type MerchantVerifyDao struct {
    db  *sqlx.DB
    sql MerchantVerifySql
}

func NewMerchantVerifyDao(db *sqlx.DB) base.MerchantVerifyDao {
    return &MerchantVerifyDao{db: db, sql: GetMerchantVerifySql(db)}
}

func (d *MerchantVerifyDao) QueryMerchant(merchantId string) (*base.MerchantVerify, error) {
    verify := &base.MerchantVerify{}
    err := d.db.NamedGet(verify, d.sql.QueryMerchantVerify(),
        iris.Map{"MerchantId": merchantId})
    return verify, err
}

func (d *MerchantVerifyDao) QueryAccessorMerchants(accessorId, merchantId string) ([]*base.MerchantVerify, error) {
    verifies := make([]*base.MerchantVerify, 0)
    err := d.db.NamedSelect(&verifies, d.sql.QueryAccessorMerchantVerifies(),
        iris.Map{"AccessorId": accessorId, "MerchantId": merchantId})
    return verifies, err
}
