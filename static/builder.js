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
    req('/monkeys').then(function(response) {
        monkeys = JSON.parse(response)
        guts = "";
        for(var i=0, monkey; monkey = monkeys[i]; i++) {
            guts += "<li><strong>" + monkey.Name + "</strong> (" + monkey.Speed.toFixed(3) + " kkps): " + monkey.Progress;
        }
        document.getElementById("results").innerHTML = guts;

        setTimeout(updateMonkeys, 2000);
    }, function(err) {
        console.log("ERROR, bailing on requests: ", err);
    });
};

updateMonkeys();



var addMonkey = function() {
    console.log("Adding...");
    req('/monkeys', "POST").then(function(response) {
        monkey = JSON.parse(response)
        console.log("MONKEY ADDED!");
        console.log(monkey);
    }, function(err) {
        console.log("Didn't work: ", err);
    });
};