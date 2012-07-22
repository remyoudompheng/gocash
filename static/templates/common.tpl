{{ define "common" }}
<!DOCTYPE html>
<html>
    <head>
        <title>{{ .Title }}</title>
        <link type="text/css" rel="stylesheet" href="https://ajax.googleapis.com/ajax/libs/jqueryui/1.8.21/themes/smoothness/jquery-ui.css">
        <link type="text/css" rel="stylesheet" href="/static/gocash.css">
        <!--
        <script src="https://ajax.googleapis.com/ajax/libs/jquery/1.7.2/jquery.min.js"></script>
        <script src="https://ajax.googleapis.com/ajax/libs/jqueryui/1.8.21/jquery-ui.min.js"></script>
        -->
        <script type="text/javascript">
            {{ template "script" . }}
        </script>
    </head>
    <body>
        {{ template "body" . }}
    </body>
</html>
{{ end }}

