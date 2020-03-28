# Simple Auth

Simple Auth is a stupidly simple authentication mechanism that was built to serve one very specific goal and flow.

A site I was consulting on wanted to host a part of their React app behind a username and password *that is shared across the organization*. They hosted their React app on Netlify. Since they didn't have access to Apache or NGINX, and there wasn't a server-side user management system, they asked for a simple authentication system to provide a roadblock to anyone trying to access that part of the site. The alternative would have been to hardcode the username and password in the minified JS.

This project was quickly written, validated, and rolled into a very cheap server. The React App now calls this server, checks the validation of the information, and then returns the appropriate HTTP codes.

## Usage

The best way to run this is to use the [Docker Image](https://hub.docker.com/r/kevineaton/simple-auth):

`docker pull kevineaton/simple-auth`

To run it, you would pass in the environment variables below.

## Environment Variables

- `SA_API_PORT`
  - The API port you want to server to listen on. Defaults to `8090`
- `SA_USERNAME`
  - The username to check against. Failure to provide this results in a panic.
- `SA_PASSWORD`
  - The password to check against. Failure to provide this results in a panic.
- `SA_TOKEN_SALT`
  - The SALT for the JWT signing. Failure to provide this results in a panic.
- `SA_RATE`
  - The rate limit. Defaults to 1000.
- `SA_TOKEN_EXPIRES_MINUTES`
  - How long the token should be valid for. Defaults to 1 day (1440).

### A Note On Encryption

As the data is passed in through the environment, it is unencrypted. As such, if someone inspects the environment, they would likely be able to sniff the information. Granted, if they have access to your environment, something else has likely gone awry, but I recommend encrypting the environment. I personally prefer [Anisble Vault](https://docs.ansible.com/ansible/latest/user_guide/vault.html) but that implementation is up to you.

## Building the Docker Image

Building the Docker image is a two step process. First, you will want to build the main image:

`docker build -t "kevineaton/simple-auth:latest" .`

Next, we want to slim it way down, so we build a slimmed down image:

`docker build -f Dockerfile-slim -t "kevineaton/simple-auth"`

## Contributing

Contributions are welcome, although not necessarily needed. This was built to solve a very specific situation in which security of the data is non-compromising. In other words, if the protected section did leak, there wouldn't be any lost personal data or IP. As such, the features list doesn't necessarily need to grow, but if you can think of a need that this tool can help with, by all means!

## Used libraries

- [Chi](https://github.com/go-chi/chi)
- [Toolbooth](https://github.com/didip/tollbooth)
- [CORS](https://github.com/goware/cors)
- [Testify](https://github.com/stretchr/testify)
