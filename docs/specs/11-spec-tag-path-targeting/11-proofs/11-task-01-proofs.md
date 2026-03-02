# 11-task-01-proofs: Implement `findRepositoriesWithFilters` and Unit Tests

## CLI Output

### TestFindRepositoriesByPath (all cases)

```
$ go test ./cmd/git-workspace/ -run TestFindRepositoriesByPath -v
=== RUN   TestFindRepositoriesByPath
=== RUN   TestFindRepositoriesByPath/Prefix_match_returns_two_repos
=== RUN   TestFindRepositoriesByPath/Prefix_match_single_repo
=== RUN   TestFindRepositoriesByPath/Substring_fallback_when_no_prefix_match
=== RUN   TestFindRepositoriesByPath/Case-sensitive_no_match_on_wrong_case
=== RUN   TestFindRepositoriesByPath/Exact_case_match_for_mixed-case_path
=== RUN   TestFindRepositoriesByPath/No_match_returns_empty
--- PASS: TestFindRepositoriesByPath (0.00s)
    --- PASS: TestFindRepositoriesByPath/Prefix_match_returns_two_repos (0.00s)
    --- PASS: TestFindRepositoriesByPath/Prefix_match_single_repo (0.00s)
    --- PASS: TestFindRepositoriesByPath/Substring_fallback_when_no_prefix_match (0.00s)
    --- PASS: TestFindRepositoriesByPath/Case-sensitive_no_match_on_wrong_case (0.00s)
    --- PASS: TestFindRepositoriesByPath/Exact_case_match_for_mixed-case_path (0.00s)
    --- PASS: TestFindRepositoriesByPath/No_match_returns_empty (0.00s)
PASS
ok  	github.com/daileyo/gws/cmd/git-workspace	0.497s
```

### TestFindRepositoriesWithFilters (all cases)

```
$ go test ./cmd/git-workspace/ -run TestFindRepositoriesWithFilters -v
=== RUN   TestFindRepositoriesWithFilters
=== RUN   TestFindRepositoriesWithFilters/Repo_filter_only
=== RUN   TestFindRepositoriesWithFilters/Path_filter_only
=== RUN   TestFindRepositoriesWithFilters/AND_logic:_intersection_of_repo_and_path
=== RUN   TestFindRepositoriesWithFilters/AND_logic:_no_intersection_returns_empty
=== RUN   TestFindRepositoriesWithFilters/Both_empty_returns_all_repos
--- PASS: TestFindRepositoriesWithFilters (0.00s)
    --- PASS: TestFindRepositoriesWithFilters/Repo_filter_only (0.00s)
    --- PASS: TestFindRepositoriesWithFilters/Path_filter_only (0.00s)
    --- PASS: TestFindRepositoriesWithFilters/AND_logic:_intersection_of_repo_and_path (0.00s)
    --- PASS: TestFindRepositoriesWithFilters/AND_logic:_no_intersection_returns_empty (0.00s)
    --- PASS: TestFindRepositoriesWithFilters/Both_empty_returns_all_repos (0.00s)
PASS
ok  	github.com/daileyo/gws/cmd/git-workspace	0.497s
```

### TestFindRepositories (no regression to existing test)

```
$ go test ./cmd/git-workspace/ -run TestFindRepositories -v
=== RUN   TestFindRepositories
=== RUN   TestFindRepositories/Find_by_partial_name_-_multiple_matches
=== RUN   TestFindRepositories/Find_by_partial_name_-_single_match
=== RUN   TestFindRepositories/Find_by_exact_path
=== RUN   TestFindRepositories/Find_by_partial_name_-_case_insensitive
=== RUN   TestFindRepositories/Find_by_exact_name_match
=== RUN   TestFindRepositories/No_matches
=== RUN   TestFindRepositories/Case_insensitive_partial_match
--- PASS: TestFindRepositories (0.00s)
    --- PASS: TestFindRepositories/Find_by_partial_name_-_multiple_matches (0.00s)
    --- PASS: TestFindRepositories/Find_by_partial_name_-_single_match (0.00s)
    --- PASS: TestFindRepositories/Find_by_exact_path (0.00s)
    --- PASS: TestFindRepositories/Find_by_partial_name_-_case_insensitive (0.00s)
    --- PASS: TestFindRepositories/Find_by_exact_name_match (0.00s)
    --- PASS: TestFindRepositories/No_matches (0.00s)
    --- PASS: TestFindRepositories/Case_insensitive_partial_match (0.00s)
PASS
ok  	github.com/daileyo/gws/cmd/git-workspace	0.497s
```

## Verification

- `TestFindRepositoriesByPath`: 6 cases covering prefix match, single prefix, substring fallback, case-sensitive rejection, mixed-case acceptance, and no-match — all PASS
- `TestFindRepositoriesWithFilters`: 5 cases covering repo-only, path-only, AND intersection, AND no-intersection, both-empty — all PASS
- `TestFindRepositories` (existing): 7 cases — all PASS, no regression
