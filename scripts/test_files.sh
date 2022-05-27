#!/usr/bin/bash

EXEC="sqlite3 /tmp/tmp.db -line"

# add_file(name, path)
function add_file {
    $EXEC "INSERT INTO files (id, name, path, created) VALUES (\"$(uuidgen)\", \"$1\", \"$2\", CURRENT_TIMESTAMP);"
}

add_file "test1" "/tmp/foo1"
add_file "test2" "/tmp/foo2"
add_file "test3" "/tmp/foo3"
add_file "test4" "/tmp/foo4"
add_file "test5" "/tmp/foo5"
add_file "test6" "/tmp/foo6"