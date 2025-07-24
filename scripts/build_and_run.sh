#!/bin/bash

# help function
show_help() {
  echo "Usage: $0 [-b <build_dir>] <project_dir> [args]"
  echo "   -b <build_dir>  Specify the output build directory (default: ./bin)"
  echo "   <project_dir>   Directory containing the Go project to build"
  echo "   [args]          Additional arguments to pass to the run command"
  echo "Builds the Go project in the specified directory."
  exit 1
}

out_build_dir="./bin"

# Check if -b option is provided
if [[ "$1" == "-b" ]]; then
  shift  # Remove the -b option
  if [ -z "$1" ]; then
    show_help
  fi
  out_build_dir="$1"
  shift  # Remove the build directory argument
fi

# Read the build dirs list from file or show usage
if [ -z "$1" ]; then
  show_help
fi

build_dir="$1"

# Trim ending slash if present
build_dir=${build_dir%/}

if [ -d "$build_dir" ]; then
  EXECUTABLE="$out_build_dir/$(basename "$build_dir")"
  echo "Building Go project in directory: $build_dir"
  go mod download
  if [ $? -ne 0 ]; then
    echo "mod download command failed: $build_dir"
    exit 1
  fi
  # Build the project or exit if it fails
  go build -o "$EXECUTABLE" "$build_dir"/
  if [ $? -ne 0 ]; then
    echo "Build failed for directory: $build_dir"
    exit 1
  fi
  chmod +x "$EXECUTABLE"
  "$EXECUTABLE" "${@:2}"  # Pass additional arguments if any
else
  echo "Directory does not exist: $build_dir"
  exit 1
fi
