#!/usr/bin/env python3
"""stagit-gen.py - static, stagit-style git pages for every public, non-fork
repo of a GitHub user. Needs only git, python3 and the gh CLI.

Usage: python3 stagit-gen.py [github-user] [out-dir] [cache-dir]
Defaults: equwal git .stagit-cache

Output layout (like stagit): OUT/index.html (repo list), and per repo
OUT/<repo>/{index.html (log), files.html, refs.html, file/<path>.html,
commit/<sha>.html}. Private repos and forks are never listed.
"""
import html, json, os, subprocess, sys, shutil
from pathlib import Path

USER = sys.argv[1] if len(sys.argv) > 1 else "equwal"
OUT = Path(sys.argv[2] if len(sys.argv) > 2 else "git")
CACHE = Path(sys.argv[3] if len(sys.argv) > 3 else ".stagit-cache")
MAX_COMMITS = 500
MAX_COMMIT_PAGES = 100
MAX_FILE = 100_000
MAX_FILES = 400
MAX_DIFF = 200_000
SITE = "Recently Written"
OWNER = "Spenser Truex"


def run(*a, cwd=None, check=True):
    r = subprocess.run(a, cwd=cwd, capture_output=True, check=check)
    return r.stdout.decode("utf-8", "replace")


def esc(s):
    return html.escape(s, quote=True)


def page(title, body, depth, header=""):
    up = "../" * depth
    return f"""<!DOCTYPE html>
<html lang="en"><head><meta charset="UTF-8" />
<meta name="viewport" content="width=device-width, initial-scale=1" />
<title>{esc(title)} - {SITE}</title>
<link rel="stylesheet" href="{up}stagit.css" />
</head><body>
<div id="top"><a href="{up}../index.html">{SITE}</a> &middot; <a href="{up}index.html">git</a></div>
{header}
<div id="content">
{body}
</div></body></html>
"""


def write(path, text):
    path.parent.mkdir(parents=True, exist_ok=True)
    path.write_text(text, encoding="utf-8")


def repo_header(name, desc, depth, url):
    up = "../" * depth
    return (
        f'<h1>{esc(name)}</h1><span class="desc">{esc(desc or "")}</span>'
        f'<p class="url">git clone {esc(url)}</p>'
        f'<p><a href="{up}{esc(name)}/index.html">Log</a> | '
        f'<a href="{up}{esc(name)}/files.html">Files</a> | '
        f'<a href="{up}{esc(name)}/refs.html">Refs</a></p><hr/>'
    )


def sync(name):
    dest = CACHE / f"{name}.git"
    url = f"https://github.com/{USER}/{name}.git"
    if dest.exists():
        run("git", "fetch", "--prune", "--tags", "origin", "+refs/heads/*:refs/heads/*", cwd=dest, check=False)
    else:
        dest.parent.mkdir(parents=True, exist_ok=True)
        run("git", "clone", "--bare", url, str(dest))
    return dest


