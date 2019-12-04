package main

import (
	"bytes"
	"html/template"
	"io"
)

const ChangeQueryKeyHtml  = `
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <title>Change Unsplash</title>
</head>
<body>
    <h1>Change Get Unsplash Query Key</h1>
    <form action="/change_query_key" method="post">
        <label for="key">query key:</label> <input type="text" id="key" name="key" value="{{ . }}" /> <br/>
        <input type="submit" value="submit" />
    </form>
</body>
</html>
`

func GetChangeIndex() (io.Reader, error) {

	tmp, err := template.New("webpage").Parse(ChangeQueryKeyHtml)
	if err != nil {
		return nil, err
	}

	var index bytes.Buffer
	if err := tmp.Execute(&index, *PhotoQueryKey); err != nil {
		return nil, err
	}
	return &index, nil
}
