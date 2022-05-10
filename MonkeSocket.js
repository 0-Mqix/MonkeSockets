class MonkeSocket {
    events 
    conn

    constructor(url) {
        if (window["WebSocket"]) {
            this.events = new Map()
            this.conn = new WebSocket("ws://" + document.location.host + url)
        
            const _this = this;
            this.conn.onmessage = (e) => {
                var messages = e.data.split('\n')
                for (let i = 0; i < messages.length; i++) {
                    
                    let data = messages[i]
                    let event = data.slice(0, data.indexOf(':') + 1)
                    let message = data.slice(data.indexOf(':') + 1)
            
                    let func = _this.events.get(event)
                    if (func != null) func(message)
                }   
            }
        }
    }

    on(event, func) {
        this.events.set(event, func)
    }

    onOpen(func) {
        this.conn.onopen = (e) => func(e)
    }

    onClose(func) {
        this.conn.onclose = (e) => func(e)
    }

    send(event, message) {
        this.conn.send(event+message)
    }
}
