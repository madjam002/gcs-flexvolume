#!/bin/bash

set -ex

# Install the gcsfuse binary onto the host
cp $(which gcsfuse) $1

# Install the driver onto the host
mkdir -p $2/awprice~gcs/
cp $(which gcsfuse-driver) $2/awprice~gcs/
