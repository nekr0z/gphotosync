## Installation

USER=$(id -u -n)

Install executable binary
`# install -m 0755 gphotosync /usr/local/bin/gphotosync`

Install systemd.service files
`# install -m 0644 system/gpohtosync@.service /etc/systemd/system/gphotosync@.service`
`# systemctl daemon-update`

Specify a backup folder in the `~/.config/gphotosync.env`
`$ echo "LIBRARY_PATH=/home/${USER}/photos" | tee /home/${USER}/.config/gphotosync.env`

Run gphotosync for first and input google access token
`/usr/local/bin/gphotosync`


Then you can enable and start the systemd.service

`systemctl start gphotosync@user`