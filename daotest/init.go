package daotest

import (
    "github.com/CharLemAznable/gokits"
    "os"
)

// add this src file for fixing test coverage.
var (
    DaoDataSet = map[string]map[string]string{
        "mysql": {
            "DriverName":      "mysql",
            "DataSourceName":  os.Getenv("MYSQL_DATA_SOURCE_NAME"),
            "InitSqlFile":     "mysql.init.sql",
            "SetupSqlFile":    "mysql.setup.sql",
            "TeardownSqlFile": "mysql.teardown.sql",
        },
        "sqlite3": {
            "DriverName":      "sqlite3",
            "DataSourceName":  os.Getenv("SQLITE_DATA_SOURCE_NAME"),
            "InitSqlFile":     "sqlite3.init.sql",
            "SetupSqlFile":    "sqlite3.setup.sql",
            "TeardownSqlFile": "sqlite3.teardown.sql",
        },
    }
    TestConfigSet = map[string]map[string]string{}

    GeneratedKeyPair, _ = gokits.GenerateRSAKeyPairDefault()
    PrivateKeyString, _ = GeneratedKeyPair.RSAPrivateKeyEncoded()
    PublicKeyString, _  = GeneratedKeyPair.RSAPublicKeyEncoded()
)
