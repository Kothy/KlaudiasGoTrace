var texts = [];
var goroutines = [];
var channels = new Map();
var jsonArray = [];
var commObjs = [];
var namesObjs = [];
var goroutinesObjs = [];
var sleepsObjs = [];
var parentsObjs = [];
var blockObjs = [];
var objects = [];
var scene;
var renderer;
var camera;
var ambientLight;
var controls;
var loader;
var font;


let spinner = null;
let division = 1;
let WINDOW_CUT = 21;
let ARROW_COLOR = "#DC143C"
let GNAME_COLOR = "#800080"
let GBODY_COLOR = "#0040ff";
let GSLEEP_COLOR = "#ffff00";
let GBLOCK_COLOR = "#fc9f1c";
let PARENT_COLOR = "#55db30";
let BG_COLOR = "#282c34";
let ONE_FRAME = 0.0000001;
let THICKNESS = 0.02;
let THICKNESS2 = 0.007;
let THICKNESS3 = THICKNESS + 0.005;
let MAXLEN = 10;
let arrowsCheck = true;
let namesCheck = true;
let sleepsCheck = false;
let blocksCheck = false;
let PARAMS = {
		Arrows: true,
		Names: true,
		Sleeps: false,
		Blocks: false,
		Background: BG_COLOR,
		GoroutinesCol: GBODY_COLOR,
		ArrowsCol: ARROW_COLOR,
		ParentsCol: PARENT_COLOR,
		NamesCol: GNAME_COLOR,
		SleepsCol: GSLEEP_COLOR,
		BlocksCol: GBLOCK_COLOR,
};


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
		this.sleeps = [];
		this.blocks = [];
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

	addSleep(item){
		this.sleeps.push(item);
	}

	addBlock(item){
		this.blocks.push(item);
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
		controls = setControls();
		setWindow();
		ambientLight = setLight();

		loader = new THREE.FontLoader();
		// font = loader.parse(cambriaMathFont);
		font = loader.parse(cambriaFont);

		createSettings();
		drawScene();
}

function createSettings(){

	const pane = new Tweakpane({title: "Settings"});
	const input = pane.addInput(PARAMS, 'Arrows');

	input.on('change', function(ev) {
			checkboxCommunication(ev.value);
	});

	const input2 = pane.addInput(PARAMS, 'Names');
	input2.on('change', function(ev) {
			checkboxNames(ev.value);
	});

	const input9 = pane.addInput(PARAMS, 'Sleeps');
	input9.on('change', function(ev) {
			checkboxSleeps(ev.value);
	});

	const input10 = pane.addInput(PARAMS, 'Blocks');
	input10.on('change', function(ev) {
			checkboxBlocks(ev.value);
	});

	const input3 = pane.addInput(PARAMS, 'Background');
	input3.on('change', function(ev) {
			scene.background = new THREE.Color(ev.value);
			BG_COLOR = ev.value;
	});

	const input4 = pane.addInput(PARAMS, 'GoroutinesCol');
	input4.on('change', function(ev) {
		for (var i = 0; i < goroutinesObjs.length; i++) {
				goroutinesObjs[i].material.color.set(ev.value);
		}
		GBODY_COLOR = ev.value;
	});

	const input5 = pane.addInput(PARAMS, 'ArrowsCol');
	input5.on('change', function(ev) {
		for (var i = 0; i < commObjs.length; i++) {
			if (commObjs[i].type == "Mesh"){
				commObjs[i].material.color.set(ev.value);
			} else{
				if (commObjs[i].children[0])
					commObjs[i].children[0].material.color.set(ev.value);

				if (commObjs[i].children[1])
					commObjs[i].children[1].material.color.set(ev.value);
			}
		}
		ARROW_COLOR = ev.value;
	});

	const input6 = pane.addInput(PARAMS, 'ParentsCol');
	input6.on('change', function(ev) {
		for (var i = 0; i < parentsObjs.length; i++) {
				parentsObjs[i].material.color.set(ev.value);
		}
		PARENT_COLOR = ev.value;

	});

	const input7 = pane.addInput(PARAMS, 'NamesCol');
	input7.on('change', function(ev) {
		for (var i = 0; i < namesObjs.length; i++) {
				namesObjs[i].material.color.set(ev.value);
		}
		GNAME_COLOR = ev.value;
	});

	const input8 = pane.addInput(PARAMS, 'SleepsCol');
	input8.on('change', function(ev) {
		for (var i = 0; i < sleepsObjs.length; i++) {
				sleepsObjs[i].material.color.set(ev.value);
		}
		GSLEEP_COLOR = ev.value;
	});

	const input11 = pane.addInput(PARAMS, 'BlocksCol');
	input11.on('change', function(ev) {
		for (var i = 0; i < blockObjs.length; i++) {
				blockObjs[i].material.color.set(ev.value);
		}
		GBLOCK_COLOR = ev.value;
	});
}

function checkboxCommunication(value){
		for (var i = 0; i < commObjs.length; i++) {
					commObjs[i].visible = value;
		}
		arrowsCheck = value;
}

