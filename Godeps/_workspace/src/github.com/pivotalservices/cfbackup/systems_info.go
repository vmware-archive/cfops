package cfbackup

// NewSystemsInfo creates a map of SystemDumps that are configured
// based on the installation settings fetched from ops manager
func NewSystemsInfo(installationSettingsFile string, sshKey string) SystemsInfo {

	configParser := NewConfigurationParser(installationSettingsFile)
	installationSettings := configParser.InstallationSettings

	const (
		defaultRemoteArchivePath = "/tmp/archive.backup"
		mysqlRemoteArchivePath   = "/var/vcap/store/mysql/archive.backup"
		nfsRemoteArchivePath     = "/var/vcap/store/shared/archive.backup"
	)

	var systemDumps = make(map[string]SystemDump)

	for _, job := range configParser.FindCFPostgresJobs() {
		switch job.Identifier {
		case "consoledb":
			systemDumps[ERConsole] = &PgInfo{
				SystemInfo: SystemInfo{
					Product:           "cf",
					Component:         "consoledb",
					Identifier:        "credentials",
					SSHPrivateKey:     sshKey,
					RemoteArchivePath: defaultRemoteArchivePath,
				},
				Database: "console",
			}
		case "ccdb":
			systemDumps[ERCc] = &PgInfo{
				SystemInfo: SystemInfo{
					Product:           "cf",
					Component:         "ccdb",
					Identifier:        "credentials",
					SSHPrivateKey:     sshKey,
					RemoteArchivePath: defaultRemoteArchivePath,
				},
				Database: "ccdb",
			}
		case "uaadb":
			systemDumps[ERUaa] = &PgInfo{
				SystemInfo: SystemInfo{
					Product:           "cf",
					Component:         "uaadb",
					Identifier:        "credentials",
					SSHPrivateKey:     sshKey,
					RemoteArchivePath: defaultRemoteArchivePath,
				},
				Database: "uaa",
			}
		}
	}
	systemDumps[ERMySQL] = &MysqlInfo{
		SystemInfo: SystemInfo{
			Product:           "cf",
			Component:         "mysql",
			Identifier:          "mysql_admin_credentials",
			SSHPrivateKey:     sshKey,
			RemoteArchivePath: mysqlRemoteArchivePath,
		},
		Database: "mysql",
	}

	systemDumps[ERDirector] = &SystemInfo{
		Product:           installationSettings.GetBoshName(),
		Component:         "director",
		Identifier:          "director_credentials",
		SSHPrivateKey:     sshKey,
		RemoteArchivePath: defaultRemoteArchivePath,
	}
	systemDumps[ERNfs] = &NfsInfo{
		SystemInfo: SystemInfo{
			Product:           "cf",
			Component:         "nfs_server",
			Identifier:          "vm_credentials",
			SSHPrivateKey:     sshKey,
			RemoteArchivePath: nfsRemoteArchivePath,
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
		if _, ok := s.SystemDumps[info]; ok {
			jobs = append(jobs, s.SystemDumps[info])
		}
	}
	return jobs
}
