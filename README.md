# teslamate-addr-fix
To fix teslamate broken addresses caused by openstreetmap unavailability

[中文说明](README_CN.md)

## Notice

**Must create a [backup](https://docs.teslamate.org/docs/maintenance/backup_restore) before doing this.**


## Pre-requisite
- You have teslamate [broken address issue](https://github.com/adriankumpf/teslamate/issues/2956)

- You have access to openstreetmap.org **via your HTTP proxy**


## Demo

Before: (no start/destination info since OSM blocked)
![Before](demo/before.jpg)

After fixed:
![After](demo/after.jpg)

## Step

- Expose your teslamate postgres port to your host. If you are using docker compose, just simply add port in the .yml file, then execute `docker-compose up -d` to recreate docker.

	```
	database:
		image: postgres:14
		restart: always
		environment:
		- POSTGRES_USER=teslamate
		- POSTGRES_PASSWORD=xxxxxxxx
		- POSTGRES_DB=teslamate

		# add this
		ports:
		- 5432:5432 

	```

- Configure and turn your HTTP proxy on. You can either set the system proxy beforehand or set the proxy at runtime. System proxy setting is like:

	```
	# Your .bashrc/.zshrc

	export all_proxy=socks5://127.0.0.1:7890
	export http_proxy=http://127.0.0.1:7890
	export https_proxy=http://127.0.0.1:7890
	```

- Run the help command by `./teslamate-addr-fix -h`. At least you should specify your teslamate psql password, otherwise it cannot connect to teslamate database. The other arguments should be specified if not same as default.

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

- Follow the instruction and start fixing. Run the program with arguments such as `./teslamate-addr-fix -password 123456` , and the log will be printed in `teslamate-addr-fix.log`

- After the program finish, check your teslamate grafana drive graph if anything correct. 

## Disclaimer

Only use this program after properly created backups, I am **not** responsible for any data loss or software failure related to this.

This project is only for study purpose, and **no web proxy (or its download link) provided**. If the network proxy is used in violation of local laws and regulations, the user is responsible for the consequences.

When you download, copy, compile or execute the source code or binary program of this project, it means that you have accepted the disclaimer as mentioned.