# Bryan Fire & Safety Website

This is the official website for Bryan Fire & Safety Inc, serving the San Jose Bay Area with fire protection services.

## Features

- Static website serving
- Contact form with email notifications
- Security measures against path traversal attacks
- Health check endpoint for monitoring

## Running Locally

### Prerequisites

- Go 1.23+ (tested with Go 1.24)
- SMTP credentials for sending emails

### Setup

1. Clone the repository
2. Copy `.env.example` to `.env` and configure your SMTP settings:
   ```bash
   cp .env.example .env
   ```

3. Edit `.env` with your actual SMTP credentials

4. Install dependencies:
   ```bash
   go mod download
   ```

5. Run the server:
   ```bash
   go run main.go
   ```

The server will start on `http://localhost:8080` for local development.

**Note:** In production, the server runs on HTTP internally (localhost:8080) while Nginx handles HTTPS externally. See "HTTPS Setup for Production Website" section below for deployment details.

## Production Deployment on Digital Ocean

### Email Configuration

For the contact form to work, you need to configure SMTP settings. Options include:

1. **Gmail SMTP** (requires App Password):
   - SMTP_HOST: smtp.gmail.com
   - SMTP_PORT: 587
   - SMTP_USER: your-gmail@gmail.com
   - SMTP_PASS: your-app-specific-password

2. **SendGrid**:
   - SMTP_HOST: smtp.sendgrid.net
   - SMTP_PORT: 587
   - SMTP_USER: apikey
   - SMTP_PASS: your-sendgrid-api-key

3. **Mailgun**, **AWS SES**, or other SMTP providers

### Email Forwarding Setup

The email forwarding from `info@bryanfire.com` → `bryanfiresafetyinc@gmail.com` should be configured at your domain registrar (Porkbun):

1. Log into Porkbun.com
2. Go to your domain settings for bryanfire.com
3. Set up email forwarding:
   - From: info@bryanfire.com
   - To: bryanfiresafetyinc@gmail.com

This is handled at the DNS level and is separate from the web server.

### Deploying to Digital Ocean

1. Create a Droplet (Ubuntu recommended)

2. Install Go on the droplet

3. Clone your repository

4. Set up environment variables:
   ```bash
   export PORT=8080
   export SMTP_HOST=smtp.gmail.com
   export SMTP_PORT=587
   export SMTP_USER=your-email@gmail.com
   export SMTP_PASS=your-password
   export TO_EMAIL=info@bryanfire.com
   export GIN_MODE=release
   ```

5. Build and run:
   ```bash
   go build -o bryanfire main.go
   ./bryanfire
   ```

6. Set up as a systemd service (recommended):

Create `/etc/systemd/system/bryanfire.service`:
```ini
[Unit]
Description=Bryan Fire Safety Website
After=network.target

[Service]
Type=simple
User=www-data
WorkingDirectory=/path/to/bryanFireSite
Environment="PORT=8080"
Environment="SMTP_HOST=smtp.gmail.com"
Environment="SMTP_PORT=587"
Environment="SMTP_USER=your-email"
Environment="SMTP_PASS=your-password"
Environment="TO_EMAIL=info@bryanfire.com"
Environment="GIN_MODE=release"
ExecStart=/path/to/bryanFireSite/bryanfire
Restart=always

[Install]
WantedBy=multi-user.target
```

7. Start the service:
   ```bash
   sudo systemctl daemon-reload
   sudo systemctl enable bryanfire
   sudo systemctl start bryanfire
   ```

8. **CRITICAL FOR PRODUCTION**: Configure HTTPS for bryanfire.com (see HTTPS Setup below)

### HTTPS Setup for Production Website

**Why HTTPS is Required:**
As a professional fire safety company serving Bay Area businesses, bryanfire.com must use HTTPS to:
- Protect customer contact form data (names, emails, phone numbers)
- Build trust with commercial clients
- Meet modern browser security standards
- Improve SEO rankings for fire safety services in San Jose

**Recommended Approach: Nginx Reverse Proxy with Certbot**

This setup keeps your Go application simple while Nginx handles HTTPS encryption.

**Step 1: Install Nginx and Certbot on your Digital Ocean droplet**
```bash
sudo apt update
sudo apt install nginx certbot python3-certbot-nginx
```

**Step 2: Create Nginx configuration for bryanfire.com**

Create file: `/etc/nginx/sites-available/bryanfire`
```nginx
server {
    server_name bryanfire.com www.bryanfire.com;
    
    location / {
        proxy_pass http://127.0.0.1:8080;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }
}
```

**Step 3: Enable the configuration**
```bash
sudo ln -s /etc/nginx/sites-available/bryanfire /etc/nginx/sites-enabled/
sudo nginx -t
sudo systemctl restart nginx
```

**Step 4: Obtain free SSL certificate from Let's Encrypt**
```bash
sudo certbot --nginx -d bryanfire.com -d www.bryanfire.com
```

Certbot will automatically:
- Obtain SSL certificates for your domain
- Update Nginx configuration to use HTTPS
- Set up automatic certificate renewal

**Step 5: Verify HTTPS is working**
- Visit https://bryanfire.com in your browser
- Check for the padlock icon showing secure connection
- Test the contact form to ensure it works over HTTPS

**Certificate Auto-Renewal:**
Certbot sets up a systemd timer to automatically renew certificates before expiration. Verify with:
```bash
sudo systemctl status certbot.timer
```

### Domain Configuration on Porkbun

1. Point your A record to your Digital Ocean droplet IP address
2. Point www subdomain to same IP (A record or CNAME)
3. Configure email forwarding: info@bryanfire.com → bryanfiresafetyinc@gmail.com
4. After DNS propagates (may take up to 48 hours), run Certbot as shown above

## API Endpoints

- `GET /` - Serves the homepage
- `POST /contact` - Handles contact form submissions
- `GET /health` - Health check endpoint
- `GET /{filename}` - Serves static assets (CSS, JS, images)

## Security Features

- Path traversal protection
- Hidden file blocking (dotfiles)
- Source code access prevention
- MIME type restrictions for static files
- Input validation on contact forms

## License

© 2024 Bryan Fire & Safety Inc. All rights reserved.
