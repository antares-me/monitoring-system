#!/bin/sh

./generator &
./monitoring &
wait -n
exit $?

