
#!/bin/bash

# This may be out of date, check at https://github.com/docker/compose/releases

DOCKER_PROJECT_USER=docker
DOCKER_COMPOSE_REPO=compose
DOCKER_COMPOSE_APIURL=https://api.github.com/repos/${DOCKER_PROJECT_USER}/${DOCKER_COMPOSE_REPO}/releases/latest
DOWNLOAD_VERSION=`curl -is ${DOCKER_COMPOSE_APIURL} | tr '\n' ' ' | sed 's/.*tag_name" *: *"\([0-9.]*\).*/\1/'`
[[ ${DOWNLOAD_VERSION} =~ ^[0-9.][0-9.]*$ ]] || DOWNLOAD_VERSION=
echo $DOWNLOAD_VERSION

DOCKER_COMPOSE_BASEURL=https://github.com/${DOCKER_PROJECT_USER}/${DOCKER_COMPOSE_REPO}/releases/download
DOCKER_COMPOSE_VERSION=${DOWNLOAD_VERSION:-1.10.0}
DOCKER_COMPOSE_BASENAME=docker-compose-`uname -s`-`uname -m`

TARGET_BIN=/usr/local/bin
TARGET_FILE=${TARGET_BIN}/docker-compose

# Be sure you have curl installed.
curl --version || __NO_CURL=1

if [ -z "${__NO_CURL}" ]; then
  echo "Curl found. Installing docker-compose."
else
  echo "Error: Can't find 'curl'" >&2
  echo "Installing Curl..."
  apt install curl
fi

curl -L ${DOCKER_COMPOSE_BASEURL}/${DOCKER_COMPOSE_VERSION}/${DOCKER_COMPOSE_BASENAME} -o $TARGET_FILE && \
chmod +x ${TARGET_FILE} && \
( echo -n "Installed " && ${TARGET_FILE} --version ) || __NOT_INSTALLED=1

if [ ! -z "${__NOT_INSTALLED}" ]; then
  echo "Error: failed to properly or fully install docker-compose" >&2
  echo "Failed to install, either the URL is wrong or you don't have permission to write/create ${TARGET_FILE}"
  exit 3
else
# run your container
apt install golang
git clone github.com/amku91/lam
docker-compose build && docker-compose up
fi