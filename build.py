import argparse
import os
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
    ("windows", "arm64"),
]


if __name__ == "__main__":
    parser = argparse.ArgumentParser(description="Auxilary build script.")
    parser.add_argument("-v", "--version", default="unspecified",
                        type=str, help="string used as a version when building binaries")
    args = parser.parse_args()
    version = args.version
    for (build_os, build_arch) in platforms:
        ext = ""
        if build_os == "windows":
            ext = ".exe"
        build_env = os.environ.copy()
        build_env["GOOS"] = build_os
        build_env["GOARCH"] = build_arch
        subprocess.run(["go", "build", "-ldflags", "-s -w -X main.version=%s" %
                        version, "-o", "bombardier-%s-%s%s" % (build_os, build_arch, ext)], env=build_env)
