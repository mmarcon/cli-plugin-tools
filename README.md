# MongoDB Tools as Atlas CLI Plugin

_*Disclaimer:* this plugin is mostly a proof of concept built for fun as a weekend project._

This project exposes MongoDB tools as an Atlas CLI plugin.
The tools experience is simplified and fully integrated into the CLI experience, removing the need of passing connections strings around.

## `restore` command

Restores an archive (created with the `dump` command or with `mongodump`) into the destination deployment.

```
$ atlas tools restore <deploymentName> --archive <path | URL> [--dbuser databaseUser] [--dbpass databasePassword] [--debug]
```

## `dump` command

Creates a dump of a running deployment.

```
$ atlas tools dump <deploymentName> --archive <path> [--dbuser databaseUser] [--dbpass databasePassword] [--db dbName] [--debug]
```

When `--db` is not passed, the plugin has the [same behavior](https://www.mongodb.com/docs/database-tools/mongodump/mongodump-behavior/#data-exclusion) as `mongodump`.