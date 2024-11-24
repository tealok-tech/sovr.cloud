document.addEventListener("DOMContentLoaded", function(event) {
	console.log("Loaded");
	var form = document.getElementById("login");
	form.addEventListener("submit", onLoginSubmit, true);
});

function onLoginSubmit(e) {
	e.preventDefault();
	console.log("hey");
}

