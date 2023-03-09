function upDatePlayerReady(params) {
    const pList = document.getElementById("salon")
    const existlist = []
    for (let [name, v] of Object.entries(params)) {
        existlist.push("player_"+name);
        const isready = v.Ready
        if (document.getElementById("player_"+name)) {
            const img = document.getElementById("i_"+name)
            img.style.display = !isready ? "none":""
            if (isready) img.setAttribute("src", "./static/img/paw.png")
        } else {
            const newP = document.createElement("div")
            newP.id = "player_"+name
            newP.classList.add("pdiv")
            pList.appendChild(newP)

            const pName = document.createElement("p")
            pName.textContent = name
            newP.appendChild(pName)

            const readycase = document.createElement("div")
            readycase.classList.add("readycase")
            newP.appendChild(readycase)

            const ready = document.createElement("img")
            ready.src = "./static/img/paw.png"
            ready.style.display = !isready ? "none":""
            ready.classList.add("readyimg")
            ready.id = "i_"+name
            readycase.appendChild(ready)
        }
    }
    Array(...pList.childNodes).forEach((n)=>{
        if (!existlist.includes(n.id)) pList.removeChild(n)
    })
}

var interval
function toggleCountGame() {
    const salon = document.getElementById("salon")
    if (interval) {
        salon.removeChild(document.getElementById("chrono"));
        clearInterval(interval);
        interval = undefined
        return
    }

    const chrono = document.createElement("div")
    chrono.id = "chrono"
    chrono.classList.add("chrono")
    salon.appendChild(chrono)
    const p = document.createElement("p")
    chrono.appendChild(p)
    var count = 10;
    const updatetime = ()=>{p.innerText = count}
    updatetime()
    interval = setInterval(()=>{
        count--
        updatetime()
        if (count === 0) {clearInterval(interval)}
    }, 1000)
}

function showgame(game) {
    console.log(game);
    document.getElementById("mainsalon").style.display = "none"
    document.getElementById("maingame").style.display = ""

    const board = document.getElementById("board")
    
    game.Cards.forEach((c) => { //add card element and img without info to preload
        const card = document.createElement("div")
        card.id = "card_id_"+ c.Id
        card.classList.add("card" ,"case")
        board.appendChild(card)
        
        const front = document.createElement("div")
        front.classList.add("front")
        card.appendChild(front)
        
        const back = document.createElement("div")
        back.classList.add("back")
        card.appendChild(back)
        const imgback = document.createElement("img")
        back.appendChild(imgback)
    });
    
    //preload image part
    let log = []
    game.Cards.forEach((c)=>{
        if (!log.includes(c.Img)) log.push(c.Img)
    })
    const shuffledArray = log.sort(() => Math.random() - 0.5);
    shuffledArray.forEach((src) => {
        const img = document.createElement("img")
        img.src = src
        img.style.display = "none"
        board.appendChild(img)
    })
    // document.getElementById()
}

var selectcard = []
function yourturn() {
    const board = document.getElementById("board")
    function getcard(event) {
        console.log(event.target);
        if (event.target.parentNode.id === null || selectcard.includes(event.target.parentNode.id)) return
        const cId = parseInt(event.target.parentNode.id.match(/\d+/g)[0])
        selectcard.push(cId)
        MessageOut.send(orderGO.SEND_CARD_BY_ID, cId)
        if (selectcard.length == 2) {
            board.removeEventListener("click", getcard)
            MessageOut.send(orderGO.CHECK_TWIN, Number(selectcard[0]), Number(selectcard[1]))
            selectcard = []
        }
    }
    board.addEventListener("click", getcard)
}

function showcard(cId, cSrc) {
    const card = document.getElementById("card_id_"+cId)
    const cardimg = card.getElementsByTagName("img")[0]
    cardimg.src = cSrc
    card.classList.toggle("flip", true)
}

function hidingcard(...cardsId) {
    cardsId.forEach((cId)=>{
        console.log(cId);
        const card = document.getElementById("card_id_"+cId)
        const cardimg = card.getElementsByTagName("img")[0]
        cardimg.src = ""
        card.classList.toggle("flip", false)    
    })
}

function updateplayerlist(plist, pplayid) {
    console.log(plist, pplayid);
    const div = document.getElementById("plist")
    for(let [pname, p] of Object.entries(plist)) {
        var pdiv = document.getElementById("p_"+pname)
        if (!pdiv) {
            const newPdiv = document.createElement("div")
            newPdiv.id = "p_"+pname
            newPdiv.classList.add("player")
            div.appendChild(newPdiv)
            pdiv = newPdiv

            const name = document.createElement("p")
            name.classList.add("name")
            name.innerText = pname
            newPdiv.appendChild(name)            
            
            const score = document.createElement("p")
            score.classList.add("score")
            newPdiv.appendChild(score)            
        }
        const score = pdiv.getElementsByClassName("score")[0]
        score.innerText = p.Score
        pdiv.classList.toggle("isPlay", p.Id === pplayid )
    };
}