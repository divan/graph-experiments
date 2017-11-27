// WebGL
var renderer = new THREE.WebGLRenderer();
renderer.setSize( window.innerWidth, window.innerHeight );
document.body.appendChild( renderer.domElement );

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

var onFrame = null;

var animate = function () {
	if (onFrame) {
		onFrame();
	}

	// frame cycle
	tbControls.update();
	renderer.render(scene, camera);
	requestAnimationFrame( animate );
};

var width = window.innerWidth;
var height = window.innerHeight;
var nodeRelSize = 1;
var nodeResolution = 8;
var warmupTicks = 0;
var cooldownTicks = 1000;
var cooldownTime = 15000; //ms

var update = function () {
	resizeCanvas();

	onFrame = null; // pause simulation

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
	graphData.nodes.forEach(node => {
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
		sphere.position.set(Math.random() * 1000, Math.random() * 1000, Math.random() * 1000);
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

update();
animate();
