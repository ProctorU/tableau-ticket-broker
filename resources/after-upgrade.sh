#!/bin/bash

systemctl --system daemon-reload >/dev/null || true
if ! systemctl is-enabled ticket-broker.service >/dev/null
then
    systemctl enable ticket-broker.service >/dev/null || true
    systemctl start ticket-broker.service >/dev/null || true
else
    systemctl restart ticket-broker.service >/dev/null || true
fi
