# SSE Basic Example (Go)

This is a simple example of **Server-Sent Events (SSE)** implemented in Go using the standard `net/http` package.

---

## What is SSE?

**Server-Sent Events (SSE)** allow a server to push real-time updates to clients over a single, long-lived HTTP connection.  
The client receives automatic event updates without needing to repeatedly poll the server.

---

## How this example works

The server listens on `/events` and sends a new message every second for 10 seconds.

---

## How to Run

1. Clone or download this folder.
2. Make sure you have Go installed (version 1.16+ recommended).
3. In In the project folder, run:

```bash
go run main.go


## How to test
Open your browser's developer console and run this JavaScript snippet:

```
const es = new EventSource("http://localhost:8080/events?channel=sports"); // spcify channel here to get only that particular channel events. 
es.onmessage = function(event) {
    console.log("Received:", event.data);
};
```


You will see messages streaming from the server every second, like:
```
Received: Message number 0 at 2025-06-08T14:00:00Z
Received: Message number 1 at 2025-06-08T14:00:01Z
...
```

##Points to remember

1. Content-Type : Tells the client this is an SSE stream.
2. Connection: Keeps the connection open for continuous streaming.
3. Normally, HTTP responses are buffered and sent only when the handler finishes or the buffer fills. SSE requires sending partial data immediately to the client. *http.Flusher* interface allows forcing the server to flush buffered data to the client right away.