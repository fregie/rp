# rp
极为简单的https反响代理服务器,适合于在有证书的情况下快速部署一个指定端口的https反代

## 用法
```bash
docker run -d \
--name rp-service \
--restart=always \
-p 4443:443 \
-e SOURCE="http://source-ip:port" \
-e CRT="-----BEGIN CERTIFICATE-----
...
-----END CERTIFICATE-----
" \
-e KEY="-----BEGIN RSA PRIVATE KEY-----
...
-----END RSA PRIVATE KEY-----
" \
fregie/rp
```