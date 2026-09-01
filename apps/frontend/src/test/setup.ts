import '@testing-library/jest-dom/vitest'

// jsdom tak menyediakan `window.matchMedia`. Sejumlah komponen (mis. sonner
// `<Toaster/>` yang di-mount di `AppShell`) memanggilnya saat render. Stub
// minimal supaya test yang menyentuh shell tak meledak.
if (!window.matchMedia) {
  window.matchMedia = (query: string): MediaQueryList => ({
    matches: false,
    media: query,
    onchange: null,
    addListener: () => {},
    removeListener: () => {},
    addEventListener: () => {},
    removeEventListener: () => {},
    dispatchEvent: () => false,
  })
}
