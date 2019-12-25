# GeoIP WebAPI on Kubernetes


## Objective
Deploy followings on Kubernetes:
* `geoipserver` -- a simple web API server that reads a GeoLite2 database binary
  and echoes out latitude and longitude information for a given IP address
  * For more details, see `geoipserver/` or https://github.com/SimonSK/geoipserver
* `geoipupdate` -- an updater that downloads a new version of the GeoLite2 database
  * For more details, see `geoipupdate/` or https://github.com/SimonSK/geoipupdate


## Software requirement
* Docker for Desktop with a local Kubernetes cluster
* Go version 1.11+ (for Go modules support)
* Tested on:
    * MacOS Docker Desktop Community version 2.1.0.5 (40693) stable
    * Kubernetes version 1.14.8
    * Go version 1.13.5


## Quickstart
* Deploy: `make deploy` or `make all`.
* Destroy all deployed resources: `make destroy-all`


## Simple test
* To check if the service is properly running:
  * Go to `localhost:8080` on browser, or
  * Run `curl localhost:8080`
  * Example:
    ```
    $ curl localhost:8080
    GeoIP WebAPI Server v0.1 on port 8080
    Description: A simple web API server that reads GeoIP2/GeoLite2 database binary and echoes out latitude and longitude information for a given IP address
    
    Database info:
      Type: GeoLite2-City
      Format Version: 2.0
      Build Timestamp: 2019-12-24T17:41:41Z
    
    Usage:
      To get all fields: /api/{ipAddress}
      To get "location" fields: /api/{ipAddress}/location
      To get GPS coordinates: /api/{ipAddress}/location/coords
    ```
* To obtain GPS coordinates of an IP address:
  * Go to `localhost:8080/api/{ipAddress}/location/coords` on browser, or
  * Run `curl localhost:8080/api/{ipAddress}/location/coords`
  * Example:
    ```
    $ curl localhost:8080/api/72.36.89.1/location/coords
    {"lat":40.1047,"lon":-88.2062}
    ```

## Configuration

### `conf/cluster.conf`
This file contains parameters used to generate Kubernetes manifests to be deployed.
* `DATABASE_DIR`: an absolute path to a directory on local host where GeoIP database files will be downloaded to
and read from
  * Current value: `/tmp/GeoIP`
  * Note: make sure this directory is allowed for file sharing on Docker for Desktop.
* `LOCALHOST_PORT`: port to bind on localhost
  * Current value: `8080`
* `NUM_SERVER_REPLICA`: number of API server replicas
  * Current value: `3`

### `conf/updater.conf`
This file contains configuration for GeoIP database updater.

For our purpose, the configured parameters (i.e. `AccountID`, `LicenseKey`, `EditionIDs`) should be left as-is.


## Deployment details

Resources are deployed in 3 steps:
1. Make local persistent volumes: `make makePV`
2. Create a cronjob for regularly updating the GeoLite2 database binary: `make scheduleUpdates`
   * To start the initial download, 1 job is manually created based on the scheduled cronjob.
3. Deploy the GeoIP API server: `make deployServer`

### Expected behaviors (assumptions)
* Database binary updates run every 3 hours on Tuesdays starting at 00:00 in UTC.
  * See `templates/scheduleUpdates.yml.tpl:L20`.
* The updater currently downloads updates with a default (invalid) `AccountID/LicenseKey` combination.
  * As of December 25th, 2019, this is still working. MaxMind update server may require a valid license in near future.
* The updater keeps old versions of database binaries.
* `TERM` signal is sent to GeoIP API server every 3 hours on Tuesdays starting at 00:30 in UTC.
  * See `misc/server-entrypoint.sh:L14`.
* The server opens the database file using memory map or loads it into memory;
  changes to the database file or the symlink that points to the file (during updates)
  should not affect the server's operations.
* GeoLite2 `City` database is used to obtain GPS coordinates.
  * The updater downloads the City database.
  * The server reads (and is intended to work only with) the City database.
