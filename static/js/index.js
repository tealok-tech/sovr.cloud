document.addEventListener("DOMContentLoaded", function(event) {
	console.log("DOM content loaded, getting started");
	document.getElementById("js_loading").style.display = "none";
	if (!window.PublicKeyCredential) {
		document.getElementById("no_webauthn").style.display = "block";
	}
	let p1 = new Promise((resolve, reject) => {
		if (!PublicKeyCredential.isUserVerifyingPlatformAuthenticatorAvailable) {
			document.getElementById("no_platform_auth_api").style.display = "block";
			reject("IsUserVerifyingPlatformAuthenticatorAvailable is not available");
		} else {
			PublicKeyCredential.isUserVerifyingPlatformAuthenticatorAvailable().then((available) => {
				console.log("isUserVerifyingPlatformAuthenticatorAvailable:", available);
				if (!available) {
					document.getElementById("no_platform_auth").style.display = "block";
				}
				resolve(available);
			});
		}
	});
	let p2 = new Promise((resolve, reject) => {
		if (!PublicKeyCredential.isConditionalMediationAvailable) {
			document.getElementById("no_conditional_mediation_api").style.display = "block";
			reject("IsConditionalMediationAvailable is not available");
		} else {
			PublicKeyCredential.isConditionalMediationAvailable().then((available) => {
				console.log("isConditionalMediationAvailable:", available);
				if (!available) {
					document.getElementById("no_conditional_mediation").style.display = "block";
				}
				resolve(available);
			});
		}
	});
	const login_form = document.getElementById("login");
	login_form.addEventListener("submit", onLoginNext, true);
	const register_form = document.getElementById("register");
	register_form.addEventListener("submit", onRegisterNext, true);
	Promise.all([p1, p2]).then((values) => {
		if(values[0] && values[1]) {
			login_form.style.display = "block";
		}
	})
});

function showRegisterForm() {
	console.log("Showing registration form");
	document.getElementById("login").style.display = "none";
	document.getElementById("register").style.display = "block";
	const login_username = document.querySelector("#login input[name='username']");
	const register_username = document.querySelector("#register input[name='username']");
	register_username.value = login_username.value;
}

function getLoginUsername() {
	var username_element = document.querySelector("#login input[name='username']");
	return username_element.value;
}

async function onLoginNext(e) {
	e.preventDefault();
	const username = getLoginUsername()
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
	const url = "/register/begin?username=" + encodeURIComponent(username);
	const response = await fetch(url);
	const json = await response.json()
	console.log(json);
}
