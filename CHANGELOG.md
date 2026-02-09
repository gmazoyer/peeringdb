# Changelog

## Version 0.1.0

### Added

- Go client for the PeeringDB API using generics, `context.Context`, and
  functional options.
- Support for all PeeringDB resource types: `Network`, `NetworkFacility`,
  `NetworkInternetExchangeLAN`, `NetworkContact`, `Organization`, `Facility`,
  `Campus`, `Carrier`, `CarrierFacility`, `InternetExchange`,
  `InternetExchangeLAN`, `InternetExchangePrefix`, and
  `InternetExchangeFacility`.
- `GetASN()` convenience function to look up a network by AS number.
- API key authentication via `WithAPIKey` option.
- Custom HTTP client support via `WithHTTPClient` option.
- Custom API URL support via `WithURL` option.
