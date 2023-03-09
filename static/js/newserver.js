const servername = document.getElementById("newservername")

// servername.onfocus = ()=>{
//     console.log("FOCCUS");
//     const onEnter = function (e) {
//         if (e.key !== "Enter") return
//         sumitserver()
//         servername.removeEventListener("keydown", onEnter)
//     }
//     servername.addEventListener("keydown", onEnter)
// }

// servername.addEventListener("focusout", )

servername.addEventListener("keydown", (e)=> {
    if (e.key !== "Enter") return
    sumitserver()
})

//newserver button
function sumitserver () {
    createServer(servername.value);
}

function createServer(sname) {
    if (sname.trim() !== "") {
        catchGoErr(fetch("http://" + document.location.host + "/newserver?servername="+sname))
    }
}

