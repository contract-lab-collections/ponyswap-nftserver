package model

import (
	"context"
	"fmt"
	"nftserver/global/setting"

	"github.com/go-redis/redis/v8"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Model struct {
	ID        uint64 `gorm:"primary_key" json:"id"`
	CreatedAt int64  `gorm:"" json:"created_at,omitempty"`
	UpdatedAt int64  `gorm:"" json:"updated_at,omitempty"`

	// DeletedOn  uint32 `json:"deleted_on"`
	// IsDel      uint8  `json:"is_del"`
}

func NewMysqlDBEngine(s *setting.DbSettingS) (*gorm.DB, error) {
	databaseSetting := s.Database
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=%s&parseTime=%t",
		databaseSetting.UserName,
		databaseSetting.Password,
		databaseSetting.Host,
		databaseSetting.DBName,
		databaseSetting.Charset,
		databaseSetting.ParseTime,
	)
	db, err := gorm.Open(mysql.New(mysql.Config{
		DSN:                       dsn,
		Conn:                      nil,
		SkipInitializeWithVersion: false,
		DefaultStringSize:         256,
		DefaultDatetimePrecision:  nil,
		DisableDatetimePrecision:  false,
		DontSupportRenameIndex:    false,
		DontSupportRenameColumn:   false,
		DontSupportForShareClause: false,
	}), &gorm.Config{
		// QueryFields: true,
	})
	//db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	//db, err := gorm.Open(databaseSetting.DBType, )
	if err != nil {
		return nil, err
	}

	// if global.IsDevMode() {
	// 	db = db.Debug()
	// }
	//db = db.Debug()

	//db.SingularTable(true)
	//db.DB().SetMaxIdleConns(databaseSetting.MaxIdleConns)
	//db.DB().SetMaxOpenConns(databaseSetting.MaxOpenConns)
	return db, nil
}

func InitRedis(s *setting.DbSettingS) (*redis.Client, error) {
	redisSetting := s.Redis
	redisClient := redis.NewClient(&redis.Options{
		Addr:     redisSetting.AddressPort,
		Password: redisSetting.Password,
		DB:       redisSetting.DefaultDB,
	})

	timeoutCtx, cancelFunc := context.WithTimeout(context.Background(), redisSetting.DialTimeout)
	defer cancelFunc()

	_, err := redisClient.Ping(timeoutCtx).Result()
	return redisClient, err
}

func DbAutoMigrate(db *gorm.DB) {
	db.AutoMigrate(
		&AssetToken{}, &AssetTokenStats{},
	)
}
