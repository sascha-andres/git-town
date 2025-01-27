@smoke
Feature: prepend a branch to a feature branch

  Background:
    Given the current branch is a feature branch "old"
    And the commits
      | BRANCH | LOCATION      | MESSAGE    |
      | old    | local, origin | old commit |
    When I run "git-town prepend parent"

  Scenario: result
    Then it runs the commands
      | BRANCH | COMMAND                        |
      | old    | git fetch --prune --tags       |
      |        | git checkout main              |
      | main   | git rebase origin/main         |
      |        | git checkout old               |
      | old    | git merge --no-edit origin/old |
      |        | git merge --no-edit main       |
      |        | git checkout -b parent main    |
    And the current branch is now "parent"
    And these commits exist now
      | BRANCH | LOCATION      | MESSAGE    |
      | old    | local, origin | old commit |
    And this lineage exists now
      | BRANCH | PARENT |
      | old    | parent |
      | parent | main   |

  Scenario: undo
    When I run "git-town undo"
    Then it runs the commands
      | BRANCH | COMMAND              |
      | parent | git checkout old     |
      | old    | git branch -D parent |
    And the current branch is now "old"
    And the initial commits exist
    And the initial lineage exists
