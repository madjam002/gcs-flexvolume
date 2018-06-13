#!/bin/bash

set -ex

# Install the gcsfuse binary onto the host
cp $(which gcsfuse) /host/bin

# Install the driver onto the host
mkdir -p /host/volume-plugin-dir/awprice~gcs/
cp $(which gcsfuse-driver) /host/volume-plugin-dir/awprice~gcs/gcs
