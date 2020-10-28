package apptest

import (
    "errors"
    "github.com/CharLemAznable/go-lannister/base"
    "github.com/CharLemAznable/sqlx"
)

type TestMerchantManageDao struct{}

func NewMerchantManageDao(_ *sqlx.DB) base.MerchantManageDao {
    return &TestMerchantManageDao{}
}

func (d *TestMerchantManageDao) QueryMerchantById(merchantId string) (*base.MerchantManage, error) {
    merchant, ok := merchantById[merchantId]
    if !ok {
        return &base.MerchantManage{}, errors.New("MerchantNotExists")
    }
    return &base.MerchantManage{
        MerchantId:   merchant["MerchantId"],
        MerchantName: merchant["MerchantName"],
        MerchantCode: merchant["MerchantCode"],
    }, nil
}

func (d *TestMerchantManageDao) QueryMerchantByCode(merchantCode string) (*base.MerchantManage, error) {
    merchant, ok := merchantByCode[merchantCode]
    if !ok {
        return &base.MerchantManage{}, errors.New("MerchantNotExists")
    }
    return &base.MerchantManage{
        MerchantId:   merchant["MerchantId"],
        MerchantName: merchant["MerchantName"],
        MerchantCode: merchant["MerchantCode"],
    }, nil
}

func (d *TestMerchantManageDao) CreateMerchant(accessorId, merchantId, merchantName, merchantCode string) (int64, error) {
    if err := accessorErrors[accessorId]; nil != err {
        return 0, err
    }
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
    if err := accessorErrors[accessorId]; nil != err {
        return 0, err
    }
    getAccessorMerchants(accessorId)[merchantId] = present
    return 1, nil
}

func (d *TestMerchantManageDao) QueryMerchants(accessorId string) ([]*base.MerchantManage, error) {
    if err := accessorErrors[accessorId]; nil != err {
        return []*base.MerchantManage{}, err
    }
    resultMap := map[string]*base.MerchantManage{}
    for _, _accessorId := range []string{accessorId, "0"} {
        merchants := getAccessorMerchants(_accessorId)
        for merchantId := range merchants {
            if merchant, ok := merchantById[merchantId]; ok {
                resultMap[merchantId] = &base.MerchantManage{
                    MerchantId:   merchant["MerchantId"],
                    MerchantName: merchant["MerchantName"],
                    MerchantCode: merchant["MerchantCode"],
                }
            }
        }
    }
    var result []*base.MerchantManage
    for _, merchant := range resultMap {
        result = append(result, merchant)
    }
    return result, nil
}

func (d *TestMerchantManageDao) QueryMerchant(accessorId, merchantId string) (*base.MerchantManage, error) {
    if err := merchantErrors[merchantId]; nil != err {
        return &base.MerchantManage{}, err
    }
    merchant := merchantById[merchantId]
    return &base.MerchantManage{
        MerchantId:   merchant["MerchantId"],
        MerchantName: merchant["MerchantName"],
        MerchantCode: merchant["MerchantCode"],
    }, nil
}

type TestMerchantVerifyDao struct{}

func NewMerchantVerifyDao(_ *sqlx.DB) base.MerchantVerifyDao {
    return &TestMerchantVerifyDao{}
}

func (d *TestMerchantVerifyDao) QueryMerchant(merchantId string) (*base.MerchantVerify, error) {
    merchant, ok := merchantById[merchantId]
    if !ok {
        return &base.MerchantVerify{}, errors.New("MerchantNotExists")
    }
    return &base.MerchantVerify{
        MerchantId: merchant["MerchantId"],
    }, nil
}

func (d *TestMerchantVerifyDao) QueryAccessorMerchants(accessorId, merchantId string) ([]*base.MerchantVerify, error) {
    if err := merchantErrors[merchantId]; nil != err {
        if "MockError" == err.Error() {
            merchantErrors[merchantId] = errors.New("SkipError")
            return []*base.MerchantVerify{}, err
        }
    }
    resultMap := map[string]*base.MerchantVerify{}
    for _, _accessorId := range []string{accessorId, "0"} {
        merchants := getAccessorMerchants(_accessorId)
        if _, ok := merchants[merchantId]; ok {
            if merchant, ok := merchantById[merchantId]; ok {
                resultMap[merchantId] = &base.MerchantVerify{
                    AccessorId: _accessorId,
                    MerchantId: merchant["MerchantId"],
                }
            }
        }
    }
    var result []*base.MerchantVerify
    for _, merchant := range resultMap {
        result = append(result, merchant)
    }
    return result, nil
}

func init() {
    base.RegisterMerchantManageDao("", NewMerchantManageDao)
    base.RegisterMerchantVerifyDao("", NewMerchantVerifyDao)
}
