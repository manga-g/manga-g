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

import itertools
import re

from .md_vars import RULES_MDV
from .settings import stg
from .utils import ddir, inmd

RE_MX = r"(?<=\{matrix.)[a-zA-Z0-9-_]+?(?=\})"

SCRIPTS = stg(None, "scripts.yml")
MD_VARS = RULES_MDV["md_vars"]

MATRIX = {}
for k, v in SCRIPTS["matrix"].items():
    MATRIX[k] = [str(i) for i in v]

PG = {}
GLOBAL = {}
for k, v in dict(
    MD_VARS["global"],
    **SCRIPTS["variables"]["global"],
    **PG
).items():
    GLOBAL[k] = str(v)

PL = {}
LOCAL = {}
for k, v in dict(
    MD_VARS["local"],
    **SCRIPTS["variables"]["local"],
    **PL
).items():
    LOCAL[k] = str(v)

def cd(*k: list[str]):
    op = []
    dls = [MATRIX[i] for i in k]
    for i in itertools.product(*dls):
        op.append(dict(zip(k, i)))
    return op

def repl(key: str, fn: str, contents: str):
    s, c = fn, contents
    lv = ddir(SCRIPTS, f"variables/local/{key}")
    rmvg = {**GLOBAL, **lv}

    for k, v in rmvg.items():
        c = c.replace(f"${{{k}}}", v)
        s = s.replace(f"${{{k}}}", v)

    return s, c

def mr(k: str, v: dict[str, str]):
    s, c = repl(k, v["path"], v["contents"])
    if mx_match:=re.findall(RE_MX, s):
        op = []
        for i in cd(*mx_match):
            _s = s
            _c = c
            for k, v in i.items():
                _s = _s.replace(f"${{matrix.{k}}}", v)
                _c = _c.replace(f"${{matrix.{k}}}", v)
            op.append([_s, _c])
        return op
    else:
        return [[s, c]]

def main():
    for k, v in SCRIPTS["scripts"].items():
        for p, c in mr(k, v):
            with open(inmd(p), "w") as f:
                f.write(c)