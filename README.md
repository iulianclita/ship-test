## Ship Test

### Description
This tool is used to distribute items amongst various pack sizes in order to comply with the following rules:
1. Only whole packs can be sent. Packs cannot be broken open.
2. Within the constraints of Rule 1 above, send out no more items than necessary to
fulfil the order.
3. Within the constraints of Rules 1 & 2 above, send out as few packs as possible to
fulfil each order.

### Usage
In order to use this tool just start the local HTTP server by running `go run main.go` from the project root.
Once the server is successfully started, navigate to `localhost:8080/ship?order_qty=520&pack_zies=100,200,300` for example and see the result in json format. The link above calculates the distribution of 500 ordered items while having at your disposal pack sizes of 100, 200 and 300 respectively. You should get back something like
```json
{
    300: 1, 
    200: 1, 
    100: 1
}
``` 
You can play with the query params however you want in order to obtain the desired values.

To stop the server, just hit `Ctrl + C` from the command line.