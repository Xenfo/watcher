# Watcher

## TODO

- Frontend
- Backend
- Logging to a file
- Rotating log files
- Graceful shutdown
  - <https://medium.com/tokopedia-engineering/gracefully-shutdown-your-go-application-9e7d5c73b5ac>
  - <https://www.rudderstack.com/blog/implementing-graceful-shutdown-in-go/>

## Usage

Create a `config.json` file and add the packages you want to watch. Then run `go run main.go`.

```json
// config.json

{
  "packages": {
    "example-package": {
      "notify": false, // unimplemented
      "currentVersion": "1.0.0",
      "targetVersion": "", // optional, if set, Watcher will check for the target version instead of the latest version
      "notes": "example note", // optional, if set, Watcher will add a note to the notification
      "includeBetas": false
    }
  }
}
```
