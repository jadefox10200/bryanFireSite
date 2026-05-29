#!/bin/bash

set -e

echo "Building Go app..."
go build -o bryanfire

echo "Installing systemd service..."
cp bryanfire.service /etc/systemd/system/bryanfire.service

echo "Reloading systemd..."
systemctl daemon-reload

echo "Enabling service..."
systemctl enable bryanfire

echo "Restarting service..."
systemctl restart bryanfire

echo "Done. Status:"
systemctl status bryanfire --no-pager
