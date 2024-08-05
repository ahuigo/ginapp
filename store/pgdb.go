package store

import (
	"fmt"
	"ginapp/conf"
	"log"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// DB 存储
type DB struct {
	*gorm.DB
	config          string
	transactionKind transactionKind
}

type transactionKind int

const (
	transactionKindNone transactionKind = iota
	transactionKindNormal
	transactionKindInternal
)

// Begin 开启一个事务
func (s *DB) Begin() (*DB, func()) {
	if s.transactionKind != transactionKindNone {
		// in transaction: keep DB, make internal transaction.
		t := &DB{
			DB:              s.DB,
			config:          s.config,
			transactionKind: transactionKindInternal,
		}
		return t, func() {}
	}

	// start normal transaction.
	transaction := &DB{
		DB:              s.DB.Begin(),
		config:          s.config,
		transactionKind: transactionKindNormal,
	}
	recovery := func() {
		if r := recover(); r != nil {
			transaction.Rollback()
			panic(r)
		}
	}
	return transaction, recovery
}

// Rollback rollback transaction
func (s *DB) Rollback() *DB {
	if s.transactionKind == transactionKindNormal {
		s.DB.Rollback()
	}
	return s
}

// Commit commit transaction
func (s *DB) Commit() *DB {
	if s.transactionKind == transactionKindNormal {
		s.DB.Commit()
	}
	return s
}
func (s *DB) Close() {
	if sqlDB, err := s.DB.DB(); err != nil {
		panic(err)
	} else {
		sqlDB.Close()
	}
}

var pgDB *DB

// GetPgDB 从配置中新建 Postgres 存储
func GetPgDB() *DB {
	if pgDB == nil {
		pgDB = newPgDBWithConfig()
	}
	return pgDB
}

// newPgDBWithConfig 从指定配置中新建 Postgres 存储
func newPgDBWithConfig() *DB {
	conf := conf.GetConf()
	pgConf := conf.PgConf
	host := pgConf.Host
	port := pgConf.Port
	user := pgConf.User
	passwd := pgConf.Password
	dbName := pgConf.Dbname
	// conn := `host=localhost user=role1 dbname=ahuigo sslmode=disable password=`
	pgConnectStr := fmt.Sprintf("host=%s port=%d user=%s dbname=%s password='%s' sslmode=disable", host, port, user, dbName, passwd)
	pgDB := newDB(pgConnectStr)
	if gin.Mode() == gin.DebugMode || gin.Mode() == gin.TestMode {
		pgDB.DB = pgDB.DB.Debug()
	}
	if conf.App.Mode == "debug" {
		pgDB.DB = pgDB.DB.Debug()
	}
	return pgDB
}

// newDB create store
func newDB(connectStr string) *DB {
	db, err := gorm.Open(postgres.Open(connectStr))
	if err != nil {
		log.Fatalf("gorm open error: %s, connectStr: %v", err, connectStr)
	}

	// use otelgorm plugin
	// if err := db.Use(otelgorm.NewPlugin()); err != nil {
	// 	panic(err)
	// }

	sqlDB, _ := db.DB()

	// SetMaxIdleConns sets the maximum number of connections in the idle connection pool.
	sqlDB.SetMaxIdleConns(10)

	// SetMaxOpenConns sets the maximum number of open connections to the database.
	sqlDB.SetMaxOpenConns(100)

	// SetConnMaxLifetime sets the maximum amount of time a connection may be reused.
	sqlDB.SetConnMaxLifetime(5 * time.Minute)

	// use callback trace db operations
	// otgorm.RegisterCallbacks(db)

	return &DB{
		DB:     db,
		config: connectStr,
	}
}
