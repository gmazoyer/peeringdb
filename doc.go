/*
Package peeringdb provides structures and functions to interact with the
PeeringDB API. The API documentation is available here:
https://www.peeringdb.com/apidocs/

The PeeringDB API is based on REST principles and returns data formatted in
JSON. This package queries the API with the correct URL and parameters, parses
the JSON response, and converts it into Go structures. Currently, this package
only supports GET requests and cannot be used to modify any PeeringDB records.

There are two levels of structures in this package. The first level represents
the top level of the JSON returned by the API. These structures are named
*Resource. They all have a Meta field containing metadata returned by the API,
and a Data field which is an array of the second level structures.

All calls to the PeeringDB API use the "depth=1" parameter. This means that
sets are expanded as integer slices instead of slices of structures, which
speeds up the API processing time. To get the structures for a given set, you
just need to iterate over the set and call the appropriate function to retrieve
structures from IDs.

For example, when requesting one or more objects from the PeeringDB API, the
response is always formatted in the same way: first comes the metadata, then
the data. The data is always in an array since it might contain more than one
object. When asking the API for a network object (called Net and represented by
the struct of the same name), this package parses the first level as a
NetResource structure. This structure contains metadata in its Meta field (if
there is any) and Net structures in the Data field (as an array).
*/
package peeringdb
