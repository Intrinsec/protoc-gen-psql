# Cascade updates and auto fill options

## Context

In our application, each time an object is updated, we want its `update_time` field to be updated to `now()` (`auto_fill` option). Moreover, for some complex objects, we want to update their `update_time` when other related objects are updated or even when relation between objects are updated.

An example is available on [tests/action.proto](tests/action.proto), here, an incident is linked to multiple actions and we want to update the incident `update_time` field when one of its action is updated (option `relay_cascade_update`). And also, when the relation between these two objects is updated. Ex.: an action is moved from one incident to another, we also want that both incidents `update_time` fields to be updated (`cascade_update_on_related_table`).

## Decision

The 3 main solutions identified to solve this problem are:

- Implement the logic directly in our application
- Recompute the value at each read access from the `update_time` fields of the various linked objects (works only with soft delete objects)
- Create PostgreSQL triggers to update the `update_time` field of the main object.

For our needs, we want to keep as much logic as possible in the database which excludes implementing it at the application level. Moreover, we want read access to our application to be as fast as possible and some of our database objects can be deleted, which excludes to compute this value at each read access.

In our case, the best solution is to create triggers. However, we don't want to have to maintain dozens of almost identical triggers. This is why we want to abstract this by generating triggers from protobuf options.

## Technical explanation

The triggers below address our need to update the `update_time` field of Incident when Incident, Action or EntityAction are updated (see. [tests/action.proto](tests/action.proto)):

- to update the incident `update_time` field when the incident is updated:

```pgsql
CREATE OR REPLACE FUNCTION fn_incident_auto_update_update_time()
  RETURNS TRIGGER
  LANGUAGE plpgsql
AS $$
BEGIN
  NEW.update_time = now();
  return NEW;
END
$$;

CREATE TRIGGER tg_incident_auto_update_update_time
  BEFORE UPDATE ON Incident
  FOR EACH ROW
  WHEN (OLD IS DISTINCT FROM NEW)
  EXECUTE FUNCTION fn_incident_auto_update_update_time();
```

- to update the incident `update_time` field when an action is updated:

```pgsql
CREATE OR REPLACE FUNCTION fn_action_incident_update_time()
  RETURNS trigger
  LANGUAGE plpgsql
AS $$
BEGIN
  IF TG_OP IN ('DELETE', 'UPDATE') THEN
    UPDATE Incident SET update_time = now() WHERE uuid = (SELECT incident_uuid FROM EntityAction WHERE action_uuid = OLD.uuid);
  END IF;

  IF TG_OP IN ('INSERT', 'UPDATE') THEN
    UPDATE Incident SET update_time = now() WHERE uuid = (SELECT incident_uuid FROM EntityAction WHERE action_uuid = NEW.uuid);
  END IF;

  RETURN NULL;
END
$$;
CREATE TRIGGER tg_action_incident_update_time AFTER UPDATE ON Action FOR EACH ROW EXECUTE FUNCTION fn_action_incident_update_time();
```

- And, another trigger on the relation table when a link between an action and an incident is created, deleted or updated:

```pgsql
CREATE OR REPLACE FUNCTION fn_entityaction_incident_update_time()
  RETURNS trigger
  LANGUAGE plpgsql
AS $$
BEGIN
  IF TG_OP IN ('DELETE', 'UPDATE') THEN
    UPDATE Incident SET update_time = now() WHERE uuid = OLD.incident_uuid;
  END IF;

  IF TG_OP IN ('INSERT', 'UPDATE') THEN
    UPDATE Incident SET update_time = now() WHERE uuid = NEW.incident_uuid;
  END IF;

  RETURN NULL;
END
$$;
CREATE TRIGGER tg_entityaction_incident_update_time AFTER UPDATE ON EntityAction FOR EACH ROW EXECUTE FUNCTION fn_action_incident_update_time();
```

In order to generate these triggers from a protobuf option, some information about the foreign keys are required (primary keys and parent tables). Below, an example of `relay_cascade_update` option with all the information needed to create a trigger on Action table to update the `update_time` field of the related tables (incident and communication):

