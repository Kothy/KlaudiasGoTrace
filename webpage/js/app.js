
var texts = [];
var goroutines = [];
var channels = new Map();
var jsonArray = [];
var commObjs = [];
var namesObjs = [];
var scene;
var renderer;
var camera;
var ambientLight;
var controls;
var drawingThread;
var loadingScreen;
var loader;
var font;

let spinner = null;
let division = 1;
let WINDOW_CUT = 21;
let ARROW_COLOR = "#DC143C"
let GNAME_COLOR = "purple"
let GBODY_COLOR = "blue";
let GREEN_COLOR = "#55db30"
let ONE_FRAME = 0.0000001;
let THICKNESS = 0.02;
let THICKNESS2 = 0.007;
let MAXLEN = 10;
let RESOURCES_LOADED = true;


class Goroutine {
	constructor(id, parentId, start) {
		this.id = id;
		this.parentId = parentId;
		this.parent = null;
		this.start = start;
		this.end = 0;
		this.vecStart = null;
		this.vecEnd = null;
		this.len = 0
		this.children = [];
		this.depth = 0;
		this.received = [];
	}
	setDepth(d) {
		this.depth = d;
	}

	setChildren(ch){
		this.children = ch;
	}

	setEnd(end){
		this.end = end;
	}

	setParrent(g){
		this.parent = g;
	}

	setLen(len){
		this.len = len;
	}

	setVecStart(start){
		this.vecStart = start;
	}

	setVecEnd(end){
		this.vecEnd = end;
	}

	addReceived(item){
		this.received.push(item);
	}
}

class Channel {
	constructor(name) {
		this.name = name;
		this.chan = [];
	}

	size(){
		return this.chan.length;
	}
	 push(item) {
		 	this.chan.push(item);
	 }

	 pop(){
		 	return this.chan.shift();
	 }

	 isEmpty() {
 		if (this.size() == 0) return true;
		return false;
 	}
}

function mainApp(){
		texts = [];
		goroutines = [];

		scene = setScene();
		renderer = setRenderer();
		camera = setCamera();
		setWindow();
		ambientLight = setLight();

		controls = setControls();

		loader = new THREE.FontLoader();
		font = loader.parse(cambriaMathFont);

		setLoadScreen();

	// 	var settings = {
	// 		  message: "dat.GUI",
	// 		  checkbox: true,
	// 		  colorA: '#FF00B4',
	// 		  colorB: '#22CBFF',
	// 		  step5: 20,
	// 		  range: 500,
	// 		  options:"Option 1",
	// 		  speed:0,
	// 		  field1: "Field 1",
	// 		  field2: "Field 2",
	// };
	//
	// 	const gui = new dat.GUI();
	//
	// 	gui.add(settings, 'checkbox').name("Zobraz komunikáciu")
	// 	.onChange(function (value) {
  // 			checkboxCommunication(null);
	// 	});
	// 	gui.add(settings, 'checkbox').name("Zobraz mená gorutín")
	// 	.onChange(function (value) {
  // 			checkboxNames(null);
	// 	});

		const PARAMS = {
  			Red: true
		};
		const PARAMS2 = {
  			Purple: true
		};

		const pane = new Tweakpane({title: "Settings"});
		const input = pane.addInput(PARAMS, 'Red');

		input.on('change', function(ev) {
  			// console.log(`change: ${ev.value}`);
				checkboxCommunication(null);
		});

		const input2 = pane.addInput(PARAMS2, 'Purple');
		input2.on('change', function(ev) {
  			// console.log(`change: ${ev.value}`);
				checkboxNames(null);
		});

		drawScene();
}

function checkboxCommunication(event){
		var check = document.getElementById("commCheckbox");
		for (var i = 0; i < commObjs.length; i++) {
					commObjs[i].visible = !commObjs[i].visible;
		}
}

function checkboxNames(event){
		var check = document.getElementById("commCheckbox");
		for (var i = 0; i < namesObjs.length; i++) {
					namesObjs[i].visible = !namesObjs[i].visible;
		}
}

