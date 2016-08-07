/*
Package peeringdb provides structures and functions to interact with the
PeeringDB API. The API documentation is available here:
https://www.peeringdb.com/apidocs/

The PeeringDB API is based on REST principles and the output is formatted in
JSON. This is what this package do. It queries the API with correct URL and
parameters. It then parses the JSON response and make it Go structures. For now
this package can only be used to do GET requests. So it cannot be used to make
any changes on any PeeringDB record.

There is two levels of structures in this package. The first level is a
representation of the first level of the JSON returned by the API. These
structures are named *Resource. They all have the same layout a Meta field,
containing metadata returned by the API, and a Data field being an array of the
second level structures.

All calls to the PeeringDB API are made using the "depth=1" parameter. This
means that sets will be expanded as interger slices instead of slice of
structures. This speeds up the API processing time. To get the structures for a
given set, you just have to iterate over the set and call the right function to
retrieve structures from IDs.

Let's take an example. When we request one or several objects from the
PeeringDB API, the response is always formatted in the same way. First comes
the metadata, then comes data. Data are always in an array since it might
contain more than one object. So when we ask the API to give us a network
object (called Net and represented by the struct of the same name), this
package will parse the first level as a NetResource structure. This structure
contains metadata in its Meta field (if their is any) and Net structures in the
Data field (being an array).
*/
package peeringdb
