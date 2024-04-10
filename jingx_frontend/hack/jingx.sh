#!/usr/bin/env bash

GENS="$1"
PROJECT="$2"
REPOSITORY="$3"
TAG="$4"

addr=""
user=""
password=""
gitInfo=""
branch=""
commitHash=""
image=""

if [ "${GENS}" = "cli" ] || grep -qw "address" <<<"${GENS}"; then
    addr=$(env | grep JINGX | grep ADDRESS | awk '{split($0,a,"="); print a[2]}')
#  echo ${addr}
fi

if [ "${GENS}" = "cli" ] || grep -qw "username" <<<"${GENS}"; then
    user=$(env | grep JINGX | grep USER | awk '{split($0,a,"="); print a[2]}')
#  echo ${user}
fi

if [ "${GENS}" = "cli" ] || grep -qw "password" <<<"${GENS}"; then
    password=$(env | grep JINGX | grep PASS | awk '{split($0,a,"="); print a[2]}')
#  echo ${password}
fi

if [ "${GENS}" = "cli" ] || grep -qw "git" <<<"${GENS}"; then
    gitInfo=$(cat .git/config | grep url | awk '{split($0,a,"/"); print a[2]}')
#  echo ${gitInfo}
fi

if [ -z "$gitInfo" ]; then
    gitInfo=$(cat .git/config | grep url | awk '{split($0,a,"/"); print a[5]}')
fi

if [ "${GENS}" = "cli" ] || grep -qw "branch" <<<"${GENS}"; then
    branch=$(env | grep -i CI_COMMIT_BRANCH | awk '{split($0,a,"="); print a[2]}')
#  echo ${branch}
fi

if [ -z "$branch" ]; then
    branch=$(git branch | grep \* | awk '{split($0,a," "); print a[2]}')
fi

if [ "${GENS}" = "cli" ] || grep -qw "commitHash" <<<"${GENS}"; then
    commitHash=$(git rev-parse HEAD)
fi

if [ -z "$commitHash" ]; then
    commitHash=$(env | grep -i CI_COMMIT_SHA | awk '{split($0,a,"="); print a[2]}')
fi

if [ "${GENS}" = "cli" ]; then
    isCi=$(env | grep -i CI_COMMIT_BRANCH | awk '{split($0,a,"="); print a[2]}')
    if [ -z "${isCi}" ]; then
        image=$(docker images --no-trunc | grep ${PROJECT} | grep ${REPOSITORY} | grep ${TAG} | awk '{print $3}' | awk '{split($0,a,":"); print a[2]}')
        jingxcli \
            -address=${addr} \
            -username=${user} \
            -password=${password} \
            -project=${PROJECT} \
            -repository=${REPOSITORY} \
            -tag=${TAG} \
            -git=${gitInfo} \
            -branch=${branch} \
            -commitHash=${commitHash} \
            -sha256=${image}
    else
        image=$(sudo docker images --no-trunc | grep ${PROJECT} | grep ${REPOSITORY} | grep ${TAG} | awk '{print $3}' | awk '{split($0,a,":"); print a[2]}')
        sudo jingxcli \
            -address=${addr} \
            -username=${user} \
            -password=${password} \
            -project=${PROJECT} \
            -repository=${REPOSITORY} \
            -tag=${TAG} \
            -git=${gitInfo} \
            -branch=${branch} \
            -commitHash=${commitHash} \
            -sha256=${image}
    fi
fi
