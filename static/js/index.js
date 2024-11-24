document.addEventListener("DOMContentLoaded", function(event) {
	console.log("Loaded");
	const login_form = document.getElementById("login");
	login_form.addEventListener("submit", onLoginNext, true);
	const register_form = document.getElementById("register");
	register_form.addEventListener("submit", onRegisterNext, true);
});

function showRegisterForm() {
	console.log("Showing registration form");
	document.getElementById("login").style.display = "none";
	document.getElementById("register").style.display = "block";
}

async function onLoginNext(e) {
	e.preventDefault();
	var username_element = document.querySelector("#login input[name='username']");
	const username = username_element.value;
	const url = "/login/begin?username=" + encodeURIComponent(username);
	const response = await fetch(url);
	if (response.status == 400) {
		showRegisterForm();
	}
	const json = await response.json()
	console.log(json);
}

async function onRegisterNext(e) {
	e.preventDefault();
	var username_element = document.querySelector("#register .username");
	console.log("register username", username_element);
}
