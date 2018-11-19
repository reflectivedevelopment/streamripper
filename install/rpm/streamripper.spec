%define debug_package %{nil}

Summary: Streamripper uses multi-threading to increase the network throughput of a TCP stream with encryption.
Name: streamripper
Provides: streamripper
Version: 2.0
Release: 1
License: MIT
Source0: %{name}-%{version}.tar.gz
#URL: 
Vendor: Kim Ebert <kim@developmint.work>
Packager: Kim Ebert <kim@developmint.work>
BuildArch: x86_64
BuildRoot: %{_builddir}/%{name}-root
#Requires: 
BuildRequires: golang

%description
Streamripper uses multi-threading to increase the network throughput of a TCP stream with encryption.

%prep
%setup -n %{name}-%{version}

%build
./configure
make

%install
make DESTDIR=%{buildroot} install

%clean
rm -rf %{buildroot}

%files
/usr/local/bin/streamripper
