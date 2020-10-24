package types

import (
    "github.com/CharLemAznable/go-lannister/elf"
    "github.com/CharLemAznable/sqlx"
    "github.com/kataras/golog"
    "sort"
    "strings"
)

type AccessorManage struct {
    BaseResp

    AccessorId      string `json:"accessorId" db:"ACCESSOR_ID"`
    AccessorName    string `json:"accessorName" db:"ACCESSOR_NAME"`
    AccessorPubKey  string `json:"accessorPubKey" db:"ACCESSOR_PUB_KEY"`   // 访问者公钥, 用于平台验证访问者发起的请求
    PayNotifyUrl    string `json:"payNotifyUrl" db:"PAY_NOTIFY_URL"`       // 支付回调地址
    RefundNotifyUrl string `json:"refundNotifyUrl" db:"REFUND_NOTIFY_URL"` // 退款回调地址
    PubKey          string `json:"pubKey" db:"PUB_KEY"`                    // 平台公钥, 用于访问者验证平台回调的请求
}

type AccessorManageDao interface {
    QueryAccessorById(accessorId string) (*AccessorManage, error)
    UpdateAccessorById(accessorId string, manage *AccessorManage) (int64, error)
    UpdateKeyPairById(accessorId, nonsense, pubKey, prvKey string)
}

type AccessorManageDaoConstructor func(db *sqlx.DB) AccessorManageDao

var accessorManageDaoConstructors = NewRegistry("AccessorManageDaoConstructor")

func RegisterAccessorManageDaoConstructor(name string, constructor AccessorManageDaoConstructor) {
    accessorManageDaoConstructors.Register(name, constructor)
}

func GetAccessorManageDao(db *sqlx.DB) AccessorManageDao {
    if nil == db {
        golog.Error("Nil sqlx.DB")
        return nil
    }
    driverName := db.DriverName()
    constructor := accessorManageDaoConstructors.Get(driverName)
    if nil == constructor {
        golog.Errorf("Unknown AccessorManageDaoConstructor"+
            " for driver %q (forgotten import?)", driverName)
        return nil
    }
    return constructor.(AccessorManageDaoConstructor)(db)
}

/****************************************************************************************************/

const (
    NONSENSE  = "nonsense"
    SIGNATURE = "signature"
)

type AccessorVerify struct {
    AccessorId     string `db:"ACCESSOR_ID"`
    AccessorPubKey string `db:"ACCESSOR_PUB_KEY"`
}

func (v *AccessorVerify) Verify(paramMap map[string]string) error {
    signature := paramMap[SIGNATURE]
    delete(paramMap, SIGNATURE)

    names := make([]string, 0)
    for name := range paramMap {
        names = append(names, name)
    }
    sort.Strings(names)

    paramPairs := make([]string, 0)
    for _, name := range names {
        value := paramMap[name]
        if "" == name || "" == value {
            continue
        }
        paramPairs = append(paramPairs, name+"="+value)
    }
    plainText := strings.Join(paramPairs, "&")

    return elf.SHA1WithRSA.VerifyBase64ByKeyString(
        plainText, signature, v.AccessorPubKey)
}

type AccessorVerifyDao interface {
    QueryAccessorById(accessorId string) (*AccessorVerify, error)
}

type AccessorVerifyDaoConstructor func(db *sqlx.DB) AccessorVerifyDao

var accessorVerifyDaoConstructors = NewRegistry("AccessorVerifyDaoConstructor")

func RegisterAccessorVerifyDaoConstructor(name string, constructor AccessorVerifyDaoConstructor) {
    accessorVerifyDaoConstructors.Register(name, constructor)
}

func GetAccessorVerifyDao(db *sqlx.DB) AccessorVerifyDao {
    if nil == db {
        golog.Error("Nil sqlx.DB")
        return nil
    }
    driverName := db.DriverName()
    constructor := accessorVerifyDaoConstructors.Get(driverName)
    if nil == constructor {
        golog.Errorf("Unknown AccessorVerifyDaoConstructor"+
            " for driver %q (forgotten import?)", driverName)
        return nil
    }
    return constructor.(AccessorVerifyDaoConstructor)(db)
}
