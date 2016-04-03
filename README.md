# Write Better

Write Better is a Golang application which takes in plaintext or an uploaded
file and returns a score with a report on how well the piece of text reads.
The lower the score the better the writing.

## Working with this project

I've included a Dockerfile to get things started that way. Once you have
Docker on your machine (not going to explain this here) then you can clone
this repo down and build it with:

    $ docker build -t write-better .

Once it's written, you can run the application with:

    $ docker run -p 8000:17644 --name wb --rm write-better

The app will be available at `127.0.0.1:8000`. When you're all done with
things, you can run:

    $ docker stop wb

This will stop the daemon and remove the container.

## Todo

- [X] Implement processors
- [X] Setup HTTP router to serve an application
- [X] Build handlers for templates, uploads, results
- [X] Update processors response to give JS ability to highlight index matches
- [X] Build text parser to highlight matches
- [X] Add popovers in frontend for messages with colors matching processor types
- [X] Finalize a website design
- [ ] Write tests
