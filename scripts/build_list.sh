#!/bin/bash

# help function
show_help() {
  echo "Usage: $0 [-b <build_output_dir>] <path_to_targets_file>"
  echo "   -b <build_output_dir>   Specify the output build directory (default: ./bin)"
  echo "   <path_to_targets_file>  File containing directories to build"
  echo "Builds Go projects in the specified directories."
  exit 1
}

build_output_dir="./bin"

# Check if -b option is provided
if [[ "$1" == "-b" ]]; then
  shift  # Remove the -b option
  if [ -z "$1" ]; then
    show_help
  fi
  build_output_dir="$1"
  shift  # Remove the build_output_dir argument
fi

# Read the build dirs list from provided targets.txt file or show usage
if [ -z "$1" ]; then
  show_help
fi

project_dirs=()
while IFS= read -r line; do
  project_dirs+=("$line")
done < "$1"

# Obtain this script's directory
script_dir=$(dirname "$0")

for dir in "${project_dirs[@]}"; do
  "$script_dir"/build.sh -b "$build_output_dir" "$dir"
done
