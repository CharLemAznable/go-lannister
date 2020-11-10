package common

import (
    "github.com/CharLemAznable/go-lannister/base"
    "github.com/CharLemAznable/sqlx"
    "github.com/kataras/iris/v12"
)

type AccessorManageSql interface {
    QueryAccessorById() string
    UpdateAccessorInfo(arg interface{}) (string, error)
    UpdateKeyPairById() string
}

var accessorManageSqlRegistry = NewSqlRegistry("AccessorManageSql")

func RegisterAccessorManageSql(name string, sql AccessorManageSql) {
    accessorManageSqlRegistry.Register(name, sql)
}

func GetAccessorManageSql(db *sqlx.DB) AccessorManageSql {
    return accessorManageSqlRegistry.GetSql(db).(AccessorManageSql)
}

type AccessorManageDao struct {
    db  *sqlx.DB
    sql AccessorManageSql
}

func NewAccessorManageDao(db *sqlx.DB) base.AccessorManageDao {
    return &AccessorManageDao{db: db, sql: GetAccessorManageSql(db)}
}

func (d *AccessorManageDao) QueryAccessor(accessorId string) (*base.AccessorManage, error) {
    manage := &base.AccessorManage{}
    err := d.db.NamedGet(manage, d.sql.QueryAccessorById(),
        iris.Map{"AccessorId": accessorId})
    return manage, err
}

func (d *AccessorManageDao) UpdateAccessor(accessorId string, manage *base.AccessorManage) (int64, error) {
    manage.AccessorId = accessorId
    result, err := d.db.DynamicNamedExec(d.sql.UpdateAccessorInfo, manage)
    if nil != err {
        return 0, err
    }
    return result.RowsAffected()
}

func (d *AccessorManageDao) UpdateKeyPair(accessorId, nonsense, pubKey, prvKey string) (int64, error) {
    result, err := d.db.NamedExec(d.sql.UpdateKeyPairById(), iris.Map{
        "AccessorId": accessorId, "Nonsense": nonsense, "PubKey": pubKey, "PrvKey": prvKey})
    if nil != err {
        return 0, err
    }
    return result.RowsAffected()
}

/****************************************************************************************************/

type AccessorVerifySql interface {
    QueryAccessorVerify() string
}

var accessorVerifySqlRegistry = NewSqlRegistry("AccessorVerifySql")

func RegisterAccessorVerifySql(name string, sql AccessorVerifySql) {
    accessorVerifySqlRegistry.Register(name, sql)
}

func GetAccessorVerifySql(db *sqlx.DB) AccessorVerifySql {
    return accessorVerifySqlRegistry.GetSql(db).(AccessorVerifySql)
}

type AccessorVerifyDao struct {
    db  *sqlx.DB
    sql AccessorVerifySql
}

func NewAccessorVerifyDao(db *sqlx.DB) base.AccessorVerifyDao {
    return &AccessorVerifyDao{db: db, sql: GetAccessorVerifySql(db)}
}

func (d *AccessorVerifyDao) QueryAccessor(accessorId string) (*base.AccessorVerify, error) {
    verify := &base.AccessorVerify{}
    err := d.db.NamedGet(verify, d.sql.QueryAccessorVerify(),
        iris.Map{"AccessorId": accessorId})
    return verify, err
}
