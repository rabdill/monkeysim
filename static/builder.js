console.log("It loaded!");

var HttpClient = function() {
    this.get = function(aUrl, aCallback) {
        var anHttpRequest = new XMLHttpRequest();
        anHttpRequest.onreadystatechange = function() { 
            if (anHttpRequest.readyState == 4 && anHttpRequest.status == 200)
                aCallback(anHttpRequest.responseText);
        }

        anHttpRequest.open( "GET", aUrl, true );            
        anHttpRequest.send( null );
    }
}

var client = new HttpClient();
client.get('/info', function(response) {
    monkeys = JSON.parse(response)
    console.log(monkeys);
    console.log(monkeys.length);
    guts = "";
    for(var i=0, monkey; monkey = monkeys[i]; i++) {
        guts += "<li><strong>" + monkey.Name + "</strong> (" + monkey.Speed + "): " + monkey.Progress;
    }
    document.getElementById("results").innerHTML = guts;
});