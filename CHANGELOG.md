# Changelog

Version v0.9.3-alpha.2
--------------------
Released: 2026-10-04

🚀 Features:
- *(packages)* Build rpms with goreleaser as well

🐛 Bug Fixes:
- *(gui)* Fix missing version information
- *(gui)* Add make target for fyne metadata
- Inject changelog during packit build

⚙️ Miscellaneous Tasks:
- Add new changelog workflow
- *(release)* Update changelog for v0.9.3-alpha.1
- *(ci)* Trigger release workflow after changelog update

Version v0.9.2
--------------------
Released: 2026-10-02

🚀 Features:
- *(container)* [**breaking**] Update image ghcr.io/heathcliff26/go-fyne-ci to v202609241228
- *(container)* [**breaking**] Update image ghcr.io/heathcliff26/go-fyne-ci to v202609301715

🐛 Bug Fixes:
- *(deps)* Update module github.com/valkey-io/valkey-go to v1.0.78

💼 Other:
- Derive build-id from semver version
- Improve mac entry experience

Version v0.9.1
--------------------
Released: 2026-09-24

💼 Other:
- Enable appstream metadata validation
- Provide additional apk per arch
- Add android installation options

Version v0.9.0
--------------------
Released: 2026-09-23

📚 Documentation:
- Add privacy policy

💼 Other:
- Add quay.io repository
- Rename Coverprofiles workflow
- Move cli binary into separate folder
- Add local tab
- Use containerized workflows
- Fix coverprofiles workflow
- Fix wrong fileending for hosts file
- *(deps)* Bump golang.org/x/image
- Fix coverprofiles permissions
- Add basic unit-tests
- Add periodic status updates
- Re-enable push to coveralls
- Implement api client
- Implement settings tab
- Build for android
- Fix wrong release artifact
- Do not create draft release by default
- Persist window state across restarts
- Do not upload gosec on release
- Build new bundle format
- Drop unused arg to keytool
- Stretch wake button over whole space
- Improve package icon quality
- Use adaptive Icon
- Improve error handling
- Move background tasks into lifecycle
- Split wakeTab update into separate functions
- Use static name for builder container
- Improve context handling for background updates
- Move status update inside fetch Hosts
- Hide window reset
- Add about button
- Update copr repo
- Update hosts in-place
- Add refresh button
- Fix js client wrong capitalization
- Add installation instructions for gui
- Use context for http calls
- Default to https schema if none specified
- Add option to import hosts from remote
- Make tabs scrollable

Version v0.8.5
--------------------
Released: 2026-09-13

🐛 Bug Fixes:
- *(deps)* Update module go.yaml.in/yaml/v3 to v3.0.5
- *(container)* Update image docker.io/library/golang to v1.26.6
- *(container)* Update image docker.io/library/golang to v1.27.1
- *(deps)* Update gomod

Version v0.8.4
--------------------
Released: 2026-07-12

🐛 Bug Fixes:
- *(container)* Update image docker.io/library/golang to v1.26.5
- *(deps)* Update gomod

Version v0.8.3
--------------------
Released: 2026-06-13

🐛 Bug Fixes:
- *(container)* Update image docker.io/library/golang to v1.26.4
- *(deps)* Update gomod

💼 Other:
- Move renovate and Link Check to tekton
- Remove obsolete Renovate status badge

Version v0.8.2
--------------------
Released: 2026-06-01

🐛 Bug Fixes:
- *(deps)* Update gomod
- *(container)* Update image docker.io/library/golang to v1.26.3

💼 Other:
- Use goreleaser to build release artifacts

Version v0.8.1
--------------------
Released: 2026-05-01

🐛 Bug Fixes:
- *(container)* Update image docker.io/library/golang to v1.26.2
- *(deps)* Update gomod

Version v0.8.0
--------------------
Released: 2026-03-21

🐛 Bug Fixes:
- *(container)* Update image docker.io/library/golang to v1.26.1
- *(deps)* Update gomod

💼 Other:
- Use user and group nobody

Version v0.7.7
--------------------
Released: 2026-03-01

🚀 Features:
- *(container)* Update image docker.io/library/golang to v1.26.0

🐛 Bug Fixes:
- *(deps)* Update gomod

💼 Other:
- Fix test no longer valid with golang 1.26

Version v0.7.6
--------------------
Released: 2026-02-01

🐛 Bug Fixes:
- *(container)* Update image docker.io/library/golang to v1.25.6
- *(deps)* Update gomod

💼 Other:
- Remove toolchain directive
- Use go.yaml.in/yaml/v3 for YAML parsing
- Remove unneeded prepare step in renovate workflow

Version v0.7.5
--------------------
Released: 2026-01-01

🐛 Bug Fixes:
- *(deps)* Update gomod
- *(container)* Update image docker.io/library/golang to v1.25.5
- *(deps)* Update gomod

Version v0.7.4
--------------------
Released: 2025-12-06

💼 Other:
- Decrease ping timeout to 200ms
- Do not display alert when successfully waking host

Version v0.7.3
--------------------
Released: 2025-12-01

🐛 Bug Fixes:
- *(container)* Update image docker.io/library/golang to v1.25.4
- *(deps)* Update gomod

💼 Other:
- Enable signing artifacts in release workflow

Version v0.7.2
--------------------
Released: 2025-11-01

