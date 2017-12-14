function req(url, method="GET") {
    return new Promise(function(resolve, reject) {
        var req = new XMLHttpRequest();
        req.onreadystatechange = function() {
            if(req.readyState == 4) {
                if(req.status >= 200 && req.status < 300) resolve(req.responseText);
                else reject(req);
            }
        }
        req.open(method, url, true);            
        req.send(null);
    });
};

function updateMonkeys() {
    getSeatedMonkeys();
    getBullpen()
};

function getSeatedMonkeys() {
    req('/seats').then(function(response) {
        monkeys = JSON.parse(response)
        guts = "";
        for(var i=0, monkey; monkey = monkeys[i]; i++) {
            guts += `<li><button onclick="stand(` + monkey.Seat + `)">stand</button>Seat ` + monkey.Seat + `: <strong>` + monkey.Name + `</strong> (` + monkey.Speed.toFixed(3) + ` kkps): ` + monkey.Progress;
        }
        document.getElementById("results").innerHTML = guts;

        setTimeout(updateMonkeys, 2000);
    }, function(err) {
        console.log("ERROR, bailing on requests: ", err);
    });
}

function getBullpen() { // monkeys not currently typing
    req('/monkeys').then(function(response) {
        monkeys = JSON.parse(response)
        console.log(monkeys)
        guts = "";
        for(var i=0, monkey; monkey = monkeys[i]; i++) {
            if(!monkey.Seated) guts += `<li><button onclick="sit(` + monkey.ID + `)">sit</button><strong>` + monkey.Name + `</strong>: ` + monkey.Progress;
        }
        document.getElementById("bullpen").innerHTML = guts;
    }, function(err) {
        console.log(`ERROR, bailing on requests: `, err);
    });
}

var stand = function(id) {
    console.log(`Telling monkey ` + id + ` to stand up.`);
    path = "/monkeys/" + id + "/stand"
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

var addSeat = function() {
    console.log("Adding...");
    req('/monkeys', "POST").then(function(response) {
        monkey = JSON.parse(response)
        console.log("MONKEY ADDED!");
        console.log(monkey);
    }, function(err) {
        console.log("Didn't work: ", err);
    });
};

updateMonkeys();