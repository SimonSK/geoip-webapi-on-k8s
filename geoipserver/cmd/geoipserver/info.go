package main

import "fmt"

const (
	name         = "GeoIP WebAPI Server"
	description  = "A simple web API server that reads GeoIP2/GeoLite2 database binary and echoes out latitude and longitude information for a given IP address"
	author       = "Simon Kim"
	versionMajor = 0
	versionMinor = 1
)

var (
	version         = fmt.Sprintf("%d.%d", versionMajor, versionMinor)
	nameWithVersion = fmt.Sprintf("%s v%s", name, version)
)
