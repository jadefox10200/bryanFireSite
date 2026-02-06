# Bryan Fire & Safety Website

This is the official website for Bryan Fire & Safety Inc, serving the San Jose Bay Area with fire protection services.

## Features

- Static website serving
- Contact form with email notifications
- Security measures against path traversal attacks
- Health check endpoint for monitoring

## Running Locally

### Prerequisites

- Go 1.23 or higher
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

The server will start on `http://localhost:8080`

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

8. Set up Nginx as reverse proxy (recommended) or configure firewall to allow port 8080

### Domain Configuration on Porkbun

1. Point your A record to your Digital Ocean droplet IP
2. Set up email forwarding as described above
3. Consider setting up SSL/TLS with Let's Encrypt

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
