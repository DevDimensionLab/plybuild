# Co-pilot
A little "go help" for the Java/Kotlin developers using Maven.

Current main capability? 
Upgrade your pom.xml dependencies to latest and greatest! 

Why?
- No installs of maven-plugins required, so if you a working in a multi-repo developer environment with lots of 2party dependencies and repos, you can easily upgrade them with ```co-pilot upgrade 2party```. 
- Brings natural semantics and support for different types of dependencies to the table: Kotlin, 2party, spring-boot (curated dependencies), (other) 3party   
- Can be used as a library for other go-projects automating the upgrade process
- Easy and fast
- Brings feature to the table, not found anywhere else, stay tuned!

Heads up!
- co-pilot rewrites your pom.xml, so make sure you have your pom.xml committed before testing out co-pilot
- known collateral damage: Even though the source code is based on generated code (from https://maven.apache.org/xsd/maven-4.0.0.xsd) 
for the intern pom-model there is a bug that removes attributes on "any-elements". But if you keep your pom.xml 
simple, no fancy plugin-configuration using attributes, there will be no known issues :-)
- start with ```co-pilot init```, verify that the rewrite is ok, commit, and from now on you will easily see the diff that co-pilot introduces with ```co-pilot upgrade <2party|3party|spring-boot|plugins|all>```
- or just use  ```co-pilot status``` (no rewrite) and manually upgrade your pom.xml based on what is reported as outdated, current option if you need to keep your pom.xml formatting
  
Requirement: https://golang.org/doc/install

## Install
```shell script
make
```

## Help
```shell script
co-pilot
```

