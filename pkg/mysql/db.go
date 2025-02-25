package mysql

import (
	"fmt"
	"goshop/configs"

	"github.com/golang/glog"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func DatabaseInit(cfg *configs.MySQLConfig) (db *gorm.DB, err error) {
	dbconn, err := gorm.Open(mysql.New(mysql.Config{
		DSN: fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=true",
			cfg.Username,
			cfg.Password,
			cfg.Host,
			cfg.Port,
			cfg.DataBase,
			cfg.Charset,
		),
		DefaultStringSize: 256,
		// unsupport before mysql-v5.6.
		DisableDatetimePrecision: true,
		// delete and renew the index when need to rename the index,
		// rename index is unsupported before mysql-v5.7 and MariaDB.
		DontSupportRenameIndex: true,
		// use 'change' to rename the column,
		// rename column is unsupported before mysql-v8.0 and MariaDB.
		DontSupportRenameColumn: true,
		// auto config according the mysql version.
		SkipInitializeWithVersion: false,
	}), &gorm.Config{})
	if err != nil {
		glog.Errorln("[LoginServer] mysql client init failed: ", err.Error())
		return nil, err
	}

	return dbconn, err
}

