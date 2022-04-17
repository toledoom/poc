## Proof of concept: gRPC simple game server
This poc is an exercise to practice a bit of one of my favourite languages (golang). It's far from perfect, but it's a start :)
This server basically simulates a battle server and offers a gRPC API to create players, start and finish a battle and get the leaderboards.
This gRPC server would act as an authoratitative entity between the real-time battle server (where the players connect and actually play) and the players' clients.

## Common requirements
- DockerEngine: 20.10.12
- Compose: 1.29.2
- GNU Make 3.81

## Instructions for Docker
- Start the server using Docker: ```make up```
- Run the gRPC client: ```make docker-run-client```

## Instructions for local development
- Install go 1.17.5
- Start the "dockerized" infrastructure: ```make docker-infra```
- Run the server from your favourite IDE (I use VSCode)
- Run the client in local: ```make run-client```

## Stop the server
- Soft stop: ```make down```
- Hard stop (with data wipe out): ```make down-prune```

## TODO
- TESTS!!! (I obviously didn't follow TDD)
- Decouple the gRPC layer from the domain layer. Currently, the gRPC server orchestrates the calls to repositories and handles a lot of logic. In case we wanted to create a HTTP API, most of the code should have to be repeated. On top of that, it makes it more difficult to unit test some of the parts.
- Authentication and authorization
- A HTTP API
- A better inversion of control (for instance, the domain event publishing)
- ...