#!/usr/bin/env bash

set -o nounset

# First, check for git in $PATH
hash git 2>/dev/null || {
  echo >&2 "Git required, not installed. Aborting."
  exit 1
}

[[ -d ".git" ]] || {
  echo >&2 "Git repo not found. Aborting."
  exit 1
}

function current_vtag() {
  local tag
  tag=$(git describe --tags --match 'v*' --abbrev=0 2>/dev/null)
  echo "$tag"
}

function short_version() {
  local vtag
  vtag=$(current_vtag)

  local n_commits

  if [[ "x$vtag" == "x" ]]; then
    n_commits=$(git rev-list --count HEAD 2>/dev/null)
    n_commits=$((++n_commits))
    echo "0.0.0.$n_commits"
  else
    n_commits=$(git rev-list --count "$vtag".. 2>/dev/null)
    local version=${vtag##v}

    OLD_IFS=$IFS && IFS="."

    local version_parts
    local last_part
    local version_parts

    version_parts="$version"
    last_part=$((${#version_parts[@]} - 1))
    version_parts["$last_part"]=$((${version_parts["$last_part"]} + 1))
    version="${version_parts[*]}.$n_commits"
    IFS=${OLD_IFS}

    echo "$version"
  fi
}

function commit_hash() {
  local commit_hash
  local vtag

  commit_hash=$(git describe --always 2>/dev/null)
  vtag=$(current_vtag)
  if [[ "x$vtag" = "x" ]]; then
      commit_hash=$(git show-ref --abbrev -s "${vtag}")
  fi
  echo "$commit_hash" | tr [a-z] [A-Z]
}
