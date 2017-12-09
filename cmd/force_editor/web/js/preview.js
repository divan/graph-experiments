// WebGL
var renderer = new THREE.WebGLRenderer();
renderer.setSize( window.innerWidth, window.innerHeight );
document.body.appendChild( renderer.domElement );

var graphData;
var positions = Array();

// request graphData and initial positions from websocket connection
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

function setGraphData(data) {
	graphData = data;

	update();
}

function updatePositions(data) {
	positions = data;

	updateGraph();
}

// Setup scene
const scene = new THREE.Scene();
scene.background = new THREE.Color(0x000011);

// Add lights
scene.add(new THREE.AmbientLight(0xbbbbbb));
scene.add(new THREE.DirectionalLight(0xffffff, 0.6));

var graphScene = new THREE.Group();
scene.add(graphScene);

// Setup camera
var camera = new THREE.PerspectiveCamera();
camera.far = 20000;

var tbControls = new THREE.TrackballControls(camera, renderer.domElement);

var animate = function () {
	// frame cycle
	tbControls.update();
	renderer.render(scene, camera);
	requestAnimationFrame( animate );
};

var width = window.innerWidth;
var height = window.innerHeight;
var nodeRelSize = 4;
var nodeResolution = 8;

var update = function () {
	resizeCanvas();

	// parse links
	graphData.links.forEach(link => {
		link.source = link["source"];
		link.target = link["target"];
	});

	// Add WebGL objects
	while (graphScene.children.length) {
		graphScene.remove(graphScene.children[0])
	} // Clear the place

	// Render nodes
	const nameAccessor = accessorFn("name");
	const valAccessor = accessorFn("size");
	const colorAccessor = accessorFn("group");
	let sphereGeometries = {}; // indexed by node value
	let sphereMaterials = {}; // indexed by color
	graphData.nodes.forEach((node, idx) => {
		const val = valAccessor(node) || 1;
		if (!sphereGeometries.hasOwnProperty(val)) {
			sphereGeometries[val] = new THREE.SphereGeometry(Math.cbrt(val) * nodeRelSize, nodeResolution, nodeResolution);
		}

		const color = colorAccessor(node);
		if (!sphereMaterials.hasOwnProperty(color)) {
			sphereMaterials[color] = new THREE.MeshLambertMaterial({
				color: /*colorStr2Hex(color || '#ffffaa')*/ '#ffffaa',
				transparent: true,
				opacity: 0.75
			});
		}

		const sphere = new THREE.Mesh(sphereGeometries[val], sphereMaterials[color]);

		sphere.name = nameAccessor(node); // Add label
		sphere.__data = node; // Attach node data

		graphScene.add(node.__sphere = sphere);
		if (positions[idx] !== undefined) {
			sphere.position.set(positions[idx].x, positions[idx].y, positions[idx].z);
		}
	});

	const linkColorAccessor = accessorFn("color");
	let lineMaterials = {}; // indexed by color
	graphData.links.forEach(link => {
		const color = linkColorAccessor(link);
		if (!lineMaterials.hasOwnProperty(color)) {
			lineMaterials[color] = new THREE.LineBasicMaterial({
				color: /*colorStr2Hex(color || '#f0f0f0')*/ '#f0f0f0',
				transparent: true,
				opacity: 0.5,
			});
		}

		const geometry = new THREE.BufferGeometry();
		geometry.addAttribute('position', new THREE.BufferAttribute(new Float32Array(2 * 3), 3));
		const lineMaterial = lineMaterials[color];
		const line = new THREE.Line(geometry, lineMaterial);

		line.renderOrder = 10; // Prevent visual glitches of dark lines on top of spheres by rendering them last

		graphScene.add(link.__line = line);
	});

	// correct camera position
	if (camera.position.x === 0 && camera.position.y === 0) {
		// If camera still in default position (not user modified)
		camera.lookAt(graphScene.position);
		camera.position.z = Math.cbrt(graphData.nodes.length) * 150;
	}

	function resizeCanvas() {
		if (width && height) {
			renderer.setSize(width, height);
			camera.aspect = width/height;
			camera.updateProjectionMatrix();
		}
	}
};

var updateGraph = function () {
	graphData.nodes.forEach((node, idx) => {
		const sphere = node.__sphere;
		if (!sphere) return;

		sphere.position.x = positions[idx].x;
		sphere.position.y = positions[idx].y || 0;
		sphere.position.z = positions[idx].z || 0;
	});


	graphData.links.forEach(link => {
		const line = link.__line;
		if (!line) return;

		linePos = line.geometry.attributes.position;

		// TODO: move this index into map/cache or even into original graph data
		let start, end;
		for (let i = 0; i < graphData.nodes.length; i++) {
			if (graphData.nodes[i].id === link.source) {
				start = i;
				break;
			}	
		}
		for (let i = 0; i < graphData.nodes.length; i++) {
			if (graphData.nodes[i].id === link["target"]) {
				end = i;
				break;
			}	
		}

		linePos.array[0] = positions[start].x;
		linePos.array[1] = positions[start].y || 0;
		linePos.array[2] = positions[start].z || 0;
		linePos.array[3] = positions[end].x;
		linePos.array[4] = positions[end].y || 0;
		linePos.array[5] = positions[end].z || 0;

		linePos.needsUpdate = true;
		line.geometry.computeBoundingSphere();
	});
};

animate();
