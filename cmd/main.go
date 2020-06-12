package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	_ "net/http/pprof"
	"os"

	jsoniter "github.com/json-iterator/go"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/sunnywalden/sync-data/config"
	"github.com/sunnywalden/sync-data/pkg/logging"
	"github.com/sunnywalden/sync-data/pkg/sync"
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

func UserList(w http.ResponseWriter, r *http.Request) {

	err := r.ParseForm()
	if err != nil {
		fmt.Errorf("%s", err.Error())
	}

	ctx, cancel := context.WithCancel(context.Background())

	defer cancel()

	users, err := sync.GetUser(ctx, configures)
	if err != nil {
		fmt.Errorf("%s", err.Error())
	} else {
		var json = jsoniter.ConfigCompatibleWithStandardLibrary
			jsonUsers, encodeErr := json.Marshal(&users)
			if encodeErr != nil {
				fmt.Errorf("%s", encodeErr)
			} else {
				fmt.Fprintf(w, "%s", jsonUsers)
			}
	}

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

	cmdLine := flag.NewFlagSet("sync-data", flag.PanicOnError)
	cmdLine.StringVar(&port, "p", "8088", "端口号，默认为8088")
	cmdLine.StringVar(&host, "h", "127.0.0.1", "主机名，默认127.0.0.1")
	cmdLine.Parse(os.Args[1:])

	flag.Usage = func() {
		log.Infof("Usage of %s:\n", "http base")
		flag.PrintDefaults()
	}

	http.HandleFunc("/", UserList)
	http.Handle("/metrics", promhttp.Handler())
	addr := host + ":" + port
	log.Printf(addr)
	err = http.ListenAndServe(addr, nil)
	if err != nil {
		log.Errorf("ERROR: ", err)
	}

}