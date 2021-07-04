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

//客户端
func main() {
	/* caCrtPath := filepath.Clean("E:/https/X509/client/ca/ca.crt")  //CA证书
	clientCrt := filepath.Clean("E:/https/X509/client/client.crt") //客户端证书
	clientKey := filepath.Clean("E:/https/X509/client/client.key") //客户端密钥 */
	caCrtPath := filepath.Clean("../cert/ca.pem")     //CA证书
	clientCrt := filepath.Clean("../cert/client.pem") //客户端证书
	clientKey := filepath.Clean("../cert/client.key") //客户端密钥

	//验证服务端证书配置
	caCrt, err := ioutil.ReadFile(caCrtPath)
	if err != nil {
		log.Println("caCertPath ReadFile err :", err)
		return
	}
	pool := x509.NewCertPool()
	pool.AppendCertsFromPEM(caCrt)
	//客户端证书与密钥
	cliCrt, err := tls.LoadX509KeyPair(clientCrt, clientKey)
	if err != nil {
		log.Println("Loadx509keypair err:", err)
		return
	}
	client := http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				RootCAs:      pool,
				Certificates: []tls.Certificate{cliCrt},
			},
		},
	}
	req, err := http.NewRequest(http.MethodGet, "https://localhost:8089/api/test/hlw", nil)
	if err != nil {
		log.Println("NewRequest Get err:", err)
		return
	}
	resp, err := client.Do(req)
	if err != nil {
		log.Println("Client Do Request err:", err)
		return
	}
	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("ReadAll resp.Body err:", err)
		return
	}
	fmt.Println(string(b))
}
