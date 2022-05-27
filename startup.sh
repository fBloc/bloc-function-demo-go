#!/bin/bash
# ----------- download relate fiels -----------
# make sure docker-compose.yml exist

filePathPrefix="https://raw.githubusercontent.com/fBloc/bloc/main/"

infraFile="docker-compose.yml"
shutdownShFile="shutdown.sh"
frontComposeFile="docker-compose-bloc-web.yml"

if [[ "$OSTYPE" == "darwin"* ]]; then
	# Mac OSX
	app_name="com.docke"
	serverComposeFile="docker-compose-bloc-server-mac.yml"
elif [[ "$OSTYPE" == "linux"* ]]; then
	# Linux
	app_name="dockerd"
	serverComposeFile="docker-compose-bloc-server-linux.yml"
else
	# tmp to just use linux. later should support windows
	echo "your os $OSTYPE maybe not supported, use linux as default!"
	exit 8
fi

necessaryFiles=($infraFile $serverComposeFile $frontComposeFile $shutdownShFile)
used_ports=(8083 8080 27017 5672 15672 9000 8086)


check_env() {
	cd $(dirname $0); pwd

	# check docker status
	# lsof -UnP |grep $app_name
	# if [ $? -eq 0 ]; then
	# 	echo "docker is ready"
	# else
	# 	echo "docker not start"
	# 	echo "you need start docker first"
	# 	exit 8
	# fi

	# check port be occupied
	for element in ${used_ports[@]}
	do
		Pid=`lsof -i:$element |grep 'LISTEN' | awk '{print $1 "  " $2}'`
		if [[ -z "$Pid" ]]
		then
			echo "bloc will use port $element"
		else
			echo "Fail: bloc is to use port $element but its in use - '$Pid'"
			exit 8
		fi
	done
}

check_conf_file() {
	for file in ${necessaryFiles[@]}
	do
		if [ -s "$file" ];then
			echo "necessary file $file ready"
			continue
		fi
			
		echo "necessary file $file missing, Try to download"
		wget -q $filePathPrefix$file -O $file

		if [ ! -s "$file" ];then
			echo "necessary file $file download file"
			exit 8
		fi

		echo "necessary file $file ready"
	done

}

check_components() {
	docker-compose up -d

	try_amount=3
	# check infra components all ready
	echo "checking whether influxDB is ready"
	while :
	do
		RESULT=$(curl -s --location --request GET 'http://localhost:8086/api/v2/setup')
		if [[ $RESULT == *"allowed"* ]]
		then
			echo "    ready!"
			break
		else
			echo "    Checking influxDB $try_amount"
			try_amount=$((try_amount - 1))
			if [[ $try_amount < 0 ]]
			then
				break
			fi
		fi
		sleep 1
	done

	echo "start check whether minio is ready"
	try_amount=3
	while :
	do
		RESULT=$(curl -s -o /dev/null -I -w "%{http_code}" 'http://localhost:9000/minio/health/live')
		if [[ $RESULT == "200" ]]
		then
			echo "    ready!"
			break
		else
			echo "    Checking minio $try_amount"
			try_amount=$((try_amount - 1))
			if [[ $try_amount < 0 ]]
			then
				break
			fi
		fi
			sleep 1
	done

	echo "start check whether rabbitMQ is ready"
	try_amount=3
	while :
	do
		RESULT=$(curl -s -o /dev/null -I -w "%{http_code}" 'http://localhost:15672/api/overview')
		if [[ $RESULT == *"401"* ]]
		then
			echo "    ready!"
			break
		else
			echo "    Checking rabbitMQ $try_amount"
			try_amount=$((try_amount - 1))
			if [[ $try_amount < 0 ]]
			then
				break
			fi
		fi
			sleep 1
	done

	echo "start check whether mongoDB is ready"
	try_amount=3
	while :
	do
		Pid=`lsof -i:27017 | awk '{print $1 "  " $2}'`
		if [[ $RESULT == "" ]]
		then
			echo "    Checking mongoDB $try_amount"
			try_amount=$((try_amount - 1))
			if [[ $try_amount < 0 ]]
			then
				break
			fi
		else
			echo "    ready!"
			break
		fi
			sleep 1
	done

	echo -e "\n↓↓↓↓↓↓ bloc-infra status ↓↓↓↓↓↓"
	docker-compose ps
	echo -e "↑↑↑↑↑↑ bloc-infra status ↑↑↑↑↑↑\n"
}

start_web() {
	# start bloc-server
	echo "Starting bloc-server"

	try_amount=3
	while [ try_amount > 0 ]
	do
		docker-compose -f "$serverComposeFile" up -d
		sleep 3
		RESULT=$(curl -s --location --request GET 'http://localhost:8080/api/v1/bloc')
		if [[ $RESULT == *"Welcome aboard!"* ]]
		then
			echo "    bloc-server is up"
			break
		else
			echo "    Checking bloc-server $try_amount"
			try_amount=$((try_amount - 1))
			if [[ $try_amount < 0 ]]
			then
				break
			fi
		fi
	done

	echo -e "\n↓↓↓↓↓↓ bloc-server status ↓↓↓↓↓↓"
	docker-compose -f "$serverComposeFile" ps
	echo -e "↑↑↑↑↑↑ bloc-server status ↑↑↑↑↑↑\n"

	# start bloc-web
	echo "Starting bloc_web, yaml file: $frontComposeFile"
	docker-compose -f "$frontComposeFile" up -d
	# server_status=`docker-compose -f "$frontComposeFile" ps | grep bloc_web`
	# if [[ $server_status == *"Up"* ]] || [[ $server_status == *"running"* ]]
	# then
	# 	echo "    bloc_web is up"
	# fi
	sleep 3

	echo -e "\n↓↓↓↓↓↓ bloc-frontend status ↓↓↓↓↓↓"
	docker-compose -f "$frontComposeFile" ps
	echo -e "↑↑↑↑↑↑ bloc-frontend status ↑↑↑↑↑↑\n"
}

check_env
check_conf_file
check_components

start_web

# Guide users to access the front-end address
echo "******************************"
echo "login user: bloc"
echo "login password: bloc"
echo "******************************"
