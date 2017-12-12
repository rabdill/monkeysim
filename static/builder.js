function get (url, cb) {
    var req = new XMLHttpRequest();
    req.onreadystatechange = function() { 
        if (req.readyState == 4 && req.status == 200)
            cb(req.responseText);
    }
    req.open("GET", url, true);            
    req.send(null);
};

function updateMonkeys() {
    get('/info', function(response) {
        monkeys = JSON.parse(response)
        console.log(monkeys);
        console.log(monkeys.length);
        guts = "";
        for(var i=0, monkey; monkey = monkeys[i]; i++) {
            guts += "<li><strong>" + monkey.Name + "</strong> (" + monkey.Speed.toFixed(3) + " kkps): " + monkey.Progress;
        }
        document.getElementById("results").innerHTML = guts;
    });
};

function sleep(ms) {
    return new Promise(resolve => setTimeout(resolve, ms));
};
async function keepGoing() {
    while(true) {
        updateMonkeys();
        await sleep(2000);
    }
};

keepGoing();