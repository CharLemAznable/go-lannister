package base

import (
    . "github.com/CharLemAznable/go-lannister/elf"
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

type MerchantManageDaoBuilder func(db *sqlx.DB) MerchantManageDao

var merchantManageDaoRegistry = NewDaoRegistry("MerchantManageDao")

func RegisterMerchantManageDao(name string, builder MerchantManageDaoBuilder) {
    merchantManageDaoRegistry.Register(name, builder)
}

func GetMerchantManageDao(db *sqlx.DB) MerchantManageDao {
    return merchantManageDaoRegistry.GetDao(db).(MerchantManageDao)
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

type MerchantVerifyDaoBuilder func(db *sqlx.DB) MerchantVerifyDao

var merchantVerifyDaoRegistry = NewDaoRegistry("MerchantVerifyDaoConstructor")

func RegisterMerchantVerifyDao(name string, builder MerchantVerifyDaoBuilder) {
    merchantVerifyDaoRegistry.Register(name, builder)
}

func GetMerchantVerifyDao(db *sqlx.DB) MerchantVerifyDao {
    return merchantVerifyDaoRegistry.GetDao(db).(MerchantVerifyDao)
}
