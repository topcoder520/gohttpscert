package main

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"path/filepath"
)

func init() {
}

//服务端
func main() {

	/* caCrtPath := filepath.Clean("E:/https/X509/server/ca/ca.crt")      //CA证书
	serverCrtPath := filepath.Clean("E:/https/X509/server/server.crt") //服务端证书
	serverKeyPath := filepath.Clean("E:/https/X509/server/server.key") //服务端密钥 */
	caCrtPath := filepath.Clean("../cert/ca.pem")         //CA证书
	serverCrtPath := filepath.Clean("../cert/server.pem") //服务端证书
	serverKeyPath := filepath.Clean("../cert/server.key") //服务端密钥

	mux := http.NewServeMux()
	mux.HandleFunc("/api/test/hlw", func(rw http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(rw, "hello wolrd")
	})
	//验证客户端证书配置

	caCrt, err := ioutil.ReadFile(caCrtPath)
	if err != nil {
		log.Println("caCertPath ReadFile err :", err)
		return
	}
	pool := x509.NewCertPool()
	pool.AppendCertsFromPEM(caCrt)
	server := &http.Server{
		Addr:    ":8089",
		Handler: middleWare(mux),
		TLSConfig: &tls.Config{
			ClientAuth: tls.RequireAndVerifyClientCert, //强制校验client端证书
			ClientCAs:  pool,                           //ca池 验证客户端证书的ca
		},
	}
	log.Println("ListenAndServeTLS start err: ", server.ListenAndServeTLS(serverCrtPath, serverKeyPath))
}

func middleWare(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		//中間件處理
		//TODO

		next.ServeHTTP(rw, r)
	})
}
