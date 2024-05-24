# Bazel Build Events

## JSON file

Generate a JSON file from the Bazel Build Event protocol. For this example we will run a test command (it also builds) and run multiple times each test.

```sh
bazel test //go/... --runs_per_test_detects_flakes -t- --runs_per_test=2 --build_event_json_file=blaze/bazelbep/readerjson/build_events.json
```
