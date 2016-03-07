package elasticruntime

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"path"

	"github.com/cloudfoundry-community/go-cfenv"
	"github.com/pivotalservices/cfbackup"
	ghttp "github.com/pivotalservices/gtils/http"
	"github.com/pivotalservices/gtils/log"
	"github.com/xchapter7x/lo"
)

// NewElasticRuntime initializes an ElasticRuntime intance
var NewElasticRuntime = func(jsonFile string, target string, sshKey string, cryptKey string) *ElasticRuntime {
	systemsInfo := cfbackup.NewSystemsInfo(jsonFile, sshKey)
	context := &ElasticRuntime{
		SSHPrivateKey:     sshKey,
		JSONFile:          jsonFile,
		BackupContext:     cfbackup.NewBackupContext(target, cfenv.CurrentEnv(), cryptKey),
		SystemsInfo:       systemsInfo,
		PersistentSystems: systemsInfo.PersistentSystems(),
	}
	return context
}

// Backup performs a backup of a Pivotal Elastic Runtime deployment
func (context *ElasticRuntime) Backup() (err error) {
	return context.backupRestore(cfbackup.ExportArchive)
}

// Restore performs a restore of a Pivotal Elastic Runtime deployment
func (context *ElasticRuntime) Restore() (err error) {
	err = context.backupRestore(cfbackup.ImportArchive)
	return
}

func (context *ElasticRuntime) backupRestore(action int) (err error) {
	var (
		ccJobs []cfbackup.CCJob
	)

	if err = context.ReadAllUserCredentials(); err == nil && context.directorCredentialsValid() {
		lo.G.Debug("Retrieving All CC VMs")
		manifest, erro := context.getManifest()
		if err != nil {
			return erro
		}
		if ccJobs, err = context.getAllCloudControllerVMs(); err == nil {
			directorInfo := context.SystemsInfo.SystemDumps[cfbackup.ERDirector]
			cloudController := cfbackup.NewCloudController(directorInfo.Get(cfbackup.SDIP), directorInfo.Get(cfbackup.SDUser), directorInfo.Get(cfbackup.SDPass), context.InstallationName, manifest, ccJobs)
			lo.G.Debug("Setting up CC jobs")
			defer cloudController.Start()
			cloudController.Stop()
		}
		lo.G.Debug("Running db action")
		if len(context.PersistentSystems) > 0 {
			err = context.RunDbAction(context.PersistentSystems, action)
			if err != nil {
				lo.G.Error("Error backing up db", err)
				err = ErrERDBBackup
			}
		} else {
			err = ErrEREmptyDBList
		}
	} else if err == nil {
		err = cfbackup.ErrERDirectorCreds
	}
	return
}

func (context *ElasticRuntime) getAllCloudControllerVMs() (ccvms []cfbackup.CCJob, err error) {

	lo.G.Debug("Entering getAllCloudControllerVMs() function")
	directorInfo := context.SystemsInfo.SystemDumps[cfbackup.ERDirector]
	connectionURL := fmt.Sprintf(ERVmsURL, directorInfo.Get(cfbackup.SDIP), context.InstallationName)
	lo.G.Debug("getAllCloudControllerVMs() function", log.Data{"connectionURL": connectionURL, "directorInfo": directorInfo})
	gateway := context.HTTPGateway
	if gateway == nil {
		gateway = ghttp.NewHttpGateway()
	}
	lo.G.Debug("Retrieving CC vms")
	if resp, err := gateway.Get(ghttp.HttpRequestEntity{
		Url:         connectionURL,
		Username:    directorInfo.Get(cfbackup.SDUser),
		Password:    directorInfo.Get(cfbackup.SDPass),
		ContentType: "application/json",
	})(); err == nil {
		var jsonObj []cfbackup.VMObject

		lo.G.Debug("Unmarshalling CC vms")
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if err = json.Unmarshal(body, &jsonObj); err == nil {
			ccvms, err = cfbackup.GetCCVMs(jsonObj)
		}
	}
	return
}

//RunDbAction - run a db action dump/import against a list of systemdump types
func (context *ElasticRuntime) RunDbAction(dbInfoList []cfbackup.SystemDump, action int) (err error) {

	for _, info := range dbInfoList {
		lo.G.Debug(fmt.Sprintf("RunDbAction info: %+v", info))

		if err = info.Error(); err == nil {
			err = context.readWriterArchive(info, context.TargetDir, action)

		} else {
			lo.G.Error("readWriterArchive err: ", err)
			break
		}
	}
	return
}

func (context *ElasticRuntime) readWriterArchive(dbInfo cfbackup.SystemDump, databaseDir string, action int) (err error) {
	filename := fmt.Sprintf(ERBackupFileFormat, dbInfo.Get(cfbackup.SDComponent))
	filepath := path.Join(databaseDir, filename)

	var pb cfbackup.PersistanceBackup

	if pb, err = dbInfo.GetPersistanceBackup(); err == nil {
		switch action {
		case cfbackup.ImportArchive:
			lo.G.Debug("Restoring %s", dbInfo.Get(cfbackup.SDComponent))
			var backupReader io.ReadCloser
			if backupReader, err = context.Reader(filepath); err == nil {
				defer backupReader.Close()
				err = pb.Import(backupReader)
				lo.G.Debug("Done restoring %s", dbInfo.Get(cfbackup.SDComponent))
			}
		case cfbackup.ExportArchive:
			lo.G.Info("Exporting %s", dbInfo.Get(cfbackup.SDComponent))
			var backupWriter io.WriteCloser
			if backupWriter, err = context.Writer(filepath); err == nil {
				defer backupWriter.Close()
				err = pb.Dump(backupWriter)
				lo.G.Debug("Done backing up %s", dbInfo.Get(cfbackup.SDComponent))
			}
		}
	}
	return
}

