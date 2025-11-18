#!/bin/bash

set -e

echo "Attempting to initialize the CockroachDB cluster..."

cockroach init --host=roach1:26357 --insecure || \
(
  echo "Initial 'init' failed. Checking if the cluster was already initialized..."

  cockroach init --host=roach1:26357 --insecure 2>&1 | \
  grep -q 'cluster has already been initialized'

  echo "Cluster is already initialized. Continuing..."
)