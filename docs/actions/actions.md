# Actions

Indexing edge cases.

- taskrc override is set / not set: solved (stderr)
- taskrc verbose uuids is set / not set: solved (_get)
- stderr vs stdout (probably the yellow)
- dependencies, updating of urgency, due
- conflicting filters and context
- grammar of filters
- virtual tags
- deleted tasks

- `task undo` exists!!

## Read

- 󰊕 List

- 󰊕 Active

- 󰊕 Contexts

- 󰊕 Context

- 󰊕 Projects

## Write

- 󰊕 SetContext
- 󰊕 SetStatus
- 󰊕 SetActive
- 󰊕 SetParent
- 󰊕 SetGrandparentAsParent
- 󰊕 Add
    - No need to mess with rc file, we can use `task _get [id].uuid`
- 󰊕 ModifyDescription
- 󰊕 Modify
- 󰊕 ModifyBatch
- 󰊕 Delete

