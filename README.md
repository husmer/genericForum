## genericForum is a web-based platform designed to facilitate user communication, allowing them to create posts, comments, associate categories with posts, like/dislike content, apply filters, and image upload (up to 20 MB)



## Built with Go and SQLite, this project was built to learn basics, session management, database manipulation, and Docker.

Run the server with command:

```console
go run .
```

Open the page in a web browser:

```console
http://localhost:8000/
```

#### Prerequisites

Docker Desktop should be installed on your machine.

Instructions how to do that:
- [Linux Mint, Ubuntu](https://docs.docker.com/desktop/install/ubuntu/)
- [MacOS](https://docs.docker.com/desktop/install/mac-install/)
- [Windows](https://docs.docker.com/desktop/install/windows-install/)

[General instructions](https://docs.docker.com/desktop/)

Build Docker Image:

```console
docker image build -t genericforum .
```

After building the image check your built image, list all of the Docker images:

```console
docker images
```

Run Docker Container:

```console
docker container run -p 8000:8000 genericforum
```

After running this command, you should be able to access the web application by navigating to the link in your web browser:
[localhost:8000](http://localhost:8000)

List all running containers:

```console
docker ps -a
```

Stop the container:

```console
docker stop <container ID>
```

If you want to remove the stopped container, you can use the docker rm command:

```console
docker rm <container ID>
```

If you want to stop and remove the container in a single command, you can use the -f option with docker rm:

```console
docker rm -f <container ID>
```

With this command you can delete all containers located on your PC

Be careful!

```console
sudo docker system prune -a
```
