....This needs to be cleaned up

# run all tests
./scripts/run_tests.sh

# run all tests on a specific distro/os
./scripts/run_tests.sh -o alpine

# run a specific test on a specific distro/os
./scripts/run_tests.sh -o alpine -t TestTaskInstallsPkgDepsCorrectly

# run a debug session for a specific test on a specific distro/os
./scripts/run_tests.sh -D -o alpine -t TestTaskInstallsPkgDepsCorrectly

# build
Include the -B flag
