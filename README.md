# XYHELPER- ACCOUNT-CHECK

用于批量检测账号，默认输出账号 refresh_token 和 access_token 

## 使用方法

首先 `go mod tidy` 安装依赖，然后

修改 `config.yaml` 文件，填入你的 `CHATPROXY`.

在程序目录下创建 `account.txt` 文件，每行一个账号，格式为 `username,password`，例如：

```
username1,password1
username2,password2
```

然后运行 `go run main.go` 即可。

## 细节描述

结果将输出到 `output.txt` 文件中，格式为 `username,password,refresh_token,access_token`，例如：

```
username1,password1,refresh_token1,access_token1
username2,password2,refresh_token2,access_token2
```

重复运行时,会跳过`output.txt`中已经存在的账号。
