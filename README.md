# Elastic Common Schema (ECS) support for uber-go/zap logger

Use this library for automatically adding a minimal set of ECS fields to your logs, when using [uber-go/zap](https://github.com/uber-go/zap).

---

**Please note** that this library is in a **beta** version and backwards-incompatible changes might be introduced in future releases. While we strive to comply to [semver](https://semver.org/), we can not guarantee to avoid breaking changes in minor releases.

---

## Documentation

Ready to get started? Documentation is at [elastic.co](https://www.elastic.co/guide/en/ecs-logging/go-zap/current/index.html).

## Test
```
go test ./...
```

## Contribute
Create a Pull Request from your own fork.

Run `mage` to update and format you changes before submitting.

Add new dependencies to the NOTICE.txt.

## License
This software is licensed under the [Apache 2 license](https://github.com/elastic/ecs-logging-go/zap/blob/master/LICENSE).
