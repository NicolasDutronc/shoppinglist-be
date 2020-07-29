# JSON encoded server sent messages

## Usage with curl

Start the server

Start a consmer :

```
curl -k https://localhost:8080/lists/1
```

Produce :

```
curl -k -d '{"name": "chocolate", "quantity": "a lot"}' https://localhost:8080/lists/1/send
```
