SHELL = /bin/bash

INIT_JOB_NAME = "initial-download"

.PHONY: makePV.yml scheduleUpdates.yml deployServer.yml

# Docker images
server-docker-image:
	@pushd ${PWD}/geoipserver >/dev/null && \
		make docker-image && \
		popd

updater-docker-image:
	@pushd ${PWD}/geoipupdate >/dev/null && \
		make docker-image && \
		popd

docker-images: server-docker-image updater-docker-image

# Manifest files
clean:
	@rm -rf *.yml

all-manifests: makePV.yml scheduleUpdates.yml deployServer.yml

makePV.yml:
	@./templates/generate-makePV-manifest.sh > makePV.yml

scheduleUpdates.yml:
	@./templates/generate-scheduleUpdates-manifest.sh > scheduleUpdates.yml

deployServer.yml:
	@./templates/generate-deployServer-manifest.sh > deployServer.yml

# Deploy on k8s
# Make local persistent volumes
makePV: makePV.yml
	@kubectl apply -f makePV.yml

# Set up a cronjob for updating GeoIP2/GeoLite2 database binaries
scheduleUpdates: scheduleUpdates.yml makePV updater-docker-image
	@kubectl apply -f scheduleUpdates.yml
	@kubectl create job --from=cronjob/geoipupdate ${INIT_JOB_NAME}

# Deploy GeoIP WebAPI service
deployServer: deployServer.yml makePV server-docker-image
	@kubectl apply -f deployServer.yml

# Deploy all
deploy: all
all: makePV scheduleUpdates deployServer

# Delete resources deployed on k8s
destroyPV: makePV.yml
	@kubectl delete -f makePV.yml

unscheduleUpdates: scheduleUpdates.yml
	@kubectl delete job ${INIT_JOB_NAME}
	@kubectl delete -f scheduleUpdates.yml

destroyServer: deployServer.yml
	@kubectl delete -f deployServer.yml

destroy-all: destroyServer unscheduleUpdates destroyPV
