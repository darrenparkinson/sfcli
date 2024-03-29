# Salesforce CLI

[![Status](https://img.shields.io/badge/status-wip-yellow)](https://github.com/darrenparkinson/sfcli) ![GitHub tag (latest SemVer)](https://img.shields.io/github/v/tag/darrenparkinson/sfcli) ![GitHub](https://img.shields.io/github/license/darrenparkinson/sfcli?color=brightgreen) [![GoDoc](https://pkg.go.dev/badge/darrenparkinson/sfcli)](https://pkg.go.dev/github.com/darrenparkinson/sfcli) [![Go Report Card](https://goreportcard.com/badge/github.com/darrenparkinson/sfcli)](https://goreportcard.com/report/github.com/darrenparkinson/sfcli)

A simple utility, written in Go, for interacting with Salesforce.  Currently only specific functionality is implemented, 
and the output is defined by current requirements, however it can be easily extended to add further capabilities.

There is already a [Salesforce CLI available on the Salesforce website](https://developer.salesforce.com/tools/sfdxcli#), but that doesn't currently support the V2 Bulk API. 

However, by way of an example, you can use the Salesforce CLI for v1 bulk upserts as follows:

```sh
$ sfdx auth:web
$ sfdx force:data:bulk:upsert -i Email -f ./examples/contacts.csv -s contact -u you@yourcompany.com -w 1
```

And you can use the existing Salesforce CLI to describe objects as follows:

```sh
sfdx force:schema:sobject:describe -s Account -u you@yourcompany.com
```

Which you could then pipe to `jq` to obtain the fields you want:

```sh
$ # To see how many fields:
$ sfdx force:schema:sobject:describe -s Account -u you@yourcompany.com | jq '.fields | length'
$ # To list specific fields you could do something like:
$ sfdx force:schema:sobject:describe -s Account -u you@yourcompany.com | jq -r '["Name","IDLookup"], (.fields[] | [.name, .idLookup]) | @tsv' | column -t
```

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

You can have different files for different environments and specify which to use with the `--config` option, e.g. `sfcli --config .sfcli.dev.yaml`

For full command help simply use:

```sh
$ sfcli --help

Salesforce CLI Utility

Usage:
  sfcli [command]

Available Commands:
  bulk        Bulk API V2 Commands
  describe    list field names for the various objects
  help        Help about any command

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
* Describe (show object fields)
  * Account 
  * Contact
  * Opportunity
* Describe other object types with `-o` option


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
  -c, --crlf              Specify CRLF Line Ending (default is LF)
  -e, --external string   External ID Field
  -f, --file string       CSV File
  -h, --help              help for upsert
  -s, --sobject string     Type of SObject for Insert, e.g. Account, Contact, Opportunity
```

### CSV Format

Use the correct column names as headers in the CSV.  These can be obtained from the "describe" endpoint for each object type.  
This CLI provides a `describe` command for objects (see below).

You can use relationship fields in the CSV so long as the `External ID` field is selected or, if it's a standard field,
its `idLookup` property is set to `true`.

For example, for an Account, you can use `Owner.Email` because `Owner` is a relationship to a "User" and the "User" 
`Email` field has its `idLookup` property set to true.  See the [Salesforce](https://developer.salesforce.com/docs/atlas.en-us.api_asynch.meta/api_asynch/relationship_fields_in_a_header_row__2_0.htm) documentation for more detail.

## Describing objects

There are some objects that have their own command, such as account, contact and opportunity.  You can also specify the object type on the command line for objects that don't have their own command. Here are some examples:

```sh
$ sfcli describe account
$ sfcli describe contact
$ sfcli describe opportunity
$ sfcli describe -o campaign
$ sfcli describe -o lead
```