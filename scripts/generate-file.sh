#!/bin/bash

targets=$(bazel query 'attr(additional_update_targets, "\[\]", kind(_write_source_file, //...))')

for target in $targets; do
  bazel run "$target"
done