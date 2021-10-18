# nimble-interview-backend

### Video Demo

https://youtu.be/C6f_cGsVIZk (only people with link can view)

### Dependencies to Run Independently
1. Docker
2. Docker Compose
3. make

### Dependencies to Run w/ Mobile App
1. All of above 
2. Ngrok

### Tests
```bash
make test
```

### Run Independently
```bash
make run
```

### Run w/ Mobile App
1. `make run` (in terminal session #1)
1. `ngrok http 1234` (in terminal session #2)
2. Follow intructions on https://github.com/Krajiyah/nimble-interview-frontend
3. Use ngrok url in prompt provided in app UI to point your app to your running instance of the backend

### Project Structure
1. Inspired by: https://github.com/golang-standards/project-layout
2. Main Entrypoint: cmd/main.go
3. Environment variables and secrets: configs (in production should be generated on fly and not committed to repo)
4. Serverside logic: internal/models, internal/routes
5. Serverside unit tests: internal/routes/*_test.go
6. CI/CD setup w/ Circle CI: .circleci/config.yml
