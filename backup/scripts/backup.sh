#!/bin/bash --login

validate_software() {
	echo "VALIDATE MANDATORY TOOLS"

	INSTALLED_PG_DUMP=`which pg_dump`
	if [ -z "$INSTALLED_PG_DUMP" ]; then
		echo "pg_dump utility not installed"
		exit 1
	fi
}

export_Encryption_key() {
  export BACKUP_DIR=$2
  export DEPLOYMENT_DIR=$3

	echo "EXPORT DB ENCRYPTION KEY"
	grep -E 'db_encryption_key' $DEPLOYMENT_DIR/cf-*.yml | cut -d ':' -f 2 | sort -u | tr -d ' ' > $BACKUP_DIR/cc_db_encryption_key.txt
}

export_db() {
	IP=$2
	USERNAME=$3
	export PGPASSWORD=$4
	PORT=$5
	DB=$6
	DB_FILE=$7

	echo "EXPORT $DB"

	pg_dump -h $IP -U $USERNAME -p $PORT $DB > $DB_FILE

}

export_mysqldb() {
	IP=$2
	USERNAME=$3
	PASSWORD=$4
	DB_FILE=$5

	echo '[mysqldump]
user='$USERNAME'
password='$PASSWORD > ~/.my.cnf

	echo "EXPORT MySQL DB"

	mysqldump -u $USERNAME -h $IP --all-databases > $DB_FILE

}

toggle_cc_job() {
	SERVER_URL=$2
	USERNAME=$3
	PASSWORD=$4
	output=`curl -v -XPUT -u "$USERNAME:$PASSWORD" $SERVER_URL --insecure -H "Content-Type:text/yaml" -i -s | grep Location: | grep Location: | cut -d ' ' -f 2`
	echo $output
}

export_nfs_server() {
	NFS_IP=$2
	NFS_SERVER_USER=$3
	NFS_SERVER_PASSWORD=$4
	NFS_DIR=$5
	echo "EXPORT NFS-SERVER"

	/usr/bin/expect -c "
		set timeout -1

		spawn scp -rp $NFS_SERVER_USER@$NFS_IP:/var/vcap/store/shared $NFS_DIR

		expect {
			-re ".*Are.*.*yes.*no.*" {
				send yes\r ;
				exp_continue
			}

			"*?assword:*" {
				send $NFS_SERVER_PASSWORD\r
			}
		}
		expect {
			"*?assword:*" {
				send $NFS_SERVER_PASSWORD\r
				interact
			}
		}

		exit
	"
}

usage() {
	echo "Usage: cfops install backup <OPS MGR HOST or IP> <SSH PASSWORD> <OPS MGR ADMIN USER> <OPS MGR ADMIN PASSWORD> <OUTPUT DIR>"
	printf "\t %s \t\t\t %s \n" "OPS MGR HOST or IP:" "OPS Manager Host or IP"
	printf "\t %s \t\t\t\t %s \n" "SSH PASSWORD:" "OPS Manager Tempest SSH Password"
	printf "\t %s \t\t\t %s \n" "OPS MGR ADMIN USER:" "OPS Manager Admin Username"
	printf "\t %s \t\t %s \n" "OPS MGR ADMIN PASSWORD:" "OPS Manager Admin Password"
	printf "\t %s \t\t\t\t %s \n" "OUTPUT DIR:" "Backup Directory"
}

$1 "$@"
