#!/bin/bash
read -p "Enter the EPrints hostname: " EP_HOST
read -p "Enter the EPrints username: " EP_USER
read -sp "Enter the EPrints password: " EP_PASSWD
read -p "Enter sample size: " SAMPLE_SIZE
eputil -auth=basic -json "https://${EP_USER}:${EP_PASSWD}@${EP_HOST}/rest/eprint/" >eprint_ids.json
cat eprint_ids.json | jsonrange -sample=${SAMPLE_SIZE} -values > eprint_ids.keys
for KEY in $(cat eprint_ids.keys); do
  echo "https://${EP_HOST}/${KEY}"
done
