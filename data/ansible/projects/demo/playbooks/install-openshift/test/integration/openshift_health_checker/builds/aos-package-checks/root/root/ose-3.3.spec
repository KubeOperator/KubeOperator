Name:           atomic-openshift
Version:        3.3
Release:        1
Summary:        package the critical aos packages

License:        NA

Source0:	http://example.com/ose.tgz
BuildArch:	noarch

%package master
Summary:        package the critical aos packages
%package node
Summary:        package the critical aos packages
%package -n openvswitch
Summary:        package the critical aos packages
Version:	2.4
%package -n docker
Summary:        package the critical aos packages
Version:	1.10

%description
Package for pretending to provide AOS

%description master
Package for pretending to provide AOS

%description node
Package for pretending to provide AOS

%description -n openvswitch
Package for pretending to provide openvswitch

%description -n docker
Package for pretending to provide docker

%prep
%setup -q


%build


%install
rm -rf $RPM_BUILD_ROOT
mkdir -p $RPM_BUILD_ROOT


%files
%files master
%files node
%files -n openvswitch
%files -n docker

%doc

%changelog
