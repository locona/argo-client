package main

import (
	"github.com/k0kubun/pp"
	"github.com/locona/argo-client/pkg/argo"
)

func main() {
	cli, _ := argo.New("argo")
	// cli.Create()
	// wi, _ := cli.Watch()
	// for {
	//
	// select {
	// case x := <-wi.ResultChan():
	// pp.Println("##########", x)
	// default:
	// }
	// }
}
