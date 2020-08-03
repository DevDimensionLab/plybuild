# Spring-boot-co-pilot


## Compilation
```shell script
go build
```

## Execution

### Help
```shell script
./spring-boot-co-pilot
```

```
Available Commands:
  download    Downloads ...
  help        Help about any command
  spring      Spring ...

```

### Download
Download functionality 
```
Available Commands:
  cli         Downloads spring-cli

```

#### CLI
```shell script
./spring-boot-co-pilot download cli
```

### Spring
Spring functionality
```
Available Commands:
  dependencies Spring dependencies
  info         Spring metadata info
  init         Spring init
  root         Spring root
```
#### Dependencies
Lists dependencies from start.spring.io
```shell script
./spring-boot-co-pilot spring dependencies
```

#### Info
Info metadata from start.spring.io, lists versions etc
```shell script
./spring-boot-co-pilot spring info
```

#### Init 
* Default webservice
```shell script
./spring-boot-co-pilot spring init
```

* custom webservice from json file
```shell script
./spring-boot-co-pilot spring init --config-file example.init.config.json
```

#### Root
Root metadata info from start.spring.io
```shell script
./spring-boot-co-pilot spring root
```
