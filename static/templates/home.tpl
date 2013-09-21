{{ define "script" }}
{{ end }}

{{ define "body" }}
<h1>Gocash: account overview</h1>

<table class="table">
<thead>
    <tr>
       <th>Account</th>
       <th>Balance</th>
    </tr>
</thead>
<tbody>
    {{ range $acct := $.Book.Accounts }}
    <tr>
        <td><a href="/account/?name={{ $acct.Name }}">{{ $acct.Name }}</a></td>
        <td class="amount">{{ index $.Book.Balance $acct }} {{ $acct.Unit }}</td>
    </tr>
    {{ end }}
</tbody>
{{ end }}
