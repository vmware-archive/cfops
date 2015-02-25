[![wercker status](https://app.wercker.com/status/d0a50d426b77a9f73da0fe4f383ad624/m/master "wercker status")](https://app.wercker.com/project/bykey/d0a50d426b77a9f73da0fe4f383ad624)

CF Ops
======

A Cloud Foundry Operations tool for IaaS installation, deployment, and management automation


### Background

CF Ops (cfops) is a command line interface (cli) tool to enable targeting a given IaaS (initially AWS) and automate the installation and deployment of Cloud Foundry.  The purpose is to reduce the complexity associated with standing up and managing a typical Cloud Foundry foundation from the command line.  The goal is to enable Cloud Foundry installations to be more easily repeatable and manageable at the IaaS level such that they are not treated as unique "snowflakes" and instead reflect the principles behind Infrastructure as Code.


### Current

This initial version *only* provides a proposed set of cli functionality along with the "help" dialog.

The project is written in "Go".

For example you can try the various commands, args and flags (and --help documentation) that are currently proposed, such as:

    $ ./cfops -help

    $ ./cfops survey -h

    $ ./cfops prepare -h

    $ ./cfops install -help

    $ ./cfops -iaas aws install add CF_FOUNDATION_I

    $ ./cfops -iaas aws install destroy CF_FOUNDATION_II

    $ ./cfops start CF_FOUNDATION_I

    $ ./cfops shutdown CF_FOUNDATION_I --force


etc.


Sample help output:

    mbp:cfops farmer$ ./cfops help
    NAME:
       cfops - Cloud Foundry Operations tool for IaaS installation, deployment, and management automation

    USAGE:
       cfops [global options] command [command options] [arguments...]

    VERSION:
       0.0.0

    COMMANDS:
       survey, sur		analyze and inspect the deployment environment
       prepare, p		prepare the deployment environment
       install, in		install cloud foundry to an iaas
       start, s		    start up an entire cloud foundry foundation
       restart, r		shutdown and restart an entire cloud foundry foundation
       shutdown, stop	shutdown and stop an entire cloud foundry foundation
       test, t		    test the Cloud Foundry deployment and underlying IaaS environment
       help, h		    Shows a list of commands or help for one command

    GLOBAL OPTIONS:
       --iaas, -i 'aws, vsphere, vcloud, openstack'	set the IaaS type to target for deployment [$CF_IAAS]
       --debug, -d 'true, false'			enable/disable debug output [$CF_TRACE]
       --lang, -l 'en, es'				language for the cfops cli responses [$CF_LANG]
       --help, -h					show help
       --version, -v				print the version

### Install

Run `./bin/build.sh` to build the binary. The binary will be built into the `./out` directory.
This will also copy the `config/assets/config.json` to your `~/.cfops` directory unless it exists.

