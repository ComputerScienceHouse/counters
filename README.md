# Counters ![Static Badge](https://img.shields.io/badge/weather-%23b16ded?style=flat&logo=github&logoColor=black&labelColor=0%2C0%2C0&link=https%3A%2F%2Fgithub.com%2FweatherGod3218%2F)
## Day's Since Last Useless House Service was Made: 0

A simple interface for logging and displaying time since a certain event has occured, allowing for CSHer's to track how long its been since a prior event. Heavily inspired off the "Days since last..." meme format.

This project uses Golang, [Gin](https://gin-gonic.com/en/), HTML/CSS, and Javascript.

This project ALSO uses a modified version of bootstrap 5! Check it out [here!](https://github.com/ComputerScienceHouse/csh-material-bootstrap/tree/bootstrap-5)

## Local Development
### Installing
1. Clone and cd into the repo: git clone https://github.com/ComputerScienceHouse/counters
>> (OPTIONAL): Make another branch if your working on a large thing!

### Setup
1. Make sure you have docker installed
>> (OPTIONAL): You can use docker compose as well!!
2. Copy the .env.template file, rename it to .env and place it in the root folder
3. Ask an RTP for counters secrets, add them to the .env accordingly

### Run

Counters is containerized through a docker file.

1. Build the docker file
```
    docker build -t Counters .
```
2. Run the newly built docker on port 8000
```
    docker run -p 8080:80 Counters
```

### Docker Compose

Counters also has support for Docker Compose, a extended version of docker that simplifies the steps.

(This is a really cool thing! If you use docker often, check it out!)
```
    docker compose up --build
```

## TODO
- [] Add Swagger
- [] Add API Integration
- [] Add Logging (similar to Votes "Actions")
- [] Time Zone support

  