var updateDelay = 1000; // how frequently the UI calls home

function req(url, method="GET", body="") {
    return new Promise(function(resolve, reject) {
        var req = new XMLHttpRequest();
        req.onreadystatechange = function() {
            if(req.readyState == 4) {
                if(req.status >= 200 && req.status < 300) resolve(req.responseText);
                else reject(req);
            }
        }
        req.open(method, url, true);            
        req.send(body);
    });
};

function updateMonkeys() {
    getSeatedMonkeys();
    getBullpen()
};

function getSeatedMonkeys() {
    req('/seats').then(function(response) {
        monkeys = JSON.parse(response)
        console.log(monkeys);
        guts = ``;
        for(var i=0, monkey; monkey = monkeys[i]; i++) {
            guts += `<tr><td>` + monkey.Seat + `<td> ` + monkey.Keyboard + `<td>`;
            if(monkey.Name) {
                guts += `<button class="btn btn-warning" onclick="stand(` + monkey.Seat + `)">` + monkey.Name + `</button>`;
            }
            
            guts += `</strong><td>`;
            
            if(monkey.Name) {
                guts += monkey.Speed.toFixed(3);
            }
            guts += `<td>` + monkey.Progress;
        }
        guts += "</table>"
        document.getElementById("results").innerHTML = guts;

        setTimeout(updateMonkeys, updateDelay);
    }, function(err) {
        console.log("ERROR, bailing on requests: ", err);
    });
}

function getBullpen() { // monkeys not currently typing
    req('/monkeys').then(function(response) {
        monkeys = JSON.parse(response)
        guts = "";
        for(var i=0, monkey; monkey = monkeys[i]; i++) {
            if(!monkey.Seated) guts += `<li><button class="btn btn-success" onclick="sit(` + monkey.ID + `)">` +  monkey.Name + `</button>`;
        }
        document.getElementById("bullpen").innerHTML = guts;
    }, function(err) {
        console.log(`ERROR, bailing on requests: `, err);
    });
}

var stand = function(id) {
    console.log(`Telling monkey in seat ` + id + ` to stand up.`);
    path = "/seats/" + id + "/stand"
    req(path, "PATCH").then(function(response) {
        console.log("MONKEY STOOD!");
    }, function(err) {
        console.log("Didn't work: ", err);
    });
};

var sit = function(id) {
    console.log(`Telling monkey ` + id + ` to sit.`);
    path = "/monkeys/" + id + "/sit"
    req(path, "PATCH").then(function(response) {
        console.log("MONKEY SAT!");
    }, function(err) {
        console.log("Didn't work: ", err);
    });
};

var addMonkey = function() {
    console.log("Adding monkey...");
    req('/monkeys', "POST").then(function(response) {
        monkey = JSON.parse(response)
        console.log("MONKEY ADDED!");
        console.log(monkey);
    }, function(err) {
        console.log("Didn't work: ", err);
    });
};

var addSeat = function(layout) {
    console.log("Adding seat...");
    req('/seats', "POST", `{"Layout": "` + layout + `"}`).then(function(response) {
        monkey = JSON.parse(response)
        console.log("SEAT ADDED!");
        console.log(monkey);
    }, function(err) {
        console.log("Didn't work: ", err);
    });
};

updateMonkeys();