#!/bin/bash
set -e


mkdir -p blaze/hello
echo 'package hello' > blaze/hello/hello.go

bazel run @rules_go//go -- mod tidy
bazel mod tidy

rm -rf blaze/hello