## Setup GitSync with Systemd

```bash
sudo cp gitsync.service /etc/systemd/system/gitsync.service
sudo systemctl daemon-reload
sudo systemctl enable --now gitsync.service
```

Ensure that the GitSync binary is installed in this path `/usr/bin/git-sync`. If its installed in a different path, you can change the path in the `ExecStart` line in the `gitsync.service` file.
