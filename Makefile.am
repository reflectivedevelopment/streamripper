
bin_PROGRAMS = streamripper
streamripper_SOURCES = ./src/streamripper.go ./src/ripper/pool.go 

streamripper$(EXEECT): ./src/streamripper.go ./src/ripper/pool.go 
	export GOPATH="$(abs_srcdir)"; go build -o $@ $<

pool_test_SOURCES = ./test/ripper/pool_test.go ./src/ripper/pool.go

pool_test$(EXEECT): ./test/ripper/pool_test.go ./src/ripper/pool.go
	export GOPATH="$(abs_srcdir)"; go test -o $@ $<

TESTS = pool_test

EXTRA_DIST = $(srcdir)/test

DEBTMP_DIR = deb_build
# Create a deb package
pkg-deb: dist
	echo tmpdir $(DEBTMP_DIR)
	-rm -r $(DEBTMP_DIR)
	mkdir -p $(DEBTMP_DIR)
	cp $(PACKAGE)-$(VERSION).tar.gz $(DEBTMP_DIR)/
	tar --directory $(DEBTMP_DIR)/ \
		-xzf $(DEBTMP_DIR)/$(PACKAGE)-$(VERSION).tar.gz
	cd $(DEBTMP_DIR)/$(PACKAGE)-$(VERSION)/; \
		echo | dh_make -y --single --copyright mit -e kim@developmint.work -f \
			../$(PACKAGE)-$(VERSION).tar.gz
	cp $(srcdir)/install/deb/changelog $(srcdir)/install/deb/control \
		$(srcdir)/install/deb/copyright $(DEBTMP_DIR)/$(PACKAGE)-$(VERSION)/debian/
	cd $(DEBTMP_DIR)/$(PACKAGE)-$(VERSION)/; \
		rm debian/README.Debian debian/*.ex debian/*.EX; \
		./configure; \
		dpkg-buildpackage -uc -us -rfakeroot;
	mv $(DEBTMP_DIR)/*.deb ./
	rm -r $(DEBTMP_DIR);

# Create an RPM package
pkg-rpm: dist
	mkdir -p $(HOME)/.rpm/RPMS/i386
	mkdir -p $(HOME)/.rpm/SRPMS
	mkdir -p $(HOME)/.rpm/BUILD
	mkdir -p $(HOME)/.rpm/SOURCES
	mkdir -p $(HOME)/.rpm/tmp
	- rm -r $(HOME)/.rpm/BUILD/$(PACKAGE)-root
	- rm -r $(HOME)/.rpm/RPMS/$(PACKAGE)-*
	- rm -r $(HOME)/.rpm/SRPMS/$(PACKAGE)-*
	- rm -r $(HOME)/.rpm/SOURCES/$(PACKAGE)-*
	cp $(PACKAGE)-$(VERSION).tar.gz $(HOME)/.rpm/SOURCES/
	rpmbuild --quiet --nodeps --define="_topdir $(HOME)/.rpm" \
	    -ba $(srcdir)/install/rpm/$(PACKAGE).spec
	mv $(HOME)/.rpm/RPMS/x86_64/$(PACKAGE)-$(VERSION)-*.x86_64.rpm ./
