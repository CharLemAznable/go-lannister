package apptest

import (
    "errors"
    . "github.com/CharLemAznable/go-lannister/base"
    "github.com/CharLemAznable/sqlx"
)

type TestMerchantManageDao struct{}

func NewMerchantManageDao(_ *sqlx.DB) MerchantManageDao {
    return &TestMerchantManageDao{}
}

func (d *TestMerchantManageDao) QueryMerchantById(merchantId string) (*MerchantManage, error) {
    merchant, ok := merchantById[merchantId]
    if !ok {
        return &MerchantManage{}, errors.New("MerchantNotExists")
    }
    return &MerchantManage{
        MerchantId:   merchant["MerchantId"],
        MerchantName: merchant["MerchantName"],
        MerchantCode: merchant["MerchantCode"],
    }, nil
}

func (d *TestMerchantManageDao) QueryMerchantByCode(merchantCode string) (*MerchantManage, error) {
    merchant, ok := merchantByCode[merchantCode]
    if !ok {
        return &MerchantManage{}, errors.New("MerchantNotExists")
    }
    return &MerchantManage{
        MerchantId:   merchant["MerchantId"],
        MerchantName: merchant["MerchantName"],
        MerchantCode: merchant["MerchantCode"],
    }, nil
}

func (d *TestMerchantManageDao) CreateMerchant(accessorId, merchantId, merchantName, merchantCode string) (int64, error) {
    merchant := map[string]string{
        "MerchantId":   merchantId,
        "MerchantName": merchantName,
        "MerchantCode": merchantCode,
    }
    merchantById[merchantId] = merchant
    merchantByCode[merchantCode] = merchant
    return 1, nil
}

func (d *TestMerchantManageDao) UpdateMerchant(accessorId, merchantId, merchantName, merchantCode string) (int64, error) {
    merchant := merchantById[merchantId]
    merchant["MerchantName"] = merchantName
    originalCode := merchant["MerchantCode"]
    merchant["MerchantCode"] = merchantCode

    delete(merchantByCode, originalCode)
    merchantByCode[merchantCode] = merchant
    return 1, nil
}

func (d *TestMerchantManageDao) UpdateAccessorMerchant(accessorId, merchantId string) (int64, error) {
    getAccessorMerchants(accessorId)[merchantId] = present
    return 0, nil
}

func (d *TestMerchantManageDao) QueryMerchants(accessorId string) ([]*MerchantManage, error) {
    resultMap := map[string]*MerchantManage{}
    for _, _accessorId := range []string{accessorId, "0"} {
        merchants := getAccessorMerchants(_accessorId)
        for merchantId := range merchants {
            if merchant, ok := merchantById[merchantId]; ok {
                resultMap[merchantId] = &MerchantManage{
                    MerchantId:   merchant["MerchantId"],
                    MerchantName: merchant["MerchantName"],
                    MerchantCode: merchant["MerchantCode"],
                }
            }
        }
    }
    var result []*MerchantManage
    for _, merchant := range resultMap {
        result = append(result, merchant)
    }
    return result, nil
}

func (d *TestMerchantManageDao) QueryMerchant(accessorId, merchantId string) (*MerchantManage, error) {
    for _, _accessorId := range []string{accessorId, "0"} {
        merchants := getAccessorMerchants(_accessorId)
        if _, ok := merchants[merchantId]; ok {
            if merchant, ok := merchantById[merchantId]; ok {
                return &MerchantManage{
                    MerchantId:   merchant["MerchantId"],
                    MerchantName: merchant["MerchantName"],
                    MerchantCode: merchant["MerchantCode"],
                }, nil
            }
        }
    }
    return nil, errors.New("MerchantNotExists")
}

type TestMerchantVerifyDao struct{}

func NewMerchantVerifyDao(_ *sqlx.DB) MerchantVerifyDao {
    return &TestMerchantVerifyDao{}
}

func (d *TestMerchantVerifyDao) QueryMerchant(merchantId string) (*MerchantVerify, error) {
    merchant, ok := merchantById[merchantId]
    if !ok {
        return &MerchantVerify{}, errors.New("MerchantNotExists")
    }
    return &MerchantVerify{
        MerchantId: merchant["MerchantId"],
    }, nil
}

func (d *TestMerchantVerifyDao) QueryAccessorMerchants(accessorId, merchantId string) ([]*MerchantVerify, error) {
    resultMap := map[string]*MerchantVerify{}
    for _, _accessorId := range []string{accessorId, "0"} {
        merchants := getAccessorMerchants(_accessorId)
        if _, ok := merchants[merchantId]; ok {
            if merchant, ok := merchantById[merchantId]; ok {
                resultMap[merchantId] = &MerchantVerify{
                    AccessorId: _accessorId,
                    MerchantId: merchant["MerchantId"],
                }
            }
        }
    }
    var result []*MerchantVerify
    for _, merchant := range resultMap {
        result = append(result, merchant)
    }
    return result, nil
}

func init() {
    RegisterMerchantManageDao("", NewMerchantManageDao)
    RegisterMerchantVerifyDao("", NewMerchantVerifyDao)
}
