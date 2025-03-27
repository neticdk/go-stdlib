#!/usr/bin/env bash

set -eo pipefail

SCRIPT_DIR=$(cd -- "$(dirname -- "${BASH_SOURCE[0]}")" &>/dev/null && pwd)

GIT_URL=$(git remote get-url origin | sed -r 's|git@([^:]*):([^/]*)/([^/]*)(.git)?|https://\1/\2/\3|' )
echo "Using git url: $GIT_URL"

# GitHub org
export OWNER=$(echo $GIT_URL | cut -d/ -f4)
echo "Using owner: $OWNER"

# GitHub repo
export REPO=$(echo $GIT_URL | cut -d/ -f5 | cut -d. -f1)
echo "Using repo: $REPO"

echo "Setting up project..."
envsubst <${SCRIPT_DIR}/../README.md | sed -e '/# '$REPO'/,$!d' >README.tmp && mv README.tmp ${SCRIPT_DIR}/../README.md

echo "Self destructing ..."
git rm ${SCRIPT_DIR}/setup-project.sh

echo "Committing changes..."
git commit -m "chore: setup project" ${SCRIPT_DIR}/..

echo "Creating some directories you might need..."
mkdir -p ${SCRIPT_DIR}/../{bin,cmd,pkg,internal}

tree ${SCRIPT_DIR}/..

echo
echo "Remember to:"
echo
echo "   - go mod init"
echo "   - update .github/workflows/* with the correct workflow triggers"
echo
