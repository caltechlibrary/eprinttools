#!/bin/bash
echo "Reloading lemurprints schema"
mysql lemurprints < srctest/lemurprints-setup-schema.sql
