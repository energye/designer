Name:           {{.PackageName}}
Version:        {{.App.Version}}
Release:        1%{?dist}
Summary:        {{.App.Desc}}
License:        {{.Linux.License}}
URL:            {{.Linux.Homepage}}
Group:          Applications/System

{{.Depends}}
Requires:       /bin/sh

%description
{{.App.Desc}}

%prep
# noop

%build
# noop

%install
rm -rf %{buildroot}
mkdir -p %{buildroot}/usr/bin
mkdir -p %{buildroot}/usr/share/applications
mkdir -p %{buildroot}/usr/share/icons
install -m 755 %{_appdir}/usr/bin/{{.PackageName}} %{buildroot}/usr/bin/{{.PackageName}}
install -m 644 %{_appdir}/usr/share/applications/{{.PackageName}}.desktop %{buildroot}/usr/share/applications/{{.PackageName}}.desktop
install -m 644 %{_appdir}/usr/share/icons/{{.PackageName}}.png %{buildroot}/usr/share/icons/{{.PackageName}}.png

%files
%attr(755, root, root) /usr/bin/{{.PackageName}}
/usr/share/applications/{{.PackageName}}.desktop
/usr/share/icons/{{.PackageName}}.png

%post
# noop

%preun
# noop
