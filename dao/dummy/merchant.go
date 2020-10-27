package dummy

import (
    "errors"
    "github.com/CharLemAznable/go-lannister/types"
    "github.com/CharLemAznable/sqlx"
)

type MerchantManageDao struct{}

func NewMerchantManageDao(_ *sqlx.DB) types.MerchantManageDao {
    return &MerchantManageDao{}
}

func (d *MerchantManageDao) QueryMerchantById(merchantId string) (*types.MerchantManage, error) {
    merchant, ok := merchantById[merchantId]
    if !ok {
        return &types.MerchantManage{}, errors.New("MerchantNotExists")
    }
    return &types.MerchantManage{
        MerchantId:   merchant["MerchantId"],
        MerchantName: merchant["MerchantName"],
        MerchantCode: merchant["MerchantCode"],
    }, nil
}

func (d *MerchantManageDao) QueryMerchantByCode(merchantCode string) (*types.MerchantManage, error) {
    merchant, ok := merchantByCode[merchantCode]
    if !ok {
        return &types.MerchantManage{}, errors.New("MerchantNotExists")
    }
    return &types.MerchantManage{
        MerchantId:   merchant["MerchantId"],
        MerchantName: merchant["MerchantName"],
        MerchantCode: merchant["MerchantCode"],
    }, nil
}

func (d *MerchantManageDao) CreateMerchant(accessorId, merchantId, merchantName, merchantCode string) (int64, error) {
    merchant := map[string]string{
        "MerchantId":   merchantId,
        "MerchantName": merchantName,
        "MerchantCode": merchantCode,
    }
    merchantById[merchantId] = merchant
    merchantByCode[merchantCode] = merchant
    return 1, nil
}

func (d *MerchantManageDao) UpdateMerchant(accessorId, merchantId, merchantName, merchantCode string) (int64, error) {
    merchant := merchantById[merchantId]
    merchant["MerchantName"] = merchantName
    originalCode := merchant["MerchantCode"]
    merchant["MerchantCode"] = merchantCode

    delete(merchantByCode, originalCode)
    merchantByCode[merchantCode] = merchant
    return 1, nil
}

func (d *MerchantManageDao) UpdateAccessorMerchant(accessorId, merchantId string) (int64, error) {
    getAccessorMerchants(accessorId)[merchantId] = present
    return 0, nil
}

func (d *MerchantManageDao) QueryMerchants(accessorId string) ([]*types.MerchantManage, error) {
    resultMap := map[string]*types.MerchantManage{}
    for _, _accessorId := range []string{accessorId, "0"} {
        merchants := getAccessorMerchants(_accessorId)
        for merchantId := range merchants {
            if merchant, ok := merchantById[merchantId]; ok {
                resultMap[merchantId] = &types.MerchantManage{
                    MerchantId:   merchant["MerchantId"],
                    MerchantName: merchant["MerchantName"],
                    MerchantCode: merchant["MerchantCode"],
                }
            }
        }
    }
    var result []*types.MerchantManage
    for _, merchant := range resultMap {
        result = append(result, merchant)
    }
    return result, nil
}

func (d *MerchantManageDao) QueryMerchant(accessorId, merchantId string) (*types.MerchantManage, error) {
    for _, _accessorId := range []string{accessorId, "0"} {
        merchants := getAccessorMerchants(_accessorId)
        if _, ok := merchants[merchantId]; ok {
            if merchant, ok := merchantById[merchantId]; ok {
                return &types.MerchantManage{
                    MerchantId:   merchant["MerchantId"],
                    MerchantName: merchant["MerchantName"],
                    MerchantCode: merchant["MerchantCode"],
                }, nil
            }
        }
    }
    return nil, errors.New("MerchantNotExists")
}

type MerchantVerifyDao struct{}

func NewMerchantVerifyDao(_ *sqlx.DB) types.MerchantVerifyDao {
    return &MerchantVerifyDao{}
}

func (d *MerchantVerifyDao) QueryMerchant(merchantId string) (*types.MerchantVerify, error) {
    merchant, ok := merchantById[merchantId]
    if !ok {
        return &types.MerchantVerify{}, errors.New("MerchantNotExists")
    }
    return &types.MerchantVerify{
        MerchantId: merchant["MerchantId"],
    }, nil
}

func (d *MerchantVerifyDao) QueryAccessorMerchants(accessorId, merchantId string) ([]*types.MerchantVerify, error) {
    resultMap := map[string]*types.MerchantVerify{}
    for _, _accessorId := range []string{accessorId, "0"} {
        merchants := getAccessorMerchants(_accessorId)
        if _, ok := merchants[merchantId]; ok {
            if merchant, ok := merchantById[merchantId]; ok {
                resultMap[merchantId] = &types.MerchantVerify{
                    AccessorId: _accessorId,
                    MerchantId: merchant["MerchantId"],
                }
            }
        }
    }
    var result []*types.MerchantVerify
    for _, merchant := range resultMap {
        result = append(result, merchant)
    }
    return result, nil
}

func init() {
    types.RegisterMerchantManageDaoConstructor("dummy", NewMerchantManageDao)
    types.RegisterMerchantVerifyDaoConstructor("dummy", NewMerchantVerifyDao)
}
