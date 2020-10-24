package types

import (
    "database/sql"
    "github.com/CharLemAznable/sqlx"
    "github.com/stretchr/testify/assert"
    "testing"
)

func TestRegisterAccessorManageDao(t *testing.T) {
    a := assert.New(t)

    a.Nil(GetAccessorManageDao(nil))

    db := sqlx.NewDb(nil, "fake")
    a.Nil(GetAccessorManageDao(db))

    flag := ""
    RegisterAccessorManageDaoConstructor("fake",
        func(db *sqlx.DB) AccessorManageDao {
            flag = "fake"
            return nil
        })
    a.Nil(GetAccessorManageDao(db))
    a.Equal("fake", flag)
}

func TestRegisterAccessorVerifyDao(t *testing.T) {
    a := assert.New(t)

    a.Nil(GetAccessorVerifyDao(nil))

    db := sqlx.NewDb(&sql.DB{}, "fake")
    a.Nil(GetAccessorVerifyDao(db))

    flag := ""
    RegisterAccessorVerifyDaoConstructor("fake",
        func(db *sqlx.DB) AccessorVerifyDao {
            flag = "fake"
            return nil
        })
    a.Nil(GetAccessorVerifyDao(db))
    a.Equal("fake", flag)
}
