import * as THREE from "three";
import { OrbitControls } from "three/addons/controls/OrbitControls.js";

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
renderer.setAnimationLoop(animate);
document.body.appendChild(renderer.domElement);
camera.position.set(0, 5, 5); // Un po' su (Y=5) e indietro (Z=5) per vedere bene il cubo 1x1

const controls = new OrbitControls(camera, renderer.domElement);

controls.target.set(0, 0, 0);

controls.enablePan = false;

controls.enableZoom = true;

controls.update();

const size = 50;
const divisions = 50;
const gridHelper = new THREE.GridHelper(size, divisions);
scene.add(gridHelper);

const geometry = new THREE.BoxGeometry(1, 1, 1);
const material = new THREE.MeshBasicMaterial({ color: 0x5f5f5f }); // Verde, per visibilità
const cube = new THREE.Mesh(geometry, material);
scene.add(cube);

function animate() {
  cube.rotation.x += 0.01;
  cube.rotation.y += 0.01;

  renderer.render(scene, camera);
}

console.log("[INFO] Tree.js loaded!");
