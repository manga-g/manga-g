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
FITNESS FOR A PARTICULAR PURPOSE AND NON-INFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
"""

import shlex
from os import makedirs, path
from subprocess import call
from typing import Any

from mako.lookup import TemplateLookup

from .settings import stg

YAHML = stg(None, "dev.yml")
PR = ["alpha", "beta", "rc"]

def ddir(d: dict[Any, Any], dir: str, de: Any={}) -> Any:
    """
    Retrieve dictionary value using recursive indexing with a string.
    ex.:
        `ddir({"data": {"attr": {"ch": 1}}}, "data/attr/ch")`
        will return `1`


    Args:
        dict (dict): Dictionary to retrieve the value from.
        dir (str): Directory of the value to be retrieved.

    Returns:
        op (Any): Retrieved value.
    """
    op = d
    for a in dir.split("/"):
        op = op.get(a)
        if not op:
            break
    return op or de

LOOKUPS = TemplateLookup(directories=ddir(YAHML, "file/templates") or [])

def srv_tpl(tn: str, lookup: TemplateLookup=LOOKUPS, **kwargs: dict[str, Any]):
    return lookup.get_template(tn).render(**kwargs)

def run(s: str):
    call(shlex.split(s))

def repl(s: str, repl_dict: dict[str, list[str]]) -> str:
    op = s
    for k, v in repl_dict.items():
        for i in v:
            op = op.replace(i, k)
    return op

def inmd(p: str, ls: list[str]=None):
    """
    "If Not `path.isdir`, Make Directories"

    Args:
        p (str): [description]
    """

    pd = path.dirname(p)
    if (pd) and (not path.isdir(pd)):
        makedirs(pd)
        if ls:
            ls.append(pd)
    return p

def rv(vls: list[str]):
    pr = ""
    if vls[3] <= 2:
        pr = f'-{PR[vls[3]]}.{vls[4]}'
    return ".".join([str(i) for i in vls[0:3]]) + pr