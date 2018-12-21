(function(){
    function channelWrite(url, data, callback) {
        fetch(url, {
            method: "POST",
            mode: "cors",
            cache: "no-cache",
            credentials: "same-origin",
            headers: {
                "Content-Type": "application/json; charset=utf-8",
            },
            redirect: "follow",
            referrer: "no-referrer",
            body: data,
        }).then(function(response){
            if (!response.ok) {
                throw new Error("Not OK response");
            }
            callback(true, null);
        }).catch(function(err){
            callback(false, JSON.stringify({
                "msg": err.toString(),
                "errObj": err,
            }));
        });
    }

    function channelRead(url, callback) {
        fetch(url, {
            method: "POST",
            mode: "cors",
            cache: "no-cache",
            credentials: "same-origin",
            headers: {
                "Content-Type": "application/json; charset=utf-8",
            },
            redirect: "follow",
            referrer: "no-referrer",
        }).then(function(response){
            if (!response.ok) {
                throw new Error("Not OK response");
            }
            return response.text()
        }).then(function(textValue) {
            callback(textValue, null);
        }).catch(function(err){
            callback(false, JSON.stringify({
                "msg": err.toString(),
                "errObj": err,
            }));
        });
    }

    window.httpChan = Object.freeze({
        read: channelRead,
        write: channelWrite,
    });
}());