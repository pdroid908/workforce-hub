# ENTERPRISE LIGHT UI/UX DESIGN SYSTEM & LAYOUT DIRECTIVES

## 🎯 Purpose & Scope
Dokumen ini menetapkan standar antarmuka visual, hirarki komponen, skema warna *Light Enterprise*, dan tata letak UI untuk dashboard data real-time berkategori **Enterprise Grade**.

---

## 🎨 Enterprise Light Color Palette (White + Emerald / Sky Gradient)

### 1. Base Canvas & Glassmorphism
* **Canvas Background:** Background putih dengan gradasi lembut dan aksen mesh terukur.
  * *Tailwind:* `bg-slate-50 bg-[radial-gradient(ellipse_80%_80%_at_50%_-20%,rgba(14,165,233,0.15),rgba(255,255,255,0))]`
  * *Alternatif Emerald:* `bg-[radial-gradient(ellipse_80%_80%_at_50%_-20%,rgba(16,185,129,0.15),rgba(255,255,255,0))]`
* **Card Surface:** `bg-white/80 backdrop-blur-xl`
* **Border System:** Boundary crisp bernuansa Slate/Sky berkontras halus agar card terlihat menonjol (*pop-out*).
  * *Tailwind:* `border border-slate-200/80 shadow-xl shadow-slate-200/50`

### 2. Typography & Contrast Rules
* **Font Family:** `Plus Jakarta Sans` / `Inter`
* **Primary Text:** `text-slate-900` (Gunakan Slate gelap, **hindari** `#000000` murni)
* **Secondary Text:** `text-slate-500`
* **Accent Color (Sky):** Primary `bg-sky-500`, Text `text-sky-600`, Border `border-sky-200`
* **Accent Color (Emerald):** Primary `bg-emerald-500`, Text `text-emerald-600`, Border `border-emerald-200`

---

## 📐 Enterprise Top-Pinned Layout Hierarchy

Semua halaman wajib mengikuti struktur hirarki tanpa melakukan *scrolling* tambahan untuk melihat hasil utama:

```text
┌──────────────────────────────────────────────────────────┐
│ 1. HEADER UTAMA & LIVE STATUS INDICATOR                  │
├──────────────────────────────────────────────────────────┤
│ 2. CONTROL BAR / SEARCH INPUT (Tab Switcher + GPS)       │
├──────────────────────────────────────────────────────────┤
│ 3. PINNED PRIMARY RESULT PANEL (Atas / Top-Anchor)       │
│    ├── Hero Card: Metric Value & Dynamic Color Badge    │
│    ├── Sub-Grid: Polutan, Suhu, Kelembapan               │
│    └── Action Card: Health & Operational Advice          │
├──────────────────────────────────────────────────────────┤
│ 4. REFERENCE & THRESHOLD GUIDE (Tabel Standar AQI Bawah) │
└──────────────────────────────────────────────────────────┘