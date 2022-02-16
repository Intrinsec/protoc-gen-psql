# protoc-gen-psql (PGP)

PGP is a protoc plugin to generate postgresql statements from protobuf files.

This project uses [protoc-gen-star](https://github.com/lyft/protoc-gen-star) to ease code generation.

## How to use it

See `./tests/asset.proto` for example on how to use it.

## Options

The plugin support the following options.

### At file level

- `initialization`: these strings will be gathered in a file called `00_init_<YOUR_PROTO_FILE>.pb.psql`.
- `finalization`: these strings will be gathered in a file called `99_final_<YOUR_PROTO_FILE>.pb.psql`.

### At message level

The `tableType` option defines what "type" of table will be contained in the message. Two values are valid:

- `DATA`: This value indicates the other psql options set in the message are describing a SQL "data" table. If this option is present, every other valid psql options in the message will be gathered in a file called `10_tables_<YOUR_PROTO_FILE>.pb.psql`.
- `RELATION`: This value indicates the other psql options set in the message are describing a "relation" table. This kind of message should probably contains constraints. If this option is present, every other valid psql options in the message will be gathered in a file called `20_relations_<YOUR_PROTO_FILE>.pb.psql`.

Other remaining options defines what statements will be contained in the table definition:

- `prefix`: these strings will be set at the beginning of the creation table definition.
- `suffix`: these strings will be set at the end of the creation table definition.
- `constraint`: this option enable defining constraints on the table. This option will handle error if the constraint already exist in the schema.
- `disabled`: this boolean indicates the message should be ignored. If other options are used in this message, they will be ignored.
- `relay_cascade_update`: create a trigger after insert, delete and update on the `source_foreign_key` parent table to update the given field on destination foreign key parent tables:
  - `source_foreign_key`: source foreign key
  - `destinations`: List of destinations to update:
    - `foreign_key`: destination foreign key
    - `field`: field to update on parent table
    - `value`: value to set. The given value can be a function (ex: now())

### At field level

- `column`: let describe how the field should be represented in the schema.
- `auto_fill_on_update`: create a trigger to update this field with the given value each time the table is updated. The given value can be a function (ex: now()).
- `cascade_update_on_related_table`: create a trigger to update other fields from another table each time the table is updated. This option must be set on the foreign key of the table to update and can be set multiple time to update various fields. This option has 2 parameters to set:
  - `field`: field to update on the parent table
  - `value`: value to set. The given value can be a function (ex: now()).

## PSQL files naming convention

Generated files are prefixes with a number to enable executing the psql statements in the right order. Code consuming these files to create/modify a schema should read them following the numbering.

## Use ALTER or not

Two modes are available and may be chosen by adding `alter` as a parameter in protoc command line like `--psql_out="alter:."`

- (default) alter = false : The psql code is generated considering the initial schema is empty (no table, etc.). This mode is great for readability and to generate code for a new database.
- alter = true : The psql code is generated to be applied on an existing database. `ALTER` statements are heavily used.

## Tests

Do a `make test-generate` to view a code generated example.

To test the generate psql code, do a `make test`.
