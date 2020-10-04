# Dockerfile References: https://docs.docker.com/engine/reference/builder/

# cd /data/workspaces/go/src/yari
# docker build --progress=plain -t=yari ./
# docker create --net="host" --name yariContainer yari
# docker run -t yari /YARI/yari
# docker logs -f <container>
# ======== explore this filesystem using bash (for example)
# docker run -t -i yari /bin/bash


# Start from the latest golang base image
FROM ubuntu

# Add Maintainer Info
LABEL maintainer="foad momtazi"

# Set the Current Working Directory inside the container
WORKDIR /YARI

# Copy the source from the current directory to the Working Directory inside the container
COPY . .
# copy the non standard libs to /usr/lib
COPY ./usr/lib/* /usr/lib/
COPY ./usr/bin/* /usr/bin/

# now remove the non standard libs from the current folder
RUN rm -rf ./usr/lib
RUN rm -rf ./usr/bin

# Command to run the executable
CMD ["yari"]
