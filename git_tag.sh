#!/bin/bash

set -e

# get latest Git tag
latest_tag=$(git tag -l "v*" | sort -V | tail -n1)
if [ -z "$latest_tag" ]; then
  latest_tag="v0.0.0"
fi
echo "Latest tag: $latest_tag"

# get latest GitHub Release version
latest_release=$(gh release view --json tagName -q .tagName 2>/dev/null || echo "v0.0.0")
echo "Latest release: $latest_release"

# calculate final version
final_version=$(printf '%s\n' "$latest_tag" "$latest_release" | sort -V | tail -n1)
echo "Final version: $final_version"

# calculate next version
IFS='.' read -r base major minor <<< "${final_version//v/}"

# Increment minor version
minor=$((minor + 1))

# Check if minor exceeds 999
if ((minor > 999)); then
    minor=0
    major=$((major + 1))
fi

# Check if major exceeds 999
if ((major > 999)); then
    major=0
    base=$((base + 1))
fi

next_version="v${base}.${major}.${minor}"
echo "Next version: $next_version"

# create new tag
git tag $next_version -m "Tag for release $next_version"
git push origin $next_version

# update latest tag
if git rev-parse latest >/dev/null 2>&1; then
    git tag -d latest
    git push origin :refs/tags/latest
fi

git tag latest -m "Update latest tag to $next_version"
git push origin latest
