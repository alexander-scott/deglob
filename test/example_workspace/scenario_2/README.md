### Scenario 2

- Exclusive glob for a single srcs/hdrs attribute on multiple cc_library targets that each match multiple cpp files
- Includes a glob in srcs in one target and a glob in hdrs in a different target
- The two targets with globs are at the top and bottom of the BUILD file, with a normal target being in the middle
- Incorporates edge case with glob target being at the very end of the file
