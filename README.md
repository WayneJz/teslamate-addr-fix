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

## Install Guide (Choose one only)

### 1. Docker Compose (Recommended)

- Ensure your HTTP proxy is set to **"allow LAN use"**, and find the **LAN IP:Port**. For example, a proxy LAN IP is set to `192.168.x.x`. Commonly **SHOULD NOT** be the localhost `127.x.x.x` or docker host `172.x.x.x`

- Modify teslamate docker compose file:

	```
	# Insert below 'database' section

	teslamate-addr-fix:
      image: waynejz/teslamate-addr-fix
      restart: always
      platform: linux/amd64
      environment:
	    - PROXY=http://192.168.0.100:7890    # Set your HTTP proxy
        - DATABASE_USER=teslamate
        - DATABASE_PASS=123456               # Copy from 'teslamate' section
        - DATABASE_NAME=teslamate
        - DATABASE_HOST=database
      depends_on:
        - database
	```

	If you have modified the other default values (such as DATABASE_HOST), then you should copy and replace them from 'teslamate' section as well.

- Then execute `docker-compose up -d` to recreate docker. This tool will run with teslamate in the same subnetwork. After serveral minutes, check your teslamate grafana drive graph if anything correct.

### 2. Binary Installation 

- Download this tool from [releases page](https://github.com/WayneJz/teslamate-addr-fix/releases). Ensure you download the right binary for your OS arch/version.

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

- After the program finished, check your teslamate grafana drive graph if anything correct.


### 3. Native Docker (Not recommended)

Of course you can use docker separately to run this tool, but since the tool has not join your teslamate subnetwork without docker compose, it could be hard and not recommended.

- To run as native docker, the preparation includes **"allow LAN use"** AND **"exposing teslamate postgres port to your host"**, see content above for details

- Pull the image and run in daemon mode:

	```
	docker pull waynejz/teslamate-addr-fix:latest

	docker run --name 'teslamate-addr-fix' --platform 'linux/amd64' \
	-e PROXY='http://192.168.0.100:7890' \ 
	-e DATABASE_USER='teslamate' \
	-e DATABASE_PASS='123456' \
	-e DATABASE_NAME='teslamate' \
	-e DATABASE_HOST='192.168.0.100' \
	-e DATABASE_PORT='5432' \
	-d waynejz/teslamate-addr-fix
	```

	Note both the `PROXY` and `DATABASE_HOST` should be LAN IPs. If encountered with syntax error, merge into one line and remove slashes and try again. 

- If run successfully, check your teslamate grafana drive graph if anything correct after serveral minutes.

## FAQ

- Q: My docker has problems when running, does the docker has other parameters to adjust?

	A: Extra parameters below can be set if necessary (DO NOT set if run properly):

	```
	- OSM_TIMEOUT=5          # timeout of openstreetmap request (default 5 seconds)
	- DATABASE_PORT=5432     # port of teslamate postgres (default 5432)
	- INTERVAL=5             # interval for running in daemon mode (default 5 minutes)
	```

	After adjust, a restart is required to take effect.

- Q: This program has no log output, is it running correctly?

	A: If no broken address data to fix, then the program will not output any logs. You can check your drive graph.

- Q: The addresses have been fixed, why I still cannot see the drive map?

	A: Drive map in Grafana is a frontend feature, so ensure your working computer (not NAS server) has access to openstreetmap, through web proxy if necessary.

## Disclaimer

Only use this program after properly created backups, I am **not** responsible for any data loss or software failure related to this.

This project is only for study purpose, and **no web proxy (or its download link) provided**. If the network proxy is used in violation of local laws and regulations, the user is responsible for the consequences.

When you download, copy, compile or execute the source code or binary program of this project, it means that you have accepted the disclaimer as mentioned.