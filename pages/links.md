---
title: Links
---

Pages that I keep. The newest are first. The list comes from
[sbm-webpublish](https://github.com/equwal/sbm-webpublish).

<div id="links"><p>Loading the links...</p></div>

<noscript><p>The list needs JavaScript. The raw list is at
<a href="links/links.json">links/links.json</a>.</p></noscript>

<script>
// Show links.json from sbm-webpublish, grouped by the first tag of each link.
fetch("links/links.json")
  .then(function (r) { if (!r.ok) throw r.status; return r.json(); })
  .then(function (links) {
    var box = document.getElementById("links");
    box.textContent = "";
    var groups = {}, order = [];
    links.forEach(function (l) {
      var g = l.tags[0] || "other";
      if (!groups[g]) { groups[g] = []; order.push(g); }
      groups[g].push(l);
    });
    order.forEach(function (g) {
      var h = document.createElement("h3");
      h.textContent = g;
      var ul = document.createElement("ul");
      groups[g].forEach(function (l) {
        var li = document.createElement("li");
        var a = document.createElement("a");
        a.href = l.url;
        a.textContent = l.title;
        li.appendChild(a);
        li.appendChild(document.createTextNode(" (" + l.host + ")"));
        ul.appendChild(li);
      });
      box.appendChild(h);
      box.appendChild(ul);
    });
    if (!links.length) box.textContent = "No links yet.";
  })
  .catch(function () {
    document.getElementById("links").textContent = "The links did not load.";
  });
</script>
