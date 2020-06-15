package main

import (
	"context"
	"flag"
	"github.com/prometheus/common/log"
	"github.com/sunnywalden/sync-data/pkg/auth"
	"github.com/sunnywalden/sync-data/pkg/logging"
	_ "net/http/pprof"
	"os"

	"github.com/gin-gonic/gin"
	ginprometheus "github.com/zsais/go-gin-prometheus"

	"github.com/sunnywalden/sync-data/apis"
	"github.com/sunnywalden/sync-data/config"
	//"github.com/sunnywalden/sync-data/pkg/auth"
	//"github.com/sunnywalden/sync-data/pkg/logging"
)


var (
	AppVersion string
	version *bool
	confPath string
	port string
	host string

	configures *config.TomlConfig
	ctx context.Context

)

func init() {
	version = flag.Bool("v", false, "show version and exit")
}

func main () {

	// get configures
	_, err := config.Init()
	if err != nil {
		panic(err)
	}

	configures = config.Conf

	if *version {
		log.Info(configures.App.Version)
		os.Exit(0)
	}

	// getting service start configure
	cmdLine := flag.NewFlagSet("sync-data", flag.PanicOnError)
	cmdLine.StringVar(&port, "p", "8090", "端口号，默认为800")
	cmdLine.StringVar(&host, "h", "127.0.0.1", "主机名，默认127.0.0.1")
	err = cmdLine.Parse(os.Args[1:])
	if err != nil {
		panic(err)
	}

	// return usage
	flag.Usage = func() {
		log.Infof("Usage of %s:\n", "http base")
		flag.PrintDefaults()
	}

	router := gin.New()

	p := ginprometheus.NewPrometheus("gin")
	p.Use(router)

	logLevel := configures.Log.Level
	router.Use(logging.Logger(logLevel))

	router.Use(gin.Recovery())

	user := router.Group("api/user")
	{
		user.Use(auth.TokenMiddleware())
		user.GET("/list", apis.UserList)
		user.GET("/", apis.User)
	}

	plat := router.Group("api/plat")
	{
		plat.POST("/register", apis.Register)
		plat.POST("/token", apis.GetToken)
	}

	addr := host + ":" + port
	err = router.Run(addr)
	if err != nil {
		log.Errorf("ERROR: ", err)
		panic(err)
	}

}