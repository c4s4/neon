# This build file should succeed

default: assert

properties:
  targets: []

targets:

  target1:
    steps:
    - print: "target 1"
    - 'targets += 1'

  target2:
    depends: "target1"
    steps:
    - print: "target 2"
    - 'targets += 2'

  target3:
    depends: ["target1", "target2"]
    steps:
    - print: "target 3"
    - 'targets += 3'

  assert:
    depends: target3
    steps:
    - if: 'targets == [1, 2, 3]'
      then:
      - print: "Success!"
      else:
      - throw: "Failure! targets = ={targets}"
