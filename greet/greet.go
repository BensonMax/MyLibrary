package main

import (
	"flag"
	"fmt"
	"net/http"

	"greet/internal/config"
	"greet/internal/handler"
	"greet/internal/svc"

	"github.com/tal-tech/go-zero/core/conf"
	"github.com/tal-tech/go-zero/core/logx"
	"github.com/tal-tech/go-zero/rest"
	"github.com/tal-tech/go-zero/rest/httpx"
)

var configFile = flag.String("f", "etc/greet-api.yaml", "the config file")

// 返回的结构体，json格式的body
type Message struct {
	Code int    `json:"code"`
	Desc string `json:"desc"`
}

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)

	ctx := svc.NewServiceContext(c)
	server := rest.MustNewServer(c.RestConf)
	defer server.Stop()

	server.Use(
		middlewareDemoFunc,
	)
	// 设置错误处理函数
	httpx.SetErrorHandler(errorHandler)
	handler.RegisterHandlers(server, ctx)

	fmt.Printf("Starting server at %s:%d...\n", c.Host, c.Port)
	server.Start()
}

// 定义错误处理函数
func errorHandler(err error) (int, interface{}) {
	return http.StatusConflict, Message{
		Code: -1,
		Desc: err.Error(),
	}
}

func middlewareDemoFunc(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logx.Info("request ... ")
		next(w, r)
		logx.Info("reponse ... ")
	}
}
