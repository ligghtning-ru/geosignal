# Contributing

Contributions are welcome when they keep the library small and predictable.

## Guidelines

- Keep behavior conservative: prefer a broader label over a precise label that
  is not supported by observations.
- Add tests for every new alias, country rule, or resolver behavior.
- Avoid provider-specific business logic.
- Avoid network calls, databases, and hidden global state.
- Keep public APIs simple.

## Development

```bash
go test ./...
go vet ./...
go test -bench=. ./...
```

## Adding Normalization Rules

New city aliases or suffixes should be backed by real examples and regression
tests. If a rule is specific to one application, prefer passing it with
`WithCityAlias`, `WithCityAliases`, `WithCountryAlias`, or `WithCitySuffix`
instead of changing the defaults.
