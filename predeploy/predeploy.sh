#!/bin/sh

# Exit on error.
set -e

usage() {
    echo 'Usage: predeploy.sh <serviceName> <namespace> <secretName>'
}

if [ "$#" -ne 3 ]; then
    usage
    exit 1
fi

service=$1
namespace=$2
secret=$3

# Create namespace if not exists.
kubectl create namespace $namespace || true

# Locate the cluster Certificate Authority for population in webhook YAML.
CA_BUNDLE=`kubectl get configmap -n kube-system extension-apiserver-authentication -o=jsonpath='{.data.client-ca-file}' | base64 | tr -d '\n'`

# Populate secrets from certificate file and key.
./generate-secret.sh $service $namespace $secret

# Replace static string with CA_BUNDLE contents.
# Add '' after -i to run this on Mac.
# sed -i '' "s/CA_BUNDLE/$CA_BUNDLE/g" webhook.yaml
sed -i "s/CA_BUNDLE/$CA_BUNDLE/g" webhook.yaml

# Deploy webhook resource.
kubectl apply -f webhook.yaml