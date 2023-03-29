# teslamate-addr-fix

[English](README.md)

本项目解决因 openstreetmap 被封禁而导致 teslamate 地址丢失的问题

## 注意

**在使用之前必须进行 [备份](https://docs.teslamate.org/docs/maintenance/backup_restore).**


## 前提条件
- 您也遇到了 teslamate [地址丢失的问题](https://github.com/adriankumpf/teslamate/issues/2956)

- 您有 **HTTP 代理** 可访问 openstreetmap.org

## 演示

修复前: (没有 start/destination 信息, 因 openstreetmap 被封禁)
![Before](demo/before.jpg)

修复后:
![After](demo/after.jpg)

## 步骤

- 将 teslamate postgres 数据库端口暴露给主机. 如果你是用 docker compose 安装的, 直接在 .yml 文件中填充端口即可, 然后执行 `docker-compose up -d` 重建 docker.

	```
	database:
		image: postgres:14
		restart: always
		environment:
		- POSTGRES_USER=teslamate
		- POSTGRES_PASSWORD=xxxxxxxx
		- POSTGRES_DB=teslamate

		# 加入这两行
		ports:
		- 5432:5432 

	```

- 配置并启动 HTTP 代理. 您可以提前配置系统代理, 或者在运行时再指定代理. 系统代理的设置类似这样:

	```
	# Your .bashrc/.zshrc

	export all_proxy=socks5://127.0.0.1:7890
	export http_proxy=http://127.0.0.1:7890
	export https_proxy=http://127.0.0.1:7890
	```

- 执行帮助指令 `./teslamate-addr-fix -h`. 您至少需要指定 teslamate 的 postgres 数据库密码, 否则程序无法连接到 teslamate 的数据库. 其他参数如果和默认值不同的话也需要指定.

	```
	Usage of ./teslamate-addr-fix:
	-db string
			teslamate psql database (default "teslamate")
	-host string
			teslamate psql host (default "127.0.0.1")
	-interval int
        	interval (minutes) for running in daemon mode
	-password string
			teslamate psql password
	-port string
			teslamate psql port (default "5432")
	-proxy string
			http proxy (default use system proxy)
	-timeout int
			timeout of openstreetmap request (default 5)
	-user string
			teslamate psql user (default "teslamate")
	```

- 根据指示填充参数然后开始修复. 指定参数执行程序, 比如 `./teslamate-addr-fix -password 123456`, 然后日志将会输出在 `teslamate-addr-fix.log`

- 当程序执行完成后, 查看 teslamate grafana 的 drive graph 检查是否修复. 

## 免责声明

仅当完整创建了备份后才能使用本程序, 本人**不对因使用本程序造成的任何数据丢失或程序错误负责**.

本项目仅供学习交流使用, **不提供任何网络代理及其下载链接**. 如违反当地法律法规使用网络代理的, 造成的后果由使用者负责.

当您下载, 复制, 编译或运行本项目的源代码或二进制程序时, 即代表您同意上述免责声明.