function checkboxNames(value){
		for (var i = 0; i < namesObjs.length; i++) {
					namesObjs[i].visible = value;
		}
		namesCheck = value;
}

function checkboxSleeps(value){
		for (var i = 0; i < sleepsObjs.length; i++) {
					sleepsObjs[i].visible = value;
		}
		sleepsCheck = value;
}

function checkboxBlocks(value){
		for (var i = 0; i < blockObjs.length; i++) {
					blockObjs[i].visible = value;
		}
		blocksCheck = value;
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
			if (text == "") return;
			var objJson = JSON.parse(text);
			jsonArray = objJson;
			Notiflix.Loading.Hourglass('Drawing...');
			resetScene();
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
    while(new Date().getTime() < now + sleepDuration){}
}

function loadJson() {

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
				} else if (obj.Command === "GoroutineSleep") {
						var g = getGoroutineById(obj.Id);
						var startSleep = obj.Time;
						var sleepDuration = obj.Duration;
						var endSleep = startSleep + sleepDuration;

						startSleep += (sleepDuration/8);
						endSleep -= (sleepDuration/8);

						var startY = calculateYFromTime(startSleep);
						var endY = calculateYFromTime(endSleep);
						g.addSleep([startY, endY]);
				} else if (obj.Command === "GoroutineBlock") {
						var g = getGoroutineById(obj.Id);
						var startBlock = obj.Time;
						var blockDuration = obj.Duration;
						var endBlock = startBlock + blockDuration;
						startBlock += (blockDuration/8);
						endBlock -= (blockDuration/8);

						var startY = calculateYFromTime(startBlock);
						var endY = calculateYFromTime(endBlock);
						g.addBlock([startY, endY]);
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
					for (var j = 0; j < g.sleeps.length; j++) {
					  g.sleeps[j][0] = g.sleeps[j][0] / div;
						g.sleeps[j][1] = g.sleeps[j][1] / div;
						if (g.sleeps[j][1] < g.vecEnd.y) g.sleeps[j][1] = g.vecEnd.y;
					}

					for (j = 0; j < g.blocks.length; j++) {
					  g.blocks[j][0] = g.blocks[j][0] / div;
						g.blocks[j][1] = g.blocks[j][1] / div;
						if (g.blocks[j][1] < g.vecEnd.y) g.blocks[j][1] = g.vecEnd.y;
					}

					console.log(g.blocks);

					if (max_len < Math.abs(g.len)) {
						max_len = Math.abs(g.len)
					}
					var children = findChildren(g);
					g.setChildren(children);
		}

		setDepths(getGoroutineById(1), 0);
		drawAllGoroutines(getGoroutineById(1));
		drawCommunication();

		checkboxSleeps(sleepsCheck);
		checkboxNames(namesCheck);
		checkboxCommunication(arrowsCheck);
		checkboxBlocks(blocksCheck);

		camera.position.set(0, -max_len, 25);

		camera.lookAt(0, -(max_len), 0);
		setControls();

		Notiflix.Loading.Remove(50);
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

		objects.push(arrowHelper);
		objects.push(arrowLabel);

		return [arrowHelper, arrowLabel];
}

