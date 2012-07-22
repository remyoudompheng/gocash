{{ define "common" }}
<!DOCTYPE html>
<html>
    <head>
        <title>{{ .Title }}</title>
        <script type="text/javascript">
            {{ template "script" . }}
        </script>
    </head>
    <body>
        {{ template "body" . }}
    </body>
</html>
{{ end }}

