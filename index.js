
var renderer = new THREE.WebGLRenderer();
renderer.setSize( window.innerWidth, window.innerHeight );
document.body.appendChild( renderer.domElement );

// Setup scene
const scene = new THREE.Scene();
scene.background = new THREE.Color(0x000011);

// Add lights
scene.add(new THREE.AmbientLight(0xbbbbbb));
scene.add(new THREE.DirectionalLight(0xffffff, 0.6));

// Setup camera
var camera = new THREE.PerspectiveCamera();
camera.far = 20000;

var tbControls = new THREE.TrackballControls(camera, renderer.domElement);

var d3ForceLayout = d3.forceSimulation()
            .force('link', d3.forceLink())
            .force('charge', d3.forceManyBody())
            .force('center', d3.forceCenter())
            .stop();

var animate = function () {
	requestAnimationFrame( animate );

	renderer.render(scene, camera);
};

animate();
