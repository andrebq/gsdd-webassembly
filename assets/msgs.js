(function(){
    var noAttrs = {};
    function asVDom(jDom) {
        if (jDom.children) {
            var childrenDom = [];
            for(var item of jDom.children) {
                childrenDom.push(asVDom(item));
            }
            return m(jDom.tag, jDom.attrs || noAttrs, childrenDom);
        } else {
            return m(jDom.tag, jDom.attrs || noAttrs, jDom.text || "");
        }
    }

    function renderJSON(jsonStr) {
        var obj = JSON.parse(jsonStr),
            root = document.body;
        
            console.time("renderLoop");
            newRoot = asVDom(obj);
            m.render(root, newRoot);
            console.timeEnd("renderLoop");

    }

    window.renderJSON = renderJSON;
}())