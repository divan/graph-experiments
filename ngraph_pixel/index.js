var createGraph = require('ngraph.graph');
var autoColorNodes = require('./colors.js');

var ajax = require('./ajax.js');
var pixel = require('ngraph.pixel');

Promise.all([
  ajax('data/data/positions.bin', { responseType: 'arraybuffer' }).then(toInt32Array),
  ajax('data/data/links.bin', { responseType: 'arraybuffer' }).then(toInt32Array),
  ajax('data/data/labels.json').then(toJson),
  ajax('data/data.json').then(toJson)
]).then(render);

function toInt32Array(oReq) {
  return new Int32Array(oReq.response);
}

function toJson(oReq) {
  return JSON.parse(oReq.responseText);
}

function render(data) {
  var positions = data[0];
  var links = data[1];
  var labels = data[2];
  var graphData = data[3];

  var graph = initGraphFromLinksAndLabels(links, labels, graphData);

  var renderer = pixel(graph, {
    node(node) {
        var props = { size: 5, color: "0xaaaaaa"};
        if (node.data.group === "coach") {
            props.size = 5;
        }

        if (node.data.color !== undefined) {
			props.color = node.data.color;
		}

        return props;
    },

    link() {
        return {
            fromColor: 0xffffff,
            toColor: 0xffffff
        };
    },

    // We need to use "dumb" links, otherwise it will be slow
    // Dumb links cannot be updated directly via properties. Have
    // to use renderer.edgeView().setFromColor(), renderer.edgeView().setToColor(), etc.
    activeLink: true
  });

  var layout = renderer.layout();
  // no need to do any layout here
  renderer.stable(true);

  // Set node positions.
  labels.forEach(function (label, index) {
    var nodeCount = index * 3;
    var x = positions[nodeCount + 0];
    var y = positions[nodeCount + 1];
    var z = positions[nodeCount + 2];

    layout.setNodePosition(label, x, y, z);
  });

  renderer.redraw();
}

function initGraphFromLinksAndLabels(links, labels, graphData) {
  var srcIndex;

  var graph =  createGraph({ uniqueLinkId: false });
  autoColorNodes(graphData.nodes);
  graphData.nodes.forEach((node, idx) => graph.addNode(node.id, node));
  links.forEach(processLink);

  return graph;

  function processLink(link) {
	console.log("Adding link", link);
    if (link < 0) {
      srcIndex = -link - 1;
    } else {
      var toNode = link - 1;
      var fromId = labels[srcIndex];
      var toId = labels[toNode];
      if (toId === undefined) {
          return
      }
      graph.addLink(fromId, toId);
    }
  }
}