```protobuf
  option (psql.relay_cascade_update) = {
      source_foreign_key : "action_uuid"
      src_table_name: "Action"
      src_pk_column_name: "uuid"
      destinations: [
          {
              foreign_key: "incident_uuid"
              dst_table_name: "Incident"
              dst_pk_column_name: "uuid"
              field: "update_time"
              value: "now()"
          },
          {
              foreign_key: "communication_uuid"
              dst_table_name: "Incident"
              dst_pk_column_name: "uuid"
              field: "update_time"
              value: "10"
          }
      ]
  };
```

However, The parent table and their primary key related to a foreign key are already defined in `psql.column` option, and it will be less error prone that `protoc-gen-psql` retrieves this information from here instead of duplicate it in cascade update options:

```protobuf
string action_uuid = 1 [
        (psql.column) = "uuid UNIQUE REFERENCES Action(id) ON DELETE CASCADE"
    ];
```

To retrieve this information from `psql.column` option, it is necessary to parse this option or rework it to ease the parsing. It is something we might end up doing one day, but, currently, `psql.column` option works well for our needs and we didn't want to add more complexity on this option.

It is also possible to retrieve this information from the information_schema table. This query retrieves the parent table name and their primary key from `action_uuid` foreign key:

```pgsql
SELECT
  kcu2.table_name,
  kcu2.column_name INTO dst_table_name,
  dst_pk_column_name
FROM
  information_schema.key_column_usage AS kcu
  INNER JOIN information_schema.referential_constraints AS rc ON rc.constraint_name = kcu.constraint_name
  INNER JOIN information_schema.key_column_usage AS kcu2 ON kcu2.constraint_name = rc.unique_constraint_name
WHERE
  kcu.table_name = 'entityaction'
  AND kcu.column_name = 'action_uuid'
  AND kcu.constraint_name LIKE '%_fkey';
```

Performance is an important requirement for our application, so, we do not want this query to be executed every time a trigger is called. That is why, we created a "shell" function that will call this query and then create the trigger and its associated function on database initialization.

We end up with a more concise option without duplicate information:

```protobuf
option (psql.relay_cascade_update) = {
    source_foreign_key : "action_uuid"
    destinations: [
        {
            foreign_key: "incident_uuid"
            field: "update_time"
            value: "now()"
        },
        {
            foreign_key: "communication_uuid"
            field: "update_time"
            value: "10"
        }
    ]
};
```

In conclusion, our solution enables us to have concise cascade update options without adding more complexity on `psql.column` option. But, to avoid sacrifying performance on trigger execution, we had to sacrifice performance at initilization, which slowed down our functional tests execution by about 30%.

### Other

### AFTER triggers

`WHEN (OLD IS DISTINCT FROM NEW)` on triggers prevents to trigger a cascade update or an auto_fill when updated object did not change. However, a `BEFORE` trigger will not be able to do this check on tables with generated columns (these columns are not yet generated on the `NEW` object on `BEFORE` trigger). On the contrary, `AFTER` triggers do not seem to limit this PostgreSQL features.

### Deferrable triggers on delete

An `AFTER` trigger called during a cascade delete raises a constraint violation errors. Ex. Delete an incident will delete the linked actions (cascade delete). However, when an action is deleted, it will trigger a cascade update to update the `update_time` of the linked incident (`relay_cascade_update` option) which is already deleted.

A deferrable trigger which will be executed at the end of the transaction seems to not have this limitation or at least it fails silently.

It is maybe linked to this explanation from [Before Delete triggers on PostgreSQL documentation](https://www.postgresql.org/docs/9.1/sql-createtrigger.html):

> SQL specifies that BEFORE DELETE triggers on cascaded deletes fire after the cascaded DELETE completes. The PostgreSQL behavior is for BEFORE DELETE to always fire before the delete action, even a cascading one. This is considered more consistent. > There is also nonstandard behavior if BEFORE triggers modify rows or prevent updates during an update that is caused by a referential action. This can lead to constraint violations or stored data that does not honor the referential constraint.

## Tips

### Format specifiers with pg_format

pg_format does not correctly reformat format specifiers (ex.: `%1$I` becomes `% 1$I`) from our generated functions. This can be fixed by adding the following placeholder argument: `pg_format --placeholder '%([0-9]+\$)?-?([0-9]+|\*|(\*[0-9]+\$))?[sIL]'`.
