- setup
```
$ cd github-network-data-server
$ go mod tidy
```

- init database
```
$ mysql -u root < ./db/init.sql
```

- run server
```
$ go run .
or
$ go build main.go && ./main
```