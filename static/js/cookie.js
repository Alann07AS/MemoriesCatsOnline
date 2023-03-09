// Fonction pour écrire un cookie
function setCookie(name, value, date) {
    var expires = "";
    if (date) {
        date.setTime(date);
        expires = "; expires=" + date.toUTCString();
    }
    document.cookie = name + "=" + (value || "")  + expires + "; SameSite=Strict;";
}

// Fonction pour lire un cookie
function getCookie(name) {
    var nameEQ = name + "=";
    var ca = document.cookie.split(';');
    for(var i=0;i < ca.length;i++) {
        var c = ca[i];
        while (c.charAt(0)==' ') c = c.substring(1,c.length);
        if (c.indexOf(nameEQ) == 0) return c.substring(nameEQ.length,c.length);
    }
    return null;
}

/**
 * @param {Promise} promise 
 */
function catchGoErr (promise) {
    promise.then((response) => {
        console.log(response);
        if (!response.ok) {
            response.text().then((data) => {
                if (data !== "") {
                    console.error(data);
                }
            });
        }
        return response;
    });
}
