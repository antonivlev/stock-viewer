# Stock viewer

Requires Go and Docker

To run:
1. Set up database (connection details are in config.json as well):
    ```bash
    $ docker run -d --name stock-app-container -p 54320:5432 postgres:11
    $ docker exec -it stock-app-container psql -U postgres -c "create database stockapp"
    ```
2. Get and run the app:
    ```bash
    $ go get github.com/antonivlev/stock-viewer
    $ cd ~/go/src/github.com/antonivlev/stock-viewer/
    $ go get
    $ go run main.go
    ```

Server API documentation: https://documenter.getpostman.com/view/6354074/SW7Z4UZo?version=latest
