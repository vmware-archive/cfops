cfdeploy
========

A Cloud Foundry Deployment tool for IaaS installation, deployment, and management automation


### Background

cfdeploy is a command line interface (cli) tool to enable targeting a given IaaS (initially AWS) and automate the installation and deployment of Cloud Foundry.  The purpose is to reduce the complexity associated with standing up a typical Cloud Foundry foundation from the command line.  The goal is to enable Cloud Foundry installations to be more easily repeatable and manageable at the IaaS level such that they are not unique "puppies" and instead reflect the principles behind Infrastructure as Code.


### Current

This initial version *only* provides a proposed set of cli functionality along with the "help" dialog.

The project is written in "Go".

For example you can try the various commands, args and flags (and --help documentation) that are currently proposed, such as:

    $ ./cfdeploy -help

    $ ./cfdeploy install -help

    $ ./cfdeploy -i aws install add FOUNDATION1

    $ ./cfdeploy -i aws install destroy FOUNDATION1

    $ ./cfdeploy start FOUNDATION1

    $ ./cfdeploy shutdown FOUNDATION1 --force


etc.

