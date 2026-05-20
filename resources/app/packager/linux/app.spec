Name:           {{.PackageName}}
Version:        {{.App.Version}}
Release:        1%{?dist}
Summary:        {{if .App.Desc}}{{.App.Desc}}{{else}}No description{{end}}
License:        {{if .Linux.License}}{{.Linux.License}}{{else}}Unknown{{end}}
{{- if .Linux.Homepage}}
URL:            {{.Linux.Homepage}}
{{- end}}
Group:          Applications/System

{{- if .Depends}}
{{.Depends}}
{{- end}}

%description
{{if .App.Desc}}{{.App.Desc}}{{else}}No description{{end}}

%prep
# noop

%build
# noop

%install
rm -rf %{buildroot}
mkdir -p %{buildroot}/usr/bin
mkdir -p %{buildroot}/usr/share/applications
mkdir -p %{buildroot}/usr/share/icons
mkdir -p %{buildroot}/usr/lib
install -m 755 %{_appdir}/usr/bin/{{.PackageName}} %{buildroot}/usr/bin/{{.PackageName}}
install -m 644 %{_appdir}/usr/share/applications/{{.PackageName}}.desktop %{buildroot}/usr/share/applications/{{.PackageName}}.desktop
install -m 644 %{_appdir}/usr/share/icons/{{.PackageName}}.png %{buildroot}/usr/share/icons/{{.PackageName}}.png
install -m 755 %{_appdir}/usr/lib/{{.LibName}} %{buildroot}/usr/lib/{{.LibName}}

%files
%attr(755, root, root) /usr/bin/{{.PackageName}}
/usr/share/applications/{{.PackageName}}.desktop
/usr/share/icons/{{.PackageName}}.png
/usr/lib/{{.LibName}}

%post
/sbin/ldconfig

%preun
/sbin/ldconfig
