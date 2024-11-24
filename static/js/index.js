document.addEventListener("DOMContentLoaded", function(event) {
	console.log("Loaded");
	var form = document.getElementById("login");
	form.addEventListener("submit", onLoginSubmit, true);
});

async function onLoginSubmit(e) {
	e.preventDefault();
	var username_element = document.getElementById("username");
	const username = username_element.value;
	const url = "/login/begin?username=" + encodeURIComponent(username);
	const response = await fetch(url);
	const json = await response.json()
	console.log(json);
}

