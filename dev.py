"""
MIT License

Copyright © 2022 whitespace_negative, https://github.com/whinee <whinyaan@protonmail.com>

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
"""

import re
import shlex
import sys
from functools import partial
from os import path
from subprocess import call

VERSIONS_NAME = {
    "user": "User",
    "dev": "Dev",
    "minor": "Minor",
    "patch": "Patch",
    "pri": "Pre-release identifier",
    "prv": "Pre-release version",
}
VERSIONS = ["user", "dev", "minor", "patch", "pri", "prv"]

def run(s: str):
    call(shlex.split(s))

def push(v: list[int]=None):
    msg = inquirer.text(message="Enter commit message", default="")
    run("git add .")
    if msg != "":
        msg = f"{msg},"
    if v:
        run(f"git commit -am '{msg}https://{GLOBAL('site')}/changelog#v{'-'.join([str(i) for i in v])}'")
    else:
        run(f"git commit -am '{msg or 'push'}'")
    run("git push")

def reset(idx: int, vls: list[int]=None):
    if vls is None:
        _vls = list(VLS)
    else:
        _vls = vls
    for i in range(idx + 1, len(VLS)):
        _vls[i] = 0
    return _vls

def _bump(v: str):
    idx = VERSIONS.index(v)
    _vls = list(VLS)
    match v:
        case "pri":
            match _vls[idx]:
                case 0:
                    _vls[3] += 1
                    _vls[4] = _vls[4] = 1
                case 3:
                    _vls[3] = _vls[4] = 0
                case _:
                    _vls = reset(idx, _vls)
                    _vls[idx] += 1
        case _:
            _vls = reset(idx, _vls)
            _vls[idx] += 1
    return _vls

def bump():
    while True:
        choices = []
        for k, v in VERSIONS_NAME.items():
            choices.append([f'{k.ljust(23)}(bump to {rv(_bump(k))})', k])
        v = inquirer.list_input(
            message=f"What version do you want to bump? (Current version: {rv(VLS)})",
            choices=choices,
        )
        _vls = _bump(v)
        print(
            f"    This will bump the version from {rv(VLS)} to {rv(_vls)} ({VERSIONS_NAME[v]} bump). "
        )
        match inquirer.list_input(
            message="Are you sure?",
            choices=[
                ["Yes", "y"],
                ["No", "n"],
                ["Cancel", "c"],
            ]
        ):
            case "y":
                with open("version", "wb") as f:
                    f.write(" ".join(_vls))
                gen_script()
                push(_vls)
                return
            case "n":
                continue
            case "c":
                break
            case _:
                pass

def vfn(answers, current):
    global vls
    x, y = re.match(r"\s*((\d+\s*){5})\s*", current).span()
    vls = current[x:y].strip().split(" ")
    if len(vls) == 5:
        vls = [int(i) for i in vls]
        return True

    raise Exception("Invalid version digits")

def set_ver():
    inquirer.text(
        message="Enter version digits seperated by spaces",
        validate = vfn
    )

    with open("version", "w") as f:
        f.write(" ".join([str(i) for i in vls]))

def main():
    match inquirer.list_input(
        message="What action do you want to take",
        choices=[
            ["Generate documentation", "docs"],
            ["Push to github", "gh"],
            ["Bump a version", "bump"],
            ["Set the version manually", "set_ver"]
        ]
    ):
        case "docs":
            from scripts import md_vars
            md_vars.main()
            run("mkdocs build")
            if inquirer.confirm(
                "Do you want to push this to github?",
                default=False
            ):
                push()
        case "gh":
            push()
        case "bump":
            bump()
        case "set_ver":
            set_ver()
        case _:
            pass

if __name__ == "__main__":
    if len(sys.argv) >= 2:
        if sys.argv[1] == "req":
            run(f'pip install --require-virtualenv pyyaml')
            import yaml
            with open('dev.yml', 'r') as f:
                y = yaml.safe_load(f)
            ls = []
            for c in y["env"]["dev"]["req"]:
                ls.append(f'-r {y["requirements"][c]}')
            run(f'pip install --require-virtualenv {" ".join(ls)}')
        elif sys.argv[1] == "docs":
            from scripts import md_vars

            md_vars.main(True)
    else:
        import httpx
        import inquirer
        import yaml

        from scripts.req import req

        try:
            VLS = req.get("https://raw.githubusercontent.com/manga-g/manga-g/main/version").text
            with open("version", "w") as f:
                f.write(VLS)
        except httpx.ConnectTimeout:
            print("Could not connect to github, Defaulting to local version")
            with open("version", "r") as f:
                VLS = [int(i) for i in f.read().split(" ")]

        from scripts.settings import stg
        from scripts.utils import ddir, rv

        YML = stg(None, "dev.yml")
        GLOBAL = partial(ddir, ddir(YML, "md_vars/global"))
        LICENSE = partial(ddir, ddir(YML, "license"))

        FN = GLOBAL("formal_name")
        ORG = GLOBAL("organization")
        USER = GLOBAL("user")

        main()