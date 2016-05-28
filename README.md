# Kupak - Kubernetes Package Manager

[![Travis Widget]][Travis]

[Travis]: https://travis-ci.org/cafebazaar/kupak
[Travis Widget]: https://travis-ci.org/cafebazaar/kupak.svg?branch=master

Kupak is package manager for installing and basic management of Kubernetes resources using a format called **pak**.

Pak is a format for parameterizing and grouping related Kubernetes resources like pods, replication controllers and services. With kupak you can install, uninstall, track and update paks in your Kubernetes cluster.

Pak supports Go text/templating for parameterization with a simple format for defining parameters and their types.

See this [repo](https://github.com/cafebazaar/paks) for some ready-to-use paks and examples.

## Features
- Simplicity
- No external database
- Tracking and listing all installed Paks
- No server-side configuration
- CLI and library

## Usage and Installation
### Prerequisite
Kupak requires a working `kubectl` installed.

### Installation
```
go get github.com/cafebaazar/kupak
cd $GOPATH/github.com/cafebaazar/kupak
glide install
go build
```

### Repositories
Any directory with an `index.yaml` file like [this](https://github.com/cafebazaar/paks/blob/master/index.yaml) can be a `kupak` repository. Repository address is the path to the directory in which the repository index file lives. The path can be an absolute or relative path in the local file system, over the internet with an `http(s)://` prefix, or hosted on github which is specified in the form of `github.com/user/repository`.

### Usage
- Browsing a repository:

      `kupak paks`
- Info of a pak in a specific repository:

      `kupak --repo=github.com/cafebazaar/paks spec elasticsearch-cluster`
- Installed paks:

      `kupak list`
- Installing a pak:

      `kupak install talkative`
- Removing a pak:

      `kupak remove redis-staging2`
- Adding a repository:

      `kupak repos add github.com/cafebazaar/booker-paks booker project paks repo`

### Writing a Pak
[See this](docs/pak.md)

## License
Copyright 2016 Hezardastan, Inc.

Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file except in compliance with the License. You may obtain a copy of the License at [http://www.apache.org/licenses/LICENSE-2.0](http://www.apache.org/licenses/LICENSE-2.0)

Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the License for the specific language governing permissions and limitations under the License.
