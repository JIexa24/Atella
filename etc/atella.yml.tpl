hostname: ""
log:
  log_file: stdout
  log_level: info
pid_file: "/usr/share/atella/atella.pid"
proc_file: "/usr/share/atella/atella.proc"
connectivity: 1
reporter:
  hex_len: 10
  message_path: "/usr/share/atella/msg"
master: false
interval: 10s
net_timeout: 2
# channels:
#   - type: tgSibnet
#     enabled: true
#     address: "127.0.0.1"
#     port: 0
#     to: ""
#   - type: mail
#     auth: false
#     enabled: true
#     address: "127.0.0.1"
#     port: 25
#     username: "user"
#     password: "password"
#     # If ended with @hostname hostname will be replace to "hostname" parameter in 
#     from: "atella@hostname"
#     to: ""
security: "CodePhrase"