def build_repo(meta):
    name, desc = meta["name"], meta.get("description")
    url = f"https://github.com/{USER}/{name}"
    try:
        g = sync(name)
    except subprocess.CalledProcessError as e:
        print(f"  skip {name}: {e.stderr.decode('utf-8','replace').strip()[:120]}")
        return None
    head = run("git", "rev-parse", "--verify", "HEAD", cwd=g, check=False).strip()
    rdir = OUT / name
    if rdir.exists():
        shutil.rmtree(rdir)
    hdr = repo_header(name, desc, 1, url)
    if not head:
        write(rdir / "index.html", page(name, "<p>Empty repository.</p>", 1, hdr))
        return {"name": name, "desc": desc, "last": ""}

    sep = "\x1f"
    fmt = sep.join(["%H", "%as", "%an", "%s", "%ad"]) + "\x1e"
    raw = run("git", "log", f"-{MAX_COMMITS}", "--date=format:%Y-%m-%d %H:%M", f"--format={fmt}", cwd=g)
    commits = [c.strip("\n").split(sep) for c in raw.split("\x1e") if c.strip()]
    last = commits[0][4] if commits else ""
    total = run("git", "rev-list", "--count", "HEAD", cwd=g).strip()

    rows = []
    for i, (h, _d, an, subj, ad) in enumerate(commits):
        link = f'<a href="commit/{h}.html">{esc(subj)}</a>' if i < MAX_COMMIT_PAGES else esc(subj)
        rows.append(f"<tr><td>{esc(ad)}</td><td>{link}</td><td>{esc(an)}</td></tr>")
    note = "" if len(commits) == int(total) else f"<p>Showing latest {len(commits)} of {total} commits.</p>"
    write(rdir / "index.html", page(name, note + "<table><thead><tr><td>Date</td><td>Commit message</td><td>Author</td></tr></thead><tbody>" + "".join(rows) + "</tbody></table>", 1, hdr))

    for h, *_ in commits[:MAX_COMMIT_PAGES]:
        meta_txt = run("git", "show", "-s", "--format=%an <%ae>%n%ad%n%n%B", "--date=iso", h, cwd=g)
        stat = run("git", "show", "--format=", "--stat", h, cwd=g)
        diff = run("git", "show", "--format=", "-p", "--no-color", h, cwd=g)
        if len(diff) > MAX_DIFF:
            diff = diff[:MAX_DIFF] + "\n[diff truncated]\n"
        d = []
        for line in diff.splitlines():
            cls = "a" if line.startswith("+") and not line.startswith("+++") else "d" if line.startswith("-") and not line.startswith("---") else "h" if line.startswith("@@") else ""
            d.append(f'<span class="{cls}">{esc(line)}</span>' if cls else esc(line))
        body = f"<pre>commit {h}\n{esc(meta_txt)}</pre><pre>{esc(stat)}</pre><pre>{chr(10).join(d)}</pre>"
        write(rdir / "commit" / f"{h}.html", page(f"{name} {h[:8]}", body, 2, repo_header(name, desc, 2, url)))

    ls = run("git", "ls-tree", "-r", "-l", "-z", "HEAD", cwd=g).split("\0")
    frows, n = [], 0
    for ent in ls:
        if not ent:
            continue
        info, path = ent.split("\t", 1)
        mode, typ, sha, size = info.split()
        size = 0 if size == "-" else int(size)
        can = typ == "blob" and mode != "160000" and size <= MAX_FILE and n < MAX_FILES
        text = None
        if can:
            blob = subprocess.run(["git", "cat-file", "blob", sha], cwd=g, capture_output=True).stdout
            if b"\0" not in blob:
                text = blob.decode("utf-8", "replace")
        if text is not None:
            n += 1
            depth = 2 + path.count("/")
            lines = text.splitlines() or [""]
            code = "\n".join(f'<a id="l{i}" href="#l{i}">{i}</a> {esc(l)}' for i, l in enumerate(lines, 1))
            fh = repo_header(name, desc, depth, url)
            write(rdir / "file" / (path + ".html"), page(f"{name}/{path}", f"<p>{esc(path)} ({size} bytes)</p><pre class=\"src\">{code}</pre>", depth, fh))
            cell = f'<a href="file/{esc(path)}.html">{esc(path)}</a>'
        else:
            cell = f'{esc(path)} <span class="dim">(not shown)</span>'
        frows.append(f"<tr><td>{mode}</td><td>{cell}</td><td class=\"num\">{size}</td></tr>")
    write(rdir / "files.html", page(name, "<table><thead><tr><td>Mode</td><td>Name</td><td>Size</td></tr></thead><tbody>" + "".join(frows) + "</tbody></table>", 1, hdr))

    refs = run("git", "for-each-ref", "--format=%(refname:short)\x1f%(objecttype)\x1f%(creatordate:short)\x1f%(objectname:short)", "refs/heads", "refs/tags", cwd=g)
    rr = "".join(
        f"<tr><td>{esc(a)}</td><td>{esc(b)}</td><td>{esc(c)}</td><td>{esc(d)}</td></tr>"
        for a, b, c, d in (l.split("\x1f") for l in refs.splitlines() if l)
    )
    write(rdir / "refs.html", page(name, f"<table><thead><tr><td>Name</td><td>Type</td><td>Date</td><td>Commit</td></tr></thead><tbody>{rr}</tbody></table>", 1, hdr))
    return {"name": name, "desc": desc, "last": last}


def main():
    repos = json.loads(run("gh", "repo", "list", USER, "--limit", "500", "--visibility", "public", "--source", "--json", "name,description,isPrivate,isFork"))
    repos = [r for r in repos if not r["isPrivate"] and not r["isFork"]]
    repos.sort(key=lambda r: r["name"].lower())
    OUT.mkdir(parents=True, exist_ok=True)
    print(f"{len(repos)} public non-fork repos")
    done = []
    for r in repos:
        print(" ", r["name"])
        m = build_repo(r)
        if m:
            done.append(m)
    rows = "".join(
        f'<tr><td><a href="{esc(m["name"])}/index.html">{esc(m["name"])}</a></td><td>{esc(m["desc"] or "")}</td><td>{esc(OWNER)}</td><td>{esc(m["last"])}</td></tr>'
        for m in done
    )
    idx = f"<h1>ye olde code shoppe</h1><p>Mirror of <a href=\"https://github.com/{USER}\">github.com/{USER}</a> (public, non-fork repos).</p><table><thead><tr><td>Name</td><td>Description</td><td>Owner</td><td>Last commit</td></tr></thead><tbody>{rows}</tbody></table>"
    write(OUT / "index.html", page("git", idx, 0))
    (OUT / "stagit.css").write_text(CSS, encoding="utf-8")
    Path("git.html").write_text('<!DOCTYPE html><meta http-equiv="refresh" content="0; url=git/index.html"><a href="git/index.html">git</a>\n', encoding="utf-8")
    print(f"wrote {len(done)} repos to {OUT}/")


CSS = """body{font-family:monospace;color:#222;background:#fff;max-width:1000px;margin:0 auto;padding:1rem}
a{color:#00e}h1{font-size:1.4em;margin:.4em 0 0}#top{border-bottom:1px solid #ccc;padding-bottom:.4em;margin-bottom:.8em}
.desc{color:#555}.url{color:#555}.dim{color:#999}table{border-collapse:collapse;width:100%}
td{padding:.15em .6em .15em 0;vertical-align:top}thead td{font-weight:bold;border-bottom:1px solid #ccc}
.num{text-align:right}pre{overflow:auto}.src a{color:#999;text-decoration:none}.a{color:#080}.d{color:#c00}.h{color:#08a}
@media(prefers-color-scheme:dark){body{background:#111;color:#ddd}a{color:#8af}.src a{color:#777}#top{border-color:#444}thead td{border-color:#444}.a{color:#6c6}.d{color:#f66}}
"""

if __name__ == "__main__":
    main()
