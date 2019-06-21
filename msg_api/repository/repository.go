package repository

import (
	"github.com/imkuqin-zw/final_consistency/basic/config"
	"github.com/imkuqin-zw/final_consistency/msg_api/models"
	"github.com/imkuqin-zw/final_consistency/msg_api/repository/mysql"
	z "github.com/imkuqin-zw/final_consistency/plugins/zap"
	"go.uber.org/zap"
	"sync"
)

var (
	repo Repository
	log  *z.Logger
	once sync.Once
)

const (
	_ = iota
	MYSQL_DB
)

type RepoConf struct {
	DbType int8 `json:"db_type"`
}

type TransactionMsg interface {
	//获取事务
	GetTransMsgByMsgId(string) (*models.TransactionMsg, error)

	//更新事务状态
	UpdateTransMsgStatusByMsgId(*models.TransactionMsg, uint8) (int64, error)

	//删除事务
	DeleteTransMsgByMsgId(string) (int64, error)

	//创建事务
	InsertTransMsg(*models.TransactionMsg) error

	//更新事务发送次数
	UpdateTransMsgSendTimesByMsgId(string) (int64, error)

	//将事务设置为死亡状态
	SetTransMsgAlreadyDeadByMsgId(string) (int64, error)
}

type Repository interface {
	TransactionMsg
}

func Init() {
	once.Do(func() {
		log = z.GetLogger()
		InitRepository()
	})
}

func InitRepository() {
	c := config.C()
	var cfg RepoConf
	if err := c.Path("repo", &cfg); err != nil {
		log.Panic("get app db_type config fault", zap.Error(err))
	}
	switch cfg.DbType {
	case MYSQL_DB:
		mysql.Init()
		repo = mysql.NewRepo()
	default:
		log.Panic("unknown db_type", zap.Any("repo_cfg", cfg))
	}
}

func GetRepo() Repository {
	return repo
}
