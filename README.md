# Forum-advanced-features

## Description

This project consists in creating a web forum that allows :

- communication between users.
- associating categories to posts.
- liking and disliking posts and comments.
- filtering posts.
- create a post containing an image.
- edit/delete users posts and comments

## Authors

Elina Razena @elinana

Litvintsev Anton @Antosha7

Alper Balaban @alpbal

Orel Margarita @maggieeagle

## Usage
  
To run the web-site on local machine:

- Download the repository

- Run with a command `./run.sh`

- Open [http://localhost:8080/](http://localhost:8080/) in browser

- Register or login with a test user(test@gmail.com, 1234)

- Stop with a command `./stop.sh`

To work with a Docker manually:

- Build image with `sudo docker build -t forum .`

- Run image with `docker container run -p 8080:8080 forum:latest`

You can login with other profiles as well:

- cowboy@gmail.com blacksheep

- snaphappyphotographer@email.com ZoomZoomZap

- rodeoqueen@email.com kingtomyqueen

## Implementation details

- SQL3

- Adaptive web-design