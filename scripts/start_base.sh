sudo /etc/init.d/fdfs_trackerd start
sudo /etc/init.d/fdfs_storaged start
sudo /usr/local/nginx/sbin/nginx

docker start mysql
docker start redis