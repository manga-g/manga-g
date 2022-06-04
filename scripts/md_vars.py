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

import base64
from datetime import date
from functools import partial

from . import docs
from .settings import stg
from .utils import ddir, rv

YML = stg(None, "dev.yml")
R_GLOBAL = ddir(YML, "md_vars/global")
GLOBAL = partial(ddir, R_GLOBAL)
LICENSE = partial(ddir, ddir(YML, "license"))

FN = GLOBAL("formal_name")
ORG = GLOBAL("organization")
USER = GLOBAL("user")

icons = ["issues", "forks", "stars", "contributors", "license", "code", "win", "macos", "linux", "twitter"]
langs = ["python", "html", "yaml"]

def b64(name: str):
    with open(f"./docs/assets/images/icons/{name}", "rb") as f:
        return base64.b64encode(f.read()).decode("utf-8")

if lc := LICENSE("cholder"):
    copyright = []
    for c, mp in lc.items():
        user = mp['user']
        for org, projects in mp["projects"].items():
            for project, pm in projects.items():
                copyright.append(
                    f"by [{c}, Github account [{user}](https://github.com/{user}) owner, {pm['year']}] as part of project [{project}](https://github.com/{org}/{project})"
                )
    if len(copyright) > 1:
        copyright[-2] += f", and {copyright[-1]}"
        del copyright[-1]
    cholder = f"""Copyright for portions of project [{FN}](https://github.com/{ORG}/{FN}) are held {', '.join(copyright)}.\n
All other copyright for project [{FN}](https://github.com/{ORG}/{FN}) are held by [Github Account [{USER}](https://github.com/{USER}) Owner, {LICENSE('year')}]."""
else:
    cholder = f"Copyright (c) 2021 Github Account [{USER}](https://github.com/{USER}) Owner"

RULES_MDV = {
    "rules": {
        "del": {},
        "repl": {},
    },
    "md_vars": {
        "global": {
            "year": str(date.today().year),
            "cholder": cholder,
            "ver": rv(docs.VLS),
            "prerel": not docs.VLS[-2] == 3,
            **{f"{i}_b64": b64(f"{i}.png") for i in icons},
            **dict(
                zip([f'v_{i}' for i in ["d", "u", "m", "p", "pri", "prv"]], docs.VLS)
            ),
            **R_GLOBAL
        },
        "local": {
            **ddir(YML, "md_vars/local"),
        },
    }
}

def main(hr: bool=False):
    docs.main(RULES_MDV, hr)