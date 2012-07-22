{{ define "script" }}
{{ end }}

{{ define "body" }}
<h1>Gocash: account overview</h1>

<table>
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
        <td>{{ with $bal := index $.Book.Balance $acct }}{{ money $bal }}{{ end }}</td>
    </tr>
    {{ end }}
</tbody>
{{ end }}
