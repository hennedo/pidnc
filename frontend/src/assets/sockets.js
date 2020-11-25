let socket
let callbacks = {}
let toSend = []
function openSocket() {
    socket = new WebSocket(`${(window.location.protocol === "http:" ? "ws:" : "wss:")}//${window.location.host}/ws`)
    socket.onopen = function() {
        while(toSend.length > 0) {
            let msg = toSend.splice(0,1)[0]
            socket.send(JSON.stringify(msg))
        }
    };
    socket.onmessage = function(e) {
        let msg;
        try {
            msg = JSON.parse(e.data)
        } catch(err) {
            console.log(err)
            console.log(e.data)
            return
        }

        if(callbacks[msg.type] && callbacks[msg.type].length > 0) {
            callbacks[msg.type].forEach((cb) => {
                cb.apply(socket, [msg.args])
            })
        }
    }
    socket.onclose = function() {
        window.setTimeout(openSocket, 200);
    }
}
openSocket();
function on(eventName, callback) {
    if(!callbacks[eventName]) {
        callbacks[eventName] = []
    }
    callbacks[eventName].push(callback)
    console.log(callbacks)
    return callback
}
function off(callback) {
    for(let i in callbacks) {
        let index = callbacks[i].indexOf(callback)
        if(index > -1) {
            callbacks[i].splice(index, 1)
        }
    }
}
function emit(eventName, data) {
    let msg = {
        type: eventName,
        args: data
    }
    if(socket.readyState === WebSocket.OPEN) {
        socket.send(JSON.stringify(msg))
    } else {
        toSend.push(msg)
    }
}
export default {on, off, emit};