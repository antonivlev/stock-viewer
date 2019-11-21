updateSearchTable();

document.querySelector("#send-symbol").onclick = () => {
    let stock = document.querySelector("#symbol-input").value
    
    console.log("sending request")

    fetch("http://localhost:3000/api/get-stock-data?symbol="+stock)
        .then(res => {
            console.log("got response")
            // TODO: need better way of checking if response contains an error
            if (res.headers.get("Content-Type") === "application/json") {
                // all good
                res.json()
                    .then(jsonData => {
                        console.log("parsed json");
                        plotData(jsonData);
                        let now = new Date().toJSON();
                        saveSearch(now, stock).then(updateSearchTable);
                    })
                    .catch(jsonError => showError("json parse error: \n"+jsonError));
            } else {
                // server returned error
                res.text().then(parsedText => showError(parsedText))
            }
        })
}

function updateSearchTable() {
    fetch("http://localhost:3000/api/get-searches").then(res => {
        // TODO: need better way of checking if response contains an error
        if (res.headers.get("Content-Type") === "application/json") {
            // all good
            res.json()
                .then(jsonData => fillInTable(jsonData))
                .catch(jsonError => showError("json parse error: \n"+jsonError));
        } else {
            // server returned error
            res.text().then(parsedText => showError(parsedText))
        }
    })
}

function fillInTable(searchList) {
    document.querySelector("tbody").innerHTML = "";
    searchList.map((search) => {
        let {SearchTime, Stock} = search;
        Stock = Stock.toUpperCase()
        document.querySelector("tbody").insertAdjacentHTML(
            "afterbegin", 
            `<tr class="stock-row">
                <td>${SearchTime}</td>
                <td class="stock-name">${Stock}</td>
                <td class="open"></td>
                <td class="high"></td>
                <td class="low"></td>
                <td class="close"></td>
                <td class="volume"></td>
            </tr>`
        );
    })
    updateCurrentStockValues();
}

// Updates on new searches. Server return empty if frequency too high,
// needs looking into
function updateCurrentStockValues() {
    document.querySelectorAll(".stock-row").forEach(row => {
        let stock = row.querySelector(".stock-name").innerText;
        // fetch and fill in values for this row
        fetch("http://localhost:3000/api/get-latest-stock-data?symbol="+stock).then(res => {
            // TODO: need better way of checking if response contains an error
            if (res.headers.get("Content-Type") === "application/json") {
                // all good
                res.json()
                    .then(jsonData => fillInStockRow(row, jsonData))
                    .catch(jsonError => showError("json parse error: \n"+jsonError));
            } else {
                // server returned error
                res.text().then(parsedText => showError(parsedText))
            }
        })
    })
}

function fillInStockRow(row, stockData) {
    row.querySelector(".open").innerText = stockData["02. open"];
    row.querySelector(".high").innerText = stockData["03. high"];
    row.querySelector(".low").innerText = stockData["04. low"];
    row.querySelector(".close").innerText = stockData["08. previous close"];
    row.querySelector(".volume").innerText = stockData["06. volume"];    
}

// TODO handle errors
function saveSearch(time, stock) {
    return fetch("http://localhost:3000/api/save-search", 
        {
            method: "POST", 
            body: JSON.stringify({
                "searchTime": time,
                "stock": stock,
            })
        }
    )
}

function plotData(datesMap) {
    let [x, open, high, low, close] = [[], [], [], [], []];
    Object.keys(datesMap).map(date => {
        x.push(date);
        open.push(datesMap[date]["1. open"]);
        high.push(datesMap[date]["2. high"]);
        low.push(datesMap[date]["3. low"]);
        close.push(datesMap[date]["4. close"]);
    });
    console.log("arranged data")

    let trace1 = {  
        // data
        x: x,
        open: open,
        high: high,
        low: low,
        close: close,     
        // config
        decreasing: {line: {color: '#7F7F7F'}},
        increasing: {line: {color: '#17BECF'}}, 
        line: {color: 'rgba(31,119,180,1)'}, 
        type: 'candlestick', 
        xaxis: 'x', 
        yaxis: 'y'
    };
      
    let data = [trace1];
      
    let layout = {
        dragmode: 'zoom', 
        margin: {r: 0, t: 0, b: 0, l: 30}, 
        showlegend: false, 
        xaxis: {
            autorange: true, 
            domain: [0, 1], 
            title: 'Date', 
            type: 'date'
        }, 
        yaxis: {
            type: 'linear'
        }
    };
    
    Plotly.newPlot('plotly-div', data, layout, {responsive: true});
    console.log("plotted!");
}

function showError(errorString) {
    let errorDiv = document.querySelector(".error-msg");
    errorDiv.innerText = errorString;
}

