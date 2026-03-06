# Task 3.0 Proof Artifacts - Verbose Levels (`-v` and `-vv`)

## Test Results

All tests pass:

```
ok  	github.com/daileyo/gws/cmd/git-workspace	0.405s
```

## Verbose Count Flag Tests

```
--- PASS: TestVerboseCountFlag/no_verbose (count=0)
--- PASS: TestVerboseCountFlag/single_-v (count=1)
--- PASS: TestVerboseCountFlag/double_-vv (count=2)
--- PASS: TestVerboseCountFlag/triple_-vvv (count=3)
```

## Verbose Level Column Override Tests

```
--- PASS: TestVerboseLevelColumnOverrides/verbose_0_-_no_overrides
  → No Show* flags set
--- PASS: TestVerboseLevelColumnOverrides/verbose_1_-_stored_data_columns
  → ShowType, ShowVisibility, ShowTags, ShowPath = true
  → ShowStatus, ShowUser, ShowRemote = false
--- PASS: TestVerboseLevelColumnOverrides/verbose_2_-_all_columns
  → All Show* flags = true
--- PASS: TestVerboseLevelColumnOverrides/verbose_1_with_filter_-_filter_preserved,_all_stored_columns_shown
  → FilterType="github" preserved, all stored columns shown
```

## Key Implementation Details

- `CountVarP` registered in Task 1.0: `-v` → 1, `-vv` → 2
- `VerboseLevel` field in `ListOptions` populated by `parseDualPurposeFlags()`
- Verbose level 1: overrides ShowType/ShowVisibility/ShowTags/ShowPath to true
- Verbose level 2: additionally overrides ShowStatus/ShowUser/ShowRemote to true
- Filter flags are preserved — verbose only affects column display, not filtering
- When verbose + column flags: verbose wins on display, filter still applies
