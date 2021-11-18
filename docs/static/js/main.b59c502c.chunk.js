(this.webpackJsonpwebapp2=this.webpackJsonpwebapp2||[]).push([[0],{15:function(e,t,a){e.exports=a(25)},20:function(e,t,a){},21:function(e,t,a){},23:function(e,t,a){},24:function(e,t,a){},25:function(e,t,a){"use strict";a.r(t);var n=a(0),r=a.n(n),i=a(13),l=a.n(i),s=(a(20),a(2));a(21);function c(e){return r.a.createElement("header",{className:"app-header"},r.a.createElement("div",{className:"icons"},r.a.createElement("div",{className:"icon github"},r.a.createElement("a",{href:"https://github.com/pedrooaugusto/steganography-png",target:"_blank"},r.a.createElement("i",{className:"fa fa-github"})," Github")),r.a.createElement("div",{className:"icon medium"},r.a.createElement("a",{href:"#dois"},r.a.createElement("i",{className:"fa fa-medium"})," Article"))),r.a.createElement("div",{className:"main-title"},r.a.createElement("h1",null,"Portable Network Graphics & Steganography"),r.a.createElement("div",{className:"subtitle"},"Hiding and retrieving secret files inside PNG files")))}var o=a(5),u=a(9),d=a(10),m=a(11),p=a(4),f=a.n(p),E=a(8);function h(e){var t=r.a.useState(v),a=Object(s.a)(t,2),n=a[0],i=a[1],l=r.a.useState(v),c=Object(s.a)(l,2),o=c[0],u=c[1],d=r.a.useState(null),m=Object(s.a)(d,2),p=m[0],h=m[1],w=r.a.useState(!1),y=Object(s.a)(w,2),N=y[0],O=y[1],T=r.a.useState(!1),I=Object(s.a)(T,2),j=(I[0],I[1]),S=r.a.useRef(null),k=function(){var t=Object(E.a)(f.a.mark((function t(){var a,n,r;return f.a.wrap((function(t){for(;;)switch(t.prev=t.next){case 0:if(O(!0),j(!0),!(null===(a=S.current)||void 0===a||null===(n=a.files)||void 0===n?void 0:n.length)){t.next=14;break}return t.t0=e,t.t1=Uint8Array,t.next=7,S.current.files[0].arrayBuffer();case 7:t.t2=t.sent,t.t3=new t.t1(t.t2),t.t0.setInputImage.call(t.t0,t.t3),h(null),r=URL.createObjectURL(S.current.files[0]),i(r),u(r);case 14:O(!1);case 15:case"end":return t.stop()}}),t)})));return function(){return t.apply(this,arguments)}}(),L=null==o||""===o;return r.a.createElement("div",{className:"config input-file"},r.a.createElement("div",{className:"title"},"Input image"),r.a.createElement("div",{className:"subtitle"},"The input file must be a PNG image, you can either load from the file system or from an external URL. This is the image in which the secret is hidden or the secret will be hidden (depending on the mode)."),r.a.createElement("div",{className:"load-url-input"},r.a.createElement("form",null,r.a.createElement("input",{list:"images",name:"url",id:"url",autoComplete:"off",placeholder:"Insert png url here",value:o,onChange:function(t){""===t.target.value&&(i(""),e.setInputImage(null)),u(t.target.value)}}),r.a.createElement("datalist",{id:"images"},r.a.createElement("option",{value:"https://raw.githubusercontent.com/pedrooaugusto/steganography-png/master/imagepack/suspicious-pitou.png"}),r.a.createElement("option",{value:"https://raw.githubusercontent.com/pedrooaugusto/steganography-png/master/imagepack/jinx.png"}),r.a.createElement("option",{value:"https://raw.githubusercontent.com/pedrooaugusto/steganography-png/master/imagepack/suspicious-bisky.png"})),r.a.createElement("button",{onClick:function(t){t.preventDefault(),O(!0),j(!0),fetch(o,{method:"GET"}).then((function(e){if(200!==e.status)throw new g("Request response was not ok:\n\t'".concat(e.statusText));var t=e.headers.get("content-type")||"";if(!t.toLocaleLowerCase().includes("png"))throw new b('Input file must be a png image!\n\tType "'.concat(t,'" is not "image/png"'));return e})).then((function(e){return e.arrayBuffer()})).then((function(t){h(null),i(o),e.setInputImage(new Uint8Array(t))})).catch((function(t){h(t.toString()),i(o),e.setInputImage(null)})).finally((function(){O(!1)}))},disabled:L},"Load"))),r.a.createElement("div",{className:"preview-img ".concat(e.empty?"empty":"")},r.a.createElement("figure",null,N?r.a.createElement("div",{className:"loading"},"Loading... ",r.a.createElement("i",{className:"fa-3x fa fa-spinner fa-spin"})):L||null==n||""===n?r.a.createElement("div",{className:"empty"},r.a.createElement("b",null,"EMPTY PREVIEW -- NO IMAGE!")):p?r.a.createElement("div",{className:"err"},r.a.createElement("span",null,p)):r.a.createElement("img",{src:n,alt:"Input preview",onLoad:function(){return j(!1)}}))),r.a.createElement("div",{className:"load-file"},r.a.createElement("label",{htmlFor:"file-upload-input-file",className:"btn"},"Or Load From File"),r.a.createElement("input",{id:"file-upload-input-file",type:"file",accept:".png",onChange:k,ref:S})))}var v=null!==""?"":"https://vignette.wikia.nocookie.net/anicrossbr/images/2/20/109_-_Neferpitou_portrait.png/revision/latest/scale-to-width-down/340?cb=20160308215759&path-prefix=pt-br",g=function(e){Object(u.a)(a,e);var t=Object(d.a)(a);function a(e){return Object(o.a)(this,a),t.call(this,"Failed to load image:\n\t"+e)}return a}(Object(m.a)(Error)),b=function(e){Object(u.a)(a,e);var t=Object(d.a)(a);function a(e){return Object(o.a)(this,a),t.call(this,"File type not supported:\n\t"+e)}return a}(Object(m.a)(Error));function w(e){return r.a.createElement("div",{className:"config mode"},r.a.createElement("div",{className:"title"},"Mode"),r.a.createElement("div",{className:"subtitle"},"You can either search for a secret hidden inside the input image or hide a new secret inside the input image."),r.a.createElement("div",{className:"opts"},r.a.createElement("label",{htmlFor:"hide"},r.a.createElement("input",{type:"radio",name:"mode",value:"hide",id:"hide",defaultChecked:!0,onClick:function(){return e.setMode("HIDE")}}),"Hide secret"),r.a.createElement("label",{htmlFor:"find"},r.a.createElement("input",{type:"radio",name:"mode",value:"find",id:"find",onClick:function(){return e.setMode("FIND")}}),"Find secret")))}function y(e){var t=e.secret,a=e.setSecret,n=e.empty,i=r.a.useRef(null),l="string"!==typeof t&&null!=(null===t||void 0===t?void 0:t.byteLength),s=l?t.slice(0,100).join(" ")+"...":t||"";return r.a.createElement("div",{className:"config secret"},r.a.createElement("div",{className:"title"},"Secret to be hidden"),r.a.createElement("div",{className:"subtitle"},"The secret can be a plain text message or a file loaded from the file system.",r.a.createElement("br",null),"Plain text messages where the first line is ",r.a.createElement("b",null,"#!HTML")," will be rendered as HTML."),r.a.createElement("div",{className:"plain-text"},r.a.createElement("textarea",{className:!n||null!=t&&""!==t?"":"empty",placeholder:"Type here the secret message to hide inside the input image",value:s.toString(),readOnly:l,disabled:l,title:l?"You cannot edit content loaded directly from a file!":"",onChange:function(e){return a(e.target.value)}}),r.a.createElement("div",{className:"footer"},r.a.createElement("div",{className:"clear-all"},r.a.createElement("button",{onClick:function(){return a(null)}},"Clear")),r.a.createElement("div",{className:"info"},"Data Length: ",l?t.length:s.length))),r.a.createElement("div",{className:"load-file"},r.a.createElement("label",{htmlFor:"file-upload-secret",className:"btn"},"Or Load From File"),r.a.createElement("input",{id:"file-upload-secret",type:"file",accept:"*",ref:i,onChange:Object(E.a)(f.a.mark((function e(){var t,n,r,l;return f.a.wrap((function(e){for(;;)switch(e.prev=e.next){case 0:if(!(null===(t=i.current)||void 0===t||null===(n=t.files)||void 0===n?void 0:n.length)){e.next=11;break}return e.t0=Uint8Array,e.next=4,i.current.files[0].arrayBuffer();case 4:e.t1=e.sent,r=new e.t0(e.t1),l=i.current.files[0].name.includes(".")&&i.current.files[0].name.split(".").pop(),r.type=i.current.files[0].type+(l?"."+l:""),a(r),e.next=12;break;case 11:a(null);case 12:case"end":return e.stop()}}),e)})))})))}function N(e){return r.a.createElement("div",{className:"config bit-loss"},r.a.createElement("div",{className:"title"},"Bit Loss"),r.a.createElement("div",{className:"subtitle"},"Given one byte of the input image, how many bits, of this byte, we should dedicate to the secret ?",r.a.createElement("br",null),"Eg: Using bitloss = 4; Would take 2 bytes of the input image to encode 1 byte of the secret."),r.a.createElement("div",{className:"bit-loss-input"},r.a.createElement("select",{name:"bitloss",defaultValue:"8",onChange:function(t){e.setBitLoss(parseInt(t.target.value))}},r.a.createElement("option",{value:"1"},"1 bit"),r.a.createElement("option",{value:"2"},"2 bits"),r.a.createElement("option",{value:"4"},"4 bits"),r.a.createElement("option",{value:"6"},"6 bits"),r.a.createElement("option",{value:"8"},"8 bits"))))}a(23);function O(e){var t=r.a.useState(!1),a=Object(s.a)(t,2),n=a[0],i=a[1];return r.a.createElement("div",{className:"config-section"},r.a.createElement("div",{className:"main-title"},r.a.createElement("div",{className:"title"},"Configuration"),r.a.createElement("div",{className:"subtitle"},"Here you can specify which parameters to apply during the proccess, such as the input image, mode (hide secret or reveal secret), bit loss and the secret to be hidden.")),r.a.createElement(h,{setInputImage:e.actions.setInputImage,empty:n&&null===e.state.imageBuf}),r.a.createElement(w,{setMode:e.actions.setMode}),"FIND"!==e.state.mode&&r.a.createElement(r.a.Fragment,null,r.a.createElement(y,{setSecret:e.actions.setDataToHide,secret:e.state.dataToHide,empty:n}),r.a.createElement(N,{setBitLoss:e.actions.setBitLoss})),r.a.createElement("div",{className:"submit-section"},r.a.createElement("button",{className:"btn",onClick:function(){return null==e.state.imageBuf?i(!0):"HIDE"!==e.state.mode||null!=e.state.dataToHide&&""!==e.state.dataToHide?(i(!1),void e.actions.startProcess()):i(!0)}},"GO!"),n&&r.a.createElement("span",{className:"subtitle"},"Pleae fill in all the required fields! ")))}function T(e){e.mode;var t=e.output;return t.err||t.loading||!t.result?"EMPTY":t.loading?"LOADING":t.err?"ERR":null}var I=function(e){if(T(e))return null;var t="HEX"===e.output.viewType,a="FIND"===e.mode;return r.a.createElement("button",{className:"btn ".concat(t?"selected":""," ").concat(a?"":"disabled"),disabled:!a,onClick:function(){return e.setOutputView("HEX")}},"Show as Hex")},j=function(e){if(T(e)||"HEX"!==e.output.viewType)return null;var t=Array.from(e.output.result).map((function(e){return e.toString(16)})).join(" ");return r.a.createElement("div",{className:"output-type hex",style:{position:"absolute"}},r.a.createElement("pre",null,t))},S=function(e){if(T(e))return null;var t="PLAIN"===e.output.viewType,a="FIND"===e.mode;return r.a.createElement("button",{className:"btn ".concat(t?"selected":""," ").concat(a?"":"disabled"),disabled:!a,onClick:function(){return e.setOutputView("PLAIN")}},"Show as Plain Text")},k=function(e){if(T(e)||"PLAIN"!==e.output.viewType)return null;var t=(new TextDecoder).decode(e.output.result);return r.a.createElement("div",{className:"output-type plain-text"},t.startsWith("#!HTML")?r.a.createElement("div",{dangerouslySetInnerHTML:{__html:t}}):r.a.createElement("pre",null,t))},L=function(e){var t,a;if(T(e))return null;var n="PNG"===e.output.viewType,i=((null===(t=e.output.dataType)||void 0===t||null===(a=t.search)||void 0===a?void 0:a.call(t,/png|gif|jpg|jpeg/gi))||-1)>=0;return r.a.createElement("button",{className:"btn ".concat(n?"selected":""," ").concat(i?"":"disabled"),disabled:!i,onClick:function(){return e.setOutputView("PNG")}},"Show as Image")},P=function(e){var t=T(e)||"PNG"!==e.output.viewType,a=e.mode,n=e.output,i=n.result,l=n.dataType,s=r.a.useMemo((function(){if(t)return null;var e="HIDE"===a?"image/png":l||"";e=e.split(".")[0];var n=new Blob([i],{type:e});return URL.createObjectURL(n)}),[t,a,i,l]);return t||null===s?null:r.a.createElement("div",{className:"output-type png"},r.a.createElement("figure",null,r.a.createElement("img",{src:s,alt:"Output file"})))},H=function(e){var t=e.subarray(0,8),a=[137,80,78,71,13,10,26,10];return t.every((function(e,t){return e===a[t]}))},D=function(e){if(T(e))return null;var t="PPNG"===e.output.viewType,a=H(e.output.result)&&("FIND"===e.mode||"HIDE"===e.mode);return r.a.createElement("button",{className:"btn ".concat(t?"selected":""," ").concat(a?"":"disabled"),disabled:!a,onClick:function(){return e.setOutputView("PPNG")}},"Show as Parsed PNG")},x=function(e){var t=T(e)||"PPNG"!==e.output.viewType,a=r.a.useState(""),n=Object(s.a)(a,2),i=n[0],l=n[1];return r.a.useEffect((function(){t||window.PNG.toString(e.output.result,(function(e,t){if(e)return console.error(e);l(t)}))}),[t]),t?null:r.a.createElement("div",{className:"output-type png-parsed"},r.a.createElement("pre",null,i))};a(24);function C(e){var t,a,n=e.state,i=n.mode,l=n.output;return r.a.createElement("div",{className:"output-section"},r.a.createElement("div",{className:"main-title"},r.a.createElement("div",{className:"title"},"Output"),r.a.createElement("div",{className:"subtitle"},"HIDE"===i&&r.a.createElement("span",null,"This is the resultant image with the secret hidden deep down in the pixels of each ",r.a.createElement("i",null,"scanline.")," Higher values for ",r.a.createElement("i",null,"bit loss "),"produces images with a high volume of noise."),"FIND"===i&&r.a.createElement("span",null,"This is what we found after looking for a hidden secret inside this image"))),r.a.createElement("div",{className:"result-file"},r.a.createElement("div",{className:"output"},r.a.createElement(M,e.state),r.a.createElement(_,e.state),r.a.createElement(j,e.state),r.a.createElement(k,e.state),r.a.createElement(P,e.state),r.a.createElement(x,e.state)),r.a.createElement("div",{className:"info"},r.a.createElement("b",null,l.result&&"HIDE"===i&&r.a.createElement("span",null,"New Image Length: ",null===l||void 0===l||null===(t=l.result)||void 0===t?void 0:t.length," bytes;"),l.result&&"FIND"===i&&r.a.createElement("span",null,"Hidden Secret Length: ",null===l||void 0===l||null===(a=l.result)||void 0===a?void 0:a.length," bytes; Hidden Secret Type: ",null===l||void 0===l?void 0:l.dataType,";"))),r.a.createElement("div",{className:"view-options"},r.a.createElement(L,Object.assign({},e.state,{setOutputView:e.actions.setOutputView})),r.a.createElement(D,Object.assign({},e.state,{setOutputView:e.actions.setOutputView})),r.a.createElement(S,Object.assign({},e.state,{setOutputView:e.actions.setOutputView})),r.a.createElement(I,Object.assign({},e.state,{setOutputView:e.actions.setOutputView})),!T(e.state)&&r.a.createElement("button",{className:"btn",onClick:function(){var e=document.createElement("a"),t=new Blob([l.result],{type:"application/octet-stream"});e.href=URL.createObjectURL(t),e.download="download-"+l.dataType,e.click()}},"Download Output"))))}var M=function(e){return"EMPTY"!==T(e)?null:r.a.createElement("div",{className:"output-type empty"},r.a.createElement("p",null,"Please, fill in the Configuration form."))},_=function(e){return"LOADING"!==T(e)?null:r.a.createElement("div",{className:"output-type loading"},r.a.createElement("h4",null,"Loading please wait..."))},F=a(1),B=a(14),G=new(function(){function e(){var t=this;Object(o.a)(this,e),this.worker=void 0,this.listeners=[],this.worker=new Worker("go-worker.js"),this.worker.onmessage=function(e){if("ErrorLoadingWorker"===e.data.type)return console.log("Killing Worker: "+e.data.error),t.worker.terminate();if("OperationResponse"===e.data.type){var a=t.listeners.find((function(t){return t[0]===e.data.id}));if(null==a)return;return t.listeners=t.listeners.filter((function(t){return t[0]!==e.data.id})),void(e.data.error?a[2](e.data.error):a[1](e.data.payload))}}}return Object(B.a)(e,[{key:"hideData",value:function(e,t,a,n){var r=this;return new Promise((function(i,l){var s=+new Date;r.listeners.push([s,i,l]),r.worker.postMessage({type:"Operation",operationName:"hideData",inputImage:e,data:t,dataType:a,bitLoss:n,id:s})}))}},{key:"revealData",value:function(e){var t=this;return new Promise((function(a,n){var r=+new Date;t.listeners.push([r,a,n]),t.worker.postMessage({type:"Operation",operationName:"revealData",inputImage:e,id:r})}))}}]),e}()),A={imageBuf:null,mode:"HIDE",dataToHide:"",bitLoss:8,output:{viewType:"PNG",result:null,err:null,loading:!1}};function R(){var e=arguments.length>0&&void 0!==arguments[0]?arguments[0]:A,t=arguments.length>1?arguments[1]:void 0;switch(t.type){case"SET_IMAGE_BUFF":return Object(F.a)(Object(F.a)({},e),{},{imageBuf:t.data,output:Object(F.a)({},A.output)});case"SET_MODE":return Object(F.a)(Object(F.a)({},e),{},{mode:t.data,output:Object(F.a)({},A.output)});case"SET_DATA_TO_HIDE":return Object(F.a)(Object(F.a)({},e),{},{dataToHide:t.data});case"SET_BITLOSS":return Object(F.a)(Object(F.a)({},e),{},{bitLoss:t.data});case"SET_OUTPUT_VIEW_TYPE":return Object(F.a)(Object(F.a)({},e),{},{output:Object(F.a)(Object(F.a)({},e.output),{},{viewType:t.data})});case"PROCCESS":return Object(F.a)(Object(F.a)({},e),{},{output:Object(F.a)(Object(F.a)({},e.output),t.data)});default:return e}}function U(e){var t=Object(s.a)(e,2),a=t[0],n=t[1];return[a,{setInputImage:function(e){n({type:"SET_IMAGE_BUFF",data:e})},setMode:function(e){n({type:"SET_MODE",data:e})},setDataToHide:function(e){n({type:"SET_DATA_TO_HIDE",data:e})},setBitLoss:function(e){n({type:"SET_BITLOSS",data:e})},setOutputView:function(e){n({type:"SET_OUTPUT_VIEW_TYPE",data:e})},startProcess:function(){n({type:"PROCCESS",data:{result:null,err:null,loading:!0}});var e,t=function(e,t){var r=arguments.length>2&&void 0!==arguments[2]?arguments[2]:"";return V(e,t,r,a.mode,n)};if("HIDE"===a.mode){var r,i=null===(r=a.dataToHide)||void 0===r?void 0:r.type;"string"===typeof a.dataToHide&&(i=a.dataToHide.startsWith("#!HTML")?"text/html.html":"text/plain.txt"),G.hideData(a.imageBuf,(e=a.dataToHide,e.buffer?e:(new TextEncoder).encode(e)),i,a.bitLoss).then((function(e){t(null,e.data,e.dataType)})).catch((function(e){t(e,new Uint8Array,"")}))}else G.revealData(a.imageBuf).then((function(e){t(null,e.data,e.dataType)})).catch((function(e){t(e,new Uint8Array,"")}));matchMedia("screen and (max-width: 860px)").matches?setTimeout((function(){return window.scrollTo(0,document.body.scrollHeight)}),100):window.scrollTo(0,0)}}]}var V=function(e,t){var a=arguments.length>2&&void 0!==arguments[2]?arguments[2]:"",n=arguments.length>3?arguments[3]:void 0,r=arguments.length>4?arguments[4]:void 0;if(e)return alert(e),r({type:"PROCCESS",data:{result:null,err:e,loading:!1}});var i=a.search(/png|gif|jpg|jpeg/gi)>=0||"HIDE"===n,l=a.search(/text/gi)>=0;r({type:"PROCCESS",data:{result:t,err:null,loading:!1,viewType:i?"PNG":l?"PLAIN":"HEX",dataType:"HIDE"===n?"image-png.png":a}})};var W=function(e){var t=r.a.useState(matchMedia(e).matches),a=Object(s.a)(t,2),n=a[0],i=a[1],l=r.a.useCallback((function(){i(matchMedia(e).matches)}),[e]);return r.a.useEffect((function(){return window.addEventListener("resize",l),function(){window.removeEventListener("resize",l)}}),[l]),n},Y=function(){var e=U(r.a.useReducer(R,A)),t=Object(s.a)(e,2),a=t[0],n=t[1],i=!W("screen and (max-width: 860px)")||a.output.result;return r.a.createElement("div",{className:"App"},r.a.createElement(c,null),r.a.createElement("main",null,r.a.createElement("section",null,r.a.createElement(O,{state:a,actions:n})),i&&r.a.createElement("section",null,r.a.createElement(C,{state:a,actions:n}))))};Boolean("localhost"===window.location.hostname||"[::1]"===window.location.hostname||window.location.hostname.match(/^127(?:\.(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)){3}$/));l.a.render(r.a.createElement(r.a.StrictMode,null,r.a.createElement(Y,null)),document.getElementById("root")),"serviceWorker"in navigator&&navigator.serviceWorker.ready.then((function(e){e.unregister()})).catch((function(e){console.error(e.message)}))}},[[15,1,2]]]);
//# sourceMappingURL=main.b59c502c.chunk.js.map