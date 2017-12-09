var ws = new WebSocket('ws://' + window.location.host + '/ws');
ws.onopen = function (event) {
  ws.send("Here's some text that the server is urgently awaiting!"); 
};

ws.onmessage = function (event) {
  console.log("WS", event.data);
}
