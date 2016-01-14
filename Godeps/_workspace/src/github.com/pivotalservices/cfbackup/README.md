# cfbackup
Cloud Foundry Backup Utilities

[![wercker status](https://app.wercker.com/status/daa1b586e39ce2801352461ca4a09078/m "wercker status")](https://app.wercker.com/project/bykey/daa1b586e39ce2801352461ca4a09078)

[![GoDoc](http://godoc.org/github.com/pivotalservices/cfbackup?status.png)](http://godoc.org/github.com/pivotalservices/cfbackup)


### this repo is meant to be included in other projects. It will provide method calls for backing up Ops Manager and Elastic Runtime.


## Running tests / build pipeline locally (docker-machine)

```

# install the wercker cli
$ curl -L https://install.wercker.com | sh

# make sure a docker host is running
$ docker-machine start default && eval $(docker-machine env default)

# run the build pipeline locally, to test your code locally
$ ./testrunner

```

## Running tests / build pipeline locally (boot2docker)

```

# install the wercker cli
$ curl -L https://install.wercker.com | sh

# make sure a docker host is running
$ boot2docker up && $(boot2docker shellinit)

# run the build pipeline locally, to test your code locally
$ ./testrunner

```
