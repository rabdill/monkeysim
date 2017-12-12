function get(url) {
    return new Promise(function(resolve, reject) {
        var req = new XMLHttpRequest();
        req.onreadystatechange = function() {
            if(req.readyState == 4) {
                if(req.status == 200) resolve(req.responseText);
                else reject(req);
            }
        }
        req.open("GET", url, true);            
        req.send(null);
    });
};

function updateMonkeys() {
    get('/info').then(function(response) {
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