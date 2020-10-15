import argparse
import subprocess

platforms = [
    ("darwin", "amd64"),
    ("darwin", "arm64"),
    ("freebsd", "386"),
    ("freebsd", "amd64"),
    ("freebsd", "arm"),
    ("linux", "386"),
    ("linux", "amd64"),
    ("linux", "arm"),
    ("linux", "arm64"),
    ("netbsd", "386"),
    ("netbsd", "amd64"),
    ("netbsd", "arm"),
    ("openbsd", "386"),
    ("openbsd", "amd64"),
    ("openbsd", "arm"),
    ("openbsd", "arm64"),
    ("windows", "386"),
    ("windows", "amd64"),
]


if __name__ == "__main__":
    parser = argparse.ArgumentParser(description="Auxilary build script.")
    parser.add_argument("-v", "--version", default="unspecified",
                        type=str, help="string used as a version when building binaries")
    args = parser.parse_args()
    version = args.version
    for (os, arch) in platforms:
        ext = ""
        if os == "windows":
            ext = ".exe"
        subprocess.run(["go", "build", "-ldflags", "-X main.version=%s" %
                        version, "-o", "bombardier-%s-%s%s" % (os, arch, ext)])
