# tasktea

TUI for `taskwarrior` using `bubbletea`, inspired by `taskwiki`.

Set filters to view tasks as trees with dependencies.

## Design issues

- How to keep reference to tree, and still use it properly with
  the update loop?
- ViewFilter, ViewSave are not really on the same level as ViewTasks, etc.
- Documentation of custom messages
- Filter does separate things for tasks vs context, etc.
- Maybe use render/actions as callbacks/cmds

## Solution

- Study [`list` bubble](/home/x/fd/code/go/references/bubbles/list)
- create `tree` bubble with navbar and statusbar
