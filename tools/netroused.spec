%global debug_package %{nil}
%global bin_name netrouse

Name:           netroused
Version:        0
Release:        %autorelease
Summary:        Wake on Lan Utility
%global package_id io.github.heathcliff26.%{name}

License:        Apache-2.0
URL:            https://github.com/heathcliff26/netrouse
Source:         %{url}/archive/refs/tags/v%{version}.tar.gz

BuildRequires: golang >= 1.27.0

%global _description %{expand:
This is a simple utility for sending Wake-On-Lan magic packet to clients in the local network.
It can be used directly via the cli, or remotely via a web interface.}

%description %{_description}

%prep
%autosetup -n %{name}-%{version} -p1

%build
export RELEASE_VERSION="%{version}-%{release}"
hack/build.sh

%install
install -D -m 0755 bin/%{bin_name} %{buildroot}%{_bindir}/%{bin_name}
install -D -m 0644 tools/%{name}.service %{buildroot}%{_prefix}/lib/systemd/system/%{name}.service
install -D -m 0644 examples/config.yaml %{buildroot}%{_sysconfdir}/netrouse/config.yaml
install -D -m 0644 %{package_id}.metainfo.xml %{buildroot}/%{_datadir}/metainfo/%{package_id}.metainfo.xml

%post
PLACEHOLDER_POSTINSTALL

%preun
PLACEHOLDER_PREREMOVE

%postun
PLACEHOLDER_POSTREMOVE

%files
%license LICENSE
%doc README.md
%{_bindir}/%{bin_name}
%{_prefix}/lib/systemd/system/%{name}.service
%dir %{_sysconfdir}/netrouse
%{_sysconfdir}/netrouse/config.yaml
%{_datadir}/metainfo/%{package_id}.metainfo.xml

%changelog
%autochangelog
