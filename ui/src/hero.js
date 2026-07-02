// WebGL hero — the "two trust domains" scene.
//
// Design intent (Igloo/Lusion depth + Raycast restraint, NOT a generic 3D
// gimmick): two point-cloud constellations — domain A (ice) and domain B
// (steel) — drift in a foggy deep-space volume, joined by one luminous,
// verifiable trust path that a signature pulse travels. Glow is done with
// additive point sprites (a soft canvas texture), not UnrealBloom — so it
// holds 60fps on mid-range laptops with no post-processing pass.
//
// Lazy-loaded (dynamic import of three) so it never blocks first paint or the
// reduced-motion path. Exposes { setScroll, pause, resume, dispose }.
import * as THREE from "three";

const POINTS_PER_CLUSTER = 800; // 1600 total — dense but cheap

function softSprite() {
  const c = document.createElement("canvas");
  c.width = c.height = 64;
  const x = c.getContext("2d");
  const g = x.createRadialGradient(32, 32, 0, 32, 32, 32);
  g.addColorStop(0, "rgba(255,255,255,1)");
  g.addColorStop(0.25, "rgba(255,255,255,.85)");
  g.addColorStop(1, "rgba(255,255,255,0)");
  x.fillStyle = g; x.fillRect(0, 0, 64, 64);
  const t = new THREE.CanvasTexture(c);
  t.colorSpace = THREE.SRGBColorSpace;
  return t;
}

