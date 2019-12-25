#!/bin/bash
WORKING_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
source ${WORKING_DIR}/../conf/cluster.conf
SERVER_ENTRYPOINT_SCRIPT="$(dirname ${WORKING_DIR})/misc/server-entrypoint.sh"
SERVER_LOOP_SCRIPT="$(dirname ${WORKING_DIR})/misc/server-loop.sh"
# Make sure the scripts are made executable
chmod +x ${SERVER_ENTRYPOINT_SCRIPT} ${SERVER_LOOP_SCRIPT}

sed \
  -e "s|LOCALHOST_PORT|${LOCALHOST_PORT}|g" \
  -e "s|NUM_SERVER_REPLICA|${NUM_SERVER_REPLICA}|g" \
  -e "s|\"SERVER_ENTRYPOINT_SCRIPT\"|\"${SERVER_ENTRYPOINT_SCRIPT}\"|g" \
  -e "s|\"SERVER_LOOP_SCRIPT\"|\"${SERVER_LOOP_SCRIPT}\"|g" \
  ${WORKING_DIR}/deployServer.yml.tpl
