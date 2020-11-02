package daotest

import (
    "github.com/CharLemAznable/sqlx"
    "io/ioutil"
    "os"
    "strings"
    "testing"

    _ "github.com/CharLemAznable/go-lannister/dao/go-sqlite3"
    _ "github.com/CharLemAznable/go-lannister/dao/mysql"
)

func TestMain(m *testing.M) {
    dbSet := map[string]*sqlx.DB{}
    for name, daoData := range DaoDataSet {
        db, err := sqlx.Open(daoData["DriverName"], daoData["DataSourceName"])
        if nil != err {
            continue
        }
        dbSet[name] = db
        TestConfigSet[name] = map[string]string{
            "DriverName":     daoData["DriverName"],
            "DataSourceName": daoData["DataSourceName"],
        }
    }

    // create tables
    for name, db := range dbSet {
        initSql, _ := ioutil.ReadFile(DaoDataSet[name]["InitSqlFile"])
        db.MustExec(string(initSql))
    }
    // setup
    for name, db := range dbSet {
        setupSql, _ := ioutil.ReadFile(DaoDataSet[name]["SetupSqlFile"])
        db.MustExec(strings.ReplaceAll(string(setupSql),
            "#PublicKeyString#", PublicKeyString))
    }

    code := m.Run()

    // teardown
    for name, db := range dbSet {
        teardownSql, _ := ioutil.ReadFile(DaoDataSet[name]["TeardownSqlFile"])
        db.MustExec(string(teardownSql))
    }

    os.Exit(code)
}
