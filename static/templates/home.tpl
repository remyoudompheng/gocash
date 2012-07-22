{{ define "script" }}
{{ end }}

{{ define "body" }}
<h1>Gocash</h1>

<p>Accounts:</p>
<ul>
{{ range $acct := $.Book.Accounts }}
<li>{{ $acct.Name }}</li>
{{ end }}
</ul>

<p>Transactions:</p>
<ul>
{{ range $trn := $.Book.Transactions }}
<li>{{ $trn.Date.Format "2006-01-02" }}: {{ $trn.Description }}
  <ul>
    {{ range $flow := $trn.Flows }}
    <li>{{ $flow.Account.Name }} {{ parsePrice $flow.Price | printf "%.2f" }}</li>
    {{ end }}
  </ul>
</li>
{{ end }}
</ul>
{{ end }}
