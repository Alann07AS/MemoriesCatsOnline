
function loadGamePage() {
    clearInterval(updateInterval)
    fetch("http://" + window.location.host + "/loadgame").then((r)=>{
        r.text().then((data)=>{
            document.body.innerHTML = data
        })
    })
}

/**@type {WebSocket} */
// var conn

if (!window["WebSocket"]) document.body.innerText = "YOUR NAVIGATOR IS NOT COMPATIBLE"

function joinGame(sName) {
    if (window["WebSocket"]) {
        const conn = new WebSocket("ws://" + window.location.host + "/joingame?servername=" + sName)
        conn.onerror = (e)=>{
            console.error(e);
        }
        conn.onclose = ()=>{
            alert("CONNECTION CLOSE")
            location.reload()       
        }
        conn.onmessage =
            /**
             * @param {MessageEvent} evt 
             */(evt)=>{
                console.log(evt.data);
                const incomingorder = JSON.parse(evt.data);
                if (incomingorder.Params === null) incomingorder.Params = [];
                actionJS.actions.get(incomingorder.Instruction)(...incomingorder.Params);
            }
        conn.onopen = ()=>{
            MessageOut.conn = conn;
            loadGamePage();
        }
    }
}

function leaveGame(sName) {
    catchGoErr(fetch("http://" + window.location.host + "/leavegame?servername=" + sName))
}

/**
 * Renseigner la con avant utilisation
 */
class MessageOut {
    /**
     * @param {number} orderGO 
     * @param  {...any} args 
     */
    static send (orderGO, ...args) {
        let newOrder = {}
        newOrder.UserName = getCookie("username")
        newOrder.Instruction = orderGO
        newOrder.Params = args
        MessageOut.conn.send(JSON.stringify(newOrder));
    }
    static conn
}


class MessageIn {
    static exec (InComing) {
        InComing.Instruction
        InComing.Params
    }
}

/**
 * Liste des fonction utilisable par le go ici en js
 */
class actionJS {
    static actions = new Map()
    static JS_UPDATE_PLAYER_READY = 1
    static JS_TOGGLE_START_CHRONO = 2
    static JS_SHOW_GAME           = 3
    static JS_SHOW_CARD           = 4
    static JS_YOUR_TURN           = 5
    static JS_UPDATE_PLIST        = 6
    static JS_HIDING_CARD         = 7
}

actionJS.actions.set(actionJS.JS_UPDATE_PLAYER_READY, (...params)=>{
    upDatePlayerReady(params[0])
})

actionJS.actions.set(actionJS.JS_TOGGLE_START_CHRONO, (...params)=>{
    toggleCountGame()
})

actionJS.actions.set(actionJS.JS_SHOW_GAME, (...params)=>{
    showgame(params[0])
})

actionJS.actions.set(actionJS.JS_SHOW_CARD, (...params)=>{
    const cId = params[0], cSrc = params[1]
    showcard(cId, cSrc)
})

actionJS.actions.set(actionJS.JS_YOUR_TURN, (...params)=>{
    yourturn()
})

actionJS.actions.set(actionJS.JS_UPDATE_PLIST, (...params)=>{
    updateplayerlist(params[0], params[1])
})

actionJS.actions.set(actionJS.JS_HIDING_CARD, (...params)=>{
    hidingcard(...params)
})

/**
 * Liste des fonction utilisable coter go depuis js
 */
class orderGO {
static TOGGLE_READY = 1
static SEND_CARD_BY_ID = 2
static CHECK_TWIN = 3
static _4 = 4
static _5 = 5
static _6 = 6
static _7 = 7
static _8 = 8
static _9 = 9
static _10 = 10
static _11 = 11
static _12 = 12
static _13 = 13
static _14 = 14
static _15 = 15
static _16 = 16
static _17 = 17
static _18 = 18
static _19 = 19
static _20 = 20
static _21 = 21
static _22 = 22
static _23 = 23
static _24 = 24
static _25 = 25
static _26 = 26
static _27 = 27
static _28 = 28
static _29 = 29
static _30 = 30
}
