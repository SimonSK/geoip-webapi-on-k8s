#!/bin/bash
WORKING_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
source ${WORKING_DIR}/../conf/cluster.conf

sed \
  -e "s|\"DATABASE_DIR\"|\"${DATABASE_DIR}\"|g" \
  ${WORKING_DIR}/makePV.yml.tpl
