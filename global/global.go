package global

import (
	"nftserver/pkg/logger"

	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
)

var DBEngine *gorm.DB
var RedisCli *redis.Client

var (
	AppSetting        *setting.AppSettingS
	DbSetting         *setting.DbSettingS
	StorageSetting    *setting.StorageSettingS
	Web3Setting       *setting.Web3SettingS
	BlockChainSetting *setting.BlockChainSettingS
	ControlorSetting  *setting.ControlorSettings
	MockSetting       *setting.MockSettingS
)

var Logger *logger.Logger
