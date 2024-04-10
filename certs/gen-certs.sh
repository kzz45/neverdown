#!/usr/bin/env bash

# Enter PEM pass phrase: neverdown
# Verifying - Enter PEM pass phrase: neverdown
# Enter pass phrase for ca.key: neverdown
# Certificate request self-signature ok
# subject=C = CN, ST = Shanghai, L = XuHui, O = NeverDown, CN = NeverDown root CA
# Enter pass phrase for ca.key: neverdown
# writing RSA key

# 创建ssl证书私钥
openssl genrsa -des3 -out ca.key 2048

# 生成CA证书 10 years
openssl req -x509 \
  -new \
  -nodes \
  -key ca.key \
  -days 3650 \
  -out ca.crt \
  -sha256 \
  -subj "/C=CN/ST=Shanghai/L=XuHui/O=NeverDown/CN=NeverDown root CA"

openssl genrsa -out server.key 2048

# 不使用CA 创建ssl证书CSR
openssl req -new -key server.key -out server.csr \
  -subj "/C=CN/ST=Shanghai/L=XuHui/O=NeverDown/CN=NeverDown root CA"

# 使用CA签署ssl证书 ssl证书有效期10年
openssl x509 -req -in server.csr -out server.crt -days 3650 \
  -CAcreateserial -CA ca.crt -CAkey ca.key \
  -CAserial serial -extfile cert.ext

openssl x509 -in server.crt -out cert.pem
openssl rsa -in server.key -text -out key.pem
