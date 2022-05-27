#!/usr/bin/bash

EXEC="sqlite3 /tmp/tmp.db -line"
BASEDIR="/tmp/srv"

mkdir --parents $BASEDIR

# add_file(name, path)
function add_file {
    #$EXEC "INSERT INTO files (id, name, path, created) VALUES (\"$(uuidgen)\", \"$1\", \"$2\", CURRENT_TIMESTAMP);"
    curl https://loripsum.net/api > "$BASEDIR/$2"
}

add_file "test1" "foo1"
add_file "test2" "foo2"
add_file "test3" "foo3"
add_file "test4" "foo4"
add_file "test5" "foo5"
add_file "test6" "foo6"