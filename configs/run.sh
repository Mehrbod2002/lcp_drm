htpasswd -c ../tmp/htpasswd admin

sudo mkdir -p /etc/nginx/auth
sudo cp /root/server/tmp/htpasswd /etc/nginx/auth/htpasswd

sudo chown www-data:www-data /etc/nginx/auth/htpasswd
sudo chmod 644 /etc/nginx/auth/htpasswd

cd frontend/ && npm install next-auth axios