### Scenario 8

- Globs in two separate attributes (srcs and hdrs) of the same target: the first glob line is replaced with a deps attribute listing all sub-targets from both globs, and the second glob line is removed.
