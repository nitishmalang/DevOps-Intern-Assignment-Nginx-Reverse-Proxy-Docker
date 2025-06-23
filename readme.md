### 🧪 **DevOps Intern Assignment: Nginx Reverse Proxy + Docker**

This project sets up uses Docker Compose to set up a system with two backend services (Golang and Python Flask) and an Nginx reverse proxy, demonstrating routing, health checks, logging and automated testing for a DevOps intern assignment.

## Setup Instructions
1. **Clone the Repository**:
   
    `git clone <your-repo-url>`
    `cd nginx-reverse-proxy-assignment`

2. **Install Prerequisites**:
   
     Docker and Docker Compose (Docker Desktop on Windows).

3. **Build and Run:**

    `docker-compose up --build`

4. **Access Services:**

     *Service 1 (Golang)*:

http://localhost:8080/service1/hello → {"message":"Hello from Service 1"}

http://localhost:8080/service1/ping → {"service":"1","status":"ok"}

   *Service 2 (Python Flask)*:

http://localhost:8080/service2/hello → {"message":"Hello from Service 2"}

http://localhost:8080/service2/ping → {"service":"2","status":"ok"}

5. **Run Tests:**

`.\test.ps1`  # Windows
`./test.sh`   # Linux

6. **View Logs:**

`docker-compose logs nginx`

7. **Stop Services:**

`docker-compose down`

8. **Bonus Features**
   
---***Health Checks:***

service_1: Verifies port 8001 with nc, compatible with distroless.

service_2: Checks /ping endpoint with curl.


---***Logging Clarity:***

Nginx logs include client IP, timestamp, request path, status code, and user agent.

Service logs show startup details (e.g., ports 8001, 8002).


---***Clean Docker Setup:***

Multi-stage builds for service_1 (minimal distroless runtime).

Lightweight images: python:3.11-slim, nginx:alpine.

Reproducible Python dependencies with uv.

Isolated bridge network (app-network).

---***Security:***

Avoided vulnerable golang:1.22-alpine (e.g., CVE-2025-22871).

Pinned uv==0.7.14 for stability.

9. **Troubleshooting:**
   
---***Port Conflicts:***

Check: `netstat -aon | findstr "8080 8001 8002"`

Resolve: `docker stop $(docker ps -q)`

***Health Check Failures:***

Verify: `curl http://localhost:8001/ping` and `curl http://localhost:8002/ping`

Check status: `docker inspect service_1 | findstr "Health"`

---***Build Issues:***

`Clear cache: docker builder prune`

Pre-pull images:
`docker pull golang:1.23`
`docker pull python:3.11-slim`
`docker pull gcr.io/distroless/base-debian12`

---***Nginx 404 Errors:***

Ensure correct paths (/service1/*, /service2/*).

Check logs: `docker-compose logs nginx`

10. **Directory Structure**

![Screenshot (907)](https://github.com/user-attachments/assets/d7bc7087-9c3c-43c1-9f7f-370271431024)







    

    
