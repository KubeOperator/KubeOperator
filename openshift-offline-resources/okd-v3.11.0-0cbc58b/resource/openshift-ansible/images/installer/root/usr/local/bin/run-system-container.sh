#!/bin/sh

export ANSIBLE_LOG_PATH=/var/log/ansible.log
exec ansible-playbook -i /etc/ansible/hosts ${OPTS} ${PLAYBOOK_FILE}
