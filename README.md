# Examitsu

## Deploy with HTTPS
Backend is at port 8080, frontend is at port 3000. The public port are only 443 and 80. We are using nginx as the reverse proxy to access brackend and frontend.

- Install nginx:
    ```
    sudo apt install nginx
    ```

- Save the nginx.conf
    ```
    sudo nano /etc/nginx/sites-available/examitsu
    sudo ln -s /etc/nginx/sites-available/examitsu /etc/nginx/sites-enabled/
    ```
- Test the configuration
    ```
    sudo nginx -t
    ```
- Reload nginx
    ```
    sudo systemctl reload nginx
    ```
- Configure SSL, we are using Let's Encrypt:
    ```
    sudo apt update
    sudo apt install certbot python3-certbot-nginx -y
    sudo certbot --nginx -d examitsu.net
    ```
- Visit https://examitsu.net
