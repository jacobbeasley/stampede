# Automatically clean migrations/schema.sql before running buffalo test
if ($args[0] -eq "test") {
    go run "$PSScriptRoot/scripts/clean_schema.go"
}

# Find the real buffalo command in PATH excluding the current directory
$realBuffalo = Get-Command buffalo -All | Where-Object { $_.Path -notlike "$PSScriptRoot*" } | Select-Object -First 1

if ($realBuffalo) {
    & $realBuffalo.Path @args
} else {
    & buffalo @args
}
