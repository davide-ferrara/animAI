import * as THREE from "three";
import { OrbitControls } from "three/addons/controls/OrbitControls.js";
import { GLTFLoader } from "three/addons/loaders/GLTFLoader.js";

var log = require("loglevel");
// Comment to disable logs
log.setLevel("trace");

const scene = new THREE.Scene();
const camera = new THREE.PerspectiveCamera(
  75,
  window.innerWidth / window.innerHeight,
  0.1,
  1000,
);

const renderer = new THREE.WebGLRenderer({ antialias: true });
renderer.setSize(window.innerWidth, window.innerHeight);
renderer.setClearColor(0x1f1f1f);
document.body.appendChild(renderer.domElement);

// Ambient light (soft general light)
const ambientLight = new THREE.AmbientLight(0xffffff, 0.5);
scene.add(ambientLight);

// Directional light (sun-like light to create shadows/definition)
const dirLight = new THREE.DirectionalLight(0xffffff, 1);
dirLight.position.set(5, 10, 7);
scene.add(dirLight);

camera.position.set(0, 5, 5);

const controls = new OrbitControls(camera, renderer.domElement);
controls.target.set(0, 0, 0);
controls.enablePan = true;
controls.enableZoom = true;

const size = 50;
const divisions = 50;
const gridHelper = new THREE.GridHelper(size, divisions);
scene.add(gridHelper);

const loader = new GLTFLoader();

(async () => {
  try {
    const gltf = await loader.loadAsync("static/human.glb");
    scene.add(gltf.scene);
    log.info("Model loaded!");
  } catch (error) {
    log.error("Error loading model:", error);
  }
})();

log.info("Scene loaded!");

function animate() {
  // Only rotate the cube if it exists, or just remove the rotation logic
  // cube.rotation.x += 0.01;

  controls.update(); // Required if enableDamping is true, harmless otherwise
  renderer.render(scene, camera);
}

// Start the loop
renderer.setAnimationLoop(animate);
