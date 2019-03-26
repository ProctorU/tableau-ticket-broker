#!/bin/bash

systemctl --system daemon-reload >/dev/null || true
systemctl enable ticket-broker.service >/dev/null || true
systemctl start ticket-broker.service >/dev/null || true
