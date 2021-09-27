#!/usr/bin/env bash
current=$(pwd)
proj=project
projpath=$current/$proj

vipconfpath="$projpath/vip.conf.d"
migrate_path="$projpath/migrate"
phoenixpath="$projpath/phoenix"
migrationspath="$projpath/migrations"
penguinpath="$projpath/penguin"
swanpath="$projpath/swan"
turkeypath="$projpath/turkey"

mkdir -p "$vipconfpath"
mkdir -p "$phoenixpath"
mkdir -p "$migrationspath"
mkdir -p "$penguinpath"
mkdir -p "$swanpath"
mkdir -p "$turkeypath"

export GOROOT=/usr/local/go
export GOPATH=$projpath
export GOPROXY=https://goproxy.cn,direct
export PATH=$GOROOT/bin:$PATH
export vip_path=$vipconfpath
export migrate_path=$migrate_path

#cp -rf $current/conf/*.json $vipconfpath
