hostname: ""
log:
  log_file: stdout
  log_level: error
reporter:
  hex_len: 10
  message_path: "/usr/share/atella/msg"
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