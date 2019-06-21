package main

import (
	"github.com/imkuqin-zw/final_consistency/basic"
	"github.com/imkuqin-zw/final_consistency/basic/common"
	"github.com/imkuqin-zw/final_consistency/basic/config"
	"github.com/imkuqin-zw/final_consistency/msg_api/handler"
	"github.com/imkuqin-zw/final_consistency/msg_api/msg_queue"
	pb "github.com/imkuqin-zw/final_consistency/msg_api/proto/transaction"
	"github.com/imkuqin-zw/final_consistency/msg_api/repository"
	"github.com/imkuqin-zw/final_consistency/msg_api/service"
	z "github.com/imkuqin-zw/final_consistency/plugins/zap"
	"github.com/micro/cli"
	"github.com/micro/go-config/source/etcd"
	"github.com/micro/go-grpc"
	"github.com/micro/go-micro"
	"github.com/micro/go-micro/registry"
	"github.com/micro/go-plugins/registry/etcdv3"
	"go.uber.org/zap"
	"strings"
	"time"
)

var (
	appName  string
	etcdAddr string
	cfg      = &msgApiCfg{}
	log      = z.GetLogger()
)

type msgApiCfg struct {
	common.AppCfg
	MsgVersion uint32 `json:"msg_version"`
}

func main() {
	svr := grpc.NewService(
		micro.Flags(
			cli.StringFlag{
				Name:        "cfg_name",
				Usage:       "service config name",
				Value:       "msg-api-srv",
				Destination: &appName,
			},
			cli.StringFlag{
				Name:        "cfg_addr",
				Usage:       "config etcd address",
				Value:       "192.168.2.118:2379",
				Destination: &etcdAddr,
			},
		),
	)
	// Initialise Cmd
	svr.Init()

	// Initialise config
	initCfg()

	micReg := etcdv3.NewRegistry(registryOptions)
	svr.Init(
		micro.Name(cfg.Name),
		micro.Version(cfg.Version),
		micro.Registry(micReg),
		micro.RegisterInterval(cfg.RegInterval),
		micro.RegisterTTL(cfg.RegTTL),
		micro.Action(func(context *cli.Context) {
			repository.Init()
			msg_queue.Init()
			service.Init(cfg.MsgVersion)
			handler.Init()
		}),
	)

	//Register Handler
	if err := pb.RegisterTransactionHandler(svr.Server(), new(handler.TransactionMsgApi)); err != nil {
		log.Panic("register transaction handler", zap.Error(err))
	}

	// Run service
	if err := svr.Run(); err != nil {
		log.Fatal("service fault", zap.Error(err))
	}
}

func registryOptions(opts *registry.Options) {
	etcdCfg := &common.EtcdCfg{}
	err := config.C().App("etcd", etcdCfg)
	if err != nil {
		log.Panic("get etcd config fault", zap.Error(err))
	}
	opts.Timeout = time.Second * 5
	opts.Addrs = etcdCfg.Addrs
}

func initCfg() {
	source := etcd.NewSource(
		etcd.WithAddress(strings.Split(etcdAddr, ",")...),
		etcd.WithPrefix("msg_api"),
	)
	basic.Init(
		config.WithSource(source),
		config.WithApp(appName),
	)
	log.Info("[initCfg] init config completed")
	initAppCfg()
	return
}

func initAppCfg() {
	err := config.C().Path("app", cfg)
	if err != nil {
		log.Panic("get app config fault", zap.Error(err))
	}
	if cfg.RegTTL <= 0 {
		cfg.RegTTL = time.Second * 15
	}
	if cfg.RegInterval > cfg.RegTTL {
		cfg.RegInterval = cfg.RegTTL - 5
		if cfg.RegInterval <= 0 {
			cfg.RegTTL = time.Second * 15
			cfg.RegInterval = time.Second * 10
		}
	}
	if cfg.MsgVersion == 0 {
		cfg.MsgVersion = 1
	}
	return
}
