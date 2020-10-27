package types

import (
    "github.com/CharLemAznable/sqlx"
)

type MerchantManage struct {
    BaseResp

    MerchantId   string `json:"merchantId"`
    MerchantName string `json:"merchantName"`
    MerchantCode string `json:"merchantCode"`
    AuthorizeAll string `json:"authorizeAll,omitempty"`
}

type MerchantManageDao interface {
    // 根据商户标识查询商户
    QueryMerchantById(merchantId string) (*MerchantManage, error)
    // 根据商户编码查询商户
    QueryMerchantByCode(merchantCode string) (*MerchantManage, error)
    // 创建商户
    CreateMerchant(accessorId, merchantId, merchantName, merchantCode string) (int64, error)
    // 更新商户
    UpdateMerchant(accessorId, merchantId, merchantName, merchantCode string) (int64, error)
    // 创建或更新访问权限
    UpdateAccessorMerchant(accessorId, merchantId string) (int64, error)

    // 查询访问者可访问的商户列表
    QueryMerchants(accessorId string) ([]*MerchantManage, error)
    // 查询指定商户信息
    QueryMerchant(accessorId, merchantId string) (*MerchantManage, error)
}

type MerchantManageDaoConstructor func(db *sqlx.DB) MerchantManageDao

var merchantManageDaoConstructors = NewDaoConstructorRegistry("MerchantManageDaoConstructor")

func RegisterMerchantManageDaoConstructor(name string, constructor MerchantManageDaoConstructor) {
    merchantManageDaoConstructors.Register(name, constructor)
}

func GetMerchantManageDao(db *sqlx.DB) MerchantManageDao {
    return merchantManageDaoConstructors.
        GetDaoConstructor(db).(MerchantManageDaoConstructor)(db)
}

type MerchantVerify struct {
    AccessorId string
    MerchantId string
}

type MerchantVerifyDao interface {
    // 查询商户是否有效
    QueryMerchant(merchantId string) (*MerchantVerify, error)
    // 查询访问者是否有权访问商户
    QueryAccessorMerchants(accessorId, merchantId string) ([]*MerchantVerify, error)
}

type MerchantVerifyDaoConstructor func(db *sqlx.DB) MerchantVerifyDao

var merchantVerifyDaoConstructors = NewDaoConstructorRegistry("MerchantVerifyDaoConstructor")

func RegisterMerchantVerifyDaoConstructor(name string, constructor MerchantVerifyDaoConstructor) {
    merchantVerifyDaoConstructors.Register(name, constructor)
}

func GetMerchantVerifyDao(db *sqlx.DB) MerchantVerifyDao {
    return merchantVerifyDaoConstructors.
        GetDaoConstructor(db).(MerchantVerifyDaoConstructor)(db)
}
