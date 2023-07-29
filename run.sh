sudo docker build -f Dockerfile -t forum . # create forum image based on Dockerfile
echo -e ""
sudo docker run --name forum -p 8080:8080 forum # create and run the container from forum image