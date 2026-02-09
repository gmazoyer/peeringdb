/*
Package peeringdb provides structures and functions to interact with the
PeeringDB API. The API documentation is available here:
https://www.peeringdb.com/apidocs/

The PeeringDB API is based on REST principles and returns data formatted in
JSON. This package queries the API with the correct URL and parameters, parses
the JSON response, and converts it into Go structures. Currently, this package
only supports GET requests and cannot be used to modify any PeeringDB records.

All calls to the PeeringDB API use the "depth=1" parameter. This means that
sets are expanded as integer slices instead of slices of structures, which
speeds up the API processing time. To get the structures for a given set, you
just need to iterate over the set and call the appropriate function to retrieve
structures from IDs.
*/
package peeringdb
