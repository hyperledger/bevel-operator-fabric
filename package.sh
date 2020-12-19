#!/bin/bash

binary=$(basename "$(dirname "${1}")")
minisign -s "/media/${USER}/kfsoftware/minisign.key" \
         -x "${binary}.minisig" \
         -Sm "${1}" < "/media/${USER}/kfsoftware/minisign-passphrase"

cp -f "${binary}.minisig" "${1}.minisig"
cp -f LICENSE "$(dirname "${1}")"
zip -r -j "${binary}.zip" "$(dirname "${1}")"

if [ -f hlf-operator_linux_amd64.minisig ]; then
    mv hlf-operator_linux_amd64.minisig hlf-operator.minisig
fi