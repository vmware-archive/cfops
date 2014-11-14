package install

import (
	"fmt"

	"github.com/cloudfoundry-community/gogobosh"
	"github.com/cloudfoundry-community/gogobosh/api"
	"github.com/cloudfoundry-community/gogobosh/local"
	"github.com/cloudfoundry-community/gogobosh/net"
	"github.com/cloudfoundry-community/gogobosh/utils"

	"github.com/pivotalservices/cfops/system"
)

type DumpCommand struct {
	CommandRunner system.CommandRunner
	Installer
}

func (cmd DumpCommand) Metadata() system.CommandMetadata {
	return system.CommandMetadata{
		Name:        "dump",
		ShortName:   "d",
		Usage:       "dump the configuration information of an existing deployment",
		Description: "dump an existing cloud foundry foundation deployment configuration from the iaas",
	}
}

func (cmd DumpCommand) HasSubcommands() bool {
	return false
}

func (cmd DumpCommand) Run(args []string) (err error) {
	utils.Logger = utils.NewLogger()

	// target := flag.String("target", "https://192.168.50.4:25555", "BOSH director host")
	// username := flag.String("username", "admin", "Login with username")
	// password := flag.String("password", "admin", "Login with password")
	// flag.Parse()

	configPath, err := local.DefaultBoshConfigPath()
	if err != nil {
		fmt.Println(err)
		return
	}
	config, err := local.LoadBoshConfig(configPath)
	if err != nil {
		fmt.Println(err)
		return
	}
	target, username, password, err := config.CurrentBoshTarget()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("Targeting %s with user %s...\n", target, username)

	director := gogobosh.NewDirector(target, username, password)
	repo := api.NewBoshDirectorRepository(&director, net.NewDirectorGateway())

	info, apiResponse := repo.GetInfo()
	if apiResponse.IsNotSuccessful() {
		fmt.Println("Could not fetch BOSH info")
		return
	}

	fmt.Println("Director")
	fmt.Printf("  Name       %s\n", info.Name)
	fmt.Printf("  URL        %s\n", info.URL)
	fmt.Printf("  Version    %s\n", info.Version)
	fmt.Printf("  User       %s\n", info.User)
	fmt.Printf("  UUID       %s\n", info.UUID)
	fmt.Printf("  CPI        %s\n", info.CPI)
	if info.DNSEnabled {
		fmt.Printf("  dns        %#v (%s)\n", info.DNSEnabled, info.DNSDomainName)
	} else {
		fmt.Printf("  dns        %#v\n", info.DNSEnabled)
	}
	if info.CompiledPackageCacheEnabled {
		fmt.Printf("  compiled_package_cache %#v (provider: %s)\n", info.CompiledPackageCacheEnabled, info.CompiledPackageCacheProvider)
	} else {
		fmt.Printf("  compiled_package_cache %#v\n", info.CompiledPackageCacheEnabled)
	}
	fmt.Printf("  snapshots  %#v\n", info.SnapshotsEnabled)
	fmt.Println("")
	fmt.Printf("%#v\n", info)
	fmt.Println("")

	stemcells, apiResponse := repo.GetStemcells()
	if apiResponse.IsNotSuccessful() {
		fmt.Println("Could not fetch BOSH stemcells")
		return
	} else {
		for _, stemcell := range stemcells {
			fmt.Printf("%#v\n", stemcell)
		}
		fmt.Println("")
	}

	releases, apiResponse := repo.GetReleases()
	if apiResponse.IsNotSuccessful() {
		fmt.Println("Could not fetch BOSH releases")
		return
	} else {
		for _, release := range releases {
			fmt.Printf("%#v\n", release)
		}
		fmt.Println("")
	}

	deployments, apiResponse := repo.GetDeployments()
	if apiResponse.IsNotSuccessful() {
		fmt.Println("Could not fetch BOSH deployments")
		return
	} else {
		for _, deployment := range deployments {
			fmt.Printf("%#v\n", deployment)
		}
		fmt.Println("")
	}

	task, apiResponse := repo.GetTaskStatus(1)
	if apiResponse.IsNotSuccessful() {
		fmt.Println("Could not fetch BOSH task 1")
		return
	} else {
		fmt.Printf("%#v\n", task)
	}

	fmt.Println("")
	fmt.Println("VMs in cf-warden deployment:")
	vms, apiResponse := repo.ListDeploymentVMs("cf-warden")
	if apiResponse.IsNotSuccessful() {
		fmt.Println("Could not get list of VM for cf-warden")
		return
	} else {
		for _, vm := range vms {
			fmt.Printf("%#v\n", vm)
		}
	}

	fmt.Println("")
	fmt.Println("VMs status in cf-warden deployment:")
	vmsStatuses, apiResponse := repo.FetchVMsStatus("cf-warden")
	if apiResponse.IsNotSuccessful() {
		fmt.Println("Could not fetch VMs status for cf-warden")
		return
	} else {
		for _, vmStatus := range vmsStatuses {
			fmt.Printf("%s/%d is %s, IPs %#v\n", vmStatus.JobName, vmStatus.Index, vmStatus.JobState, vmStatus.IPs)
		}
	}
	return
}
