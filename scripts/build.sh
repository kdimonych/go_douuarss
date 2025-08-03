#!/bin/bash

# help function
show_help() {
  echo "Usage: $0 [-b <build_output_dir>] <project_dir1> <project_dir2> ..."
  echo "   -b <build_output_dir>               Specify the output build directory (default: ./bin)"
  echo "   <project_dir1> <project_dir2> ...   Directories containing Go projects to build"
  echo "Builds Go projects in the specified directories."
  exit 1
}

# Read the build dirs list from file or show usage
if [ -z "$1" ]; then
  show_help
fi

build_output_dir="./bin"

# Check if -b option is provided
if [[ "$1" == "-b" ]]; then
  shift  # Remove the -b option
  if [ -z "$1" ]; then
    show_help
  fi
  build_output_dir="$1"
  shift  # Remove the build directory argument
fi

project_dirs=("$@")

for dir in "${project_dirs[@]}"; do
  # Trim ending slash if present
  dir=${dir%/}
  if [ -d "$dir" ]; then
    echo "Building Go project in directory: $dir"
    go mod download
    go build -o "$build_output_dir/$(basename "$dir")" "$dir"/
  fi
done
