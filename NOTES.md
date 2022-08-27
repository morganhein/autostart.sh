### TODO
* If someone has apt installed on a non-debian-like machine, we don't want to detect apt exists and try to use that. How?
* Handle cyclic dependencies
* Way to set a prioritized list of install targets, for example `--installers=gvm,brew,npm`
  * This is done now in the general.installer_preferences setting
* Ability to use templates when creating links
* To support the above, some way to drive the configuration of the links
* Cleanup readme and comments
* Ability to add sections of code to pre-existing files (like sourcing aliases in .bashrc etc)
* Support gvm, npm, etc.
* Use go-releaser to add pre-built binaries https://goreleaser.com/
* Full in-environment tests using docker for every supported environment
* Loads a default.toml from /usr/share/envy/default.toml that has all the default configuration. Can be overridden by ~/home/<user>/.config/envy/default.toml also existing.

### Scenarios

Scenario: A user wants to install a package on a machine, and does not care what installer is used.
Context: The package the user wants to install is defined in their recipe, with supported installers on the local machine
Result: They can type `envy install <pkg>` and it will install using the preferred installer

Scenario: A user wants to install a single task as defined in a recipe
Context: The task, installer, and packages are defined in the recipe. It has no dependencies, no pre/post-install steps, or any sync actions.
Result: By typing `envy task <task>` the task will be run, and the package(s) will be installed

Scenario: A user wants to install many tasks
Context: The tasks, installer, and packages are defined in the recipe. One "entrypoint" task must have the other tasks as dependencies.
Result: By typing `envy task <entrypoint task>` the main task will be run, and therefore the children tasks as well

