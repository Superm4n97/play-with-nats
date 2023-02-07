# nats
## run
Run the consumer for this_is_nats server:
```bash
export SERVER=????????
```
### credentials
for credentials export `CRED` nats credential file path or provide nats `USER` and `PASS`

```bash
export CRED=????????
```
or
```bash
export USER=????????
export PASS=????????
```
then run the consumer
```bash
go run consumer.go
```

nats server examples: 
```bash
export natsAddress = "localhost:30222"
export natsServerAddress = "nats://127.0.0.1:30222"
export natsServerAddress = "this-is-nats.appscode.ninja:4222"
export natsServerAddress = "nats://nats.appscode.ninja:4222"
```