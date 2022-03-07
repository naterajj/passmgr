# passmgr - password manager for the CLI - v0.0.1

## Usage

```
passmgr --search <host> | --update <id> | --insert | --import-csv </path/to/csv> | --config </path/to/config> | --help | --version
```

## Configuration

Configuration is stored as single level object/dictionary in JSON format. It has a default path of:

```
$HOME/.passmgr_config
```

The default path can be overriden with the `--config` CLI option.

### Configuration options:

|**key**|**value**|
|-------|---------|
|"dbfile"|"/path/to/database"|
|"enforce-db-permissions"|`true` or `false`|

## Database Schema

The database is automatically initialized by `passmgr` using the following DDL:

`SQL
CREATE TABLE passwords (
  id         INTEGER PRIMARY KEY AUTOINCREMENT,
  host       TEXT NOT NULL,
  url        TEXT NOT NULL,
  username   TEXT,
  password   BLOB NOT NULL,
  created_at TEXT NOT NULL,
  updated_at TEXT NOT NULL
);
`

The password is stored as a binary blob using AES-256 GCM with a `nonce` that is stored along with the ciphertext.
