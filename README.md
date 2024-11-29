# Examitsu

## Description
A simple MCQ test platform, made in Indonesian language, initially made to help my high school to run an offline paperless exam session, so the students can open this website, but not other sites. Basically, a google form like, but specially made to help exams.

## Features
- Manage exams.
- Authorize student's exam session, so, only one student can do one exam session at one time.
- Easily manage questions and answers.
- Easily get the scores.

## Demo
https://youtu.be/wRo37EnK3Js

## Deployment

### Architecture
- Worker, mysql, and redis in one VM --> "private" VM, private facing.
- Backend and frontend in another VM --> "public" VM, public facing.
- In this manner, if something happened with the "public" VM, we can just have a backup VM to replace it.

### Deploy!
- Replace all the environment variables secret, and prepare GCS service account.
- Clone the source code into both VMs.
- Private VM:
    ```
    docker compose up backend-mysql
    docker compose up backend-redis
    docker compose up backend-worker
    ```
- Public VM:
    ```
    docker compose up backend
    docker compose up frontend
    ```

### HTTPS
Backend is at port 8080, frontend is at port 3000. The public port are only 443 and 80. We are using nginx as the reverse proxy to access backend and frontend. In public VM, do the following (for example, I use `examitsu.net` as the domain):

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
