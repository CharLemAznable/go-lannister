package apptest

import (
    "errors"
    "github.com/CharLemAznable/go-lannister/base"
    "github.com/CharLemAznable/sqlx"
)

type AccessorManageDao struct{}

func NewAccessorManageDao(_ *sqlx.DB) base.AccessorManageDao {
    return &AccessorManageDao{}
}

func (d *AccessorManageDao) QueryAccessor(accessorId string) (*base.AccessorManage, error) {
    if err := accessorErrors[accessorId]; nil != err {
        return &base.AccessorManage{}, err
    }
    accessor := accessors[accessorId]
    return &base.AccessorManage{
        AccessorId:      accessor["AccessorId"],
        AccessorName:    accessor["AccessorName"],
        AccessorPubKey:  accessor["AccessorPubKey"],
        PayNotifyUrl:    accessor["PayNotifyUrl"],
        RefundNotifyUrl: accessor["RefundNotifyUrl"],
    }, nil
}

func (d *AccessorManageDao) UpdateAccessor(accessorId string, manage *base.AccessorManage) (int64, error) {
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

func (d *AccessorManageDao) UpdateKeyPair(accessorId, nonsense, pubKey, prvKey string) error {
    if err := accessorErrors[accessorId]; nil != err {
        return err
    }
    accessor := accessors[accessorId]
    accessor["PubKey"] = pubKey
    return nil
}

type AccessorVerifyDao struct{}

func NewAccessorVerifyDao(_ *sqlx.DB) base.AccessorVerifyDao {
    return &AccessorVerifyDao{}
}

func (d *AccessorVerifyDao) QueryAccessorById(accessorId string) (*base.AccessorVerify, error) {
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
