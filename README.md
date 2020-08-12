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
  init        initializes a project
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

## Init
Initializes a projects and writes the output to a pom.xml.new file. This does not upgrade or touch the contents of the file in any way, only rearranges the file according to the co-pilot standard.
```shell script
co-pilot init
```
* Custom target
```shell script
co-pilot init --target /path/to/folder
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
Upgrade functionalities. Every available command writes the output to a new `pom.xml.new` file instead of overwriting the original pom.xml file.
```
Available Commands:
  deps         upgrade dependencies to project
  kotlin       upgrade kotlin version in project
  spring-boot  upgrade spring-boot to the latest version
```


### Dependencies
Upgrades the dependencies of the current project
* Current directory
```shell script
co-pilot upgrade deps
```

* Custom target
```shell script
co-pilot upgrade deps --target /path/to/folder
```

### Kotlin
Upgrades the version information found under `<properties><kotlin.version>...</properties>` to the latest kotlin version found for kotlin-stdlib-jdk8 on maven.org
```shell script
co-pilot upgrade kotlin
```

* Custom target
```shell script
co-pilot upgrade kotlin --target /path/to/folder
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
