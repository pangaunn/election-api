## Election API

### How to run with docker

`$ docker-compose up`

### How to run on local

- Run redis you can run it via docker `$ make local-redis`.
- Copy .env.example and edit it according to your redis server.
- `$ make run`

After you running API you can access playground via http:/localhost:3000/pg

### Note on Mutation vote
You must supply header with Authorization:{IDcard} in order to call it. For example Authorization:"1234567890123"

