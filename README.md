# Albatross Client Go

A go client to interact with the albatross API

## Installation

```
go get -u github.com/gojekfarm/albatross-client-go
```

## Quickstart

```go
client := api.NewClient("http://localhost:8080")
```

The default client configures a default timeout without any retry configuration. To set custom configuration, you can pass the appropriate config option methods to NewClient.

```go
client := api.NewClient(
	"http://localhost:8080",
	config.WithTimeout(10*time.Second),
	config.WithRetry(&config.Retry{
		RetryCount: 5,
		Backoff: 500 * time.Millisecond,
	}),
	config.WithLogger(logger),
)
```

You can provide a custom logger for the client. The custom logger must implement the logger interface, defined under `logger/interface.go`.

### Install

```go

flags := flags.InstallFlags{
	CommonFlags: flags.CommonFlags{
		Namespace: "namespace",
	},
}

status, err := client.Install(
	context.Background(),
	"testrelease",
	"stable/chart",
	api.Values{"override": "some"},
	flags,
)

```

### Upgrade

```go

flags := flags.UpgradeFlags{
	Install: true,
	Version: "xyz",
	CommonFlags: flags.CommonFlags{
		Namespace: "namespace",
	},
}

status, err := client.Upgrade(
	context.Background(),
	"testrelease",
	"stable/chart",
	api.Values{"override": "another"},
	flags,
)

```

### List

```go

flags := flags.ListFlags{
	Deployed: true,
	CommonFlags: flags.CommonFlags{
		Namespace: "namespace",
	},
}

releases, err := client.List(context.Background(), flags)

```

### Uninstall

```go 

flags := flags.UninstallFlags{
	CommonFlags: flags.CommonFlags{
		Namespace: "namespace"
	}
}
 
release, err := client.Uninstall(context.Background(), name, flags)

```

## Status

The project is under development, and the API is subject to breaking changes.

## License

```
Copyright 2020 GO-JEK Tech

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
```


