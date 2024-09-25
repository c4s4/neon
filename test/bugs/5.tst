# Bug 5: when calling super after a call, super target is wrong
# This build must fail with message: "target 'bug-5' not found in parent build files"

default: bug-5

targets:

  test:
    steps:
    - pass:

  bug-5:
    steps:
    - call: test
    - super:
