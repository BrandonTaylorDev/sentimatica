# Sentimatica
## Edge
Frontend microservice that receives requests and publishes them on the queue cluster.

## Worker
Backend microservice that receives messages on the queue cluster, handles the request, then returns a response.
