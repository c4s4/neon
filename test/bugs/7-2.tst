# This build file should fail: test1 -> test2 -> test1 (infinite loop)

default: test1

targets:

  test1:
    steps:
    - print: "test1"
    - call: test2

  test2:
    steps:
    - print: "test2"
    - call: test1
