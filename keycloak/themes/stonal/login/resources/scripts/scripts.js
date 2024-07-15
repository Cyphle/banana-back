function togglePassword(resourcesPath, passwordElementId = "password", imgElementId = "vi") {
    const passwordElement = document.getElementById(passwordElementId);
    const imgElement = document.getElementById(imgElementId);
    if (passwordElement.type === "password") {
        passwordElement.type = "text";
        imgElement.src = resourcesPath + "/img/eye.png";
    } else {
        passwordElement.type = "password";
        imgElement.src = resourcesPath + "/img/eye-off.png";
    }
}
