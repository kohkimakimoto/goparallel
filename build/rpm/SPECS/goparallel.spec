Name:           goparallel
Version:        0.1.0
Release:        1.el%{rhel}
Summary:        Execute commands in parallel.

Group:          Development/Tools
License:        MIT
Source0:        %{name}_linux_amd64.zip
BuildRoot:      %(mktemp -ud %{_tmppath}/%{name}-%{version}-%{release}-XXXXXX)

%description
Execute commands in parallel.

%prep
%setup -q -c

%install
mkdir -p %{buildroot}/%{_bindir}
cp goparallel %{buildroot}/%{_bindir}

%pre

%post

%preun

%clean
rm -rf %{buildroot}

%files
%defattr(-,root,root,-)
%attr(755, root, root) %{_bindir}/goparallel

%doc
