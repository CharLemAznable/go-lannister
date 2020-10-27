package apptest

import (
    "errors"
    . "github.com/CharLemAznable/go-lannister/base"
    "github.com/CharLemAznable/sqlx"
)

type TestAccessorManageDao struct{}

func NewAccessorManageDao(_ *sqlx.DB) AccessorManageDao {
    return &TestAccessorManageDao{}
}

func (d *TestAccessorManageDao) QueryAccessor(accessorId string) (*AccessorManage, error) {
    if err := accessorErrors[accessorId]; nil != err {
        return &AccessorManage{}, err
    }
    accessor := accessors[accessorId]
    return &AccessorManage{
        AccessorId:     accessor["AccessorId"],
        AccessorName:   accessor["AccessorName"],
        AccessorPubKey: accessor["AccessorPubKey"],
    }, nil
}

func (d *TestAccessorManageDao) UpdateAccessor(accessorId string, manage *AccessorManage) (int64, error) {
    if err := accessorErrors[accessorId]; nil != err {
        return 0, err
    }
    accessor := accessors[accessorId]
    if "" != manage.AccessorName {
        accessor["AccessorName"] = manage.AccessorName
    }
    if "" != manage.AccessorPubKey {
        accessor["AccessorPubKey"] = manage.AccessorPubKey
    }
    if "" != manage.PayNotifyUrl {
        accessor["PayNotifyUrl"] = manage.PayNotifyUrl
    }
    if "" != manage.RefundNotifyUrl {
        accessor["RefundNotifyUrl"] = manage.RefundNotifyUrl
    }
    return 1, nil
}

func (d *TestAccessorManageDao) UpdateKeyPair(accessorId, nonsense, pubKey, prvKey string) error {
    if err := accessorErrors[accessorId]; nil != err {
        return err
    }
    accessor := accessors[accessorId]
    accessor["PubKey"] = pubKey
    return nil
}

type TestAccessorVerifyDao struct{}

func NewAccessorVerifyDao(_ *sqlx.DB) AccessorVerifyDao {
    return &TestAccessorVerifyDao{}
}

func (d *TestAccessorVerifyDao) QueryAccessorById(accessorId string) (*AccessorVerify, error) {
    accessor, ok := accessors[accessorId]
    if !ok {
        return &AccessorVerify{}, errors.New("AccessorNotExists")
    }
    return &AccessorVerify{
        AccessorId:     accessor["AccessorId"],
        AccessorPubKey: accessor["AccessorPubKey"],
    }, nil
}

func init() {
    RegisterAccessorManageDao("", NewAccessorManageDao)
    RegisterAccessorVerifyDao("", NewAccessorVerifyDao)
}
