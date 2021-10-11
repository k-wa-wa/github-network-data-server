- setup
```
$ cd github-network-server
$ go mod tidy
```

- init database
```
$ mysql -u root < ./db/init.sql
```

- run server
```
$ go run .
```