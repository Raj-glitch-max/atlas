// Smooth-scroll + scene choreography (lazy-loaded, motion-only).
//   • Lenis inertia scroll, driven by GSAP's ticker and synced to ScrollTrigger
//   • hero "camera beat" scrub + render pause when offscreen
//   • scroll-progress rail, cursor spotlight, active-section nav
//   • offset-aware anchor smooth-scroll
//   • reveal-on-scroll, metric count-up, card pointer-glow
// All progressive enhancement — if this never runs, content is already visible.
import Lenis from "lenis";
import { gsap } from "gsap";
import { ScrollTrigger } from "gsap/ScrollTrigger";
import { runVerify } from "./verify.js";

gsap.registerPlugin(ScrollTrigger);

export function initScroll(hero) {
  /* ---------- Lenis smooth scroll, on the GSAP ticker ---------- */
  const lenis = new Lenis({ lerp: 0.1, smoothWheel: true, wheelMultiplier: 1, touchMultiplier: 1.4 });
  lenis.on("scroll", ScrollTrigger.update);
  gsap.ticker.add((time) => lenis.raf(time * 1000));
  gsap.ticker.lagSmoothing(0);

  /* ---------- scroll-progress rail ---------- */
  const prog = document.createElement("div"); prog.id = "progress"; document.body.appendChild(prog);
  lenis.on("scroll", ({ scroll, limit }) => { prog.style.transform = `scaleX(${limit ? scroll / limit : 0})`; });

  /* ---------- cursor spotlight (fine pointers only) ---------- */
  if (matchMedia("(pointer:fine)").matches) {
    const spot = document.createElement("div"); spot.id = "spot"; document.body.appendChild(spot);
    let sx = 0, sy = 0, queued = false;
    addEventListener("pointermove", (e) => {
      sx = e.clientX; sy = e.clientY;
      if (!queued) { queued = true; requestAnimationFrame(() => { spot.style.setProperty("--mx", sx + "px"); spot.style.setProperty("--my", sy + "px"); queued = false; }); }
    }, { passive: true });
    requestAnimationFrame(() => spot.classList.add("on"));
  }

  /* ---------- offset-aware anchor smooth-scroll ---------- */
  document.querySelectorAll('a[href^="#"]').forEach((a) => {
    a.addEventListener("click", (e) => {
      const id = a.getAttribute("href"); if (!id || id.length < 2) return;
      const el = document.querySelector(id); if (!el) return;
      e.preventDefault(); lenis.scrollTo(el, { offset: -64, duration: 1.1 });
    });
  });

  /* ---------- hero camera beat + render pause ---------- */
  const heroEl = document.querySelector(".hero"), canvas = document.getElementById("field");
  if (hero && heroEl) {
    ScrollTrigger.create({
      trigger: heroEl, start: "top top", end: "bottom top", scrub: true,
      onUpdate: (self) => { hero.setScroll(self.progress); if (canvas) canvas.style.opacity = String(1 - self.progress * 0.92); },
      onLeave: () => hero.pause(), onEnterBack: () => hero.resume(),
    });
  }

  /* ---------- reveal-on-scroll ---------- */
  document.querySelectorAll(".reveal").forEach((el) => {
    ScrollTrigger.create({ trigger: el, start: "top 90%", once: true, onEnter: () => el.classList.add("in") });
  });

  /* ---------- active-section nav ---------- */
  document.querySelectorAll(".nav .links a").forEach((a) => {
    const sec = document.querySelector(a.getAttribute("href")); if (!sec) return;
    ScrollTrigger.create({
      trigger: sec, start: "top 45%", end: "bottom 45%",
      onToggle: (self) => { if (self.isActive) { document.querySelectorAll(".nav .links a").forEach((x) => x.classList.remove("active")); a.classList.add("active"); } },
    });
  });

  /* ---------- metric count-up ---------- */
  document.querySelectorAll(".metric .n[data-to]").forEach((el) => {
    const to = parseFloat(el.dataset.to); if (isNaN(to)) return;
    ScrollTrigger.create({
      trigger: el, start: "top 92%", once: true,
      onEnter: () => { const o = { v: 0 }; gsap.to(o, { v: to, duration: 1.1, ease: "power2.out", onUpdate: () => { el.textContent = Math.round(o.v); } }); },
    });
  });

  /* ---------- card pointer-tracking glow ---------- */
  document.querySelectorAll(".card").forEach((card) => {
    card.addEventListener("pointermove", (e) => {
      const r = card.getBoundingClientRect();
      card.style.setProperty("--cx", ((e.clientX - r.left) / r.width * 100) + "%");
      card.style.setProperty("--cy", ((e.clientY - r.top) / r.height * 100) + "%");
    }, { passive: true });
  });

  /* ---------- auto-run the verify cinematic once in view ---------- */
  const cns = document.querySelector("#verify .console");
  if (cns) ScrollTrigger.create({ trigger: cns, start: "top 65%", once: true, onEnter: () => setTimeout(runVerify, 350) });

  ScrollTrigger.refresh();
}
