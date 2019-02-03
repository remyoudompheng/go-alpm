package alpm

import (
	"bytes"
	"reflect"
	"testing"
)

const pacmanConf = `
#
# GENERAL OPTIONS
#
[options]
RootDir     = /
DBPath      = /var/lib/pacman/
CacheDir    = /var/cache/pacman/pkg/ /other/cachedir
LogFile     = /var/log/pacman.log
GPGDir      = /etc/pacman.d/gnupg/
HoldPkg     = pacman glibc
# If upgrades are available for these packages they will be asked for first
SyncFirst   = pacman
#XferCommand = /usr/bin/curl -C - -f %u > %o
XferCommand = /usr/bin/wget --passive-ftp -c -O %o %u
CleanMethod = KeepInstalled
Architecture = x86_64

# Pacman won't upgrade packages listed in IgnorePkg and members of IgnoreGroup
IgnorePkg   = hello world
IgnoreGroup = kde

NoUpgrade   = kernel26
NoExtract   =

# Misc options
UseSyslog
#UseDelta
TotalDownload
CheckSpace
#VerbosePkgLists
ILoveCandy
# By default, pacman accepts packages signed by keys that its local keyring
# trusts (see pacman-key and its man page), as well as unsigned packages.
SigLevel    = Required DatabaseOptional
LocalFileSigLevel = Optional
RemoteFileSigLevel = Required

[core]
SigLevel = Required
Server = ftp://ftp.example.com/foobar/$repo/os/$arch/

[custom]
SigLevel = Optional TrustAll
Server = file:///home/custompkgs
`

var pacmanConfRef = PacmanConfig{
	CacheDir:    []string{"/var/cache/pacman/pkg/", "/other/cachedir"},
	HoldPkg:     []string{"pacman", "glibc"},
	SyncFirst:   []string{"pacman"},
	IgnorePkg:   []string{"hello", "world"},
	IgnoreGroup: []string{"kde"},
	NoUpgrade:   []string{"kernel26"},
	NoExtract:   nil,

	RootDir:      "/",
	DBPath:       "/var/lib/pacman/",
	GPGDir:       "/etc/pacman.d/gnupg/",
	LogFile:      "/var/log/pacman.log",
	Architecture: "x86_64",
	XferCommand:  "/usr/bin/wget --passive-ftp -c -O %o %u",
	CleanMethod:  "KeepInstalled",

	Options: ConfUseSyslog | ConfTotalDownload | ConfCheckSpace | ConfILoveCandy,

	SigLevel: SigPackage | SigDatabase | SigDatabaseOptional,
	LocalFileSigLevel: SigPackage | SigPackageOptional |
		SigDatabase | SigDatabaseOptional,
	RemoteFileSigLevel: SigPackage | SigDatabase,

	Repos: []RepoConfig{
		{Name: "core", Servers: []string{"ftp://ftp.example.com/foobar/$repo/os/$arch/"},
			SigLevel: SigPackage | SigDatabase},
		{Name: "custom", Servers: []string{"file:///home/custompkgs"},
			SigLevel: SigPackage | SigPackageOptional |
				SigPackageMarginalOk | SigPackageUnknownOk |
				SigDatabase | SigDatabaseOptional |
				SigDatabaseMarginalOk | SigDatabaseUnknownOk},
	},
}

func detailedDeepEqual(t *testing.T, x, y interface{}) {
	v := reflect.ValueOf(x)
	w := reflect.ValueOf(y)
	if v.Type() != w.Type() {
		t.Errorf("differing types %T vs. %T", x, y)
		return
	}
	for i := 0; i < v.NumField(); i++ {
		v_fld := v.Field(i).Interface()
		w_fld := w.Field(i).Interface()
		if v.Type().Field(i).Name == "Repos" {
			repos1 := v_fld.([]RepoConfig)
			repos2 := w_fld.([]RepoConfig)
			if len(repos1) != len(repos2) {
				t.Errorf("repos length mismatch: %+v vs %+v", repos1, repos2)
			}
			for i := 0; i < len(repos1) && i < len(repos2); i++ {
				detailedDeepEqual(t, repos1[i], repos2[i])
			}
		}
		if !reflect.DeepEqual(v_fld, w_fld) {
			t.Errorf("field %s differs: got %#v, expected %#v",
				v.Type().Field(i).Name, v_fld, w_fld)
		}
	}
}

func TestPacmanConfigParser(t *testing.T) {
	buf := bytes.NewBufferString(pacmanConf)
	conf, err := ParseConfig(buf)
	if err != nil {
		t.Error(err)
	}

	detailedDeepEqual(t, conf, pacmanConfRef)
}
