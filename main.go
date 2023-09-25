/*
 * @Author: Vincent Young
 * @Date: 2023-09-25 11:56:51
 * @LastEditors: Vincent Young
 * @LastEditTime: 2023-09-25 12:00:07
 * @FilePath: /simpleproxy/main.go
 * @Telegram: https://t.me/missuo
 * @GitHub: https://github.com/missuo
 * 
 * Copyright © 2023 by Vincent, All Rights Reserved. 
 */

package main

import (
	"flag"
	"fmt"
	"github.com/elazarl/goproxy"
	"log"
	"net/http"
)

func main() {
	// Define command-line flags
	usernamePtr := flag.String("username", "", "Username for the proxy")
	passwordPtr := flag.String("password", "", "Password for the proxy")
	portPtr := flag.Int("port", 0, "Port to listen on")

	// Parse command-line flags
	flag.Parse()

	// Validate command-line inputs
	if *usernamePtr == "" || *passwordPtr == "" || *portPtr == 0 {
		log.Fatal("Username, password, and port must be provided")
	}

	proxy := goproxy.NewProxyHttpServer()

	proxy.OnRequest().HandleConnect(goproxy.AlwaysMitm)
	proxy.OnRequest().DoFunc(func(r *http.Request, ctx *goproxy.ProxyCtx) (*http.Request, *http.Response) {
		user, pass, ok := r.BasicAuth()
		if !ok || user != *usernamePtr || pass != *passwordPtr {
			return r, goproxy.NewResponse(r,
				goproxy.ContentTypeText, http.StatusUnauthorized,
				"Unauthorized")
		}
		return r, nil
	})

	// Start the proxy server
	log.Printf("Starting proxy server on :%d", *portPtr)
	err := http.ListenAndServe(fmt.Sprintf(":%d", *portPtr), proxy)
	if err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