export async function initHero(canvas) {
  const reduce = matchMedia("(prefers-reduced-motion: reduce)").matches;
  const host = canvas.parentElement;
  let W = host.clientWidth, H = host.clientHeight;

  const renderer = new THREE.WebGLRenderer({ canvas, antialias: true, alpha: false, powerPreference: "high-performance" });
  renderer.setPixelRatio(Math.min(devicePixelRatio || 1, 1.8));
  renderer.setSize(W, H, false);
  renderer.setClearColor(0x07080a, 1);

  const scene = new THREE.Scene();
  scene.fog = new THREE.FogExp2(0x07080a, 0.115);
  const camera = new THREE.PerspectiveCamera(52, W / H, 0.1, 100);
  camera.position.set(0, 0, 6.4);

  const group = new THREE.Group();
  scene.add(group);

  // ---- point-cloud shader (twinkle + size attenuation) ----
  // Sprite size scales with viewport height and is clamped in the shader:
  // GPUs clamp gl_PointSize at wildly different maxima (Intel ~255, NVIDIA
  // ~2047, SwiftShader higher), so unclamped sizes render as delicate points
  // on one machine and a blown-out white mass on another.
  const sizeScale = () => 42 * Math.min(1.5, Math.max(0.75, H / 900));
  const uniforms = { uTime: { value: 0 }, uPix: { value: renderer.getPixelRatio() }, uScale: { value: sizeScale() }, uTex: { value: softSprite() } };
  const material = new THREE.ShaderMaterial({
    uniforms,
    transparent: true,
    depthWrite: false,
    blending: THREE.AdditiveBlending,
    vertexShader: `
      uniform float uTime; uniform float uPix; uniform float uScale;
      attribute float aSize; attribute float aPhase; attribute vec3 aColor;
      varying vec3 vColor; varying float vTw;
      void main(){
        vColor = aColor;
        vec3 p = position;
        p.x += sin(uTime*0.30 + aPhase)*0.07;
        p.y += cos(uTime*0.26 + aPhase*1.3)*0.07;
        vec4 mv = modelViewMatrix * vec4(p,1.0);
        vTw = 0.55 + 0.45*sin(uTime*1.1 + aPhase*3.0);
        gl_PointSize = min(aSize * uPix * (uScale / -mv.z), 44.0 * uPix);
        gl_Position = projectionMatrix * mv;
      }`,
    fragmentShader: `
      uniform sampler2D uTex;
      varying vec3 vColor; varying float vTw;
      void main(){
        float a = texture2D(uTex, gl_PointCoord).a;
        if(a < 0.02) discard;
        gl_FragColor = vec4(vColor*(0.6+0.7*vTw), a*vTw*0.85);
      }`,
  });

  const ICE = new THREE.Color(0x7dd3fc), STEEL = new THREE.Color(0xaeb9cc);
  function cluster(cx, tint) {
    const n = POINTS_PER_CLUSTER;
    const pos = new Float32Array(n * 3), col = new Float32Array(n * 3), size = new Float32Array(n), ph = new Float32Array(n);
    for (let i = 0; i < n; i++) {
      // ellipsoidal shell, denser toward the core
      const u = Math.random() * 2 - 1, th = Math.random() * Math.PI * 2, rr = Math.pow(Math.random(), 0.6);
      const sp = Math.sqrt(1 - u * u);
      pos[i * 3] = cx + sp * Math.cos(th) * rr * 2.6;
      pos[i * 3 + 1] = u * rr * 1.9;
      pos[i * 3 + 2] = sp * Math.sin(th) * rr * 2.2 - 0.5;
      const c = tint.clone().offsetHSL(0, 0, (Math.random() - 0.5) * 0.18);
      col[i * 3] = c.r; col[i * 3 + 1] = c.g; col[i * 3 + 2] = c.b;
      size[i] = 1.6 + Math.pow(Math.random(), 2.2) * 6.4;
      ph[i] = Math.random() * 6.28;
    }
    const g = new THREE.BufferGeometry();
    g.setAttribute("position", new THREE.BufferAttribute(pos, 3));
    g.setAttribute("aColor", new THREE.BufferAttribute(col, 3));
    g.setAttribute("aSize", new THREE.BufferAttribute(size, 1));
    g.setAttribute("aPhase", new THREE.BufferAttribute(ph, 1));
    return new THREE.Points(g, material);
  }
  group.add(cluster(-2.7, ICE));
  group.add(cluster(2.7, STEEL));

  // ---- the trust path (issuer → delegate → verifier → decision) ----
  const path = new THREE.CatmullRomCurve3([
    new THREE.Vector3(-2.7, -0.5, 0.4),
    new THREE.Vector3(-0.9, 0.7, 0.9),
    new THREE.Vector3(1.1, -0.4, 0.9),
    new THREE.Vector3(2.7, 0.5, 0.4),
  ]);
  const pts = path.getPoints(120);
  const lineGeo = new THREE.BufferGeometry().setFromPoints(pts);
  const line = new THREE.Line(lineGeo, new THREE.LineBasicMaterial({ color: 0x7dd3fc, transparent: true, opacity: 0.22, blending: THREE.AdditiveBlending, depthWrite: false }));
  group.add(line);

  // hub markers + a traveling signature pulse, as additive sprites
  const sprMap = uniforms.uTex.value;
  const mkSprite = (color, scale) => {
    const s = new THREE.Sprite(new THREE.SpriteMaterial({ map: sprMap, color, transparent: true, blending: THREE.AdditiveBlending, depthWrite: false }));
    s.scale.setScalar(scale); return s;
  };
  const hubPts = [path.getPointAt(0), path.getPointAt(0.34), path.getPointAt(0.68), path.getPointAt(1)];
  const hubColors = [0x7dd3fc, 0x7dd3fc, 0xaeb9cc, 0xaeb9cc];
  hubPts.forEach((p, i) => { const s = mkSprite(hubColors[i], 0.5); s.position.copy(p); group.add(s); });
  const pulse = mkSprite(0xffffff, 0.55); group.add(pulse);

  // ---- interaction + loop ----
  const pointer = { x: 0, y: 0, ex: 0, ey: 0 };
  const onMove = (e) => { pointer.x = e.clientX / innerWidth - 0.5; pointer.y = e.clientY / innerHeight - 0.5; };
  addEventListener("pointermove", onMove, { passive: true });

  let scroll = 0, running = true, raf = 0, last = performance.now();
  const clock = new THREE.Clock();

  function resize() {
    W = host.clientWidth; H = host.clientHeight;
    renderer.setSize(W, H, false); camera.aspect = W / H; camera.updateProjectionMatrix();
    uniforms.uPix.value = renderer.getPixelRatio();
    uniforms.uScale.value = sizeScale();
  }
  addEventListener("resize", resize);

  function frame() {
    if (!running) return;
    raf = requestAnimationFrame(frame);
    const now = performance.now();
    // frame-rate guard: skip heavy work if a frame ran long (keeps mid-range smooth)
    last = now;
    const el = clock.getElapsedTime();
    uniforms.uTime.value = el;

    // ambient rotation + pointer parallax (eased)
    pointer.ex += (pointer.x - pointer.ex) * 0.05; pointer.ey += (pointer.y - pointer.ey) * 0.05;
    group.rotation.y = Math.sin(el * 0.06) * 0.22 + pointer.ex * 0.5 - scroll * 0.35;
    group.rotation.x = pointer.ey * 0.28 + scroll * 0.18;

    // camera beat: dolly in + drop as the hero scrolls away
    camera.position.z = 6.4 - scroll * 2.4;
    camera.position.y = scroll * 0.9;
    camera.lookAt(0, 0, 0);

    // signature pulse travels the trust path
    const tt = (el * 0.11) % 1;
    pulse.position.copy(path.getPointAt(tt));
    const beat = 0.4 + 0.3 * Math.sin(el * 4);
    pulse.scale.setScalar(0.45 + beat * 0.25);
    line.material.opacity = 0.16 + 0.12 * (0.5 + 0.5 * Math.sin(el * 0.8));

    renderer.render(scene, camera);
  }
  if (!reduce) { running = true; frame(); }
  else { renderer.render(scene, camera); running = false; } // one static frame

  return {
    setScroll(p) { scroll = Math.max(0, Math.min(1, p)); if (reduce) renderer.render(scene, camera); },
    pause() { if (running) { running = false; cancelAnimationFrame(raf); } },
    resume() { if (!running && !reduce) { running = true; last = performance.now(); frame(); } },
    dispose() { this.pause(); removeEventListener("pointermove", onMove); removeEventListener("resize", resize); renderer.dispose(); },
  };
}
