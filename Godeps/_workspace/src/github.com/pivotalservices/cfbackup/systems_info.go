package cfbackup

// NewSystemsInfo creates a map of SystemDumps that are configured
// based on the installation settings fetched from ops manager
func NewSystemsInfo(installationSettingsFile string, sshKey string) SystemsInfo {
	var (
		uaadbInfo = &PgInfo{
			SystemInfo: SystemInfo{
				Product:       "cf",
				Component:     "uaadb",
				Identity:      "root",
				SSHPrivateKey: sshKey,
			},
			Database: "uaa",
		}
		consoledbInfo = &PgInfo{
			SystemInfo: SystemInfo{
				Product:       "cf",
				Component:     "consoledb",
				Identity:      "root",
				SSHPrivateKey: sshKey,
			},
			Database: "console",
		}
		ccdbInfo = &PgInfo{
			SystemInfo: SystemInfo{
				Product:       "cf",
				Component:     "ccdb",
				Identity:      "admin",
				SSHPrivateKey: sshKey,
			},
			Database: "ccdb",
		}
		mysqldbInfo = &MysqlInfo{
			SystemInfo: SystemInfo{
				Product:       "cf",
				Component:     "mysql",
				Identity:      "root",
				SSHPrivateKey: sshKey,
			},
			Database: "mysql",
		}
		directorInfo = &SystemInfo{
			Product:       BoshName(),
			Component:     "director",
			Identity:      "director",
			SSHPrivateKey: sshKey,
		}
		nfsInfo = &NfsInfo{
			SystemInfo: SystemInfo{
				Product:       "cf",
				Component:     "nfs_server",
				Identity:      "vcap",
				SSHPrivateKey: sshKey,
			},
		}
	)
	return SystemsInfo{
		SystemDumps: map[string]SystemDump{
			ERDirector: directorInfo,
			ERConsole:  consoledbInfo,
			ERUaa:      uaadbInfo,
			ERCc:       ccdbInfo,
			ERMySQL:    mysqldbInfo,
			ERNfs:      nfsInfo,
		},
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
