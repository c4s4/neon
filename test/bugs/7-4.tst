# This build file should fail: test1 -> test2 -> test1

default: test1

targets:

  test1:
    steps:
    - call: test2
    - print: "test1"

  test2:
    steps:
    - print: "test2"
    - call: test1
