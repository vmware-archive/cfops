package cfbackup

import "fmt"

// NewSystemsInfo creates a map of SystemDumps that are configured
// based on the installation settings fetched from ops manager
func NewSystemsInfo(installationSettingsFile string, sshKey string) SystemsInfo {

	configParser := NewConfigurationParser(installationSettingsFile)

	var systemDumps = make(map[string]SystemDump)

	for _, job := range configParser.FindCFPostgresJobs() {
		switch job.Identifier {
		case "consoledb":
			systemDumps[ERConsole] = &PgInfo{
				SystemInfo: SystemInfo{
					Product:       "cf",
					Component:     "consoledb",
					Identity:      "root",
					SSHPrivateKey: sshKey,
				},
				Database: "console",
			}
		case "ccdb":
			systemDumps[ERCc] = &PgInfo{
				SystemInfo: SystemInfo{
					Product:       "cf",
					Component:     "ccdb",
					Identity:      "admin",
					SSHPrivateKey: sshKey,
				},
				Database: "ccdb",
			}
		case "uaadb":
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
// jobs that need to be backed up
func (s SystemsInfo) PersistentSystems() []SystemDump {
	ps := []string{ERCc, ERUaa, ERConsole, ERNfs, ERMySQL}
	jobs := []SystemDump{}

	for _, info := range ps {
		fmt.Printf("info:%+v\n", s.SystemDumps[info])
		if _, ok := s.SystemDumps[info]; ok {
			jobs = append(jobs, s.SystemDumps[info])
		}
	}
	return jobs
}
