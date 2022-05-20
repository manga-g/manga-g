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

import os
import re
import shutil
from os import listdir, path
from os.path import dirname as dn
from pathlib import Path
from typing import Any, Dict, List

import frontmatter
import yaml
from mako.template import Template

from .settings import stg, wr_stg
from .utils import ddir, inmd, repl, stg

with open("version", "r") as f:
    VLS = [int(i) for i in f.read().split(" ")]
v_ud = '/'.join([str(i) for i in VLS[0:2]])

RE_MD_TPL = r"(?<=# {key} start\n).+(?=\n\s+# {key} end)"

YML = stg(None, "dev.yml")

DOCS = ddir(YML, "docs")
MAKO = ddir(YML, "mako")
RULES = ddir(YML, "rules")
MD_VARS_YML = ddir(YML, "md_vars")

RMVC = ddir(MD_VARS_YML, "global")
IDF = Path(f'./{ddir(YML, "docs/input")}')

GEN = {
    "docs": [[], []],
    "mako": [[], []],
}

class Constants:
    pass

def dd(od: Dict[str, List[str]], *dicts: List[Dict[str, List[str]]]) -> Dict[str, List[str]]:
    for d in dicts:
        for a, v in d.items():
            od[a] = [*(od.get(a, []) or []), *v]
    return od

def rules_fn(rules: Dict[Any, Any]) -> Dict[str, List[str]]:
    return dd({"": ddir(rules, "del", [])}, ddir(rules, "repl"))

def docs_dir(mn: str, absolute: bool=True, api: bool=False) -> str:
    mls = mn.split(".")
    rel_ls = [
        *[str(i) for i in VLS[0:2]],
        *mls[:-1],
        f"{mls[-1]}.md"
    ]
    if api:
        rel_ls.insert(2, "api")
        rel_ls.insert(0, "docs")
    rel = path.join(*rel_ls)

    if absolute:
        return abs
    else:
        return rel

def yield_text(mod):
    yield mod.name, mod.text()
    sm = {}
    for submod in mod.submodules():
        sm[submod.name] = docs_dir(submod.name, api=True)
        yield from yield_text(submod)

    if sm:
        if sum:=mod.supermodule:
            sum = f"\n\n## Super-module\n- [{sum.name}]({docs_dir(sum.name, False)})\n"
        else:
            sum = ""

        smls = []
        for k, v in sm.items():
            v = "/".join(v.split("/")[5:])
            smls.append(f'- [{k}]({v})')
        sm = "\n\n## Sub-modules\n\n{}\n".format("\n".join(smls))

        idx = """# {}{}{}""".format(
            mod.name,
            sum,
            sm,
        )

        idx_path = docs_dir(mod.name, api=True)
        with open(idx_path, "w") as f:
            f.write(idx)

def del_gen():
    try:
        for _, v in stg("generated", "docs/_meta.yml").items():
            for i in v["folders"]:
                if path.isdir(i):
                    shutil.rmtree(i)
            for i in v["files"]:
                if path.isfile(i):
                    os.remove(i)
    except TypeError:
        shutil.copy("docs/_meta.yml.bak", "docs/_meta.yml")
        del_gen()

def main(rmv: Dict[Any, Any]={}, hr :bool=False):
    docs_pdir = DOCS["op"]
    rmv_r = ddir(rmv, "rules")
    rmv_mv = ddir(rmv, "md_vars")
    MVC = dict(RMVC, **ddir(rmv_mv, "global"))

    del_gen()

    for rip in list(IDF.rglob("*.ymd")):
        out = path.join(
            docs_pdir,
            *rip.parts[1:-1],
            f"{rip.stem}.md"
        )
        GEN["docs"][1].append(out)

        rf = frontmatter.load(rip)
        md = repl(rf.content, dd(rules_fn(RULES), rules_fn(rmv_r)))

        d = {
            **MVC,
            **ddir(
                MD_VARS_YML,
                f"local/{rip.stem}"
            ),
            **ddir(
                rmv_mv,
                f"local/{rip.stem}"
            )
        }
        for k, v in {k: str(v) for k, v in d.items()}.items():
            md = md.replace(f"${{{k}}}", v)

        if title:=rf.get("title"):
            if link:=rf.get("link"):
                md = """<h1 align="center" style="font-weight: bold">
    <a target="_blank" href="{}">{}</a>
</h1>\n\n{}\n""".format(link, title, md)
            else:
                md = """<h1 align="center" style="font-weight: bold">
    {}
</h1>\n\n{}\n""".format(title, md)

        with open(inmd(out, GEN["docs"][0]), "w") as f:
            f.write(md)

    for pip in list(IDF.rglob("*.mako")):
        ip = str(pip)
        with open(ip, "r") as f:
            iprl = f.readline()
        if iprl != "## ignore\n":
            op = path.join(
                docs_pdir,
                *pip.parts[1:-1],
                f"{pip.stem}.md"
            )
            mytemplate = Template(filename=ip)
            tpl_rd = mytemplate.render(
                **{
                    "cwd": Path(dn(ip)),
                    "VARS": {
                        **MVC,
                        **ddir(
                            MD_VARS_YML,
                            f"local/{pip.stem}"
                        ),
                        **ddir(
                            rmv_mv,
                            f"local/{pip.stem}"
                        )
                    },
                }
            )
            with open(inmd(op, GEN["mako"][0]), "w") as f:
                f.write(tpl_rd)
            GEN["mako"][1].append(op)

    shutil.copy("docs/_meta.yml", "docs/_meta.yml.bak")
    for k, v in GEN.items():
        for key, i in zip(["folders", "files"], v):
            wr_stg(f"generated/{k}/{key}", list(set(i)), "docs/_meta.yml")

    if not hr:
        base = path.join("raw_docs", "docs")

        ndd = {}
        for u in listdir(base):
            for d in listdir(path.join(base, u)):
                ndd[f"{u}.{d}"] = f"docs/{u}/{d}/index.html"
        lk = list(ndd.keys())[-1]
        ndd[f'{lk} (Current)'] = ndd.pop(lk)
        nd = yaml.dump(ndd, default_flow_style=False)
        nd = "\n".join([f'    - {i}' for i in nd.strip().split("\n")][::-1])

        

        with open("mkdocs.yml", "r") as f:
            mkdocs = f.read()

        mkdocs = re.sub(RE_MD_TPL.format(key="nav docs"), nd, mkdocs, 1, flags=re.DOTALL)

        with open("mkdocs.yml", "w") as f:
            f.write(mkdocs)
        shutil.copy("mkdocs.yml", "mkdocs.yml.bak")