function setLoadScreen() {
	loadingScreen = {
		scene: new THREE.Scene(),
		camera: new THREE.PerspectiveCamera(75, width / height, 0.1, 1000),
		light: new THREE.AmbientLight(0xffffff,1.0),
		box: new THREE.Mesh(
			new THREE.TextGeometry("Drawing",{font: font,size: 1,height: 0.5}),
			new THREE.MeshBasicMaterial({color: GREEN_COLOR})
		)
	};
	loadingScreen.scene.background = new THREE.Color("rgb(255, 255, 255)");
	loadingScreen.camera.position.z = 6;
	loadingScreen.camera.lookAt(0, -5, 0);
	loadingScreen.box.geometry.center();

	// loadingScreen.box.position.set(0,0,0);
	loadingScreen.camera.lookAt(loadingScreen.box.position);
	loadingScreen.scene.add(loadingScreen.box);
}

mainApp();

function getChannel(name) {
		if (channels.has(name)) {
				return channels.get(name);
		} else {
				return createChannel(name);
		}
}

function createChannel(name) {
		var chan = new Channel(name);
		channels.set(name, chan);
		return chan;
}

function openFile(event) {
		var input = event.target;

		var reader = new FileReader();
		reader.onload = function() {
			var text = reader.result;
			var objJson = JSON.parse(text);
			jsonArray = objJson;
			drawText("Drawing", font, 0, 0, 0, 3, 0.5, GNAME_COLOR);
			Notiflix.Notify.Success('Drawing');

			setTimeout(loadJson, 100);
		};

		reader.readAsText(input.files[0]);
}

function resetScene() {
		clearScene(scene);
		setLight();
}

function clearScene(obj) {
  while (obj.children.length > 0){
    clearScene(obj.children[0]);
    obj.remove(obj.children[0]);
  }
  if (obj.geometry) obj.geometry.dispose();

  if (obj.material){
    Object.keys(obj.material).forEach(prop => {
      if (!obj.material[prop])
        return;
      if (obj.material[prop] !== null && typeof obj.material[prop].dispose === 'function')
        obj.material[prop].dispose();
    })
    obj.material.dispose();
  }
}

function calculateYFromTime(time) {
		var result;
		result = time * ONE_FRAME;
		return -result;
}

function calculateYFromTimeWDiv(time) {
		var result;
		result = time * ONE_FRAME;
		return -result / division;
}

function sleep(sleepDuration){
    var now = new Date().getTime();
    while(new Date().getTime() < now + sleepDuration){ /* do nothing */ }
}

function loadJson() {
			// RESOURCES_LOADED = false;
			goroutines = [];
			channels = new Map();
			texts = [];
			division = 1;

			var timeDiff;
			var max_len = 0;

			for (var i = 0; i < jsonArray.length; i++) {
				var obj = jsonArray[i];
				if (obj.Command === "GoroutineStart"){
					goroutines.push(new Goroutine(obj.Id, obj.ParentId, obj.Time));
				}
				else if (obj.Command === "GoroutineEnd"){
					var g = getGoroutineById(obj.Id);
					g.setEnd(obj.Time);
					startY = calculateYFromTime(g.start);
					endY = calculateYFromTime(g.end);
					g.vecEnd = new THREE.Vector3(0, endY, 0);
					g.vecStart = new THREE.Vector3(0, startY, 0);
					g.len = Math.abs(endY - startY);

					if (max_len < Math.abs(g.len)) {
						max_len = Math.abs(g.len)
					}
				}
				else if (obj.Command === "GoroutineSend"){
						var chan = getChannel(obj.Channel);
						var g2 = getGoroutineById(obj.Id);
						var value = obj.Value;
						chan.push([g2, value, obj.Time]);
				}
				else if (obj.Command === "GoroutineReceive"){
						var chan = getChannel(obj.Channel);
						var g2 = getGoroutineById(obj.Id);
						var value = obj.Value;
						var message = chan.pop();
						g2.addReceived([message[2], obj.Time, value, message[0]]);
				}
			}
			var div = 0;
			if ((max_len / MAXLEN) < 1) {
					div = 1;
			} else {
				div = (max_len / MAXLEN);
			}
			division = div;
			max_len = 0;
			for (var i = 0; i < goroutines.length; i++) {
					var g = goroutines[i];
					parent = getGoroutineById(g.parentId);
					g.parent = parent;

					g.vecEnd.y = g.vecEnd.y / div;
					g.vecStart.y = g.vecStart.y / div;
					g.len = g.vecEnd.y;
					g.len = Math.abs(g.vecEnd.y - g.vecStart.y);
					if (max_len < Math.abs(g.len)) {
						max_len = Math.abs(g.len)
					}
					var children = findChildren(g);
					g.setChildren(children);
		}
		setTimeout(function(){setDepths(getGoroutineById(1), 0);}, 100);

		// for (var i = 0; i < goroutines.length; i++) {
		// 		var g = goroutines[i];
		// 		// console.log(g);
		// }

		setTimeout(function(){resetScene();}, 100);
		setTimeout(function(){drawAllGoroutines(getGoroutineById(1));}, 100);
		setTimeout(function(){drawCommunication();}, 100);

		var mainGoroutine = goroutines[0];
		var start = mainGoroutine.vecStart;
		var end = mainGoroutine.vecEnd;

		camera.position.x = 0;
		camera.position.y = -(max_len/2) * 2.5;
		camera.position.z = -(max_len * 1.5);

		camera.lookAt(0, ((start.y + end.y)/2), -1);
		//camera.lookAt(0, 0, 0);
		setControls();
}

