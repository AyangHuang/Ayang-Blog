package main

import (
	"autodeploy/deploy"
	"net/http"
)

func main() {
	http.HandleFunc("/autodeploy", deploy.AutoDeploy)
	_ = http.ListenAndServe(":1314", nil)
}
