const serverlist = document.getElementById("serverlist")

var lastupdate = new Date()
updateserverlist()

function updateserverlist() {
    fetch("http://" + window.location.host + "/getservelist").then((r)=>{
        r.json().then((j)=>{
            let existlist = Array(...serverlist.getElementsByClassName("splayers")).map((v)=>v.id.split("_")[0])
            let incomingSname = []
            Object.entries(j).forEach(element => {
                const s = element[1]
                incomingSname.push(s.name)
                // const existS = existlist.find((v)=> v.innerText === s.name)
                if (!existlist.includes(s.name)) {
                    const sname = document.createElement("div")
                    sname.onclick = ()=>{
                        joinGame(s.name)
                    }
                    sname.classList.add("sname", "s_"+s.name)
                    sname.innerText = s.name
                    serverlist.appendChild(sname);
                    
                    const stype = document.createElement("div")
                    stype.classList.add("s_"+s.name)
                    stype.innerText = "MemorieGame"
                    serverlist.appendChild(stype);

                    const splayers = document.createElement("div")
                    splayers.classList.add("splayers", "s_"+s.name)
                    splayers.id = s.name+"_players"
                    splayers.innerText = Object.entries(s.users).length
                    serverlist.appendChild(splayers);
                    // lineclass = !lineclass
                } else {
                    const players = document.getElementById(s.name+"_players")
                    players.innerText = Object.entries(s.users).length
                }
            });

            let difference = existlist.filter(x => !incomingSname.includes(x));
            difference.forEach((sname)=> {
                const el = Array(...serverlist.getElementsByClassName("s_"+sname))
                el.forEach(div => {
                    serverlist.removeChild(div)
                })
            })
            
            existlist = Array(...serverlist.getElementsByClassName("splayers")).map((v)=>v.id.split("_")[0])
            let lineclass = false
            existlist.forEach((sname)=>{
                const el = Array(...serverlist.getElementsByClassName("s_"+sname))
                el.forEach(div => {
                    div.classList.toggle("lp", lineclass)
                    div.classList.toggle("lip", !lineclass)
                })
                lineclass = !lineclass
            })
        })
    })
}

const updateInterval = setInterval(()=>{
    checkForUpdate()
}, 3000)

function checkForUpdate() {
    fetch("http://" + window.location.host + "/checkforupdate").then((r)=>{
        r.text().then((t)=>{
            const dateParts = t.split(" ");
            const lu = new Date(dateParts[0] + "T" + dateParts[1] + "Z");
            if (lu > lastupdate) {
                console.log("UPDATE");
                lastupdate = lu
                updateserverlist()
            }
        })
    })
}
