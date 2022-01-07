### TODO
* If someone has apt installed on a non-debian-like machine, we don't want to detect apt exists and try to use that. How?
* Handle cyclic dependencies
* Way to set a prioritized list of install targets, for example `--installers=gvm,brew,npm`
* Ability to use templates when creating links
* To support the above, some way to drive the configuration of the links
* Cleanup readme and comments
* Ability to add sections of code to pre-existing files (like sourcing aliases in .bashrc etc)
* Support gvm, npm, etc.
* Use go-releaser to add pre-built binaries https://goreleaser.com/
* Full in-environment tests using docker for every supported environment
* Loads a default.toml from /usr/share/autostart/default.toml that has all the default configuration. Can be overridden by ~/home/<user>/.config/autostart/default.toml also existing.

### Names
Maybe "shoelace" aka "bootstrap" for users