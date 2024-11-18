package rpc

import "github.com/lesismal/arpc"

func StartRPC(addr string) {
	srv := arpc.NewServer()
	srv.Handler.Handle("/", func(ctx *arpc.Context) {})

	srv.Run(addr)
}
