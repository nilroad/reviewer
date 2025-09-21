#!/usr/bin/env bash
set -Eeuo pipefail

export CMD=$1
export BOOTSTRAP_FILE=${BOOTSTRAP_FILE:-bootstrap.sh}
export APP_ENV="${APP_ENV:-local}"
export APP_NAME="${APP_NAME:-ozone}"

secret_file_path="/vault/secrets/secrets"
if [ -e "$secret_file_path" ]; then
    source $secret_file_path
    echo "source $secret_file_path" >> /etc/bash.bashrc
fi
bash ./deploy/docker/general/scripts/persist_env.sh


################################################################
#  Trying to test connection to all related stateful services  #
#  so we can sure all related services are up and running      #
################################################################
/app/cmd healthcheck


case $APP_ENV in

  production)
    DOCKER_DIR_FOR_ENV="production"
    ;;

  staging)
    DOCKER_DIR_FOR_ENV="staging"
    ;;

  testing)
    DOCKER_DIR_FOR_ENV="testing"
    ;;

  develop)
    DOCKER_DIR_FOR_ENV="develop"
    ;;

  *)
    DOCKER_DIR_FOR_ENV="local"
    ;;
esac

/app/cmd version
/app/cmd mysql:migrate
/app/cmd rabbitmq:setup
/app/cmd mysql:seed --app-name $APP_NAME --app-env $APP_ENV
echo "Preparing to run $CMD"
if [ ! -f ./deploy/docker/$DOCKER_DIR_FOR_ENV/$CMD/$BOOTSTRAP_FILE ]; then
    echo "No bootstrap file found for ./deploy/docker/$DOCKER_DIR_FOR_ENV/$CMD/$BOOTSTRAP_FILE "
    exit 1
else
    echo "Found bootstrap file for $CMD"
    bash ./deploy/docker/$DOCKER_DIR_FOR_ENV/$CMD/$BOOTSTRAP_FILE
fi