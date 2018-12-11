package jsGuard

const (
	ver3TplData = `<!DOCTYPE html>
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
			<div class="loading-cell py4">
				<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 32 32" width="40" height="40" fill="#ff6700">
					<path transform="translate(2)" d="M0 12 V20 H4 V12z">
						<animate attributeName="d" values="M0 12 V20 H4 V12z; M0 4 V28 H4 V4z; M0 12 V20 H4 V12z; M0 12 V20 H4 V12z" dur="1.2s" repeatCount="indefinite"
							begin="0" keytimes="0;.2;.5;1" keySplines="0.2 0.2 0.4 0.8;0.2 0.6 0.4 0.8;0.2 0.8 0.4 0.8" calcMode="spline"
						/>
					</path>
					<path transform="translate(8)" d="M0 12 V20 H4 V12z">
						<animate attributeName="d" values="M0 12 V20 H4 V12z; M0 4 V28 H4 V4z; M0 12 V20 H4 V12z; M0 12 V20 H4 V12z" dur="1.2s" repeatCount="indefinite"
							begin="0.2" keytimes="0;.2;.5;1" keySplines="0.2 0.2 0.4 0.8;0.2 0.6 0.4 0.8;0.2 0.8 0.4 0.8" calcMode="spline"
						/>
					</path>
					<path transform="translate(14)" d="M0 12 V20 H4 V12z">
						<animate attributeName="d" values="M0 12 V20 H4 V12z; M0 4 V28 H4 V4z; M0 12 V20 H4 V12z; M0 12 V20 H4 V12z" dur="1.2s" repeatCount="indefinite"
							begin="0.4" keytimes="0;.2;.5;1" keySplines="0.2 0.2 0.4 0.8;0.2 0.6 0.4 0.8;0.2 0.8 0.4 0.8" calcMode="spline"
						/>
					</path>
					<path transform="translate(20)" d="M0 12 V20 H4 V12z">
						<animate attributeName="d" values="M0 12 V20 H4 V12z; M0 4 V28 H4 V4z; M0 12 V20 H4 V12z; M0 12 V20 H4 V12z" dur="1.2s" repeatCount="indefinite"
							begin="0.6" keytimes="0;.2;.5;1" keySplines="0.2 0.2 0.4 0.8;0.2 0.6 0.4 0.8;0.2 0.8 0.4 0.8" calcMode="spline"
						/>
					</path>
					<path transform="translate(26)" d="M0 12 V20 H4 V12z">
						<animate attributeName="d" values="M0 12 V20 H4 V12z; M0 4 V28 H4 V4z; M0 12 V20 H4 V12z; M0 12 V20 H4 V12z" dur="1.2s" repeatCount="indefinite"
							begin="0.8" keytimes="0;.2;.5;1" keySplines="0.2 0.2 0.4 0.8;0.2 0.6 0.4 0.8;0.2 0.8 0.4 0.8" calcMode="spline"
						/>
					</path>
				</svg>
				<h2 id="header-text" class="font-text"></h2>
				<center>
					<div class="g-recaptcha" data-sitekey="6LcA6FgUAAAAAEsItyrJ0Y9FAiNiGQvBwacUI_WO"></div>
					<a id="valert-btn" class="alert-btn">立即进入</a>
				</center>
				<h5 class="font-copyright">@ XGuard Powered by grapeGuard</h5>
			</div>
			</div>
	
			<script src='https://www.google.com/recaptcha/api.js'></script>
			<script>
	
				function showMsg(e) {
					document.getElementById("header-text").innerHTML = e;
				}
			</script>
			<script>
				showMsg("为您选择最优线路，请稍后...")
			</script>
	</body>
	
	</html>
	`
)
