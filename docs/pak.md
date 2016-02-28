# Pak Format Specification
A Pak is composed of these parts:
- Pak metadata file stored as `pak.yaml` or `pak.json` which contains general information and path to resource templates
- One or more Kubernetes resource files which supports parameterization

## Pak Metadata File
This file called `pak.yaml` or `pak.json` and contains following information:
- **name**: Name of the Pak
- **version**: SemVer 2 compatible version
- **url**: Absolute and canonical URL for the Pak
- **description (optional)**: Multi-line strings that describes the Pak
- **tags (optional)**: List of related tags which used for searching and grouping in UIs
- **icon (optional)**: An absolute URL to an image used in GUIs
- **resources**: List of absolute or relative paths to kubernetes resources that need to be installed (see _Resources_ below)
- **properties**: List of parameters and their specification used for parameterizing resources (see _Properties_ below)

For example this a Pak file for Redis packages:

```yaml
name: redis
url: paks/redis-1.0/redis.yaml
version: "1.0"
description: |
  Reliable, Scalable Redis on Kubernetes

resources:
  - rc.yaml
  - svc.yaml

properties:
  - name: port
    description: Listening port for redis
    type: int
    default: 6379
```
