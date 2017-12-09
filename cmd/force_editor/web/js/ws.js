var ws = new WebSocket('ws://' + window.location.host + '/ws');
ws.onopen = function (event) {
	ws.send('{"cmd": "init"}'); 
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
	}
}
