#!/bin/bash
WORKING_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
UPDATER_CONF="$(dirname ${WORKING_DIR})/conf/updater.conf"

sed \
  -e "s|\"UPDATER_CONF\"|\"${UPDATER_CONF}\"|g" \
  ${WORKING_DIR}/scheduleUpdates.yml.tpl
