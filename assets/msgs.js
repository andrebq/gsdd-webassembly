(function(){
    var root = document.body;
    function asVDom(jDom) {
        if (jDom.children) {
            var childrenDom = [];
            for(var item of jDom.children) {
                childrenDom.push(asVDom(item));
            }
            return m(jDom.tag, jDom.attrs, childrenDom);
        } else {
            return m(jDom.tag, jDom.attrs, jDom.text || "");
        }
    }

    function renderJSON(jsonStr) {
        var obj = JSON.parse(jsonStr);

        var newRoot = asVDom(obj);

        render(root, newRoot);
    }

    window.renderJSON = renderJSON;
}())