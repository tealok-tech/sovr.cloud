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

function getRegisterUsername() {
	var username_element = document.querySelector("#register input[name='username']");
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
	var username = getRegisterUsername()
	console.log("register username", username);
	var displayname_element = document.querySelector("#register input[name='displayname']");
	const url = "/register/begin?displayname=" + encodeURIComponent(displayname_element.value) + "&" + "username=" + encodeURIComponent(username);
	const response = await fetch(url);
	const publicKeyCredentialCreationOptions = await response.json()
	// Decode our URL-encoded base64 data
	publicKeyCredentialCreationOptions.publicKey.challenge = urlEncodedBase64ToArrayBuffer(publicKeyCredentialCreationOptions.publicKey.challenge);
	publicKeyCredentialCreationOptions.publicKey.user.id = urlEncodedBase64ToArrayBuffer(publicKeyCredentialCreationOptions.publicKey.user.id);
	await createPublicKey(publicKeyCredentialCreationOptions);
}
async function createPublicKey(options) {
	const credential = await navigator.credentials.create(options);
	console.log("New credential", credential);
	var username = getRegisterUsername()
	const url = "/register/finish?username=" + encodeURIComponent(username);
	const response = await fetch(url, {
		body: JSON.stringify({
			authenticatorAttachment: credential.authenticatorAttachment,
			id: credential.id,
			rawID: arrayBufferToUrlEncodedBase64(credential.rawID),
			response: {
				attestationObject: arrayBufferToUrlEncodedBase64(credential.response.attestationObject),
				clientDataJson: arrayBufferToUrlEncodedBase64(credential.response.clientDataJson)
			},
			type: credential.type
		}),
		method: "POST",
	})
	if (!response.ok) {
		console.error("Failed to finish registration", response)
	}
	const json = await response.json()
	console.log("Registration response", json)
}

function arrayBufferToUrlEncodedBase64(buffer) {
    // Convert ArrayBuffer to base64
    const bytes = new Uint8Array(buffer);
    const binaryString = Array.from(bytes).map(byte => String.fromCharCode(byte)).join('');
    const base64 = btoa(binaryString);
    
    // Make base64 URL-safe and URL encode
    return encodeURIComponent(base64.replace(/\+/g, '-').replace(/\//g, '_'));
}

function urlEncodedBase64ToArrayBuffer(base64) {
    const decodedBase64 = decodeURIComponent(base64.replace(/-/g, '+').replace(/_/g, '/'));

    const binaryString = atob(decodedBase64);
    const bytes = new Uint8Array(binaryString.length);

    for (let i = 0; i < binaryString.length; i++) {
        bytes[i] = binaryString.charCodeAt(i);
    }

    return bytes.buffer;
}
