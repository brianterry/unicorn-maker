![cobra logo](https://cloud.githubusercontent.com/assets/173412/10886352/ad566232-814f-11e5-9cd0-aa101788c117.png)

The unicorn-maker project is a repository that contains an end-to-end example of an AWS Cloudformation resource provider.

The CloudFormation Command Line Interface (CLI) is an open-source tool that enables you to develop and test AWS and third-party resources and register them for use in AWS CloudFormation. Begin by going to the documentation guide and setting up your build environment.

[![](https://img.shields.io/github/workflow/status/spf13/cobra/Test?longCache=tru&label=Test&logo=github%20actions&logoColor=fff)](https://github.com/spf13/cobra/actions?query=workflow%3ATest)
[![Build Status](https://travis-ci.org/spf13/cobra.svg "Travis CI status")](https://travis-ci.org/spf13/cobra)
[![GoDoc](https://godoc.org/github.com/spf13/cobra?status.svg)](https://godoc.org/github.com/spf13/cobra)
[![Go Report Card](https://goreportcard.com/badge/github.com/spf13/cobra)](https://goreportcard.com/report/github.com/spf13/cobra)
[![Slack](https://img.shields.io/badge/Slack-cobra-brightgreen)](https://gophers.slack.com/archives/CD3LP1199)


# Overview
Welcome to the unicorn-maker project!


What is a CloudFormation resource provider? Excellent question! AWS CloudFormation introduced a set of capabilities that made it easy to model and automate third-party resources such as SaaS monitoring or incident management tools with infrastructure-as-code benefits.

With this launch, you can use AWS CloudFormation as a single tool to automate the provisioning of your infrastructure and application resources, whether AWS or third party, without the need for custom scripts or manual processes. You can now create your own private AWS CloudFormation resource providers, share them with the open-source community, and leverage third-party providers developed by others.

Cool, right?  How do I get started? Wow, you are full of great questions. I built this project to help you get started. In this repository, you will find an example of an AWS Cloudformation resource provider that you can use as an example.

# Getting Started

First, start by installing the AWS Cloudformation CLI and the language plugins.

Although not necessary, I recommend creating a Python virtual environment. It makes getting started a little easier.

    $ python3 -m venv env
    $ source source env/bin/activate


Now that you have your enviroment setup, clone the is repo.

    git clone https://github.com/brianterry/unicorn-maker.git

## Choose your path
What's great about creating an AWS CloudFormation provider is you can write it in JAVA, Go, Python, or Javascript.

In this repo, you will find a folder that contains an example resource built in the following languages:

[Go](https://github.com/brianterry/unicorn-maker/tree/master/go)

[Python](https://github.com/brianterry/unicorn-maker/tree/master/python)

[TypeScript](https://github.com/brianterry/unicorn-maker/tree/master/typescript)

Java (Comming soon)

No matter what path you choose, the resource design is the same. That way, you can use this project as a "rosetta stone."

For example, if you are good at Go and want to learn how to create a Python provider, compare the projects.


## Backend service

A flag is a way to modify the behavior of a command. Cobra supports
fully POSIX-compliant flags as well as the Go [flag package](https://golang.org/pkg/flag/).
A Cobra command can define flags that persist through to children commands
and flags that are only available to that command.

In the example above, 'port' is the flag.

Flag functionality is provided by the [pflag
library](https://github.com/spf13/pflag), a fork of the flag standard library
which maintains the same interface while adding POSIX compliance.




# License

Cobra is released under the Apache 2.0 license. See [LICENSE.txt](https://github.com/spf13/cobra/blob/master/LICENSE.txt)



