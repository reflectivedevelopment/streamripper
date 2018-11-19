
bin_PROGRAMS = streamripper
streamripper_SOURCES = streamripper.go
streamripper: streamripper.go
	go build -o $@ $<


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

