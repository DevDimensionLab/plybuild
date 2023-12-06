# ply
A little "go help" for the Java/Kotlin developers using Maven.

Current main capability? 
Upgrade your pom.xml dependencies to latest and greatest! 

Why?
- No installs of maven-plugins required, so if you a working in a multi-repo developer environment with lots of 2party dependencies and repos, you can easily upgrade them with `ply upgrade 2party`. 
- Brings natural semantics and support for different types of dependencies to the table: Kotlin, 2party, spring-boot (curated dependencies), (other) 3party   
- Can be used as a library for other go-projects automating the upgrade process
- Easy and fast
- Brings feature to the table, not found anywhere else, stay tuned!

Heads up!
- ply rewrites your pom.xml, so make sure you have your pom.xml committed before testing out ply
- start with `ply format pom`, verify that the rewrite of the pom.xml is ok, commit, and from now on you will easily see the diff that ply introduces with ```ply upgrade <2party|3party|spring-boot|plugins|all>```
- or just use  `ply status` (no rewrite) and manually upgrade your pom.xml based on what is reported as outdated, current option if you need to keep your pom.xml formatting
  
Requirement: https://golang.org/doc/install

```shell script
  _____  _       
 |  __ \| |      
 | |__) | |_   _ 
 |  ___/| | | | |
 | |    | | |_| |
 |_|    |_|\__, |
            __/ |
           |___/ 

Usage:
  ply [command]

Available Commands:
  bitbucket   Bitbucket functionality
  clean       Clean files and folder in a project
  completion  Generate the autocompletion script for the specified shell
  diagrams    Various tools for generating diagrams
  doc         Opens documentation in default browser
  examples    Examples found in cloud-config
  format      Format functionality for a project
  generate    Initializes a maven project with ply files and formatting
  git         Git commands
  help        Help about any command
  info        Prints info on spring version, dependencies etc
  init        Initializes a maven project with ply files and formatting
  install     Various install options for generating autocompletion etc
  lint        Linting commands
  maven       Run maven (mvn) commands
  merge       Merge functionalities for files to a project
  profiles    Manage profiles settings for ply
  query       Query dependencies in a project
  status      Status functionality for a project
  tips        Use tips to learn information faster
  upgrade     Upgrade options

Flags:
      --debug   turn on debug output
  -h, --help    help for ply
      --json    turn on json output logging

Additional help topics:
  ply about      About ply

Use "ply [command] --help" for more information about a command.
```

## Install
```shell script
make
```

## Help
```shell script
ply
```

