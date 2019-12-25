#!/bin/sh
# Assume the last argument is an absolute path to the database binary to be read
for FILENAME in "$@"; do true; done
# Check every 10 seconds if the file exists until found.
while [ ! -f "${FILENAME}" ]
do
  sleep 10
done

PROCESS_NAME="geoipserver"
COMMAND="${PROCESS_NAME} $@"

# Terminate the process approximately 30 min after database update attempts on every Tuesday
crontab -l | { cat; echo "30 */3 * * 2 pkill -15 ${PROCESS_NAME}"; } | crontab -
crond -b

# Run loop
exec /loop.sh ${COMMAND}
