var graph = require('../index.js');

var ws = new WebSocket('ws://' + window.location.host + '/ws');

// request graphData and initial positions from websocket connection
ws.onopen = function (event) {
	ws.send('{"cmd": "init"}'); 
};

ws.onmessage = function (event) {
	let msg = JSON.parse(event.data);
	switch(msg.type) {
		case "graph":
			graph.setGraphData(msg.graph);
			break;
		case "positions":
			console.log("Updating positions...");
			graph.updatePositions(msg.positions);
			break;
	}
}

// refresh restarts propagation animation.
function refresh() {
    ws.send('{"cmd": "refresh"}');
}
// js functions after browserify cannot be accessed from html,
// so instead of using onclick="refresh()" we need to attach listener
// here.
// Did I already say that whole frontend ecosystem is a one giant
// museum of hacks for hacks on top of hacks?
var refreshButton = document.getElementById('refreshButton');
refreshButton.addEventListener('click', refresh);

module.exports = { ws };
