# gohttpscert
#https双向验证 自签名证书生成
首先需要下载 OpenSSL [http://slproweb.com/products/Win32OpenSSL.html](http://slproweb.com/products/Win32OpenSSL.html) 
## 第一种方法：GO1.15版本以下证书生成
**go 1.15 版本开始废弃 CommonName**   
**前提：需要设置环境变量 GODEBUG 为 x509ignoreCN=0**   
### 1、建立我们自己的CA

需要生成一个CA私钥和一个CA的数字证书:

	openssl genrsa -out ca.key 2048
	openssl req -x509 -new -nodes -key ca.key -subj "/CN=localhost" -days 5000 -out ca.crt

如果报以下错误：
	
	Subject does not start with '/'.
	problems making Certificate Request
该错误是由Git for Windows中MinGW/MSYS模块的路径转换机制引起的。
解决方案：将-subj参数中第一个“/”改为“//”，其余“/”改为“\”，如下：

	"//CN=localhost"

###2、接下来，生成server端的私钥，生成数字证书请求，并用我们的ca私钥签发server的数字证书

	openssl genrsa -out server.key 2048
	openssl req -new -key server.key -subj "/CN=localhost" -out server.csr
	openssl x509 -req -in server.csr -CA ca.crt -CAkey ca.key -CAcreateserial -out server.crt -days 5000


现在我们的工作目录下有如下一些私钥和证书文件：
CA:
    私钥文件 ca.key
    数字证书 ca.crt

Server:
    私钥文件 server.key
    数字证书 server.crt

###3、生成客户端的私钥与证书

	openssl genrsa -out client.key 2048
	openssl req -new -key client.key -subj "/CN=localhost" -out client.csr
golang的tls需要校验ExtKeyUsage，所以要在生成client.crt时指定extKeyUsage   
创建文件client.ext  
client.ext内容：   
extendedKeyUsage=clientAuth

	openssl x509 -req -in client.csr -CA ca.crt -CAkey ca.key -CAcreateserial -extfile client.ext -out client.crt -days 5000





## 第二种方法：GO1.15版本及1.5版本以上的证书生成
因为 go 1.15 版本开始废弃 CommonName，因此推荐使用 SAN 证书。    
下面就介绍一下SAN证书生成    

###第1步：生成 CA 根证书

	openssl genrsa -out ca.key 2048
	openssl req -new -x509 -days 3650 -key ca.key -out ca.pem    
填写信息:  

	You are about to be asked to enter information that will be incorporated
	into your certificate request.
	What you are about to enter is what is called a Distinguished Name or a DN.
	There are quite a few fields but you can leave some blank
	For some fields there will be a default value,
	If you enter '.', the field will be left blank.
	-----
	Country Name (2 letter code) [AU]:cn
	State or Province Name (full name) [Some-State]:shanghai
	Locality Name (eg, city) []:shanghai
	Organization Name (eg, company) [Internet Widgits Pty Ltd]:custer
	Organizational Unit Name (eg, section) []:custer
	Common Name (e.g. server FQDN or YOUR name) []:localhost
	Email Address []:



### 第2步：用 openssl 生成 ca 和双方 SAN 证书。
准备默认 OpenSSL 配置文件于当前目录    
linux系统在 : /etc/pki/tls/openssl.cnf    
Mac系统在: /System/Library/OpenSSL/openssl.cnf    
Windows：安装目录下 openssl.cfg 比如 D:\Program Files\OpenSSL-Win64\bin\openssl.cfg   
   
1：拷贝配置文件到项目 然后修改

2：找到 [ CA_default ],打开 copy_extensions = copy

3：找到[ req ],打开 req_extensions = v3_req # The extensions to add to a certificate request

4：找到[ v3_req ],添加 subjectAltName = @alt_names

5：添加新的标签 [ alt_names ] , 和标签字段(建议在openssl.cfg文件最末尾添加在标签)

	[ alt_names ]
	DNS.1 = localhost
	DNS.2 = *.custer.fun
这里填入需要加入到 Subject Alternative Names 段落中的域名名称，可以写入多个。
### 第3步：生成服务端证书

	openssl genpkey -algorithm RSA -out server.key
 
	openssl req -new -nodes -key server.key -out server.csr -days 3650 -subj "/C=cn/OU=custer/O=custer/CN=localhost" -config ./openssl.cfg -extensions v3_req
 
	openssl x509 -req -days 3650 -in server.csr -out server.pem -CA ca.pem -CAkey ca.key -CAcreateserial -extfile ./openssl.cfg -extensions v3_req


server.csr是上面生成的证书请求文件。ca.pem/ca.key是CA证书文件和key，用来对server.csr进行签名认证。这两个文件在之前生成的。

###第4步：生成客户端证书

	
	openssl genpkey -algorithm RSA -out client.key
 
	openssl req -new -nodes -key client.key -out client.csr -days 3650 -subj "/C=cn/OU=custer/O=custer/CN=localhost" -config ./openssl.cfg -extensions v3_req
 
	openssl x509 -req -days 3650 -in client.csr -out client.pem -CA ca.pem -CAkey ca.key -CAcreateserial -extfile ./openssl.cfg -extensions v3_req

现在 Go 1.15 以上版本的 GRPC 通信，这样就完成了使用自签CA、Server、Client证书和双向认证


###最后

如果出现创建Server证书请求出现错误：

	2604:error:08064066:object identifier routines:OBJ_create:oid exists:crypto\objects\obj_dat.c:698: error in req

创建Client证书请求出现错误

	16524:error:08064066:object identifier routines:OBJ_create:oid exists:crypto\objects\obj_dat.c:698:error in req
 
**解决办法：关闭PowerShell 重新进入OpenSSL问题解决。**




文章参考：


[https://blog.csdn.net/ma_jiang/article/details/111950872](https://blog.csdn.net/ma_jiang/article/details/111950872)   
[https://studygolang.com/articles/9267](https://studygolang.com/articles/9267)



