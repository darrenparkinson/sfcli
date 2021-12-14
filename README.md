# Salesforce CLI

[![Status](https://img.shields.io/badge/status-wip-yellow)](https://github.com/darrenparkinson/sfcli) ![GitHub tag (latest SemVer)](https://img.shields.io/github/v/tag/darrenparkinson/sfcli) ![GitHub](https://img.shields.io/github/license/darrenparkinson/sfcli?color=brightgreen) [![GoDoc](https://pkg.go.dev/badge/darrenparkinson/sfcli)](https://pkg.go.dev/github.com/darrenparkinson/sfcli) [![Go Report Card](https://goreportcard.com/badge/github.com/darrenparkinson/sfcli)](https://goreportcard.com/report/github.com/darrenparkinson/sfcli)

A simple utility, written in Go, for interacting with Salesforce.  Currently only specific functionality is implemented, 
and the output is defined by current requirements, however it can be easily extended to add further capabilities.

There is already a [Salesforce CLI available on the Salesforce website](https://developer.salesforce.com/tools/sfdxcli#), but that doesn't currently support the V2 Bulk API.

## Authentication

Currently only the "username-password authorization flow" is supported.  Moving forward, the aim would be to support other authorization flows.

You can follow the [Quick Start](https://developer.salesforce.com/docs/atlas.en-us.api_rest.meta/api_rest/quickstart_oauth.htm) for setting up this authorization method.

The following environment variables are required for authentication:

* `CLIENT_ID` - consumer key from salesforce with access to the APIs
* `CLIENT_SECRET` - associated secret for the consumer key
* `USERNAME` - username with access to salesforce
* `PASSWORD` - associated password for that user
* `BASEURL` - base url for the salesforce tenant, e.g.  `https://mycompany--uat.my.salesforce.com`

You can also provide these in a `.sfcli` yaml file in your home directory or the directory in which you are running the command, e.g.:

```yaml
CLIENT_ID: fjafafhalsdjfhaksj§hdf
CLIENT_SECRET: afjf34kjhsgljdgnajk
PASSWORD: supersecretpassword
USERNAME: me@mycompany.com
BASEURL: https://mycompany--uat.my.salesforce.com
```

You can have different files for different environments and specify which to use with the `--config` option, e.g. `sfcli --config .sfcli.dev`

For full command help simply use:

```sh
$ sfcli --help

Salesforce CLI Utility

Usage:
  sfcli [command]

Available Commands:
  accounts      account related commands
  bulk          bulk API V2 Commands
  completion    generate the autocompletion script for the specified shell
  contacts      contact related commands
  help          Help about any command
  opportunities opportunity related commands

Flags:
      --config string   config file (default is $HOME/.sfcli.yaml)
  -h, --help            help for sfcli

Use "sfcli [command] --help" for more information about a command.
```

## Supported Features

The following capabilities are currently available with this tool:

* Bulk Uploads
  * List Bulk Upload Jobs
  * Show Bulk Upload Job Status
  * Create a Bulk Insert Job
  * Create a Bulk Upsert Job
* Accounts
  * Describe Accounts (show account fields)
  * List Accounts
* Contacts
  * Describe Contacts
  * List Contacts
* Opportunities
  * Describe Contacts
  * List Contacts


## Bulk Uploads

Bulk uploads are achieved with the `sfcli bulk` command:

* `bulk insert` : Bulk Insert a CSV File
* `bulk list` : List the last 1000 bulk jobs
* `bulk status` : Get the status of a specific job
* `bulk upsert` : Bulk Upsert a CSV File

Each of the commands supports various flags as required which can be displayed within the help, e.g.:

```sh
$ sfcli bulk upsert --help

Bulk Upsert a CSV File

Usage:
  sfcli bulk upsert [flags]

Flags:
  -e, --external string   External ID Field
  -f, --file string       CSV File
  -h, --help              help for upsert
  -o, --object string     Type of Object for Insert, e.g. Account, Contact, Opportunity
```

### CSV Format

Use the correct column names as headers in the CSV.  These can be obtained from the "describe" endpoint for each object type.  
This CLI provides a `describe` command for some objects.

You can use relationship fields in the CSV so long as the `External ID` field is selected or, if it's a standard field,
its `idLookup` property is set to `true`.

For example, for an Account, you can use `Owner.Email` because `Owner` is a relationship to a "User" and the "User" 
`Email` field has its `idLookup` property set to true.  See the [Salesforce](https://developer.salesforce.com/docs/atlas.en-us.api_asynch.meta/api_asynch/relationship_fields_in_a_header_row__2_0.htm) documentation for more detail.