# Live Subdomain Finder

![Project Logo/Image]

## Overview

The Live Subdomain Finder is a robust tool designed to simplify the discovery and monitoring of live subdomains for a specified domain. It provides a user-friendly solution for identifying active subdomains associated with a target domain, making it invaluable for security assessments, bug bounty programs, or general web presence analysis.

## Features

- **Live Subdomain Discovery:** Quickly enumerate live subdomains associated with a target domain.

- **Continuous Monitoring:** Monitor the availability of discovered subdomains for real-time updates.

- **Security Assessments:** Enhance penetration testing efforts by identifying potential attack surfaces.

- **Bug Bounty Programs:** Expand your bug bounty program scope by discovering additional attack vectors.

## Requirements

Ensure you have the following dependencies installed before using the Live Subdomain Finder:

- **amass:** (https://github.com/OWASP/Amass)
- **subfinder:**(https://github.com/projectdiscovery/subfinder)

## Usage

To run the Live Subdomain Finder, execute the following command:

```bash
git clone https://github.com/hacker301et/live-sub.git
cd live-sub
go run cmd/main.go
```

To run the Live Subdomain Finder Using Docker, execute the following command:

```bash
git clone https://github.com/hacker301et/live-sub.git
cd live-sub
docker build -t live-sub .
docker run --rm -it live-sub
```



