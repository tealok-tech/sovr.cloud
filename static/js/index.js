document.addEventListener("DOMContentLoaded", function(event) {
	console.log("Loaded");
	var form = document.getElementById("login");
	form.addEventListener("submit", onLoginSubmit, true);
});

async function onLoginSubmit(e) {
	e.preventDefault();
	const url = "/login/begin";
	const response = await fetch(url);
	const json = await response.json()
	console.log(json);
}

