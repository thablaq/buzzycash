#!/bin/bash

DUMP_FILE="schema.sql"
OUT_DIR="./internal/db/migrations"

mkdir -p "$OUT_DIR"

UP_FILE="$OUT_DIR/000001_init.up.sql"
DOWN_FILE="$OUT_DIR/000001_init.down.sql"

# Clear previous content
> "$UP_FILE"
> "$DOWN_FILE"

TABLES=()

# Step 1: Extract CREATE TABLE, ALTER TABLE (constraints), CREATE INDEX
awk '
  BEGIN { in_table=0 }

  /^CREATE TABLE/ {
    in_table=1
    table=$3
    gsub(/"/, "", table)
    tables[i++] = table
  }

  in_table {
    print >> upfile
    if (/;/) { in_table=0 }
  }

  /^ALTER TABLE/ {
    print >> upfile
  }

  /^CREATE INDEX/ {
    print >> upfile
  }

  END {
    # Print table names to stdout for shell
    for (t in tables) print tables[t]
  }
' upfile="$UP_FILE" "$DUMP_FILE" > /tmp/tables_list.txt

# Step 2: Generate DOWN file (drop tables in reverse order)
t_count=$(wc -l < /tmp/tables_list.txt)
for ((i=t_count; i>=1; i--)); do
  table=$(sed -n "${i}p" /tmp/tables_list.txt)
  echo "DROP TABLE IF EXISTS \"$table\" CASCADE;" >> "$DOWN_FILE"
done

echo "Migration files created:"
echo "  UP:   $UP_FILE"
echo "  DOWN: $DOWN_FILE"
