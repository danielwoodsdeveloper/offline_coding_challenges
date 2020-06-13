# Offline Coding Challenges

Complete coding challenges. Anywhere. Anytime.

## Status

*Still being developed.*

To do:
* Refactor code
* More graceful error handling
* Add more tests

## Building and Running

`docker build . -t offline-coding-challengs`

`docker run -p 8080:8080 -v /var/run/docker.sock:/var/run/docker.sock -v cc-data:/app offline-coding-challengs`