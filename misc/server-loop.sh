#!/bin/sh
COMMAND="$@"
while :
do

  ${COMMAND} >&2
done
