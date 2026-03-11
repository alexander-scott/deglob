### Scenario 7

- Multiple glob calls combined on the same attribute line (glob(["file_a*.h"]) + glob(["file_b*.h"])): all matched files from every glob call are collected and each gets its own sub-target listed under a single deps attribute.
