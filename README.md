# Loudspeaker Runtime

[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](https://opensource.org/licenses/MIT)
[![Docker](https://img.shields.io/docker/v/masanetes/loudspeaker-runtime/v0.0.1?color=blue&logo=docker)](https://hub.docker.com/repository/docker/masanetes/loudspeaker-runtime)
[![Go Reference](https://pkg.go.dev/badge/github.com/masanetes/loudspeaker-runtime.svg)](https://pkg.go.dev/github.com/masanetes/loudspeaker-runtime)
[![Test](https://github.com/masanetes/loudspeaker-runtime/actions/workflows/test.yaml/badge.svg)](https://github.com/masanetes/loudspeaker-runtime/actions/workflows/test.yaml)
[![report](https://goreportcard.com/badge/github.com/masanetes/loudspeaker-runtime)](https://goreportcard.com/report/github.com/masanetes/loudspeaker-runtime)
[![codecov](https://codecov.io/gh/masanetes/loudspeaker-runtime/branch/master/graph/badge.svg?token=9HT5CC8XDK)](https://codecov.io/gh/masanetes/loudspeaker-runtime)

Get the kubernetes event from the WatchAPI and send it to the target.

```mermaid
flowchart LR
Loudspeaker(loudspeaker) -->|Watch Event| KubeAPI[KubeAPI]  
Loudspeaker(loudspeaker) -->|Watch Configmaps| KubeAPI[KubeAPI]  
Loudspeaker -->|Events| C[Listener1]
```

# Settings

## Required Environment

|Environment|Details|
|-|-|
|Type|Listener Type. This will change the format of the credentials to be read.|
|ConfigmapName|The name of the configmaps to load the configuration from.|
|Namespace|Specify the namespace of configmaps to be monitored by WatchAPI.|

## Configmaps data format

```yaml
observes:
  - namespace: "default"
    ignoresReasons:
      - ""
    involvedObjectNames: 
      - "sample-cronjob"
    involvedObjectKinds:
      - "Cronjob"
    eventTypes:
      - "Warning"
```

Refer to the CRD API specification for a description of each field.
https://pkg.go.dev/github.com/masanetes/loudspeaker@v0.0.6/api/v1alpha1
