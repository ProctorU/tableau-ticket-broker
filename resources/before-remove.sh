#!/bin/bash

systemctl stop ticket-broker.service >/dev/null || true
systemctl disable ticket-broker.service >/dev/null || true
systemctl --system daemon-reload >/dev/null || true
