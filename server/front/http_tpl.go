package front

const front_tpl = `
<!DOCTYPE html>
<html>
<head>
<title>BAR</title>
</head>
<body>
<h1>Welcome to BAR!</h1>
<p>BAR is simple BLOB vendoring system.
	Visit <a href="https://github.com/akaspin/bar">github repo</a> for details.</p>

<h2>Windows users</h2>
<p>To export specs without git download
	<a href="{{.Info.HTTPEndpoint}}/win/bar-export.bat"><code>bar-export.bat</code></a>
	and save in root of the working tree.</p>
<p>This script is no-brain solution to export bar specs. It automatically
download <code>bar.exe</code> if it is not found and upload BLOBs and spec
to bard</p>
<p>You can also download <a href="{{.Info.HTTPEndpoint}}/win/bar.exe"><code>bar.exe</code></a> and save
	it beside <code>bar-export.bat</code> or somewhere in PATH.</p>

<h2>Git users</h2>
<p>To install bar into git repository use <code>bar git install</code></p>
<pre>
$ bar git install --endpoint={{.Info.JoinRPCEndpoints}}
</pre>
</body>
</html>
`

const spec_tpl string = `
<!DOCTYPE html>
<html>
<head>
<title>SPEC {{.ID}}</title>
</head>
<body>
<h1>:-)</h1>
<pre>
{{.ID}}
</pre>
<p>Send it to someone who knows what to do with it.</p>
<h2>Windows users</h2>
<p>To import spec download
	<a href="{{.Info.HTTPEndpoint}}/win/bar-import/{{.ID}}/bar-import-{{.ShortID}}.bat"><code>bar-import-{{.ShortID}}.bat</code></a>,
	save in root of the working tree and run.</p>
<p>This script is no-brain solution to import spec {{.ShortID}}. It automatically
download <code>bar.exe</code> if it is not found and download spec and BLOBs from bard</p>
<p>You can also download <a href="{{.Info.HTTPEndpoint}}/win/bar.exe"><code>bar.exe</code></a> and save
	it root of the working tree or somewhere in PATH.</p>
</body>
</html>
`
