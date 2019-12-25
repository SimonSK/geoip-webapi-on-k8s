# GeoIP WebAPI Server
```
DESCRIPTION:
  A simple web API server that reads GeoIP2/GeoLite2 database binary and echoes out latitude and longitude information for a given IP address

USAGE:
  geoipserver [options] path-to-database-binary

OPTIONS:
  --port value, -p value  Network listening port (default: 8080)
  --verbosity value       Logging verbosity: 0=panic, 1=fatal, 2=error, 3=warn, 4=info, 5=debug, 6=trace (default: 4)
  --help, -h              show help
  --version, -v           print the version
```


## How to run
1. Build binary: `make server`
2. Run: `./bin/geoipserver path-to-database-binary`

To build a Docker image: `make docker-image`


## Test
* Go to a default page to check if the service is properly running:
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

## Default pages
`GET /` and `GET /api` requests are routed to the server's default page.
This page presents following information:
* A brief description of the server
* Details of currently loaded database
* API usage


## API functions (`pkg/webapi/handle.go`)
* The API is intended to provide lookups based on IPv4/IPv6 addresses.
* All successful lookup results are returned as `application/json`.
* User-provided `{ipAddress}` key is expected to be 39-character long at max, including `.` for IPv4 or `:` for IPv6.

### `GET /api/{ipAddress}`
* Description: returns all fields in `record` object of the given IP address
* Supported database type: all

### `GET /api/{ipAddress}/location`
* Description: returns all fields in `record["location"]` of the given IP address
* Supported database type: GeoLite2-City

### `GET /api/{ipAddress}/location/coords`
* Description: returns `latitude` and `longitude` fields in `record["location"]` of the given IP address
* Supported database type: GeoLite2-City


## Error handling
The server uses HTTP status code in response to indicate error types.

* `400 BadRequest`: provided key is not parsed as a valid IP address
* `404 NotFound`: requested resource/path is not found on the server
* `500 InternalServerError`: error occurred while obtaining requested information on the given IP address
* `501 NotImplemented`: requested HTTP method is not implemented on the server


## Logging
Events are logged in 7 levels:
* 0=panic
* 1=fatal
* 2=error
* 3=warn
* 4=info (default)
* 5=debug
* 6=trace

Shown below is the overall format of log messages. Timestamp is in RFC3339 format.
```
timestamp [log level] [details represented in a list of key/value pairs] message
```

Example:
```
$ ./bin/geoipserver --verbosity 6 /tmp/GeoIP/GeoLite2-City.mmdb
2019-12-25T16:42:38-06:00 [DEBU] starting GeoIP WebAPI Server v0.1
2019-12-25T16:42:38-06:00 [DEBU] [databaseBinaryFilepath=/tmp/GeoIP/GeoLite2-City.mmdb listenPort=8080] server configs
2019-12-25T16:42:38-06:00 [DEBU] opening /tmp/GeoIP/GeoLite2-City.mmdb
2019-12-25T16:42:38-06:00 [DEBU] [type=GeoLite2-City formatVersion=2.0 buildTimestamp=2019-12-24T11:41:41-06:00] database info
GeoIP WebAPI Server v0.1 on port 8080
Description: A simple web API server that reads GeoIP2/GeoLite2 database binary and echoes out latitude and longitude information for a given IP address

Database info:
	Type: GeoLite2-City
	Format Version: 2.0
	Build Timestamp: 2019-12-24T11:41:41-06:00

Usage:
	To get all fields: /api/{ipAddress}
	To get "location" fields: /api/{ipAddress}/location
	To get GPS coordinates: /api/{ipAddress}/location/coords
2019-12-25T16:42:43-06:00 [INFO] [method=GET route=details remoteAddr=127.0.0.1:62296 query=72.36.89.1] req
2019-12-25T16:42:43-06:00 [INFO] [status=200 route=details remoteAddr=127.0.0.1:62296] resp
2019-12-25T16:42:44-06:00 [INFO] [method=GET route=details remoteAddr=127.0.0.1:62297 query=0.0.0.0] req
2019-12-25T16:42:44-06:00 [INFO] [status=200 route=details remoteAddr=127.0.0.1:62297] resp
2019-12-25T16:43:02-06:00 [INFO] [method=GET route=details remoteAddr=127.0.0.1:62303 query=0.0] req
2019-12-25T16:43:02-06:00 [DEBU] unable to parse 0.0 as an IP address
2019-12-25T16:43:02-06:00 [INFO] [status=400 route=details remoteAddr=127.0.0.1:62303] resp
2019-12-25T16:44:50-06:00 [INFO] [method=POST route=details remoteAddr=127.0.0.1:62370 query=0.0.0.0] req
2019-12-25T16:44:50-06:00 [INFO] [status=501 route=details remoteAddr=127.0.0.1:62370] resp
```


## Note
* oschwald/maxminddb-golang is used to to read and look up records.
* The current API is intended to work with GeoLite2 City database.
  * The server will load other database types (e.g. Country) fine, but the API will mostly likely return null with 500 status code.
