package repository

import (
	"github.com/imkuqin-zw/final_consistency/msg_api/models"
	"github.com/imkuqin-zw/final_consistency/msg_api/repository/mysql"
	"go.uber.org/zap"
	"shop/basic/config"
	z "shop/plugins/zap"
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
	GetTransMsgByMsgId(string) (*models.TransactionMsg, error)
	UpdateTransMsgStatusByMsgId(*models.TransactionMsg, uint8) (int64, error)
	InsertTransMsg(*models.TransactionMsg) error
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
