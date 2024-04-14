ACTION=$1
SERVICE="./main"
SERVICE_PID=$(ps -aux | awk '$11=="'$SERVICE'" {print $2}')

if [ -z $ACTION ]
then
	echo "Command not found"
	exit;
fi

if ! [ -x $SERVICE ]
then
	echo "Service execute permission not granted"
	exit;
fi

if [ $ACTION = 'start' ]
then
	if ! [ -z $SERVICE_PID ]
	then
		echo "Service PID:{$SERVICE_PID}"
	else
		nohup $SERVICE > log/$(date +%F-%H-%M-%S).log 2>&1 </dev/null &
	fi
elif [ $ACTION = 'restart' ]
then
	if  ! [ -z $SERVICE_PID ]
	then
		kill -s KILL $SERVICE_PID
	fi
	nohup $SERVICE > log/$(date +%F-%H-%M-%S).log 2>&1 </dev/null &
elif [ $ACTION == 'stop' ]
then
	if  ! [ -z $SERVICE_PID ]
        then
                kill -s KILL $SERVICE_PID
	else
		echo "Service PID not found"
        fi

fi