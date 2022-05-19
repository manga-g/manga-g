<%
    import re
    from os import path

    RE_H1 = r"# ([^\n]*)(\n[^#]+)(.*?)(?=\n#\s|$)"
    RE_H2 = r"## ([^\n]*)\s*(.*?)(?=\n##\s|$)"
    RE_LC = r"- (.*?)(?=- |$)"

    PR = ["alpha", "beta", "rc"]

    def rv(vls: list[str | int]):
        pr = ""
        vls = [int(i) for i in vls]
        if vls[4] <= 2:
            pr = f'-{PR[vls[4]]}.{vls[5]}'
        return ".".join([str(i) for i in vls[0:4]]) + pr

    with open(path.join(cwd, "changelog.mmd"), "r") as f:
        CL = f.read()

    d_op = {}
    for vb in re.findall(RE_H1, CL, re.DOTALL):
        ver, desc, tls = vb
        vls = ver.split(" ")
        d_op[rv(vls)] = ov = {
            "vls": vls,
            "anchor": ver.replace(" ", "-"),
            "desc": desc.strip(),
            "changes": {}
        }
        for tch in re.findall(RE_H2, tls, re.DOTALL):
            t, chls = tch
            ov["changes"][t] = ovt = []
            for ch in re.findall(RE_LC, chls, re.DOTALL):
                ovt.append(ch)

    md_op = ""
    for k, v in d_op.items():
        md_op += f'\n\n<h2 id="{v["anchor"]}">{k}</h2>'
        if desc:=v["desc"]:
            md_op += f'\n\n{desc}\n\n'
        else:
            md_op += ''
        for t, chls in v["changes"].items():
            md_op += f'\n\n<h3 id="{v["anchor"]}-{t.lower()}">{t}</h3>\n'
            for ch in chls:
                md_op += f'\n- {ch}'
%>
<h1 align="center" style="font-weight: bold">
    Changelog
</h1>${md_op}