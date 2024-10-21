#!/bin/bash

GITLAB_PAT="${1}"
GITLAB_USERNAME="${2}"

# Variables
GITLAB_API_URL="https://gitlab.alpinelinux.org/api/v4"
FORK_PROJECT_PATH="${GITLAB_USERNAME}/aports"
ORIGINAL_PROJECT_PATH="alpine/aports"
FORKED_PROJECT_URL="https://gitlab.alpinelinux.org/${FORKED_PROJECT_PATH}.git"
ORIGINAL_PROJECT_URL="https://gitlab.alpinelinux.org/${ORIGINAL_PROJECT_PATH}.git"

CLONE_DIR="aports"
APKBUILD_PATH="testing/tenv/APKBUILD"
NEW_PKGVER="3.0.0"
NEW_SHA512SUM="NEW_SHA512SUM_VALUE"

# Step 1: Check if fork exists, and delete if it does
EXISTING_FORK=$(curl -s --header "PRIVATE-TOKEN: ${GITLAB_PAT}" "${GITLAB_API_URL}/projects?search=aports" | jq -r ".[] | select(.path_with_namespace == \"${FORK_PROJECT_PATH}\") | .id")

if [ -n "${EXISTING_FORK}" ]; then
    echo "Fork found. Deleting..."
    curl --request DELETE --header "PRIVATE-TOKEN: ${GITLAB_PAT}" "${GITLAB_API_URL}/projects/${EXISTING_FORK}"
else
    echo "No fork found."
fi

# Step 2: Create a new fork of the original project
echo "Forking the repository..."
curl --request POST --header "PRIVATE-TOKEN: ${GITLAB_PAT}" "${GITLAB_API_URL}/projects/$(echo ${ORIGINAL_PROJECT_PATH} | sed 's/\//%2F/g')/fork" > /dev/null

# Wait for fork to be created
sleep 5

# Step 3: Clone the forked repository
echo "Cloning the forked repository..."
git clone "$FORKED_PROJECT_URL"
cd "$CLONE_DIR"

# Step 4: Modify the APKBUILD file
echo "Modifying the APKBUILD file..."
sed -i "s/pkgver=.*/pkgver=$NEW_PKGVER/" $APKBUILD_PATH
sed -i "s/sha512sums=.*/sha512sums=\"$NEW_SHA512SUM\"/" $APKBUILD_PATH

# Step 5: Git add, commit, and push
echo "Adding, committing, and pushing changes..."
git checkout -b update-tenv-pkgver
git add $APKBUILD_PATH
git commit -m "Update tenv to pkgver $NEW_PKGVER"
git push origin update-tenv-pkgver

# Step 6: Open a merge request from fork to original repo
echo "Opening a merge request..."
curl --request POST --header "PRIVATE-TOKEN: ${GITLAB_PAT}" \
    --data "source_branch=update-tenv-pkgver" \
    --data "target_branch=master" \
    --data "title=Update tenv to $NEW_PKGVER" \
    --data "target_project_id=$(curl -s --header "PRIVATE-TOKEN: $GITLAB_PERSONAL_TOKEN" "$GITLAB_API_URL/projects?search=aports" | jq -r ".[] | select(.path_with_namespace == \"$ORIGINAL_PROJECT_PATH\") | .id")" \
    "$GITLAB_API_URL/projects/$EXISTING_FORK/merge_requests"

echo "Merge request created!"
