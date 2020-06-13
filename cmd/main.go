package main

import (
	"context"
	"flag"
	_ "net/http/pprof"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/sunnywalden/sync-data/apis"
	"github.com/sunnywalden/sync-data/config"
	"github.com/sunnywalden/sync-data/pkg/logging"
)


var (
	AppVersion string
	version *bool
	confPath string
	port string
	host string

	configures *config.TomlConfig
	ctx context.Context

	log = logging.GetLogger()
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
	cmdLine.StringVar(&port, "p", "8088", "端口号，默认为8088")
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

	router := gin.Default()

	user := router.Group("user")
	{
		user.GET("/users", apis.UserList)
		user.GET("/user", apis.User)
	}

	plat := router.Group("plat")
	{
		plat.POST("/register", apis.Register)
		plat.POST("/token", apis.GetToken)
	}

	router.GET("/metrics", promhttp.Handler())

	addr := host + ":" + port
	err = router.Run(addr)
	if err != nil {
		log.Errorf("ERROR: ", err)
		panic(err)
	}

	// api route register
	//http.HandleFunc("/token", apis.GetToken)
	//http.HandleFunc("/register", apis.Register)
	//http.HandleFunc("/", apis.UserList)
	//http.HandleFunc("/user", apis.User)
	//http.Handle("/metrics", promhttp.Handler())
	//
	//// service start
	//addr := host + ":" + port
	//log.Printf(addr)
	//err = http.ListenAndServe(addr,nil)
	//if err != nil {
	//	log.Errorf("ERROR: ", err)
	//}

}