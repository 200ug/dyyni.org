## Deployment

1. Allow rootless binding to 443: `echo 'net.ipv4.ip_unprivileged_port_start=443' | sudo tee -a /etc/sysctl.conf && sudo sysctl -p`
2. Upload the backend files to remote: `./scripts/upload.sh username@host`
3. SSH into remote and deploy: `./scripts/deploy_backend.sh` (builds the image and starts the container)
4. Push to `master` or trigger the Actions runner manually to build and push the frontend to CF Pages

Rootless Podman deployment likely requires the following packages: `podman podman-compose slirp4netns fuse-overlayfs uidmap dbus-user-session catatonit` (and enabling lingering: `sudo loginctl enable-lingering <username>`. Firewalling should be configured to block 8081 to not expose the local healthcheck endpoint just in case.

Notably in this case the incoming traffic filtering (CF IPs only) is done on application level, but it could just as well be done on OS level with e.g. `ufw`.

---

###### Mirrors: [Codeberg](https://codeberg.org/2ug/dyyni.org) / [Github](https://github.com/200ug/dyyni.org)
