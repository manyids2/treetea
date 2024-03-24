---
title: "Setup"
---

## Fist time

Set some location:

```bash
mkdir /tmp/task
export TASKDATA=/tmp/task
export TASKRC=/tmp/task/taskrc
```

Run taskwarrior for first time.
Prompted for initializing config file.
Note that warnings are printed to `stderr` (in yellow as opposed to white).

```bash
task
```

```stdout
A configuration file could not be found in /home/x

Would you like a sample /tmp/task/taskrc created, so Taskwarrior can proceed? (yes/no)
```

```stderr
TASKRC override: /tmp/task/taskrc
TASKDATA override: /tmp/task
No matches.
Recently upgraded to 2.6.0. Please run 'task news' to read highlights about the new release.
```

Check created file:

```bash
$EDITOR $TASKRC
```

```rc
# [Created by task 2.6.2 3/24/2024 22:13:59]
# Taskwarrior program configuration file.
# For more documentation, see https://taskwarrior.org or try 'man task', 'man task-color',
# 'man task-sync' or 'man taskrc'

# Here is an example of entries that use the default, override and blank values
#   variable=foo   -- By specifying a value, this overrides the default
#   variable=      -- By specifying no value, this means no default
#   #variable=foo  -- By commenting out the line, or deleting it, this uses the default

# You can also refence environment variables:
#   variable=$HOME/task
#   variable=$VALUE

# Use the command 'task show' to see all defaults and overrides

# Files
data.location=/tmp/task

# To use the default location of the XDG directories,
# move this configuration file from ~/.taskrc to ~/.config/task/taskrc and uncomment below

#data.location=~/.local/share/task
#hooks.location=~/.config/task/hooks

# Color theme (uncomment one to use)
#include light-16.theme
#include light-256.theme
#include dark-16.theme
#include dark-256.theme
#include dark-red-256.theme
#include dark-green-256.theme
#include dark-blue-256.theme
#include dark-violets-256.theme
#include dark-yellow-green.theme
#include dark-gray-256.theme
#include dark-gray-blue-256.theme
#include solarized-dark-256.theme
#include solarized-light-256.theme
#include no-color.theme
```

Check tasks, there should be none:

```bash
task
```

```stdout

```

```stderr
TASKRC override: /tmp/task/taskrc
TASKDATA override: /tmp/task
No matches.
```

Add line to taskrc to not show override message:

```rc
verbose=blank,footnote,label,new-id,affected,edit,special,project,sync,unwait
```

Next add a task and check (`--` used to preserve spaces and avoid taskwarrior grammar),
usual warnings on stderr:

```bash
task add -- task a
```

```stdout
Created task 1.
```

task

````

```stdout
ID Age  Description Urg
 1 2s   task a         0

1 task
````

Since id is more of a human readable id, and not a primary key,
we want to print the UUID instead when a new task is created.

So in the verbose, add `new-uuid` for an actual primary key. 

NOTE: specifying both `new-id` and `new-uuid` will print UUID, order does not matter.

```rc
verbose=blank,footnote,label,new-id,new-uuid,affected,edit,special,project,sync,unwait
```

NOTE: better to use `task _get [id].uuid`, so humans can use it as well.

NOTE : Projects can be renamed. ( works by renaming projects of all tasks basically, so covered )
