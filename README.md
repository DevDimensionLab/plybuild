# Co-pilot


## Compilation
```shell script
go build
```

## Help
```shell script
co-pilot
```

```
Available Commands:
  download    Downloads options
  help        Help about any command
  maven       maven options
  spring      Spring boot tools
  upgrade     Upgrade options
```

## Download
Download functionality 
```
Available Commands:
  cli         Downloads spring-cli

```

### CLI
Downloads the spring cli to a target folder in the current directory
```shell script
co-pilot download cli
```

## Maven
Maven helper funcionality
```
Available Commands:
  repositories list repositories
```

### Repositories
List the repositories found in settings.xml or default
```shell script
co-pilot maven repositories
```

## Spring
Spring functionality
```
Available Commands:
  init        Spring init
  status      Spring status
```

### Init 
Creates a simple webservice using start.spring.io

* Default webservice
```shell script
co-pilot spring init
```

* Custom webservice from json file
```shell script
co-pilot spring init --config-file example.init.config.json
```

### Status
Status gets last default version from start.spring.io
```shell script
co-pilot spring status
```


## Upgrade
Upgrade functionality
```
Available Commands:
  dependencies upgrade dependencies to project
  spring-boot  upgrade spring-boot to the latest version
```


### Dependencies
Upgrades the dependencies of the current project
* Current directory
```shell script
co-pilot upgrade dependencies
```

* Custom target
```shell script
co-pilot upgrade dependencies --target /path/to/folder
```

### Spring Boot
Upgrades the pom.xml file found in directory to newest version of spring boot

* Current directory
```shell script
co-pilot upgrade spring-boot
```

* Custom target
```shell script
co-pilot upgrade spring-boot --target /path/to/folder
```
