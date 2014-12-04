# Go Go BOSH - BOSH client API for golang applications

This project is a golang library for applications wanting to talk to a BOSH/MicroBOSH or bosh-lite.

<a href='http://www.babygopher.com'><img src='https://raw2.github.com/drnic/babygopher-site/gh-pages/images/babygopher-badge.png' ></a>

* [![GoDoc](https://godoc.org/github.com/cloudfoundry-community/gogobosh?status.png)](https://godoc.org/github.com/cloudfoundry-community/gogobosh)
* Test status [![Build Status](https://travis-ci.org/cloudfoundry-community/gogobosh.svg)](https://travis-ci.org/cloudfoundry-community/gogobosh)


## API

The following client functions are available, as a subset of the full BOSH Director API.

* repo.GetInfo()
* repo.GetStemcells()
* repo.DeleteStemcell("bosh-stemcell", "993")
* repo.GetReleases()
* repo.DeleteReleases("cf")
* repo.DeleteRelease("cf", "144")
* repo.GetDeployments()
* repo.ListDeploymentVMs("cf-warden")
* repo.FetchVMsStatus("cf-warden")
* repo.GetTaskStatus(123)

### API compatibility

Note: Development is currently being done against bosh-lite v147.

The BOSH Core team nor its Product Managers do not claim that a BOSH director has a public API; and they want to make changes to the API in the future. It will may be tricky for golang apps to support different BOSH APIs. We'll figure this out as we go.

The best way to describe the API support in this library is to document what version of bosh-lite is being tested against, the date that it was published. Hopefully bosh-lite is always approximately parallel (via rebasing) in its API with the main BOSH project; and the same timestamps can map to the continuously delivered releases of BOSH & its RubyGems.

Trying to write a client library for an API without any versioning strategy could get messy for client applications. Please write your own integration tests that work against running BOSHes that you'll use in production.

If you are using this library, or the Ruby library within the `bosh_cli` rubygem, or talking directly with the BOSH director API - please announce yourself on the bosh-users google group and/or to the PM of BOSH. This way they can be aware of who many be affected by API changes.

## Install

```
go get github.com/cloudfoundry-community/gogobosh
````

## Documentation

The documentation is published to [https://godoc.org/github.com/cloudfoundry-community/gogobosh](https://godoc.org/github.com/cloudfoundry-community/gogobosh).

Also, view the documentation locally with:

```
godoc -goroot=$GOPATH github.com/cloudfoundry-community/gogobosh
```

### Use

There is an extensive [example application](https://github.com/cloudfoundry-community/gogobosh/blob/master/example/bosh-lite-example.go) showing usage of many of the read-only functions.

As a short getting started guide:

``` golang
package main

import (
  "github.com/cloudfoundry-community/gogobosh"
  "github.com/cloudfoundry-community/gogobosh/api"
  "github.com/cloudfoundry-community/gogobosh/net"
  "github.com/cloudfoundry-community/gogobosh/utils"
  "fmt"
  "flag"
)

func main() {
  utils.Logger = utils.NewLogger()

  target := flag.String("target", "https://192.168.50.4:25555", "BOSH director host")
  username := flag.String("username", "admin", "Login with username")
  password := flag.String("password", "admin", "Login with password")
  flag.Parse()

  director := models.NewDirector(*target, *username, *password)
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
}
```

You can automatically detect the current director target, and username/password, from the BOSH CLI's `~/.bosh_config` file as well (see `example/current_target.go`)

```golang
package main

import (
  "github.com/cloudfoundry-community/gogobosh"
  "github.com/cloudfoundry-community/gogobosh/local"
)

func main() {
  configPath, err := local.DefaultBoshConfigPath()
  config, err := local.LoadBoshConfig(configPath)
  target, username, password, err := config.CurrentBoshTarget()
  director := gogobosh.NewDirector(target, username, password)
```

## Tests

Tests are all local currently; and do not test against a running bosh or bosh-lite. I'd like to at least do integration tests against a bosh-lite in future.
