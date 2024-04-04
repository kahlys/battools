#!/bin/bash
set -e

bazel run @rules_go//go -- mod tidy
bazel mod tidy