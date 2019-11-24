# Stock viewer

Requires Docker (or postgres)

To run:
1. Set up database (connection details are in ```config.json``` as well):
    ```bash
    $ docker run -d --name stock-app-container -p 54320:5432 postgres:11
    $ docker exec -it stock-app-container psql -U postgres -c "create database stockapp"
    ```
2. Get and run the app
    ```bash
    $ git clone https://github.com/antonivlev/stock-viewer.git
    $ cd stock-viewer
    $ ./stock-viewer
    ```

3. (Optional) If you have Go, run in development mode:
    ```bash
    $ go get github.com/antonivlev/stock-viewer
    $ cd ~/go/src/github.com/antonivlev/stock-viewer/
    $ go get
    $ go run main.go
    ```

Server API documentation: https://documenter.getpostman.com/view/6354074/SW7Z4UZo?version=latest

Play around with ```cacheExpirationInterval``` (which is in seconds) in ```config.json```! The server logs any cache saving and deleting.

---------
My approach was as follows:

1. On paper, mock out design, and API; identify data flow between frontend components
2. Implement/loosely test API with Postman
3. Implement frontend


There are a few major things I would have added with more time:

- Testing: in particular package tests for the backend. Simple expected output vs actual output comparisons, in particular for things like database interaction.
- The remaining stretch goals. Clicking in history table to see previous graph is a simple api call, followed by ```plotData```. The Ag-grid looks like it can be implemented with [this library.](https://www.ag-grid.com/javascript-grid/)
- Migrate frontend to a framework (I would choose Svelte). It would "componentise" the frontend and simplify data binding. 

I have left some thoughts on how to improve code readability/reliability as TODOs in comments.