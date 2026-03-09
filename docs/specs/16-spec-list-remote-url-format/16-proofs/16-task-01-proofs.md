# Task 1.0 Proof Artifacts - URL Formatting Utility

## Test Results

```
$ go test ./internal/git/ -run TestFormatRemoteURL -v

=== RUN   TestFormatRemoteURL
=== RUN   TestFormatRemoteURL/SSH_GitHub
=== RUN   TestFormatRemoteURL/SSH_GitLab_with_subgroup
=== RUN   TestFormatRemoteURL/HTTPS_without_user_info
=== RUN   TestFormatRemoteURL/HTTPS_with_user_info
=== RUN   TestFormatRemoteURL/HTTPS_with_user_and_password
=== RUN   TestFormatRemoteURL/Azure_DevOps_SSH
=== RUN   TestFormatRemoteURL/Azure_DevOps_HTTPS_with_user_info
=== RUN   TestFormatRemoteURL/file_protocol_unchanged
=== RUN   TestFormatRemoteURL/empty_string_unchanged
=== RUN   TestFormatRemoteURL/SSH_protocol_URL_unchanged
=== RUN   TestFormatRemoteURL/HTTP_with_user_info
--- PASS: TestFormatRemoteURL (0.00s)
    --- PASS: TestFormatRemoteURL/SSH_GitHub (0.00s)
    --- PASS: TestFormatRemoteURL/SSH_GitLab_with_subgroup (0.00s)
    --- PASS: TestFormatRemoteURL/HTTPS_without_user_info (0.00s)
    --- PASS: TestFormatRemoteURL/HTTPS_with_user_info (0.00s)
    --- PASS: TestFormatRemoteURL/HTTPS_with_user_and_password (0.00s)
    --- PASS: TestFormatRemoteURL/Azure_DevOps_SSH (0.00s)
    --- PASS: TestFormatRemoteURL/Azure_DevOps_HTTPS_with_user_info (0.00s)
    --- PASS: TestFormatRemoteURL/file_protocol_unchanged (0.00s)
    --- PASS: TestFormatRemoteURL/empty_string_unchanged (0.00s)
    --- PASS: TestFormatRemoteURL/SSH_protocol_URL_unchanged (0.00s)
    --- PASS: TestFormatRemoteURL/HTTP_with_user_info (0.00s)
PASS
ok  	github.com/daileyo/gws/internal/git	0.607s
```

## Verification

All 11 test cases pass covering:
- SSH → HTTPS conversion (GitHub, GitLab with subgroups)
- HTTPS passthrough (no user info)
- HTTPS user info stripping (user only, user:pass)
- Azure DevOps SSH → HTTPS conversion
- Azure DevOps HTTPS user info stripping
- Unformattable URL passthrough (file://, empty string, ssh:// protocol)
- HTTP with user info stripping
