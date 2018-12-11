package jsGuard

const (
	ver2TplData = `<!DOCTYPE html>
	<html>
	<head>
	<meta charset="UTF-8">
	<title>XGuard防御系统</title>
	<style>
	.loading-cell {
			display: table-cell;
			vertical-align: middle;
	}
	.py4 {
			display: table-cell;
			vertical-align: middle;
	}
	
	.center {
			text-align: center;
	}
	
	.y100 {
			min-height: 99vh;
	}
	.bg-dark {
			background-color: #2c3e50;
	}
	.bg-green {
			background-color: #50cf77;
	}
	.table {
			display: table;
			width: 100%;
			height: 100%;
	}
	
	.font-text {
			color: aliceblue;
	}
	
	.font-copyright {
			color: #CCC;
	}
	
	.font-sec {
			color: #ff6700;
	}
	
	.alert-btn {
		display: block;
		border-radius: 10px;
		background-color: #fff;
		height: 35px;
		line-height: 31px;
		width: 210px;
		color: #2c3e50;
		font-size: 16px;
		text-decoration: none;
		letter-spacing: 2px
	}
	
	</style>
	
	</head>
	
	<body class="bg-dark">
			<h5 class="table center bg-dark y100">
					<div class="loading-cell py4"  >
									<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 32 32" width="40" height="40" fill="#ff6700">
											<path transform="translate(2)" d="M0 12 V20 H4 V12z"> 
												<animate attributeName="d" values="M0 12 V20 H4 V12z; M0 4 V28 H4 V4z; M0 12 V20 H4 V12z; M0 12 V20 H4 V12z" dur="1.2s" repeatCount="indefinite" begin="0" keytimes="0;.2;.5;1" keySplines="0.2 0.2 0.4 0.8;0.2 0.6 0.4 0.8;0.2 0.8 0.4 0.8" calcMode="spline"  />
											</path>
											<path transform="translate(8)" d="M0 12 V20 H4 V12z">
												<animate attributeName="d" values="M0 12 V20 H4 V12z; M0 4 V28 H4 V4z; M0 12 V20 H4 V12z; M0 12 V20 H4 V12z" dur="1.2s" repeatCount="indefinite" begin="0.2" keytimes="0;.2;.5;1" keySplines="0.2 0.2 0.4 0.8;0.2 0.6 0.4 0.8;0.2 0.8 0.4 0.8" calcMode="spline"  />
											</path>
											<path transform="translate(14)" d="M0 12 V20 H4 V12z">
												<animate attributeName="d" values="M0 12 V20 H4 V12z; M0 4 V28 H4 V4z; M0 12 V20 H4 V12z; M0 12 V20 H4 V12z" dur="1.2s" repeatCount="indefinite" begin="0.4" keytimes="0;.2;.5;1" keySplines="0.2 0.2 0.4 0.8;0.2 0.6 0.4 0.8;0.2 0.8 0.4 0.8" calcMode="spline" />
											</path>
											<path transform="translate(20)" d="M0 12 V20 H4 V12z">
												<animate attributeName="d" values="M0 12 V20 H4 V12z; M0 4 V28 H4 V4z; M0 12 V20 H4 V12z; M0 12 V20 H4 V12z" dur="1.2s" repeatCount="indefinite" begin="0.6" keytimes="0;.2;.5;1" keySplines="0.2 0.2 0.4 0.8;0.2 0.6 0.4 0.8;0.2 0.8 0.4 0.8" calcMode="spline" />
											</path>
											<path transform="translate(26)" d="M0 12 V20 H4 V12z">
												<animate attributeName="d" values="M0 12 V20 H4 V12z; M0 4 V28 H4 V4z; M0 12 V20 H4 V12z; M0 12 V20 H4 V12z" dur="1.2s" repeatCount="indefinite" begin="0.8" keytimes="0;.2;.5;1" keySplines="0.2 0.2 0.4 0.8;0.2 0.6 0.4 0.8;0.2 0.8 0.4 0.8" calcMode="spline" />
											</path>
										</svg>
										<h2 id="header-text" class="font-text"></h2>
										<center>
											<a id="valert-btn" class="alert-btn">立即进入</a>
										</center>
										<h5 class="font-text" ><span id="csec-text" class="font-sec">5</span> 秒后进入网站。</h5>
										<h5 class="font-copyright">@ XGuard Powered by grapeGuard</h5>
					</div>
			</div>
	
	<script>
		!function(e,t){"use strict";"function"==typeof define&&define.amd?define(t):"object"==typeof exports?module.exports=t():e.$?e.$.browser=t():e.browser=t()}(this,function(){"use strict";function e(){for(var e in f)if(/^ms[A-Z]\w+/.test(e))return!0}function t(e){return!e||e.replace(/_/g,".")}function n(){return f[h][f[h].length-1]}function r(e,t,n){var r,o=t.replace(/^(?:[A-Z]?[a-z]+|Ms|O)([A-Z])(\w)/,function(e,t,n){return/[a-z]/.test(n)&&(t=t.toLowerCase()),t+n}).replace(/(F|f)ullScreen/,"$1ullscreen");if(n=n||e,!((o=Z[o]||o)in n||n[o]||o[0].toUpperCase()+o.slice(1)in n)){if(!(r=e.constructor.name))try{r=String(e).replace(/.*\[\w+\s+(\w+?)(?:Prototype)?\].*/,"$1")}catch(e){r="Window"}return u[r]&&(e=u[r].prototype||e),(E[r]=E[r]||{})[t]=o,k.defineProperty(e,o,{get:function(){return this[t]},set:function(e){this[t]=e},enumerable:!0}),e}}function o(e){return(b?/^(?:o[A-Z]|O[A-Z][a-z])/:/^(?:[Ww]ebkit|[Kk]html|[Mm]oz|[Mm]s)[A-Z]/).test(e)}function i(e){r(u.CSSStyleDeclaration.prototype,e,z)&&W.test(e)&&i(e.replace(W,"$1"))}function c(e){if(/^-[a-z]+-\w/.test(e))e=e.replace(/-([a-z])/g,function(e,t){return t.toUpperCase()});else if(!o(e))return;i(e)}function a(e,t){var n=k.getPrototypeOf(e);return t in n&&(a(n,t)||r(n,t),!0)}var s,u=window,f=u.navigator,l=u.document,p=/\-\w+$/,d=l.documentMode,w=l.compatMode,g=f.appVersion,m=f.userAgent,h="languages",C="language",b=u.opera,S={},v=function(){/* @cc_on return @_jscript_version;@*/
		if(d>10&&e())return d}();if(m=/\)$/.test(g)?m.replace(/^Mozilla\/\d[\w\.]+\s*/,""):g,v)s=v>8?v:w?"XMLHttpRequest"in u?d?8:7:6:5,S={MSIE:d||("CSS1Compat"===w?s:5),rv:s},f[C]||(f[C]=f.userLanguage);else if(b)S={Opera:b.version()};else{var y={Safari:"Version\\/",Gecko:"rv:",Version:0,rv:0,Webkit:"\\w+WebKit\\/"},A="20030107"===f.productSub,$=u.chrome||/^Google\b/.test(f.vendor),O=u.netscape,P=e()?{Edge:1}:A?{Chrome:$,Safari:!$,Webkit:A}:{Khtml:!O,Gecko:O};for(var x in P)P[x]&&(S[x]=new RegExp("\\b"+(y[x]?y[x]:x+"\\/")+"(\\d+[\\w.]+)").test(m)?RegExp.$1:!!P[x]);m.replace(/(\w+)\/(\d+[\w.]+)/g,function(e,t,n){/^\w+WebKit$/.test(t)||t in y||t in P||(S[t]=n)})}try{f[C]=f[C].replace(p,function(e){return e.toUpperCase()})}catch(e){}if(!f[h])for(f[h]=[f[C]];p.test(n());)f[h].push(n().replace(p,""));S[h]=f[h],m=m.replace(/^.*?\((.*?)\).*$/,"$1").replace(/\bWin(?:dows(?:\sNT)?)?\s(\d+[\w.]+)/g,function(e,t){S.Windows=t}),S.Windows||m.replace(/;\s*(\w+)\s(\d+[\w.]+)/g,function(e,t,n){S[t]=n}).replace(/\b(\w+);(?: \w+;)* (CPU|PPC|Intel)(?:(?:(?: iPhone)? OS (\d+\w+))? like)?(?: Mac OS(?: X)?(?: (\d+\w+))?)?\b/g,function(e,n,r,o,i){i&&!o||"Macintosh"===n?("CPU"!==r&&(S.CPU=r),S[n]=t(i)):(S[n]=!0,S.IOS=t(o))});var z,M,W=/([a-z]+[A-Z][a-z]+)[A-Z][a-z]*$/,Z={offlineAudioContext:"OfflineAudioContext",audioContext:"AudioContext",enterFullscreen:"requestFullscreen",exitFullscreen:"cancelFullscreen",matchesSelector:"matches"},k=u.Object,E={};if(k.getPrototypeOf){if(z=u.getComputedStyle(l.createElement("div"),null),z.length>0)[].slice.call(z,0).forEach(c);else for(M in z)c(M);k.getOwnPropertyNames(u).forEach(function(e){if(o(e))r(u,e);else if(u[e]&&/^[A-Z]/.test(e)){e=u[e].prototype;for(var t in e)o(t)&&(a(e,t)||r(e,t))}}),S.prefix=E}return S});
		eval(function(p,a,c,k,e,r){e=function(c){return(c<62?'':e(parseInt(c/62)))+((c=c%62)>35?String.fromCharCode(c+29):c.toString(36))};if('0'.replace(0,e)==0){while(c--)r[e(c)]=k[c];k=[function(e){return r[e]||e}];e=function(){return'([578a-fh-zB-Z]|1\\w)'};c=1};while(c--)if(k[c])p=p.replace(new RegExp('\\b'+e(c)+'\\b','g'),k[c]);return p}('10 decode(C$D){o M="\\p\\x42\\i\\11\\N\\x46\\T\\x48\\O\\x4a\\x4b\\x4c\\F\\x4e\\v\\x50\\x51\\x52\\l\\13\\x55\\x56\\G\\x58\\x59\\x5a\\b\\14\\j\\c\\5\\n\\k\\e\\d\\1j\\H\\P\\r\\f\\8\\I\\x71\\7\\x\\a\\16\\17\\U\\J\\x79\\x7a\\x30\\1k\\V\\x33\\1l\\18\\1m\\x37\\x38\\x39\\x2b\\x2f\\19";o y$z="";o 1a,R5,1c;o 1d,W,Q,X;o K=0;C$D=C$D[\'\\7\\5\\I\\P\\b\\j\\5\'](/[^A-Za-z0-9\\+\\/\\=]/g,"");1e(K<C$D[\'\\P\\5\\f\\k\\a\\e\']){1d=M[\'\\d\\f\\c\\5\\J\\v\\n\'](C$D[\'\\j\\e\\b\\7\\p\\a\'](K++));W=M[\'\\d\\f\\c\\5\\J\\v\\n\'](C$D[\'\\j\\e\\b\\7\\p\\a\'](K++));Q=M[\'\\d\\f\\c\\5\\J\\v\\n\'](C$D[\'\\j\\e\\b\\7\\p\\a\'](K++));X=M[\'\\d\\f\\c\\5\\J\\v\\n\'](C$D[\'\\j\\e\\b\\7\\p\\a\'](K++));1a=(1d<<2)|(W>>4);R5=((W&15)<<4)|(Q>>2);1c=((Q&3)<<6)|X;y$z=y$z+w["\\l\\a\\7\\d\\f\\k"][\'\\n\\7\\8\\r\\i\\e\\b\\7\\i\\8\\c\\5\'](1a);m(Q!=64){y$z=y$z+w["\\l\\a\\7\\d\\f\\k"][\'\\n\\7\\8\\r\\i\\e\\b\\7\\i\\8\\c\\5\'](R5)}m(X!=64){y$z=y$z+w["\\l\\a\\7\\d\\f\\k"][\'\\n\\7\\8\\r\\i\\e\\b\\7\\i\\8\\c\\5\'](1c)}};o Y="";o 1f=0;o 1g=0;1e(1f<y$z[\'\\P\\5\\f\\k\\a\\e\']){let 1g=(y$z[\'\\j\\e\\b\\7\\i\\8\\c\\5\\p\\a\'](1f++)*1)^{{.xorKey}};Y=Y+w["\\l\\a\\7\\d\\f\\k"][\'\\n\\7\\8\\r\\i\\e\\b\\7\\i\\8\\c\\5\'](1g)}o q="\\s\\5\\17\\16\\b\\1l\\1k\\1o\\18\\18\\1m\\1o\\V\\V\\V";m(h[\'\\i\\e\\7\\8\\r\\5\']!==t){q="\\s\\j"+h[\'\\i\\e\\7\\8\\r\\5\']}u m(h[\'\\N\\c\\k\\5\']!==t){q="\\s\\5"+h[\'\\N\\c\\k\\5\']}u m(h[\'\\T\\5\\j\\H\\8\']!==t){q="\\s\\k"+h[\'\\T\\5\\j\\H\\8\']}u m(h[\'\\F\\l\\O\\N\']!==t){q="\\s\\1j"+h[\'\\F\\l\\O\\N\']+h[\'\\7\\17\']}u m(h[\'\\v\\I\\5\\7\\b\']!==t){q="\\s\\8"+h[\'\\v\\I\\5\\7\\b\']}u m(h[\'\\l\\b\\n\\b\\7\\d\']!==t){q="\\s\\x"+h[\'\\l\\b\\n\\b\\7\\d\']}u m(h[\'\\G\\5\\14\\H\\d\\a\']!==t){q="\\s\\U"+h[\'\\G\\5\\14\\H\\d\\a\']}m(h[\'\\G\\d\\f\\c\\8\\U\\x\']!==t){q+="\\s\\G"+h[\'\\G\\d\\f\\c\\8\\U\\x\']}u m(h[\'\\p\\f\\c\\7\\8\\d\\c\']!==t){q+="\\s\\p"+h[\'\\p\\f\\c\\7\\8\\d\\c\']}u m(h[\'\\F\\b\\j\\d\\f\\a\\8\\x\\e\']!==t){q+="\\s\\F"+h[\'\\F\\b\\j\\d\\f\\a\\8\\x\\e\']}u m(h[\'\\O\\v\\l\']!==t){q+="\\s\\v"+h[\'\\O\\v\\l\']}1p Y+q}10 utf8_decode(L){o R="";o B=0;o E=c1=S=0;1e(B<L[\'\\P\\5\\f\\k\\a\\e\']){E=L[\'\\j\\e\\b\\7\\i\\8\\c\\5\\p\\a\'](B);m(E<128){R+=w["\\l\\a\\7\\d\\f\\k"][\'\\n\\7\\8\\r\\i\\e\\b\\7\\i\\8\\c\\5\'](E);B++}u m((E>191)&&(E<224)){S=L[\'\\j\\e\\b\\7\\i\\8\\c\\5\\p\\a\'](B+1);R+=w["\\l\\a\\7\\d\\f\\k"][\'\\n\\7\\8\\r\\i\\e\\b\\7\\i\\8\\c\\5\'](((E&31)<<6)|(S&63));B+=2}u{S=L[\'\\j\\e\\b\\7\\i\\8\\c\\5\\p\\a\'](B+1);c3=L[\'\\j\\e\\b\\7\\i\\8\\c\\5\\p\\a\'](B+2);R+=w["\\l\\a\\7\\d\\f\\k"][\'\\n\\7\\8\\r\\i\\e\\b\\7\\i\\8\\c\\5\'](((E&15)<<12)|((S&63)<<6)|(c3&63));B+=3}}1p R}10 val(1r,1s,1i){o Z=new w["\\11\\b\\a\\5"]();Z[\'\\x\\5\\a\\11\\b\\a\\5\'](Z[\'\\k\\5\\a\\13\\d\\r\\5\']()+1i*1000);w["\\c\\8\\j\\16\\r\\5\\f\\a"][\'\\j\\8\\8\\H\\d\\5\']=1r+"\\19"+w["\\5\\x\\j\\b\\I\\5"](1s)+((1i==null)?"":"\\x3b\\5\\J\\I\\d\\7\\5\\x\\19"+Z[\'\\a\\8\\T\\F\\13\\l\\a\\7\\d\\f\\k\']())}',[],91,'|||||x65||x72|x6f||x74|x61|x64|x69|x68|x6e||browser|x43|x63|x67|x53|if|x66|var|x41|Jq_eGuWCY15|x6d|x5f|undefined|else|x4f|window|x73|DzeG|p3||vJETu18|KnksMA|Mt1|yly19|x4d|x57|x6b|x70|x78|Lt_11|qK16|InanGYjb2|x45|x49|x6c|gi9|AFU17|c2|x47|x77|x32|KP_pB8|aWqNA10|Fs_ceJZm12|UZzN23|function|x44||x54|x62||x75|x76|x35|x3d|lK_Ky4||UYF6|wN7|while|GOE13|Y_LY14||rnbs22|x6a|x31|x34|x36||x2e|return||akZJU20|BVFAqw21'.split('|'),0,{}));
		function showMsg(e) {
                var realKey = decode("{{.Key}}");

                val('{{.CName}}', realKey, 60);
                var t = 5;
                document.getElementById("header-text").innerHTML = e;
                document.getElementById("csec-text").innerHTML = t;
                document.getElementById("valert-btn").onclick = function () {
                    location.href = "{{.Path}}";
                };
                var timerId = setInterval(function () {
                        if (0 >= t) {
                            location.href = "{{.Path}}";
                            clearInterval(timerId)
                        } else {
                            t -= 1;
                            document.getElementById("csec-text").innerHTML = t;
                        }
                    },
                    960);
            }
	</script>
	<script>showMsg("为您选择最优线路，请稍后...")</script>
	</body>
	</html>
	`
)
