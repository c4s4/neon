# This build file should succeed and print: test1 test1

default: test2

targets:

  test1:
    steps:
    - print: "test1"

  test2:
    depends: test1
    steps:
    - call: test1
