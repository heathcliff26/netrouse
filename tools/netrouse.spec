%global debug_package %{nil}

Name:           netrouse
Version:        0
Release:        %autorelease
Summary:        Wake on Lan Utility

%global package_id io.github.heathcliff26.%{name}

License:        Apache-2.0
URL:            https://github.com/heathcliff26/netrouse
Source:         %{url}/archive/refs/tags/v%{version}.tar.gz

BuildRequires: golang >= 1.27
BuildRequires: gcc libXcursor-devel libXrandr-devel mesa-libGL-devel libXi-devel libXinerama-devel libXxf86vm-devel libxkbcommon-devel wayland-devel
BuildRequires: gettext-envsubst

%global _description %{expand:
This is a simple utility for sending Wake-On-Lan magic packet to clients in the local network.
It can be used directly via the cli, or remotely via a web interface.}

%description %{_description}

%prep
%autosetup -n %{name}-%{version} -p1

%build
export RELEASE_VERSION="%{version}-%{release}"
make build-gui

%install
install -D -m 755 bin/%{name}-gui %{buildroot}/%{_bindir}/%{name}-gui
install -D -m 644 packages/%{package_id}.desktop %{buildroot}/%{_datadir}/applications/%{package_id}.desktop
install -D -m 644 packages/%{package_id}.png %{buildroot}/%{_datadir}/icons/hicolor/256x256/apps/%{package_id}.png
install -D -m 644 packages/%{package_id}.svg %{buildroot}/%{_datadir}/icons/hicolor/scalable/apps/%{package_id}.svg
install -D -m 644 %{package_id}.metainfo.xml %{buildroot}/%{_datadir}/metainfo/%{package_id}.metainfo.xml

%files
%license LICENSE
%doc README.md
%{_bindir}/%{name}-gui
%{_datadir}/applications/%{package_id}.desktop
%{_datadir}/icons/hicolor/256x256/apps/%{package_id}.png
%{_datadir}/icons/hicolor/scalable/apps/%{package_id}.svg
%{_datadir}/metainfo/%{package_id}.metainfo.xml

%changelog
%autochangelog
