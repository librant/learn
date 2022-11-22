var ABDecode= function(ab) {
    let ret = '';
    for (let i = 0; i < ab.length; ) {
        let code = 0;
        if((ab[i] & 0x80) === 0){
            code = ab[i++];
        }else if((ab[i] & 0xe0) === 0xc0){
            code = (ab[i++] & 0x1f) << 6;
            code |= (ab[i++] & 0x3f);
        }else if((ab[i] & 0xf0) === 0xe0){
            code = (ab[i++] & 0x0f) << 12;
            code |= (ab[i++] & 0x3f) << 6;
            code |= (ab[i++] & 0x3f);
        }else if((ab[i] & 0xf8) === 0xf0){
            code = (ab[i++] & 0x07) << 18;
            code |= (ab[i++] & 0x3f) << 12;
            code |= (ab[i++] & 0x3f) << 6;
            code |= (ab[i++] & 0x3f);
        }else{
            i++;
            code = '?';
        }
        ret = ret + String.fromCharCode(code);
    }
    return ret;
};

var str2AB = function(str) {
    let utf8str = '';
    for (let i = 0; i < str.length; i++) {
        let c = str.charCodeAt(i);
        let bytesLeft;
        if (c <= 0x7F) {
            utf8str += str.charAt(i);
            continue;
        } else if (c <= 0x7FF) {
            utf8str += String.fromCharCode(0xC0 | (c >>> 6));
            bytesLeft = 1;
        } else if (c <= 0xFFFF) {
            utf8str += String.fromCharCode(0xE0 | (c >>> 12));
            bytesLeft = 2;
        } else {
            utf8str += String.fromCharCode(0xF0 | (c >>> 18));
            bytesLeft = 3;
        }
        while (bytesLeft > 0) {
            bytesLeft--;
            utf8str += String.fromCharCode(0x80 | ((c >>> (6 * bytesLeft)) & 0x3F));
        }
    }
    let ret = new Uint8Array(utf8str.length);
    for(let i = 0; i<utf8str.length; i++){
        ret[i] = utf8str.charCodeAt(i);
    }
    return ret;
};

(function() {
    let httpsEnabled = window.location.protocol === "https:";
    let wsURL = (httpsEnabled ? 'wss://' : 'ws://') + window.location.host + '/login?context=' + context +
        '&namespace=' + namespace + '&pod=' + pod + "&container=" + container;

    let openWs = function() {
        let ws = new WebSocket(wsURL);
        ws.binaryType = "arraybuffer";

        let term;
        let ideaTimer;

        ws.onopen = function(event) {
            console.log("websocket connection opened");
            ideaTimer = setInterval(sendPing, 30 * 1000, ws);
            term = new Terminal({
                cursorBlink:  true
            });

            term.on('resize', function(size) {
                setTimeout(function() {
                    term.showOverlay(size.cols + 'x' + size.rows);
                }, 500);
                term.fit()
                let msg = {type: "resize", rows: size.rows, cols: size.cols}
                ws.send(JSON.stringify(msg))
            });

            term.on("data", function(data) {
                ws.send(str2AB(data));
            });

            term.on('open', function() {
                window.addEventListener('resize', function(event) {
                    term.fit();
                });
                term.fit();
                term.focus();
            });

            term.open(document.getElementById("terminal-container"))
        };

        ws.onmessage = function(event) {
            let dataArray = new Uint8Array(event.data);
            term.write(ABDecode(dataArray));
        };

        ws.onclose = function(event) {
            if (term) {
                term.off('data');
                term.off('resize');
                term.showOverlay("kubectl exec connection closed", null);
            }
            clearInterval(ideaTimer);
        };
    };

    let sendPing = function(ws) {
        let newDate = new Date();
        let msg = {type: "ping", input: newDate.toISOString()}
        ws.send(JSON.stringify(msg));
    }

    openWs();
})();
