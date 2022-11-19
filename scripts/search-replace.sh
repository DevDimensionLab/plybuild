#!/usr/bin/env bash

set -e pipefail

pomFile="pom.xml"

while IFS= read -r input; do
  type=$(echo "$input" | jq -r .type)
  versionIsProperty=$(echo "$input" | jq -r .versionIsProperty)
  artifactId=$(echo "$input" | jq -r .artifactId)
  groupId=$(echo "$input" | jq -r .groupId)
  oldVersion=$(echo "$input" | jq -r .oldVersion)
  newVersion=$(echo "$input" | jq -r .newVersion)

  if [[ "$type" != "outdated"* ]]; then
      continue
  fi

  echo "=> [$type] $artifactId:$groupId ($oldVersion -> $newVersion) in $pomFile"
  if [[ ("$versionIsProperty" == "true") ]]; then
      versionValue=$(echo "$input" | jq -r .versionValue | sed -E -e 's/\${|}//g')
      #echo "s/<${versionValue}>${oldVersion}<\/${versionValue}>/<${versionValue}>${newVersion}<\/${versionValue}>/g"
      sed -i '' "s/<${versionValue}>${oldVersion}<\/${versionValue}>/<${versionValue}>${newVersion}<\/${versionValue}>/g" "$pomFile"
  else
      while read -r match; do
          line=$(echo "$match" | awk {'print $1'} | sed 's/://g')
          #echo "${line}s/${oldVersion}/${newVersion}/g" "$pomFile"
          gsed -i "${line}s/${oldVersion}/${newVersion}/g" "$pomFile"
      done <<< "$(grep -n "$oldVersion" -C 2 "$pomFile" | grep -C 2 "$artifactId" | grep -e "^[0-9]\+:")"
  fi
done



