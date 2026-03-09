# Task 2.0 Proof Artifacts - Flag Registration and Display Integration

## Test Results

```
$ go test ./cmd/git-workspace/ -v -run "TestFilterFlagsOnListCmd|TestListCmdFlagStackingWithRemoteRaw"

=== RUN   TestFilterFlagsOnListCmd
=== RUN   TestFilterFlagsOnListCmd/remote-raw
--- PASS: TestFilterFlagsOnListCmd (0.00s)
    --- PASS: TestFilterFlagsOnListCmd/remote-raw (0.00s)
=== RUN   TestListCmdFlagStackingWithRemoteRaw
--- PASS: TestListCmdFlagStackingWithRemoteRaw (0.00s)
PASS
```

## CLI Output - Formatted Remote (`-r`)

```
$ gws list -r | head -10

NAME                                                     REMOTE
-------------------------------------------------------  ------------------------------------------------------------------------------------------
bgs-feedback-service                                     https://github-aa/AAInternal/bgs-feedback-service.git
bgs-prime                                                https://github-aa/AAInternal/bgs-prime.git
gws                                                      https://github.com/daileyo/gws
```

## CLI Output - Raw Remote (`-R`)

```
$ gws list -R | head -10

NAME                                                     REMOTE
-------------------------------------------------------  ---------------------------------------------------------------------------------------------------------------------------------------------------------
bgs-feedback-service                                     git@github-aa:AAInternal/bgs-feedback-service.git
bgs-prime                                                git@github-aa:AAInternal/bgs-prime.git
gws                                                      git@github.com:daileyo/gws
```

## CLI Output - Override (`-rR`)

```
$ gws list -rR | head -5

NAME                                                     REMOTE
-------------------------------------------------------  ---------------------------------------------------------------------------------------------------------------------------------------------------------
bgs-feedback-service                                     git@github-aa:AAInternal/bgs-feedback-service.git
```

Raw URLs shown when both flags used — `-R` overrides `-r`.

## JSON Output - Formatted

```
$ gws list -r -o json | python3 -c "..." | head

[{"name": "bgs-feedback-service", "remote_url": "https://github-aa/AAInternal/bgs-feedback-service.git"}]
```

## JSON Output - Raw

```
$ gws list -R -o json | python3 -c "..." | head

[{"name": "bgs-feedback-service", "remote_url": "git@github-aa:AAInternal/bgs-feedback-service.git"}]
```
