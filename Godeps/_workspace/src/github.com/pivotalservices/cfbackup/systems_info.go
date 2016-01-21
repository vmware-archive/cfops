package cfbackup

import "fmt"

// NewSystemsInfo creates a map of SystemDumps that are configured
// based on the installation settings fetched from ops manager
func NewSystemsInfo(installationSettingsFile string, sshKey string) SystemsInfo {

	fmt.Println("we have an sys info installation settings file %s", installationSettingsFile)

	configParser := NewConfigurationParser(installationSettingsFile)
	installationSettings := configParser.installationSettings

	fmt.Println("we have a some installationSettings  %v", installationSettings)
	//director config
	//blobstore type
	// db type

	// ccdb backup postgres if instance_count >=1 else will be backedf up by HAMySQL
	// uaadb backup postgres if instance_count >=1 else will be backedf up by HAMySQL
	// consoleDB backup postgres if instance_count >=1 else will be backedf up by HAMySQL
	// mysql backup postgres if instance_count >=1 else will be backedf up by HAMySQL
	var systemDumps = make(map[string]SystemDump)

	// Search for CF and check DB jobs for postgres
	for _, product := range installationSettings.Products {
		identifier := product.Identifer
		if identifier == "cf" {
			fmt.Println("we have a cf product")
			for _, job := range product.Jobs { //jobs are release jobs like nats/ccdb

				if isPostgres(job.Identifier, job.Instances) { // job id is db

					fmt.Println("we have some psql")

					systemDumps[ERConsole] = &PgInfo{
						SystemInfo: SystemInfo{
							Product:       "cf",
							Component:     "consoledb",
							Identity:      "root",
							SSHPrivateKey: sshKey,
						},
						Database: "console",
					}

					systemDumps[ERCc] = &PgInfo{
						SystemInfo: SystemInfo{
							Product:       "cf",
							Component:     "ccdb",
							Identity:      "admin",
							SSHPrivateKey: sshKey,
						},
						Database: "ccdb",
					}

					systemDumps[ERUaa] = &PgInfo{
						SystemInfo: SystemInfo{
							Product:       "cf",
							Component:     "uaadb",
							Identity:      "root",
							SSHPrivateKey: sshKey,
						},
						Database: "uaa",
					}
				}
			}
		}
	}
	systemDumps[ERMySQL] = &MysqlInfo{
		SystemInfo: SystemInfo{
			Product:       "cf",
			Component:     "mysql",
			Identity:      "root",
			SSHPrivateKey: sshKey,
		},
		Database: "mysql",
	}
	systemDumps[ERDirector] = &SystemInfo{
		Product:       BoshName(),
		Component:     "director",
		Identity:      "director",
		SSHPrivateKey: sshKey,
	}
	systemDumps[ERNfs] = &NfsInfo{
		SystemInfo: SystemInfo{
			Product:       "cf",
			Component:     "nfs_server",
			Identity:      "vcap",
			SSHPrivateKey: sshKey,
		},
	}

	return SystemsInfo{
		SystemDumps: systemDumps,
	}
}

// PersistentSystems returns a slice of all the
// configured SystemDump for an installation
func (s SystemsInfo) PersistentSystems() []SystemDump {
	v := make([]SystemDump, len(s.SystemDumps))
	idx := 0
	for _, value := range s.SystemDumps {
		v[idx] = value
		idx++
	}
	return v
}

func isPostgres(jobdb string, instances []Instances) bool {
	//Create a slice containing all pg dbs
	pgdbs := []string{"ccdb", "uaadb", "consoledb"}

	for _, pgdb := range pgdbs {
		if pgdb == jobdb {
			//Check to see in say uaadb has a local db instance
			for _, instances := range instances {
				val := instances.Value
				if val >= 1 {
					return true
				}
			}
		}
	}
	return false
}
