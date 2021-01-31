function f() {
	console.log("pracujem");
	f();
}

self.onmessage = (m) => {
		if (m.data.subject == "sprava") {
			postMessage("sprava prijata");
		}
};