//ReadAllUserCredentials - get all user creds from the installation json
func (context *ElasticRuntime) ReadAllUserCredentials() (err error) {
	configParser := cfbackup.NewConfigurationParser(context.JSONFile)
	installationSettings := configParser.InstallationSettings
	err = context.assignCredentialsAndInstallationName(installationSettings)
	return
}

func (context *ElasticRuntime) assignCredentialsAndInstallationName(installationSettings cfbackup.InstallationSettings) (err error) {

	if err = context.assignCredentials(installationSettings); err == nil {
		context.InstallationName, err = context.getDeploymentName(installationSettings)
	}
	return
}

func (context *ElasticRuntime) assignCredentials(installationSettings cfbackup.InstallationSettings) (err error) {

	for name, sysInfo := range context.SystemsInfo.SystemDumps {
		var (
			userID string
			ip     string
			pass   string
		)
		productName := sysInfo.Get(cfbackup.SDProduct)
		jobName := sysInfo.Get(cfbackup.SDComponent)
		identifier := sysInfo.Get(cfbackup.SDIdentifier)

		if userID, pass, ip, err = context.getVMUserIDPasswordAndIP(installationSettings, productName, jobName); err == nil {
			sysInfo.Set(cfbackup.SDIP, ip)
			sysInfo.Set(cfbackup.SDVcapPass, pass)
			sysInfo.Set(cfbackup.SDVcapUser, userID)
			if identifier == "vm_credentials" {
				sysInfo.Set(cfbackup.SDUser, userID)
				sysInfo.Set(cfbackup.SDPass, pass)
			} else if userID, pass, err = context.getUserIDPasswordForIdentifier(installationSettings, productName, jobName, identifier); err == nil {
				sysInfo.Set(cfbackup.SDUser, userID)
				sysInfo.Set(cfbackup.SDPass, pass)
			}
		}
		context.SystemsInfo.SystemDumps[name] = sysInfo
	}
	return
}

func (context *ElasticRuntime) directorCredentialsValid() (ok bool) {
	var directorInfo cfbackup.SystemDump

	if directorInfo, ok = context.SystemsInfo.SystemDumps[cfbackup.ERDirector]; ok {
		connectionURL := fmt.Sprintf(cfbackup.ERDirectorInfoURL, directorInfo.Get(cfbackup.SDIP))
		gateway := context.HTTPGateway
		userId := directorInfo.Get(cfbackup.SDUser)
		password := directorInfo.Get(cfbackup.SDPass)
		if gateway == nil {
			gateway = ghttp.NewHttpGateway()
		}
		if _, err := gateway.Get(ghttp.HttpRequestEntity{
			Url:         connectionURL,
			Username:    userId,
			Password:    password,
			ContentType: "application/json",
		})(); err == nil {
			ok = true
		} else {
			ok = false
			lo.G.Debug("Error connecting to director using %s, UserId-[%s] and Password[%s]", connectionURL, userId, password)
		}
	}
	return
}

func (context *ElasticRuntime) getManifest() (manifest string, err error) {
	directorInfo, _ := context.SystemsInfo.SystemDumps[cfbackup.ERDirector]
	director := cfbackup.NewDirector(directorInfo.Get(cfbackup.SDIP), directorInfo.Get(cfbackup.SDUser), directorInfo.Get(cfbackup.SDPass), 25555)
	mfs, err := director.GetDeploymentManifest(context.InstallationName)
	if err != nil {
		return
	}
	data, err := ioutil.ReadAll(mfs)
	if err != nil {
		return
	}
	return string(data), nil
}

func (context *ElasticRuntime) getDeploymentName(installationSettings cfbackup.InstallationSettings) (deploymentName string, err error) {
	var product cfbackup.Products
	if product, err = installationSettings.FindByProductID("cf"); err == nil {
		deploymentName = product.InstallationName
	}
	return
}

func (context *ElasticRuntime) getUserIDPasswordForIdentifier(installationSettings cfbackup.InstallationSettings, product, component, identifier string) (userID, password string, err error) {
	var propertyMap map[string]string
	if propertyMap, err = installationSettings.FindPropertyValues(product, component, identifier); err == nil {
		userID = propertyMap["identity"]
		password = propertyMap["password"]
	}
	return
}

func (context *ElasticRuntime) getVMUserIDPasswordAndIP(installationSettings cfbackup.InstallationSettings, product, component string) (userID, password, ip string, err error) {
	var ips []string
	if ips, err = installationSettings.FindIPsByProductAndJob(product, component); err == nil {
		if len(ips) > 0 {
			ip = ips[0]
		} else {
			err = fmt.Errorf("No IPs found for %s, %s", product, component)
		}
		var vmCredential cfbackup.VMCredentials
		if vmCredential, err = installationSettings.FindVMCredentialsByProductAndJob(product, component); err == nil {
			userID = vmCredential.UserID
			password = vmCredential.Password
		}
	}
	return
}
