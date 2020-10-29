package mysql

import (
    . "github.com/CharLemAznable/go-lannister/elf"
    "github.com/CharLemAznable/sqlx"
    "github.com/kataras/golog"
    "io/ioutil"
    "os"
    "strings"
    "testing"
)

var (
    mysqlDataSourceName = os.Getenv("MYSQL_DATA_SOURCE_NAME")

    GeneratedKeyPair, _ = GenerateKeyPairDefault()
    PrivateKeyString, _ = GeneratedKeyPair.PrivateKeyEncoded()
    PublicKeyString, _  = GeneratedKeyPair.PublicKeyEncoded()
)

func TestMain(m *testing.M) {
    db, err := sqlx.Open("mysql", mysqlDataSourceName)
    if nil != err {
        golog.Error(err)
        os.Exit(1)
    }

    // create tables
    databaseSetupSql, _ := ioutil.ReadFile("database.setup.sql")
    db.MustExec(string(databaseSetupSql))

    // prepare test data
    databaseTestSql, _ := ioutil.ReadFile("database.test.sql")
    db.MustExec(strings.ReplaceAll(string(databaseTestSql),
        "#PublicKeyString#", PublicKeyString))

    code := m.Run()

    // drop tables
    databaseDropSql, _ := ioutil.ReadFile("database.teardown.sql")
    db.MustExec(string(databaseDropSql))

    os.Exit(code)
}
