# 🐦‍⬛ Notifier

A Go-based notification service for handling email notifications via Brevo.

### Overview

Notifier is a lightweight notification service built with Go that provides a simple API for sending email notifications. It leverages the Brevo platform for reliable email delivery and includes monitoring capabilities through New Relic.

### 🏗️ Features

- **Email Notifications**: Send transactional and marketing emails via Brevo
- **RESTful API**: Built with Echo framework for high performance
- **Monitoring**: Integrated with New Relic for performance tracking and observability
- **Lightweight**: Minimal dependencies, fast startup time

### Tech Stack

- **Language**: Go
- **Web Framework**: [Echo](https://echo.labstack.com/)
- **Email Service**: [Brevo](https://www.brevo.com/)
- **Monitoring**: [New Relic](https://newrelic.com/)

### Getting Started

#### Prerequisites

- Go 1.x or higher
- Brevo API key
- New Relic license key (optional, for monitoring)

#### 🛶 Installation

```bash
git clone git@github.com:melih-gulerb/notifier.git
cd notifier
go mod download
```

```bash
go run main.go
```

```bash
.env:
  BREVO_API_KEY="your_brevo_api_key"
  FROM_MAIL="from_mail"
  NEW_RELIC_LICENSE_KEY="your_new_relic_license_key"
  NEW_RELIC_APP_NAME="new_relic_app_name"
```



