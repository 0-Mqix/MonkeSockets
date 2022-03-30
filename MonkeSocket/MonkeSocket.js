class MonkeSocket {
    events = new Map();
    conn;
    url;

    constructor(url) {
        if (window["WebSocket"]) {
            this.url
            this.conn = new WebSocket("ws://" + document.location.host + url)
        
            const _this = this;
            this.conn.onmessage = function (e) {
                var messages = e.data.split('\n')
                for (var i = 0; i < messages.length; i++) {
                    var split = messages[i].split(":", 2)
                    var func = _this.events.get(split[0]+":")
                    if (func != null) { func(split[1]); }
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
        conn.send(event+message);
    }
}

export default MonkeSocket