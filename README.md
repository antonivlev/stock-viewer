# Stock viewer

Requires Go and Postgres

To run:
1. Create a postgres database with the following details:
```host=localhost port=5432 user=postgres dbname=stockapp password=12345```
2. ```bash
    $ go get github.com/antonivlev/stock-viewer
    $ cd ~/go/src/github.com/antonivlev/stock-viewer/
    $ go get
    $ go run main.go
    ```

Server API documentation: https://documenter.getpostman.com/view/6354074/SW7Z4UZo?version=latest
