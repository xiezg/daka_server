module daka

go 1.15

require (
	github.com/gomarkdown/markdown v0.0.0-20210208175418-bda154fe17d8
	github.com/gorilla/mux v1.8.0
	github.com/tencentcloud/tencentcloud-sdk-go v1.0.104
	github.com/xiezg/glog v0.0.0-20200512024805-77553a7bf27a
	github.com/xiezg/go-jsonify v0.0.0-20200411230712-1b0b52358430
	github.com/xiezg/muggle v0.0.0-20200808060153-865c4e13ebbc
	golang.org/x/lint v0.0.0-20201208152925-83fdc39ff7b5 // indirect
	golang.org/x/tools v0.1.0 // indirect
)

replace github.com/xiezg/muggle => ../muggle

replace github.com/xiezg/glog => ../glog
