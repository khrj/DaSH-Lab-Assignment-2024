docker pull hello-world
docker image save hello-world > dump.tar
mkdir dump
tar xvf dump.tar --directory=./dump