function setDepths(g, d) {
	g.setDepth(depth(g, d));
	for (var i = 0; i < g.children.length; i++) {
			var child = g.children[i];
			setDepths(child, 0);
	}
}

function depth(g, d){
		if (g.children.length == 0) {
			return d;
		}
		depths = [];
		for (var i = 0; i < g.children.length; i++) {
				var child = g.children[i];
				depths.push(depth(child, d + 1));
		}
		return Math.max(...depths);
}

function findChildren(goroutine) {
		var children = [];

		gid = goroutine.id;
		for (var i = 0; i < goroutines.length; i++) {
				if (goroutines[i].parentId == gid){
					children.push(goroutines[i]);
				}
		}
		return children;
}

function toRadians(degrees) {
		return degrees * (Math.PI/180);
}

Math.radians = function(degrees) {
	return degrees * Math.PI / 180;
}

Math.degrees = function(radians) {
	return radians * 180 / Math.PI;
}

function getMidPoint(a, b) {
		return new THREE.Vector3((a.x + b.x)/2, (a.y + b.y)/2, (a.z + b.z)/2);
}

function drawArrowWithText(origin, tip, color, textt) {
		const length = origin.distanceTo(tip);
		const midpoint = getMidPoint(tip, origin);

		var direction = new THREE.Vector3().subVectors(tip, origin);

		var arrowHelper = new THREE.ArrowHelper(direction.clone().normalize(),
		origin, direction.length(), color, 0.55);

		var arrowLabel = createText(textt, font, midpoint.x, midpoint.y + 0.1, midpoint.z, 0.2, 0.003, color);

		arrowLabel.setRotationFromEuler(arrowHelper.rotation);

		texts.push(arrowLabel);
		scene.add(arrowHelper);
		scene.add(arrowLabel);

		return [arrowHelper, arrowLabel];
}

function drawAllGoroutines(g) {
		// clearScene();
		var vecStart = g.vecStart;
		var vecEnd = g.setVecEnd;
		var x = vecStart.x;
		var y1 = vecStart.y;
		var y2 = vecEnd.y;
		var z = vecStart.z;
		var name = g.id == 1 ? 'main' : "#" + g.id;
		drawLineWithThickness(g.vecStart, g.vecEnd, THICKNESS, GBODY_COLOR);
		var nameObj = drawText(name, font, g.vecStart.x, g.vecStart.y + 0.35,
					g.vecStart.z, 0.35, 0.003, GNAME_COLOR);

		namesObjs.push(nameObj);

		var deg = 1 + Math.floor(Math.random() * 359);
		for (var i = 0; i < g.children.length; i++) {
				childrenNum = g.children.length
				if (childrenNum == 2) {
					childrenNum = 3;
				}
				deg += (360 / childrenNum);
				var mul = 1;
				if (g.children[i].children.length > 0) {
					mul = 1.8;
				}
				var x2 = (x + (Math.cos(Math.radians(deg)) * 2 * g.depth * mul));
				var z2 = (z + (Math.sin(Math.radians(deg)) * 2 * g.depth * mul));

				g.children[i].vecStart.x = x2;
				g.children[i].vecStart.z = z2;
				g.children[i].vecEnd.x = x2;
				g.children[i].vecEnd.z = z2;

				var yy = g.children[i].vecStart.y;
				var yy2 = g.children[i].vecEnd.y;
				var lVec = new THREE.Vector3(x ,yy, z);
				var lVec2 = new THREE.Vector3(x ,yy2, z);

				// drawSimpleLine(lVec, g.children[i].vecStart, "green");
				// drawSimpleLine(lVec2, g.children[i].vecEnd, "green");

				drawLineWithThickness(lVec, g.children[i].vecStart, THICKNESS2, GREEN_COLOR);
				drawLineWithThickness(lVec2, g.children[i].vecEnd, THICKNESS2, GREEN_COLOR);

		}
		for (var i = 0; i < g.children.length; i++) {
				drawAllGoroutines(g.children[i]);
		}
}

