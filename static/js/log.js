const backpanel = document.getElementById("backpanel")

const removelogpanel = ()=>{backpanel.style.display = "none"}
const addlogpanel = ()=>{backpanel.style.display = ""}

fetch("http://"+ window.location.host + "/checkuser?username=" + getCookie("username") ).then(/**@param p @type {Promise}}*/(p)=>{
    p.text().
    then((check)=> {
        console.log("check", check);
        if (getCookie("username") === null || check !== "") {
            addlogpanel()
            const input = document.getElementById("usernameinput")
            input.focus()
            const savecookie = (e)=>{
                if (e.key !== "Enter") return
                fetch("http://"+ window.location.host + "/newuser?username=" + input.value ).then(/**@param p @type {Promise}}*/(p)=>{
                    p.text().
                    then((t)=> {
                        if (t !== "") {input.placeholder = t; input.value = ""; return}
                        const d = new Date()
                        d.setDate(d.getDate() + 1)
                        d.setHours(6,0,0,0) //reset le cookie user tout les jours a 6h du matin heure locale a l'utilisateur
                        setCookie("username", input.value, d)
                        removelogpanel()
                        document.getElementById("username").innerText = getCookie("username")
                        window.removeEventListener("keydown", savecookie)
                        window.removeEventListener("click", focus)
                    } )
                })
                    
            }
            window.addEventListener("keydown", savecookie)
            const focus = ()=>{
                input.focus()
            }
            window.addEventListener("click", focus)
        } else {
            backpanel.innerHTML = ""
            document.getElementById("username").innerText += " " + getCookie("username")
        }
    })
})