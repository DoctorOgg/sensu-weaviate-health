[![Sensu Bonsai Asset](https://img.shields.io/badge/Bonsai-Download%20Me-brightgreen.svg?colorB=89C967&logo=sensu)](https://bonsai.sensu.io/assets/DoctorOgg/sensu-weaviate-health)
![goreleaser](https://github.com/DoctorOgg/sensu-weaviate-health/workflows/goreleaser/badge.svg)

# Sensu Weaviate Health Check

## Table of Contents

- [Overview](#overview)
- [Files](#files)
- [Usage examples](#usage-examples)
- [Configuration](#configuration)
  - [Asset registration](#asset-registration)
  - [Check definition](#check-definition)
- [Installation from source](#installation-from-source)
- [Additional notes](#additional-notes)
- [Contributing](#contributing)

## Overview

This is a [Sensu Check][6] that checks the health of a Weaviate instance by call the `/v1/nodes` endpoint.

The check will return a CRITICAL status if the Weaviate if the node is UNHEALTHY. And a WARNING status if the node is UNAVAILABLE but not ready.

## Files

- bin/check-weaviate-health

## Usage examples

```bash
sensu-weaviate-health -u http://localhost:8080
```

## Configuration

### Asset registration

[Sensu Assets][10] are the best way to make use of this plugin. If you're not using an asset, please
consider doing so! If you're using sensuctl 5.13 with Sensu Backend 5.13 or later, you can use the
following command to add the asset:

```
sensuctl asset add DoctorOgg/sensu-weaviate-health
```

If you're using an earlier version of sensuctl, you can find the asset on the [Bonsai Asset Index][https://bonsai.sensu.io/assets/DoctorOgg/sensu-weaviate-health].

### Check definition

```yml
---
type: CheckConfig
api_version: core/v2
metadata:
  name: https://github.com/DoctorOgg/sensu-weaviate-health
  namespace: default
spec:
  command: sensu-weaviate-health -u http://localhost:8080
  subscriptions:
  - system
  runtime_assets:
  - DoctorOgg/sensu-weaviate-health
```

## Installation from source

The preferred way of installing and deploying this plugin is to use it as an Asset. If you would
like to compile and install the plugin from source or contribute to it, download the latest version
or create an executable script from this source.

From the local path of the <https://github.com/DoctorOgg/sensu-weaviate-health> repository:

```
go build
```

## Contributing

For more information about contributing to this plugin, see [Contributing][1].

[1]: https://github.com/sensu/sensu-go/blob/master/CONTRIBUTING.md
[6]: https://docs.sensu.io/sensu-go/latest/reference/checks/
[10]: https://docs.sensu.io/sensu-go/latest/reference/assets/