function drawCommunication() {
	for (var i = 0; i < goroutines.length; i++) {
			var g = goroutines[i];
			for (var j = 0; j < g.received.length; j++) {
				// received[0] - time of send
				// received[1] - time of receive
				// received[2] - received value
				// received[3] - send by Goroutine
				var timeSend = g.received[j][0];
				var timeReceived = g.received[j][1];
				var recValue = g.received[j][2];
				var sendG = g.received[j][3];
				var originY = calculateYFromTimeWDiv(timeSend);
				var tipY = calculateYFromTimeWDiv(timeReceived);
				var arrOrigin = new THREE.Vector3(sendG.vecStart.x, tipY, sendG.vecStart.z);
				var arrTip = new THREE.Vector3(g.vecStart.x, tipY, g.vecStart.z);
				if (sendG.id != g.id) {
						var objs;
						objs = drawArrowWithText(arrOrigin, arrTip, ARROW_COLOR, recValue);
						commObjs.push(objs[0]);
						commObjs.push(objs[1]);
				} else {
						var obj;
						obj = drawText(recValue, font, arrTip.x, tipY, arrTip.z, 0.2, 0.003, ARROW_COLOR);
						commObjs.push(obj);
				}

			}
	}
}

function getGoroutineById(id) {
	for (var i = 0; i < goroutines.length; i++) {
		if (goroutines[i].id === id){
			return goroutines[i];
		}
	}
	return null;
}

// function drawLineWithThickness(startV, endV, thick, color) {
// 		let len = startV.distanceTo(endV);
// 		let x1 = startV.x;
// 		let y1 = startV.y;
// 		let z1 = startV.z;
// 		let x2 = endV.x;
// 		let y2 = endV.y;
// 		let z2 = endV.z;
// 		let midX = (x1 + x2) / 2;
// 		let midY = (y1 + y2) / 2;
// 		let midZ = (z1 + z2) / 2;
// 		const geometry = new THREE.CylinderGeometry(thick, thick, len, 6);
// 		const material = new THREE.MeshBasicMaterial({color: color});
// 		const cylinder = new THREE.Mesh(geometry, material);
//
// 		cylinder.position.x = midX;
// 		cylinder.position.y = midY;
// 		cylinder.position.z = midZ;
//
// 		scene.add(cylinder);
// }

function drawLineWithThickness(pointX, pointY, thick, color) {
		var direction = new THREE.Vector3().subVectors(pointY, pointX);
	  var orientation = new THREE.Matrix4();
	  orientation.lookAt(pointX, pointY, new THREE.Object3D().up);
	  orientation.multiply(new THREE.Matrix4().set(1, 0, 0, 0,
	                0, 0, 1, 0,
	                0, -1, 0, 0,
	                0, 0, 0, 1));
	  var edgeGeometry = new THREE.CylinderGeometry(thick, thick, direction.length(), 6);
		var material = new THREE.MeshBasicMaterial({color: color});
	  var edge = new THREE.Mesh(edgeGeometry, material);
	  edge.applyMatrix4(orientation);

	  edge.position.x = (pointY.x + pointX.x) / 2;
	  edge.position.y = (pointY.y + pointX.y) / 2;
	  edge.position.z = (pointY.z + pointX.z) / 2;

		scene.add(edge);
	  return edge;
}

function drawSimpleLine(start, end, color){
		const material = new THREE.LineBasicMaterial({
				color: color
		});

		const points = [];
		points.push(start);
		points.push(end);

		const geometry = new THREE.BufferGeometry().setFromPoints(points);

		const line = new THREE.Line(geometry, material);
		scene.add(line);
		return line;
}

function random(min, max) {
	return Math.floor(Math.random() * (max - min) + min);
}

