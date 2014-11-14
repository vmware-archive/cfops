#!/bin/bash --login

validate_software() {
	echo "VALIDATE MANDATORY TOOLS"

	INSTALLED_BOSH=`which bosh`
	if [ -z "$INSTALLED_BOSH" ]; then
		echo "BOSH CLI not installed"
		exit 1
	fi

	INSTALLED_PG_DUMP=`which pg_dump`
	if [ -z "$INSTALLED_PG_DUMP" ]; then
		echo "pg_dump utility not installed"
		exit 1
	fi

	INSTALLED_JAVA=`which java`
	if [ -z "$INSTALLED_JAVA" ]; then
		echo "Java JRE is missing"
		exit 1
	fi
}

copy_deployment_files() {

  export OPS_MANAGER_HOST=$1
  export OPS_MGR_SSH_PASSWORD=$2
  export DEPLOYMENT_DIR=$6

	echo "COPY DEPLOYMENT MANIFEST"

	/usr/bin/expect -c "
		set timeout -1

		spawn scp tempest@$OPS_MANAGER_HOST:/var/tempest/workspaces/default/deployments/*.yml $DEPLOYMENT_DIR

		expect {
			-re ".*Are.*.*yes.*no.*" {
				send yes\r ;
				exp_continue
			}

			"*?assword:*" {
				send $OPS_MGR_SSH_PASSWORD\r
			}
		}
		expect {
			"*?assword:*" {
				send $OPS_MGR_SSH_PASSWORD\r
				interact
			}
		}

		exit
	"

	echo "COPY MICRO-BOSH DEPLOYMENT MANIFEST"
	/usr/bin/expect -c "
		set timeout -1

		spawn scp tempest@$OPS_MANAGER_HOST:/var/tempest/workspaces/default/deployments/micro/*.yml $DEPLOYMENT_DIR

		expect {
			-re ".*Are.*.*yes.*no.*" {
				send yes\r ;
				exp_continue
			}

			"*?assword:*" {
				send $OPS_MGR_SSH_PASSWORD\r
			}
		}
		expect {
			"*?assword:*" {
				send $OPS_MGR_SSH_PASSWORD\r
				interact
			}
		}

		exit
	"
}

export_Encryption_key() {
  export BACKUP_DIR=$5
  export DEPLOYMENT_DIR=$6

	echo "EXPORT DB ENCRYPTION KEY"
	grep -E 'db_encryption_key' $DEPLOYMENT_DIR/cf-*.yml | cut -d ':' -f 2 | sort -u | tr -d ' ' > $BACKUP_DIR/cc_db_encryption_key.txt
}

export_installation_settings() {
  export OPS_MANAGER_HOST=$1
  export OPS_MGR_ADMIN_USERNAME=$3
  export OPS_MGR_ADMIN_PASSWORD=$4
  export BACKUP_DIR=$5

	CONNECTION_URL=https://$OPS_MANAGER_HOST/api/installation_settings

	echo "EXPORT INSTALLATION FILES FROM " $CONNECTION_URL

	curl "$CONNECTION_URL" -X GET -u $OPS_MGR_ADMIN_USERNAME:$OPS_MGR_ADMIN_PASSWORD --insecure -k -o $BACKUP_DIR/installation.json
}

fetch_bosh_connection_parameters() {
  export BACKUP_DIR=$5
	echo "GATHER BOSH DIRECTOR CONNECTION PARAMETERS"

	output=`sh appassembler/bin/app $BACKUP_DIR/installation.yml microbosh director director`

	export DIRECTOR_USERNAME=`echo $output | cut -d '|' -f 1`
	export DIRECTOR_PASSWORD=`echo $output | cut -d '|' -f 2`
	export BOSH_DIRECTOR_IP=`echo $output | cut -d '|' -f 3`

}

bosh_login() {
  echo "GATHER BOSH DIRECTOR CONNECTION PARAMETERS"

  export BOSH_DIRECTOR_IP=$1
  export DIRECTOR_USERNAME=$2
  export DIRECTOR_PASSWORD=$3

	echo "BOSH LOGIN $DIRECTOR_USERNAME $DIRECTOR_PASSWORD "
	rm -rf ~/.bosh_config

	bosh target $BOSH_DIRECTOR_IP << EOF
	$DIRECTOR_USERNAME
	$DIRECTOR_PASSWORD
EOF

	bosh login $DIRECTOR_USERNAME $DIRECTOR_PASSWORD
}

verify_deployment_backedUp() {
  export BACKUP_DIR=$5
	echo "VERIFY CF DEPLOYMENT MANIFEST"
	export CF_DEPLOYMENT_NAME=`bosh deployments | grep "cf-" | cut -d '|' -f 2 | tr -s ' ' | grep "cf-" | tr -d ' '`
	export CF_DEPLOYMENT_FILE_NAME=$CF_DEPLOYMENT_NAME.yml

	echo "FILES LOOKING FOR $CF_DEPLOYMENT_NAME $CF_DEPLOYMENT_FILE_NAME"

	if [ -f $BACKUP_DIR/$CF_DEPLOYMENT_FILE_NAME ]; then
		echo "file exists"
	else
		echo "file does not exist"
		bosh download manifest $CF_DEPLOYMENT_NAME $BACKUP_DIR/$CF_DEPLOYMENT_FILE_NAME
	fi

	echo $CF_DEPLOYMENT_FILE_NAME
}

bosh_status() {
  export BACKUP_DIR=$5
	export CF_DEPLOYMENT_FILE_NAME=`find $BACKUP_DIR -name "cf-*.yml" -maxdepth 1`
	echo "EXECUTE BOSH STATUS"

	bosh status > $BACKUP_DIR/bosh_status.txt
	export BOSH_UUID=`grep UUID $BACKUP_DIR/bosh_status.txt | cut -d 'D' -f 2 | tr -d ' ' | sort -u`

	export UUID_EXISTS=`grep -Fxq $BOSH_UUID $BACKUP_DIR/$CF_DEPLOYMENT_FILE_NAME`
	if [[ -z $UUID_EXISTS ]]; then
		echo "UUID Matches"
	else
		echo "UUID Mismatch"
		exit 1
	fi

	rm -rf $BACKUP_DIR/bosh_status.txt
}

set_bosh_deployment() {
  export BACKUP_DIR=$5
	export CF_DEPLOYMENT_FILE_NAME=`find $BACKUP_DIR -name "cf-*.yml" -maxdepth 1`
  echo "SET THE BOSH DEPLOYMENT"
	bosh deployment $CF_DEPLOYMENT_FILE_NAME
}

export_cloud_controller_bosh_vms() {
	export BACKUP_DIR=$5
  echo "EXPORT BOSH VMS"
	OUTPUT=`bosh vms | grep "cloud_controller-*" | cut -d '|' -f 2 | tr -d ' '`
	echo $OUTPUT > $BACKUP_DIR/bosh-vms.txt
}

stop_cloud_controller() {
	export BACKUP_DIR=$5
	echo "STOPPING CLOUD CONTROLLER"
	OUTPUT=`cat $BACKUP_DIR/bosh-vms.txt`

	for word in $OUTPUT
	do
		JOB=`echo $word | cut -d '/' -f 1`
		INDEX=`echo $word | cut -d '/' -f 2`

		/usr/bin/expect -c "
			set timeout -1

			spawn bosh stop $JOB $INDEX --force

			expect {
				-re ".*continue.*" {
					send yes\r ;
					interact
					sleep 30
				}
			}

			exit
		"
	done
}

export_db() {
	IP=$1
	USERNAME=$2
	export PGPASSWORD=$3
	PORT=$4
	DB=$5
	DB_FILE=$6

	echo "EXPORT $DB"

	pg_dump -h $IP -U $USERNAME -p $PORT $DB > $DB_FILE

}

export_nfs_server() {
	NFS_IP=$1
	NFS_SERVER_USER=$2
	NFS_SERVER_PASSWORD=$3
	NFS_DIR=$4
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

start_cloud_controller() {
	export BACKUP_DIR=$5
	echo "STARTING CLOUD CONTROLLER"
	OUTPUT=`cat $BACKUP_DIR/bosh-vms.txt`

	for word in $OUTPUT
	do
		JOB=`echo $word | cut -d '/' -f 1`
		INDEX=`echo $word | cut -d '/' -f 2`

		/usr/bin/expect -c "
			set timeout -1

			spawn bosh start $JOB $INDEX --force

			expect {
				-re ".*continue.*" {
					send yes\r ;
					interact
					sleep 30
				}
			}

			exit
		"
	done

}

export_installation() {
	if [[ "Y" = "$COMPLETE_BACKUP" || "y" = "$COMPLETE_BACKUP" ]]; then
		CONNECTION_URL=https://$OPS_MANAGER_HOST/api/installation_asset_collection

		echo "EXPORT INSTALLATION FILES FROM " $CONNECTION_URL

		curl "$CONNECTION_URL" -X GET -u $OPS_MGR_ADMIN_USERNAME:$OPS_MGR_ADMIN_PASSWORD --insecure -k -o $WORK_DIR/installation.zip
	fi
}

execute() {
	validate_software
	copy_deployment_files
	export_Encryption_key
	export_installation_settings
	fetch_bosh_connection_parameters
	bosh_login
	verify_deployment_backedUp
	bosh_status
	set_bosh_deployment
	export_cloud_controller_bosh_vms
	stop_cloud_controller
	export_cc_db
	export_uaadb
	export_consoledb
	export_nfs_server
	start_cloud_controller
	export_installation
}

usage() {
	echo "Usage: cfops install backup <OPS MGR HOST or IP> <SSH PASSWORD> <OPS MGR ADMIN USER> <OPS MGR ADMIN PASSWORD> <OUTPUT DIR>"
	printf "\t %s \t\t\t %s \n" "OPS MGR HOST or IP:" "OPS Manager Host or IP"
	printf "\t %s \t\t\t\t %s \n" "SSH PASSWORD:" "OPS Manager Tempest SSH Password"
	printf "\t %s \t\t\t %s \n" "OPS MGR ADMIN USER:" "OPS Manager Admin Username"
	printf "\t %s \t\t %s \n" "OPS MGR ADMIN PASSWORD:" "OPS Manager Admin Password"
	printf "\t %s \t\t\t\t %s \n" "OUTPUT DIR:" "Backup Directory"
}

$1 $2 $3 $4 $5 $6 $7 $8 $9
