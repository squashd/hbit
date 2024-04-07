# hbit

hbit is a gamified task management application, built with Go and RabbitMQ
on the backend and NextJS and Tailwind on the frontend.

## Getting started

Why would you? This is a personal (learning) project and a knockoff of a fully
fledged application. Nonetheless, should you choose to run this, you will need
to take these steps

1. Copy the example env file

```bash
mv .env.example .env
```

2. Populate the database environment variables—several!— with your serverless
SQLite service of choice and generate a new JWT secret (copy the command below).

```bash
openssl rand -base64 32
```

3. Fire off the build script
```bash
./scripts/buildscript; docker-compose up --build -d
```

4. Marvel at the creation.

5. Ask why there's no frontend here (it's hideous and non-functional)

## Roadmap

- [x] Set up basic services (auth, rpg, task)
- [x] Pub/sub messaging
- [x] Realize mistakes were made
- [x] Set up reverse proxy / API gateway
- [x] Orchestrate tasks with RPG resolution and registration (kinda)
- [ ] Implement quest mechanics
- [ ] Set up a show (?)
- [ ] Make a purdy frontend
- [ ] ???
- [ ] Profit

## Motivation

hbit is a [Habitica](https://habitica.com/) clone—a task and habit-building
application with gamified rpg elements. It started as an exploration of 
resource driven design. Not for any particular architectural fit, just as a
point of reference of what an end application could or should look like.

I quickly realized I'd made a poor decision and moved on to trying an
event-based microservices architecture... which I did poorly. The RPG service
became wholly dependent on the task service as a driver of events which meant
I had to asynchronously update the user. For this I decided to use a websocket
connection (?!?). But what better way to learn than to fail.

The architecture as of now (2024-07-04) consists of an API gateway with an
auth service and orchestration between registration and task resolutions (do/undo).
