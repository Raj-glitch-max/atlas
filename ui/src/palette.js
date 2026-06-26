// Command palette (⌘K) — Raycast metaphor. Pure DOM, no deps.
const $ = (s, r = document) => r.querySelector(s);
const $$ = (s, r = document) => [...r.querySelectorAll(s)];

export function initPalette(cmds) {
  const scrim = $("#palScrim"), input = $("#palInput"), list = $("#palList");
  if (!scrim) return;
  let sel = 0, filtered = cmds;

  const open = () => { scrim.classList.add("show"); input.value = ""; render(""); input.focus(); };
  const close = () => scrim.classList.remove("show");
  const isOpen = () => scrim.classList.contains("show");

  function render(q) {
    filtered = cmds.filter((c) => c.label.toLowerCase().includes(q.toLowerCase()));
    sel = 0;
    list.innerHTML = filtered.map((c, i) =>
      `<div class="it ${i ? "" : "sel"}" data-i="${i}"><svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.8">${c.icon}</svg>${c.label}<span class="pk">${c.key || ""}</span></div>`
    ).join("") || `<div class="it" style="color:var(--fg4)">No matches</div>`;
    $$(".it", list).forEach((it) => it.addEventListener("click", () => { const c = filtered[+it.dataset.i]; if (c) c.run(); close(); }));
  }
  const move = (d) => { if (!filtered.length) return; sel = (sel + d + filtered.length) % filtered.length; $$(".it", list).forEach((it, i) => it.classList.toggle("sel", i === sel)); };

  $("#cmdk")?.addEventListener("click", open);
  scrim.addEventListener("click", (e) => { if (e.target === scrim) close(); });
  input.addEventListener("input", (e) => render(e.target.value));
  addEventListener("keydown", (e) => {
    if ((e.metaKey || e.ctrlKey) && e.key.toLowerCase() === "k") { e.preventDefault(); isOpen() ? close() : open(); }
    else if (isOpen()) {
      if (e.key === "Escape") close();
      else if (e.key === "ArrowDown") { e.preventDefault(); move(1); }
      else if (e.key === "ArrowUp") { e.preventDefault(); move(-1); }
      else if (e.key === "Enter") { const c = filtered[sel]; if (c) c.run(); close(); }
    }
  });
}
