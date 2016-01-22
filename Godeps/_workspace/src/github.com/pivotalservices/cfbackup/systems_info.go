package cfbackup

import "github.com/xchapter7x/lo"

// NewSystemsInfo creates a map of SystemDumps that are configured
// based on the installation settings fetched from ops manager
func NewSystemsInfo(installationSettingsFile string, sshKey string) SystemsInfo {

	lo.G.Debugf("we have an sys info installation settings file %s", installationSettingsFile)

	configParser := NewConfigurationParser(installationSettingsFile)
	installationSettings := configParser.installationSettings

	lo.G.Debugf("we have a some installationSettings  %v", installationSettings)

	var systemDumps = make(map[string]SystemDump)
	dumps := []SystemDump{}

	for _, product := range installationSettings.Products {
		identifier := product.Identifer
		if identifier == "cf" {
			for _, job := range product.Jobs {

				if isPostgres(job.Identifier, job.Instances) {

					systemDumps[ERConsole] = &PgInfo{
						SystemInfo: SystemInfo{
							Product:       "cf",
							Component:     "consoledb",
							Identity:      "root",
							SSHPrivateKey: sshKey,
						},
						Database: "console",
					}
					dumps = append(dumps, systemDumps[ERConsole])

					systemDumps[ERCc] = &PgInfo{
						SystemInfo: SystemInfo{
							Product:       "cf",
							Component:     "ccdb",
							Identity:      "admin",
							SSHPrivateKey: sshKey,
						},
						Database: "ccdb",
					}
					dumps = append(dumps, systemDumps[ERCc])

					systemDumps[ERUaa] = &PgInfo{
						SystemInfo: SystemInfo{
							Product:       "cf",
							Component:     "uaadb",
							Identity:      "root",
							SSHPrivateKey: sshKey,
						},
						Database: "uaa",
					}
					dumps = append(dumps, systemDumps[ERUaa])
					break
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
	dumps = append(dumps, systemDumps[ERMySQL])
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
	dumps = append(dumps, systemDumps[ERNfs])

	return SystemsInfo{
		SystemDumps: systemDumps,
		Dumps:       dumps,
	}
}

// PersistentSystems returns a slice of all the
// configured SystemDump for an installation
func (s SystemsInfo) PersistentSystems() []SystemDump {
	return s.Dumps
}

func isPostgres(jobdb string, instances []Instances) bool {
	pgdbs := []string{"ccdb", "uaadb", "consoledb"}

	for _, pgdb := range pgdbs {
		if pgdb == jobdb {
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