function drawText(text, font, x, y, z, size, height, color) {
		var geometryText = new THREE.TextGeometry(text, {
				font: font, size: size, height: height});

		geometryText.center();
		var materialText = new THREE.MeshLambertMaterial({color: color});
		var meshText = new THREE.Mesh(geometryText, materialText);

		// objects.push(materialText);
		// objects.push(meshText);
		// objects.push(geometryText);

		meshText.position.y = y;
		meshText.position.x = x;
		meshText.position.z = z;
		scene.add(meshText);

		texts.push(meshText);
		return meshText;
}

function drawTextSpin(text, size, height, color){
		var geometryText = new THREE.TextGeometry(text,
			{
				font: font,
				size: 1,
				height: 0.5,
			});

		geometryText.center();
		var materialText = new THREE.MeshBasicMaterial({color: color});
		var meshText = new THREE.Mesh(geometryText, materialText);

		meshText.position.y = 0;
		meshText.position.x = 0;
		meshText.position.z = 0;
		// scene.add(meshText);

		return meshText;
}

function createText(text, font, x, y, z, size, height, color){
		var geometryText = new THREE.TextGeometry(text, {
				font: font, size: size, height: height});

		geometryText.center();
		var materialText = new THREE.MeshLambertMaterial({color: color});
		var meshText = new THREE.Mesh(geometryText, materialText);

		meshText.position.y = y;
		meshText.position.x = x;
		meshText.position.z = z;
		return meshText;
}

function setControls(){
	return new THREE.OrbitControls(camera,renderer.domElement);
}

function drawScene(){

	var update = function(){
		// console.log("Kreslim scenu");
		//object.rotation.x += 0.01;
		//object.rotation.y += 0.001;

		texts.forEach((item, i) => {
			item.rotation.x = camera.rotation.x;
			item.rotation.y = camera.rotation.y;
			item.rotation.z = camera.rotation.z;
		});

	};
	// draw scene
	var render = function(){

		renderer.render(scene,camera);
	};
	// run game loop (update, renderer, repeat)
	var GameLoop = function(){
		requestAnimationFrame(GameLoop);
		update();
		render();

	};
	GameLoop();
}

function checkLoading() {
	if (RESOURCES_LOADED == false){
		requestAnimationFrame(drawScene);
		loadingScreen.box.rotation.x += 0.01;

		renderer.render(loadingScreen.scene, loadingScreen.camera);
		return true;
	}
	return false;
}

function setScene(){
	var scene = new THREE.Scene();
	scene.background = new THREE.Color("#282c34");
	// "#3cded3" - modrozelena
	// "white" - biela
	// "#cccccc" - slabo siva
	// "#282c34" - tmavo siva
	return scene;
}

function setRenderer(){
	var renderer = new THREE.WebGLRenderer();
	renderer.setSize(window.innerWidth, window.innerHeight - WINDOW_CUT);
	document.body.appendChild(renderer.domElement);
	return renderer
}

function setCamera(){
	width = window.innerWidth;
	height = window.innerHeight - WINDOW_CUT;
	var camera = new THREE.PerspectiveCamera(75, width / height, 0.1, 1000);
	camera.position.z = 6;
	camera.position.y = -4;
	camera.lookAt(0, 0, 6);
	return camera;
}

function setCameraPos(z) {
		camera.position.z = z;
}

function printCamera(){
		var lookAtVector = new THREE.Vector3(0,0, -1);
		lookAtVector.applyQuaternion(camera.quaternion);
		console.log(lookAtVector);
}

function setWindow(){

	var width = window.innerWidth;
	var height = window.innerHeight - WINDOW_CUT;

	window.addEventListener("resize", function(){
		var width = window.innerWidth;
		var height = window.innerHeight - WINDOW_CUT;
		renderer.setSize(width,height);
		camera.aspect = width / height;
		camera.updateProjectionMatrix();
	});
}

function setLight(){
	var ambientLight = new THREE.AmbientLight(0xffffff,1.0);
	scene.add(ambientLight);
	return ambientLight;
}

function getRandomColor() {
	var letters = '0123456789ABCDEF';
	var color = '#';
	for (var i = 0; i < 6; i++) {
		//color += letters[Math.floor(Math.random() * 16)];
		color += letters[Math.floor(Math.random() * 16)];
	}
	return color;
}
