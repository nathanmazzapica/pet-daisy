## Refactoring Time

I need to refactor the entire server to use the new dependency injection model.

```go
type Server struct {
    Store   *db.UserStore
    Hub     *Hub
}
```
The server will be moved into a struct that accepts a db.UserStore.

This is also a good opportunity to generally clean up this folder. It's pretty messy right now

### Questions to ponder...

1. Q: would hub work as the server struct?
   A: no, I don't think so. Server should have hub in it.
2. ...