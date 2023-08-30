# Modal

Modals are simply a dimension to a world in which components can return values without any setters, as simple as

```go

name := c.modal.Show(&NameSelection{})

// or

task := c.modal.ShowAsync(&NameSelection{})
// Compute something
name := task.Result()

```
