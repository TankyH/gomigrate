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

git clone git@gitlab.xinghuolive.com:birds-backend/migrations.git "$migrationspath"
git clone git@gitlab.xinghuolive.com:birds-backend/phoenix.git "$phoenixpath"
git clone git@gitlab.xinghuolive.com:birds-backend/swan.git "$swanpath"
git clone git@gitlab.xinghuolive.com:birds-backend/penguin.git "$penguinpath"
git clone git@gitlab.xinghuolive.com:birds-backend/turkey.git "$turkeypath"
git clone git@gitlab.xinghuolive.com:birds-backend-limit/vip.conf.d.git "$vipconfpath"

#cp -rf $current/conf/*.json $vipconfpath
