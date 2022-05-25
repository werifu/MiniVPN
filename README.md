# MiniVPN
VPN Lab of seedlab (Go.ver)

## 使用指南
__请先启动服务器再启动客户端__
### 服务端配置
1. 将example.config.json改名为config.json，放在server文件同目录下
2. 准备服务器私钥与证书（证书需可信CA签名）
3. 配置config.json文件
4. ./server启动
### 客户端配置
1. 将example.config.json改名为config.json，放在client文件同目录下
2. 准备客户端私钥与证书（证书需可信CA签名）
3. 配置config.json文件
4. ./client启动

### ⚠️注意事项
如果使用自己的自签名CA，需在client机器上信任自己的CA证书
```shell
# ubuntu
sudo cp 你的ca.crt /usr/local/share/ca-certificates/
sudo update-ca-certificates
```
