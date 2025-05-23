package mysql

import (
	"database/sql"
	"fmt"
	"goshop/configs"
	"time"

	"github.com/golang/glog"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/plugin/dbresolver"
)

type ConnectionPoolConfig struct {
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifeTime time.Duration
	ConnMaxIdleTime time.Duration
}

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
		glog.Errorln("[Mysql] client init failed: ", err.Error())
		return nil, err
	}

	return dbconn, err
}

func DBClusterInit(cfg *configs.MySQLClusterConfig) (*gorm.DB, error) {
	var (
		masterDSN          = cfg.Master.GetDSN()
		replicas           = cfg.Replicas
		replicaDSNs        []string
		replicasDialectors []gorm.Dialector
	)
	db, err := gorm.Open(
		mysql.Open(masterDSN),
		&gorm.Config{},
	)
	if err != nil {
		glog.Errorln("[Mysql] cluster master init failed: ", err.Error())
		return nil, err
	}
	sqldb, err := db.DB()
	if err != nil {
		glog.Errorln("[Mysql] cluster master init failed: ", err.Error())
		return nil, err
	}
	setConnectionPool(sqldb, &ConnectionPoolConfig{
		MaxOpenConns:    20,
		MaxIdleConns:    5,
		ConnMaxLifeTime: 2 * time.Hour,
		ConnMaxIdleTime: 30 * time.Minute,
	})

	for _, repl := range replicas {
		replicaDSN := repl.GetDSN()
		if replicaDSN == "" {
			glog.Warning("[Mysql] replica DSN is empty, skipping: ", repl)
			continue
		}
		replicaDSNs = append(replicaDSNs, replicaDSN)
		replicasDialectors = append(replicasDialectors, mysql.Open(replicaDSN))
	}

	err = db.Use(dbresolver.Register(dbresolver.Config{
		// use `db1` as sources, `db2`, `db3` as replicas
		Sources:  []gorm.Dialector{mysql.Open(masterDSN)},
		Replicas: replicasDialectors,
		// sources/replicas load balancing policy
		Policy: dbresolver.RandomPolicy{},
		// print sources/replicas mode in logger
		TraceResolverMode: true,
	}))
	if err != nil {
		glog.Errorln("[Mysql] cluster replicas init failed: ", err.Error())
		return nil, err
	}

	if len(replicasDialectors) > 0 {
		resolverDb, err := db.DB()
		if err != nil {
			glog.Errorln("[Mysql] cluster replicas init failed: ", err.Error())
			return nil, err
		}
		setConnectionPool(resolverDb, &ConnectionPoolConfig{
			MaxOpenConns:    50,
			MaxIdleConns:    10,
			ConnMaxLifeTime: 2 * time.Hour,
			ConnMaxIdleTime: 30 * time.Minute,
		})
	}

	glog.Infof("[Mysql] cluster init success: Master=[%s], Replicas=[%v]\n", masterDSN, replicaDSNs)

	return db, nil
}

func setConnectionPool(db *sql.DB, cfg *ConnectionPoolConfig) {
	db.SetMaxOpenConns(cfg.MaxOpenConns)
	db.SetMaxIdleConns(cfg.MaxIdleConns)
	db.SetConnMaxLifetime(cfg.ConnMaxLifeTime)
	db.SetConnMaxIdleTime(cfg.ConnMaxIdleTime)
}
