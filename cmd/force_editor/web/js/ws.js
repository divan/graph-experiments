var ws = new WebSocket('ws://' + window.location.host + '/ws');

var prev = function(e) {
	e.preventDefault();
	ws.send('{"cmd": "prev"}');
}

var next = function(e) {
	e.preventDefault();
	ws.send('{"cmd": "next"}');
}

var calc = function(e) {
	e.preventDefault();
	ws.send('{"cmd": "calc"}');
}

// request graphData and initial positions from websocket connection
ws.onopen = function (event) {
	ws.send('{"cmd": "init"}'); 
	ws.send('{"cmd": "calc"}');
};


ws.onmessage = function (event) {
	let msg = JSON.parse(event.data);
	switch(msg.type) {
		case "graph":
			setGraphData(msg.graph);
			break;
		case "positions":
			updatePositions(msg.positions);
			break;
		case "forces":
			updateForces(msg.forces);
			break;
	}
}

