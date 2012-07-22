{{ define "script" }}
{{ end }}

{{ define "body" }}
<h1>Account {{ .Account.Name }}</h1>

<p>Current balance: {{ index .Book.Balance .Account }} {{ .Account.Unit }}</p>

<h2>Transactions</h2>

<table class="ui-widget ui-widget-content">
    {{ $flows := index .Book.Flows .Account }}
    {{ $balance := cumul $flows }}
    <thead>
    <tr class="ui-widget-header">
        <th>Date</th>
        <th>Description</th>
        <th>Amount</th>
        <th>Balance</th>
    </tr>
    </thead>
    <tbody>
    {{ range $i, $flow := $flows }}
    <tr>
        <td>{{ $flow.Parent.Date.Format "2006-01-02" }}</td>
        <td>{{ $flow.Parent.Description }}</td>
        <td class="amount">{{ $flow.Price }}</td>
        <td class="amount">{{ index $balance $i }} {{ .Account.Unit}}</td>
    </tr>
    {{ end }}
    </tbody>
</table>

{{ end }}
