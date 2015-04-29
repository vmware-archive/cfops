[![wercker status](https://app.wercker.com/status/d0a50d426b77a9f73da0fe4f383ad624/m/master "wercker status")](https://app.wercker.com/project/bykey/d0a50d426b77a9f73da0fe4f383ad624)

[![GoDoc](http://godoc.org/github.com/pivotalservices/cfops?status.png)](http://godoc.org/github.com/pivotalservices/cfops)


CF Ops
======

A Cloud Foundry Operations tool for IaaS installation, deployment, and management automation


### Background

CF Ops (cfops) is a command line interface (cli) tool to enable targeting a given IaaS (initially AWS) and automate the installation and deployment of Cloud Foundry.  The purpose is to reduce the complexity associated with standing up and managing a typical Cloud Foundry foundation from the command line.  The goal is to enable Cloud Foundry installations to be more easily repeatable and manageable at the IaaS level such that they are not treated as unique "snowflakes" and instead reflect the principles behind Infrastructure as Code.


### Current

This initial version *only* provides backup and restore.

The project is written in "Go".

For example you can try the various commands, args and flags (and --help documentation) that are currently proposed, such as:

    $ ./cfops -help

    $ ./cfops backup

    $ ./cfops restore

etc.


Sample help output:

    $ ./cfops
    NAME:
       cfops - Cloud Foundry Operations tool for IaaS installation, deployment, and management automation
    
    USAGE:
       cfops [global options] command [command options] [arguments...]
    
    VERSION:
       0.0.0
    
    COMMANDS:
       backup, b	backup a Cloud Foundry deployment
       help, h	Shows a list of commands or help for one command
    
    GLOBAL OPTIONS:
       --help, -h		show help
       --version, -v	print the version
    
    
    
    $ ./cfops help backup
    NAME:
       backup - backup a Cloud Foundry deployment
    
    USAGE:
       command backup [command options] [arguments...]
    
    DESCRIPTION:
       Backup a Cloud Foundry deployment, including Ops Manager configuration, databases, and blob store
    
    OPTIONS:
       --hostname, --host 		hostname for Ops Manager
       --username, -u 		username for Ops Manager
       --password, -p 		password for Ops Manager
       --tempestpassword, --tp 	password for the Ops Manager tempest user
       --destination, -d 		directory where the Cloud Foundry backup should be stored

### Install

download latest version for your system. details here:
https://github.com/pivotalservices/cfops/wiki

