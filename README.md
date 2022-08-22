# dockerhub-ls
Dockerhub-ls allows you to get the name of all images (with tags) of a user on Dockerhub. You can also batch run this by passing a list of users to dockerhub-ls.

## Usage

Getting all images from one user:
```
echo netflixoss | dockerhub-ls
```

Getting images from many users:
```
cat users.txt | dockerhub-ls
```

## Installation

First, you'll need to [install go](https://golang.org/doc/install).

Then run this command to download + compile dockerhub-ls:
```
go install github.com/edivangalindo/dockerhub-ls@latest
```

You can now run `~/go/bin/dockerhub-ls`. If you'd like to just run `dockerhub-ls` without the full path, you'll need to `export PATH="/go/bin/:$PATH"`. You can also add this line to your `~/.bashrc` file if you'd like this to persist.
