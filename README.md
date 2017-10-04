## Description

I had to write an application that reads `\n`-separated URLs from standard input, loads specified pages and counts occurrences of the word `"Go"` in response body.

Processing must be parallel with at most `maxWorkers`, but you can't use a pool of workers: if there's currently no jobs, no worker routines should be running.

It's a weird case for Go (a pool of workers would be the best solution, I think), but the code gets sort of interesting with this setup.

## Run

Run with:

```bash
echo -e 'https://golang.com\nhttps://golang.com\nhttps://pepeissad.com' | go run main.go
Failed to get https://pepeissad.com: Get https://pepeissad.com: dial tcp: lookup pepeissad.com on 10.19.211.11:53: no such host
https://golang.com: 9
https://golang.com: 9
18

```

One url is broken, an error will be logged.