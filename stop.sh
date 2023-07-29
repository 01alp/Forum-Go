docker stop forum
docker rm forum
docker image prune -a
echo "Forum container stopped and deleted, unused images cleaned"