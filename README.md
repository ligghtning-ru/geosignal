# geosignal

[![ci](https://github.com/ligghtning-ru/geosignal/actions/workflows/ci.yml/badge.svg)](https://github.com/ligghtning-ru/geosignal/actions/workflows/ci.yml)
[![Go Reference](https://pkg.go.dev/badge/github.com/ligghtning-ru/geosignal.svg)](https://pkg.go.dev/github.com/ligghtning-ru/geosignal)
[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)

`geosignal` is a small Go library for turning conflicting IP geolocation
observations into honest user-facing labels.

IP geolocation is not a single truth. Different providers may report nearby
districts, old city spellings, city-level drift, or entirely different
countries. This package helps backend services avoid confidently showing a
precise city when the evidence does not support it.

## Status

`geosignal` is intentionally small and stable. The public API is designed around
plain structs, deterministic output, and no external services.

## What It Does

- Normalizes common city aliases and spellings.
- Collapses nearby city observations into an area label.
- Keeps far same-country city disagreement at country level.
- Treats cross-country disagreement as ambiguous.
- Exposes agreement metadata for UI, logs, and debugging.
- Includes timezone compatibility helpers for IANA aliases and same-offset zones.

## Design Goals

- Be conservative by default.
- Keep all logic offline and deterministic.
- Make disagreement visible instead of hiding it behind a single provider.
- Keep provider-specific policy outside the package.
- Prefer a small API over a configurable framework.

## Non-Goals

- No external API calls.
- No IP reputation, proxy, VPN, or fraud verdicts.
- No provider scoring.
- No MaxMind/IP2Location/IPInfo clients.
- No database or cache layer.

## Install

```bash
go get github.com/ligghtning-ru/geosignal
```

## Quick Start

```go
package main

import (
	"fmt"

	"github.com/ligghtning-ru/geosignal"
)

func main() {
	display := geosignal.Resolve([]geosignal.Observation{
		{Source: "provider_a", City: "Lviv", Country: "UA", Lat: 49.8397, Lng: 24.0297, HasCoord: true},
		{Source: "provider_b", City: "Vynnyky", Country: "UA", Lat: 49.8156, Lng: 24.1297, HasCoord: true},
	}, geosignal.Observation{City: "Lviv", Country: "UA"})

	fmt.Println(display.Label)
	// Lviv area, UA
}
```

## Examples

### Nearby Cities

```go
display := geosignal.Resolve([]geosignal.Observation{
	{Source: "a", City: "Lviv", Country: "UA", Lat: 49.8397, Lng: 24.0297, HasCoord: true},
	{Source: "b", City: "Vynnyky", Country: "UA", Lat: 49.8156, Lng: 24.1297, HasCoord: true},
}, geosignal.Observation{City: "Lviv", Country: "UA"})

fmt.Println(display.Label)
// Lviv area, UA
```

### Far Same-Country Disagreement

```go
display := geosignal.Resolve([]geosignal.Observation{
	{Source: "a", City: "Dallas", Country: "US", Lat: 32.7767, Lng: -96.7970, HasCoord: true},
	{Source: "b", City: "San Diego", Country: "US", Lat: 32.7157, Lng: -117.1611, HasCoord: true},
	{Source: "c", City: "New Brunswick", Country: "US", Lat: 40.4862, Lng: -74.4518, HasCoord: true},
}, geosignal.Observation{City: "Dallas", Country: "US"})

fmt.Println(display.Kind, display.Label)
// ambiguous_city US
```

### Cross-Country Disagreement

```go
display := geosignal.Resolve([]geosignal.Observation{
	{Source: "a", City: "São Paulo", Country: "BR"},
	{Source: "b", City: "Los Angeles", Country: "US"},
}, geosignal.Observation{City: "São Paulo", Country: "BR"})

fmt.Println(display.Kind, display.Label)
// ambiguous_country Location varies by source
```

### Agreement Metadata

```go
obs := []geosignal.Observation{
	{Source: "a", City: "Los Angeles", Country: "US", Count: 2},
	{Source: "b", City: "Brooklyn", Country: "US", Count: 1},
}

display := geosignal.Resolve(obs, geosignal.Observation{Country: "US"})
agreement := geosignal.AgreementForDisplay(display, obs)

fmt.Println(display.Label)
fmt.Println(agreement.ObservedCities)
```

## API Overview

```go
display := geosignal.Resolve(observations, primary)

display, agreement := geosignal.ResolveWithAgreement(observations, primary)

ok := geosignal.TimezonesCompatible("Europe/Kiev", "Europe/Kyiv", "UA", time.Now())
```

Configuration is passed explicitly with options:

```go
display := geosignal.Resolve(
	observations,
	primary,
	geosignal.WithCityAlias("SPB", "Saint Petersburg"),
	geosignal.WithRegionRadiusKM(35),
)
```

## Display Kinds

| Kind | Meaning |
| --- | --- |
| `single_source` | One stable location, or country-only data. |
| `region_cluster` | Different city names are geographically close enough to show an area. |
| `plurality_same_country` | One city has a strong same-country plurality. |
| `ambiguous_city` | Country is stable, but there is no safe city winner. |
| `ambiguous_country` | Sources disagree on the country. |

## Timezones

```go
ok := geosignal.TimezonesCompatible(
	"Europe/Kiev",
	"Europe/Kyiv",
	"UA",
	time.Now(),
)

fmt.Println(ok)
// true
```

Compatibility is broader than string equality. It handles known IANA aliases,
same-offset DST signatures, and country-level civil timezone variants.

## Why This Exists

Many systems turn provider output into a confident city too early. That creates
bad UX and misleading decisions. `geosignal` is intentionally conservative: it
prefers a slightly broader label over a precise label that is not supported by
the observations.

## Development

```bash
go test ./...
go vet ./...
go test -run TestDoesNotExist -bench . -benchmem -count=1
```

Tests live next to the package source, which is the standard Go layout for small
libraries. Examples are also included as GoDoc examples, so they render on
pkg.go.dev.

## Package Layout

The repository keeps one import path:

```go
import "github.com/ligghtning-ru/geosignal"
```

Source files are split by responsibility (`resolve`, `normalize`, `agreement`,
`timezone`, `distance`, `options`, `types`) while staying in one package. This
keeps the public API compact and idiomatic for Go users.

## License

MIT