🐛 Bug Fixes:
- *(deps)* Update gomod
- *(container)* Update image docker.io/library/golang to v1.25.3

Version v0.7.1
--------------------
Released: 2025-10-01

🐛 Bug Fixes:
- *(deps)* Update module github.com/stretchr/testify to v1.11.1
- *(deps)* Update module github.com/spf13/cobra to v1.10.0
- *(deps)* Update module github.com/spf13/cobra to v1.10.1
- *(container)* Update image docker.io/library/golang to v1.25.1
- *(deps)* Update gomod

Version v0.7.0
--------------------
Released: 2025-08-21

🚀 Features:
- *(container)* Update image docker.io/library/golang to v1.25.0

🐛 Bug Fixes:
- *(deps)* Update gomod

💼 Other:
- Add address to host
- Change AddHost to accept only json body
- Add address for hosts
- Fix missing permission in coverprofile

Version v0.6.2
--------------------
Released: 2025-08-01

🐛 Bug Fixes:
- *(deps)* Update module sigs.k8s.io/yaml to v1.6.0

💼 Other:
- Optimize GetHosts

Version v0.6.1
--------------------
Released: 2025-07-22

💼 Other:
- Add timeout to valkey operations
- Add get hosts endpoint
- Always render index.html instead of server side caching

Version v0.6.0
--------------------
Released: 2025-07-21

🐛 Bug Fixes:
- *(deps)* Update gomod
- *(container)* Update image docker.io/library/golang to v1.24.5

💼 Other:
- Add option to seed hosts on startup
- Fix missing newlines in credits section
- Return 403 forbidden when changing hosts in readonly

Version v0.5.3
--------------------
Released: 2025-07-01

🐛 Bug Fixes:
- *(container)* Update image docker.io/library/golang to v1.24.4
- *(deps)* Update gomod

Version v0.5.2
--------------------
Released: 2025-06-01

🐛 Bug Fixes:
- *(container)* Update image docker.io/library/golang to v1.24.3
- *(deps)* Update gomod

💼 Other:
- Tag release to ensure correct version in binary
- Update release job to target renamed workflow

Version v0.5.1
--------------------
Released: 2025-05-01

🐛 Bug Fixes:
- *(deps)* Update gomod

💼 Other:
- Add gosec to improve code quality
- Change the help target to dynamically generate it's output
- Move target description to above them
- Update the favicon to support different sizes
- Fix release workflow permissions

Version v0.5.0
--------------------
Released: 2025-04-24

💼 Other:
- Add unit-test for cli errors
- Add unit-tests for storage failures
- Expand testsuite to cover more scenarios
- Replace assert with require where it improves readability

Version v0.4.1
--------------------
Released: 2025-04-21

💼 Other:
- Add run command
- Update favicon to include background
- Add help target to list available commands
- Enhance accessibility by adding aria-labels
- Fix footer text too far apart

Version v0.4.0
--------------------
Released: 2025-04-20

💼 Other:
- Change default options for release job
- Combine app name and version in footer
- Enhance footer layout with span instead of free text
- Small layout fix
- Enable using valkey as storage backend
- Expand unit-tests to cover storage

Version v0.3.1
--------------------
Released: 2025-04-20

💼 Other:
- Fix missing data volume

Version v0.3.0
--------------------
Released: 2025-04-19

🐛 Bug Fixes:
- *(container)* Update image docker.io/library/golang to v1.24.2

💼 Other:
- Add favicon for website
- Combine embedded dirs under asset FS
- Add inactive elements for adding/removing hosts
- Move hosts storage into separate package
- Implement unit-tests for storage and file backend
- Move api under v1
- Implement hosts endpoint
- Move add hosts button to bottom of card
- Add footer with app name and version

Version v0.2.1
--------------------
Released: 2025-04-01

🐛 Bug Fixes:
- *(deps)* Update module github.com/heathcliff26/simple-fileserver to v1.2.3

💼 Other:
- Fix position of title to close to top

Version v0.2.0
--------------------
Released: 2025-03-15

🐛 Bug Fixes:
- *(deps)* Update module github.com/heathcliff26/simple-fileserver to v1.2.2
- *(container)* Update image docker.io/library/golang to v1.24.1

💼 Other:
- Improve style of displayed mac addresses
- Switch to bootstrap CSS
- Add option to enter custom mac address
- Improve alerts by adding MAC and name to output
- Implement theme change based on user preferences
- Use bootstrap alerts instead of standard javascript alerts.

Version v0.1.2
--------------------
Released: 2025-03-01

🐛 Bug Fixes:
- *(deps)* Update module github.com/spf13/cobra to v1.9.1

💼 Other:
- Add unit-tests
- Fix broken docker hub link
- Use the version from build info
- Replace abandoned yaml module with sig-kubernetes version

Version v0.1.1
--------------------
Released: 2025-02-01

🐛 Bug Fixes:
- *(container)* Update image docker.io/library/golang to v1.23.5
- *(deps)* Update module github.com/heathcliff26/simple-fileserver to v1.2.1

💼 Other:
- Build new Version on any change to static files
- Add caching headers for static files

Version v0.1.0
--------------------
Released: 2025-01-04

🐛 Bug Fixes:
- *(container)* Update image docker.io/library/golang to v1.23.4

💼 Other:
- Fix editorconfig failures
- Add example configuration file

