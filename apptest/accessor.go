package apptest

import (
    "errors"
    "github.com/CharLemAznable/go-lannister/base"
    "github.com/CharLemAznable/sqlx"
)

type TestAccessorManageDao struct{}

func NewAccessorManageDao(_ *sqlx.DB) base.AccessorManageDao {
    return &TestAccessorManageDao{}
}

func (d *TestAccessorManageDao) QueryAccessor(accessorId string) (*base.AccessorManage, error) {
    if err := accessorErrors[accessorId]; nil != err {
        return &base.AccessorManage{}, err
    }
    accessor := accessors[accessorId]
    return &base.AccessorManage{
        AccessorId:     accessor["AccessorId"],
        AccessorName:   accessor["AccessorName"],
        AccessorPubKey: accessor["AccessorPubKey"],
    }, nil
}

func (d *TestAccessorManageDao) UpdateAccessor(accessorId string, manage *base.AccessorManage) (int64, error) {
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

func NewAccessorVerifyDao(_ *sqlx.DB) base.AccessorVerifyDao {
    return &TestAccessorVerifyDao{}
}

func (d *TestAccessorVerifyDao) QueryAccessorById(accessorId string) (*base.AccessorVerify, error) {
    accessor, ok := accessors[accessorId]
    if !ok {
        return &base.AccessorVerify{}, errors.New("AccessorNotExists")
    }
    return &base.AccessorVerify{
        AccessorId:     accessor["AccessorId"],
        AccessorPubKey: accessor["AccessorPubKey"],
    }, nil
}

func init() {
    base.RegisterAccessorManageDao("", NewAccessorManageDao)
    base.RegisterAccessorVerifyDao("", NewAccessorVerifyDao)
}
