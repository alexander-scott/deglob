### Scenario 9

- A combination of a glob and an explicit file list for the same attribute (glob(["file*.h"]) + ["explicit.h"]): the glob is expanded into sub-targets under deps, and the explicit files are preserved as a separate attribute line.
