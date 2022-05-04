## Election API

### How to run with docker-compose

`$ docker-compose up`

### How to run on local

- Run redis. you can run it via docker `$ make local-redis`.
- Copy .env.example and edit it according to your redis server.
- `$ make run`

After you running API you can access playground via [http:/localhost:3000/pg](http:/localhost:3000/pg)

### Normal flow
You have to open an election with mutaion `open` then you can vote with your valid IDCard. 

### How to access GraphiQL document
- open [http:/localhost:3000/pg](http:/localhost:3000/pg)
- click on the top right `< Docs` as show in picture below ![access document](graphiql.jpg)

