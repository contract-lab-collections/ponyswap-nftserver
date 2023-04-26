package main

import (
	"fmt"
	"log"
	"net/http"
	"nftserver/global/setting"
	"nftserver/pkg/logger"
	"os"
	"os/signal"
	"path"

	"runtime"
	"syscall"

	"github.com/gin-gonic/gin"
	"github.com/urfave/cli/v2"
	"github.com/zeromicro/go-zero/rest/router"
	"gopkg.in/natefinch/lumberjack.v2"
)

var Version = "v1.0.0"

func main() {
	if err := setupApp().Run(os.Args); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func setupApp() *cli.App {
	app := cli.NewApp()
	app.Usage = "api server"
	app.Action = startService
	app.Version = Version
	app.Flags = []cli.Flag{}
	app.Commands = []cli.Command{}
	app.Before = func(context *cli.Context) error {
		runtime.GOMAXPROCS(runtime.NumCPU())
		return nil
	}
	return app
}

func startService(ctx *cli.Context) {
	var err error
	err = handleInitSetting()
	if err != nil {
		log.Fatalf("startService err: %v", err)
		return
	}

	gin.SetMode(global.AppSetting.Server.RunMode)

	storageSet := global.StorageSetting
	os.MkdirAll("storage/logs", 075)
	os.MkdirAll(path.Join(storageSet.BasePath, storageSet.Images.StorePath), 075)
	os.MkdirAll(path.Join(storageSet.BasePath, storageSet.Resources.StorePath), 075)

	routerHandler := router.NewRouter()

	s := &http.Server{
		Addr:           ":" + global.AppSetting.Server.HttpPort,
		Handler:        routerHandler,
		ReadTimeout:    global.AppSetting.Server.ReadTimeout,
		WriteTimeout:   global.AppSetting.Server.WriteTimeout,
		MaxHeaderBytes: 1 << 20,
	}

	err = s.ListenAndServe()
	if err != nil {
		log.Fatalf("Server err: %v", err)
		return
	}

	signalHandle()
}

func handleInitSetting() error {
	err := setupSetting()
	if err != nil {
		log.Fatalf("init.setupSetting err: %v", err)
		return err
	}

	err = setupLogger()
	if err != nil {
		log.Fatalf("init.setupLogger err: %v", err)
		return err
	}

	err = setupRedis()
	if err != nil {
		log.Fatalf("init.setupRedis err: %v", err)
		return err
	}

	err = setupDBEngine()
	if err != nil {
		log.Fatalf("init.setupDBEngine err: %v", err)
		return err
	}

	return nil
}

func setupSetting() error {
	pSetting, err := setting.NewSetting()
	if err != nil {
		return err
	}

	err = pSetting.ReadSection("App", &global.AppSetting)
	if err != nil {
		return err
	}

	err = pSetting.ReadSection("Db", &global.DbSetting)
	if err != nil {
		return err
	}

	err = pSetting.ReadSection("Storage", &global.StorageSetting)
	if err != nil {
		return err
	}

	err = pSetting.ReadSection("Web3", &global.Web3Setting)
	if err != nil {
		return err
	}

	err = pSetting.ReadSection("BlockChain", &global.BlockChainSetting)
	if err != nil {
		return err
	}

	err = pSetting.ReadSection("Controlor", &global.ControlorSetting)
	if err != nil {
		return err
	}

	if global.IsDevMode() {
		err = pSetting.ReadSection("Mock", &global.MockSetting)
		if err != nil {
			return err
		}
	}

	return nil
}

func setupLogger() error {
	global.Logger = logger.NewLogger(&lumberjack.Logger{
		Filename:  global.AppSetting.Log.SavePath + "/" + global.AppSetting.Log.FileName + global.AppSetting.Log.FileExt,
		MaxSize:   500, // Size of a single file
		MaxAge:    10,  // Number of backup files
		LocalTime: true,
	}, "", log.LstdFlags).WithCaller(2)
	return nil
}

func setupRedis() error {
	var err error
	global.RedisCli, err = model.InitRedis(global.DbSetting)
	if err != nil {
		panic("redis init err! " + err.Error())
	}
	return err
}

func setupDBEngine() error {
	var err error
	global.DBEngine, err = model.NewMysqlDBEngine(global.DbSetting)
	if err != nil {
		panic("database init err! " + err.Error())
	}

	model.DbAutoMigrate(global.DBEngine)
	model.InitDbData(global.DBEngine)

	return nil
}

func signalHandle() {
	var (
		ch = make(chan os.Signal, 1)
	)
	signal.Notify(ch, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
	for {
		si := <-ch
		switch si {
		case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT:
			//deal with some connect
			err := global.RedisCli.Close()
			if err != nil {
				fmt.Println("get a signal exit redis err:", err.Error())
			}
			fmt.Println("get a signal: stop the nft server process", si.String())
			return
		case syscall.SIGHUP:
		default:
			return
		}
	}
}
