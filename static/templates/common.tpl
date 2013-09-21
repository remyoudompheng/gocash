{{ define "common" }}
<!DOCTYPE html>
<html>
    <head>
        <title>{{ .Title }}</title>
        <link type="text/css" rel="stylesheet" href="/libs/bootstrap/css/bootstrap.min.css">
        <link type="text/css" rel="stylesheet" href="/static/gocash.css">
        <script src="/libs/jquery.min.js">
        <script src="/libs/bootstrap/js/bootstrap.min.js">
        <script type="text/javascript">
            {{ template "script" . }}
        </script>
    </head>
    <body>
        <nav class="navbar navbar-default" role="navigation">
            <div class="navbar-header">
                <a class="navbar-brand">Gocash</a>
            </div>
        </nav>
        <div class="container">
        {{ template "body" . }}
        </div>
    </body>
</html>
{{ end }}