function drawAllGoroutines(g) {
		var vecStart = g.vecStart;
		var vecEnd = g.setVecEnd;
		var x = vecStart.x;
		var y1 = vecStart.y;
		var y2 = vecEnd.y;
		var z = vecStart.z;
		var name = g.id == 1 ? 'main' : "#" + g.id;
		var line = drawLineWithThickness(g.vecStart, g.vecEnd, THICKNESS, GBODY_COLOR);
		goroutinesObjs.push(line);
		var nameObj = drawText(name, font, g.vecStart.x, g.vecStart.y + 0.35,
					g.vecStart.z, 0.35, 0.003, GNAME_COLOR);

		g.sleeps.forEach((item, i) => {
			var start = new THREE.Vector3(g.vecStart.x, item[0], g.vecStart.z);
			var end = new THREE.Vector3(g.vecEnd.x, item[1], g.vecEnd.z);
			var sleep = drawLineWithThickness(start, end, THICKNESS3, GSLEEP_COLOR);
			sleepsObjs.push(sleep);
		});

		g.blocks.forEach((item, i) => {
			var start = new THREE.Vector3(g.vecStart.x, item[0], g.vecStart.z);
			var end = new THREE.Vector3(g.vecEnd.x, item[1], g.vecEnd.z);
			var block = drawLineWithThickness(start, end, THICKNESS3, GBLOCK_COLOR);
			blockObjs.push(block);
		});

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

				var line1 = drawLineWithThickness(lVec, g.children[i].vecStart, THICKNESS2, PARENT_COLOR);
				var line2 = drawLineWithThickness(lVec2, g.children[i].vecEnd, THICKNESS2, PARENT_COLOR);
				parentsObjs.push(line1);
				parentsObjs.push(line2);

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

function drawLineWithThickness(pointX, pointY, thick, color) {
		var direction = new THREE.Vector3().subVectors(pointY, pointX);
	  var orientation = new THREE.Matrix4();
	  orientation.lookAt(pointX, pointY, new THREE.Object3D().up);
	  orientation.multiply(new THREE.Matrix4().set(1, 0, 0, 0,
	                0, 0, 1, 0,
	                0, -1, 0, 0,
	                0, 0, 0, 1));
	  var edgeGeometry = new THREE.CylinderGeometry(thick, thick, direction.length(), 6);
		var material = new THREE.MeshBasicMaterial({
			color: color,
		});
	  var edge = new THREE.Mesh(edgeGeometry, material);
	  edge.applyMatrix4(orientation);

	  edge.position.x = (pointY.x + pointX.x) / 2;
	  edge.position.y = (pointY.y + pointX.y) / 2;
	  edge.position.z = (pointY.z + pointX.z) / 2;

		scene.add(edge);
		objects.push(edge);
	  return edge;
}

function drawSimpleLine(start, end, color){
		const material = new THREE.LineBasicMaterial({
				color: color,
		});

		const points = [];
		points.push(start);
		points.push(end);

		const geometry = new THREE.BufferGeometry().setFromPoints(points);

		const line = new THREE.Line(geometry, material);
		objects.push(line);
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
		var materialText = new THREE.MeshBasicMaterial({
			color: color,
		});
		var meshText = new THREE.Mesh(geometryText, materialText);

		// objects.push(materialText);
		// objects.push(meshText);
		// objects.push(geometryText);

		meshText.position.y = y;
		meshText.position.x = x;
		meshText.position.z = z;
		scene.add(meshText);
		objects.push(meshText);
		texts.push(meshText);
		return meshText;
}

// function drawTextSpin(text, size, height, color){
// 		var geometryText = new THREE.TextGeometry(text,
// 			{
// 				font: font,
// 				size: 1,
// 				height: 0.5,
// 			});
//
// 		geometryText.center();
// 		var materialText = new THREE.MeshBasicMaterial({color: color});
// 		var meshText = new THREE.Mesh(geometryText, materialText);
//
// 		meshText.position.y = 0;
// 		meshText.position.x = 0;
// 		meshText.position.z = 0;
// 		// scene.add(meshText);
//
// 		return meshText;
// }

function createText(text, font, x, y, z, size, height, color){
		var geometryText = new THREE.TextGeometry(text, {
				font: font, size: size, height: height});

		geometryText.center();
		var materialText = new THREE.MeshBasicMaterial({
			color: color,
		});
		var meshText = new THREE.Mesh(geometryText, materialText);

		meshText.position.y = y;
		meshText.position.x = x;
		meshText.position.z = z;
		return meshText;
}

function setControls(){
	var orbControls;
	orbControls = new THREE.OrbitControls(camera, renderer.domElement);
	return orbControls;
}

function drawScene(){

	var update = function(){

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

function setScene(){
	var scene = new THREE.Scene();
	scene.background = new THREE.Color(BG_COLOR);
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
	var camera = new THREE.PerspectiveCamera(60, width / height, 0.1, 1000);
	// camera.position.z = -6;
	// camera.position.y = -4;
	// camera.lookAt(0, 0, 6);
	return camera;
}

function setWindow(){

	var width = window.innerWidth;
	var height = window.innerHeight - WINDOW_CUT;

	window.addEventListener("resize", function(){
		var width = window.innerWidth;
		var height = window.innerHeight - WINDOW_CUT;
		renderer.setSize(width,height);
	});
}

function setLight(){
	var ambientLight = new THREE.AmbientLight(0xffffff,1.0);
	scene.add(ambientLight);
	return ambientLight;

	// const color = 0xFFFFFF;
	// const intensity = 4;
	// const light = new THREE.DirectionalLight(color, intensity);
	// light.position.set(0, 10, 0);
	// light.target.position.set(0, -10, 0);
	// scene.add(light);
	// scene.add(light.target);
	//
	// const light1 = new THREE.DirectionalLight(color, intensity);
	// light1.position.set(0, 1, 0);
	// light1.target.position.set(1, 1, 1);
	// scene.add(light1);
	// scene.add(light1.target);
	//
	// const light2 = new THREE.DirectionalLight(color, intensity);
	// light2.position.set(0, -1, 0);
	// light2.target.position.set(-1, 1, -1);
	// scene.add(light2);
	// scene.add(light1.target);

	// const color = 0xFFFFFF;
  //   const intensity = 1;
  //   const light = new THREE.PointLight(color, intensity);
  //   light.position.set(0, 5, 0);
	//
	// 	const helper = new THREE.PointLightHelper(light);
  //   scene.add(helper);
	//
	// 	const light2 = new THREE.PointLight(color, intensity);
  //   light2.position.set(-2, 20, -3);
	//
	// 	const light3 = new THREE.PointLight(color, intensity);
  //   light3.position.set(2, 20, 3);
	//
	// 	const light4 = new THREE.PointLight(color, intensity);
  //   light4.position.set(2, -20, 3);
	//
	// 	const light5 = new THREE.PointLight(color, intensity);
  //   light5.position.set(-2, -20, -3);
	//
	// 	const light6 = new THREE.PointLight(color, 3);
  //   light6.position.set(0, 0, 0);
	//
  //   scene.add(light, light2, light3, light4, light5, light6);

	return light;
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
