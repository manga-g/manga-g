// This is only for the Inputs to work dynamically

const types = ['primary', 'secondary'];
const hue = {
    "complement": 180,
    "triad-1": 120,
    "triad-2": -120,
}

var style = getComputedStyle(document.body)
style.getProperty = style.getPropertyValue
style.setProperty = function (name, value) {
    document.documentElement.style.setProperty(name, value);
}

function* objkv(obj) {
    let propKeys = Reflect.ownKeys(obj);

    for (let propKey of propKeys) {
        // `yield` returns a value and then pauses
        // the generator. Later, execution continues
        // where it was previously paused.
        yield [propKey, obj[propKey]];
    }
}

// could pass in an array of specific stylesheets for optimization
function getAllCSSVariableNames(styleSheets = document.styleSheets) {
    var cssVars = [];
    // loop each stylesheet
    for (var i = 0; i < styleSheets.length; i++) {
        // loop stylesheet's cssRules
        try { // try/catch used because 'hasOwnProperty' doesn't work
            for (var j = 0; j < styleSheets[i].cssRules.length; j++) {
                try {
                    // loop stylesheet's cssRules' style (property names)
                    for (var k = 0; k < styleSheets[i].cssRules[j].style.length; k++) {
                        let name = styleSheets[i].cssRules[j].style[k];
                        // test name for css variable signiture and uniqueness
                        if (name.startsWith('--') && cssVars.indexOf(name) == -1) {
                            cssVars.push(name);
                        }
                    }
                } catch (error) { }
            }
        } catch (error) { }
    }
    return cssVars;
}

var tsLight = {};

getAllCSSVariableNames().forEach((i) => {
    var m = /--ts-(light|dark)-([0-9]+)/g.exec(i);
    if (m && m[0] == i) {
        if (!tsLight[m[1]]) {
            tsLight[m[1]] = {};
        }
        tsLight[m[1]][m[2]] = style.getProperty(i);
    }
});

types.forEach((i) => {
    var [h, s, l] = style.getPropertyValue(`--${i}-color`).slice(4, -1).split(", ");

    var h = parseInt(h);
    var lp = parseFloat(l);

    style.setProperty(`--${i}-color-h`, h);
    style.setProperty(`--${i}-color-s`, s);
    style.setProperty(`--${i}-color-l`, l);
    style.setProperty(`--${i}-color`, `hsl(${h}, ${s}, ${l})`);

    // hue
    Object.keys(tsLight).forEach(function (k, v) {
        var hc = h + v;
        style.setProperty(`--${i}-${k}-h`, hc);
        style.setProperty(`--${i}-${k}`, `hsl(${hc}, ${s}, ${l}%)`);
    });

    // lightness
    for (let [k, v] of objkv(tsLight)) {
        for (let [vk, vv] of objkv(v)) {
            var lpl = lp + parseFloat(vv);
            style.setProperty(`--${i}-color-${k}-${vk}-l`, lpl);
            style.setProperty(`--${i}-color-${k}-${vk}`, `hsl(${h}, ${s}, ${lpl}%)`);
        }
    }
});

const body = document.querySelector("body");
document.querySelector("#banner").outerHTML = `<img id="banner" src="${bannerTheme(body.getAttribute('data-md-color-scheme'))}">`;
const banner = document.querySelector("#banner")
const bodySetAttribute = body.setAttribute;

function bannerTheme(theme) {
    if (theme == "wh_dark") {
        console.log("dark");
        return 'assets/images/icons/console/wh.gif';
    } else if (theme == "wh_white") {
        console.log("white");
        return 'assets/images/icons/console/bl.gif';
    }
}

body.setAttribute = (key, value) => {
    if (key == "data-md-color-scheme") {
        banner.src = bannerTheme(value)
        console.log(bannerTheme(value))
        bodySetAttribute.call(body, key, value);
    }
};