package main

import (
	"time"

	"github.com/jessevdk/go-assets"
)

var _Assets6748cf6e1affb3a363c1f9ced1881af87163013e = "const monthNames = [\n  \"January\", \"February\", \"March\", \"April\", \"May\", \"June\",\n  \"July\", \"August\", \"September\", \"October\", \"November\", \"December\"\n];\n\nfunction startTime() {\n  var today = new Date();\n  var h = today.getHours();\n  var m = today.getMinutes();\n  var M = today.getMonth();\n  var d = today.getDate();\n  m = checkTime(m);\n  document.getElementById('clock').innerHTML = h + \" \" + m;\n  document.getElementById('date').innerHTML = monthNames[M] + \" \" + d;\n  var t = setTimeout(startTime, 500);\n}\n\nfunction checkTime(i) {\n  if (i < 10) {\n    i = \"0\" + i\n  };\n  return i;\n}\n\nstartTime();\n"
var _Assetsfe42054a804c4bd857a5b5fec3bf398b4659e8ce = ":root {\n  --nord0: #2e3440;\n  --nord1: #3b4252;\n  --nord2: #434c5e;\n  --nord3: #4c566a;\n  --nord4: #d8dee9;\n  --nord5: #e5e9f0;\n  --nord6: #eceff4;\n  --nord7: #8fbcbb;\n  --nord8: #88c0d0;\n  --nord9: #81a1c1;\n  --nord10: #5e81ac;\n  --nord11: #bf616a;\n  --nord12: #d08770;\n  --nord13: #ebcb8b;\n  --nord14: #a3be8c;\n  --nord15: #b48ead;\n}\n\n* {\n  color: #fff;\n  font-family: monospace;\n  font-weight: normal;\n  margin: 0;\n  padding: 0;\n}\n\nhtml {\n  background-color: var(--nord0);\n  height: 100%;\n}\n\nbody {\n  height: 100%;\n  margin: 0;\n  padding: 0;\n  display: flex;\n  align-content: center;\n  justify-content: center;\n  align-items: flex-start;\n  flex-wrap: wrap;\n  text-align: center;\n}\n\nform {\n  margin: 10px;\n}\n\ninput {\n  border: none;\n  padding: 10px;\n  width: 50%;\n  background-color: var(--nord1);\n}\n\n.item-12 {\n  padding: 10px;\n  width: 100%;\n}\n\n.item-6 {\n  width: 50%;\n}\n\n.item-3 {\n  width: 25%;\n}\n\n.item-2 {\n  width: 16.6%;\n}\n\n#clock {\n  font-size: 10rem;\n}\n"
var _Assets950a25363eee220f7d8ce234bcc0b349e4ea9072 = "<!DOCTYPE html>\n<html>\n\t<head>\n\t\t<link rel=\"stylesheet\" href=\"/assets/styles.css\">\n\t</head>\n\n\t<body>\n\t\t<h1 id=\"clock\" class=\"item-12\"></h1>\n\t\t<h1 id=\"date\" class=\"item-12\"></h1>\n\t\t<form class=\"item-12\" name=\"x\" action=\"//duckduckgo.com/\">\n\t\t\t<input type=\"hidden\" value=\"1\" name=\"kh\"></input>\n\t\t\t<input type=\"hidden\" value=\"1\" name=\"kn\"></input>\n\t\t\t<input type=\"hidden\" value=\"1\" name=\"kac\"></input>\n\t\t\t<input type=\"search\" placeholder=\"DuckDuckGo\" name=\"q\"></input>\n\t\t</form>\n\t</body>\n\n\t<script src=\"/assets/date.js\"></script>\n</html>\n"

// Assets returns go-assets FileSystem
var Assets = assets.NewFileSystem(map[string][]string{"/": []string{"assets", "html"}, "/assets": []string{"date.js", "styles.css"}, "/html": []string{"index.html"}}, map[string]*assets.File{
	"/": &assets.File{
		Path:     "/",
		FileMode: 0x800001fd,
		Mtime:    time.Unix(1575919881, 1575919881224250198),
		Data:     nil,
	}, "/assets": &assets.File{
		Path:     "/assets",
		FileMode: 0x800001fd,
		Mtime:    time.Unix(1575919868, 1575919868752666510),
		Data:     nil,
	}, "/assets/date.js": &assets.File{
		Path:     "/assets/date.js",
		FileMode: 0x1b4,
		Mtime:    time.Unix(1575912349, 1575912349215457794),
		Data:     []byte(_Assets6748cf6e1affb3a363c1f9ced1881af87163013e),
	}, "/assets/styles.css": &assets.File{
		Path:     "/assets/styles.css",
		FileMode: 0x1b4,
		Mtime:    time.Unix(1575912159, 1575912159111660396),
		Data:     []byte(_Assetsfe42054a804c4bd857a5b5fec3bf398b4659e8ce),
	}, "/html": &assets.File{
		Path:     "/html",
		FileMode: 0x800001fd,
		Mtime:    time.Unix(1575919868, 1575919868752666510),
		Data:     nil,
	}, "/html/index.html": &assets.File{
		Path:     "/html/index.html",
		FileMode: 0x1b4,
		Mtime:    time.Unix(1575919868, 1575919868758000000),
		Data:     []byte(_Assets950a25363eee220f7d8ce234bcc0b349e4ea9072),
	}}, "")
