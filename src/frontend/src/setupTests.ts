import '@testing-library/jest-dom'
import { vi } from 'vitest'

// Mock window.matchMedia
Object.defineProperty(window, 'matchMedia', {
  writable: true,
  value: vi.fn().mockImplementation((query) => ({
    matches: false,
    media: query,
    onchange: null,
    addListener: vi.fn(),
    removeListener: vi.fn(),
    addEventListener: vi.fn(),
    removeEventListener: vi.fn(),
    dispatchEvent: vi.fn(),
  })),
})

// Add custom matchers
expect.extend({
  toBeInTheDocument(received) {
    if (!(received instanceof HTMLElement)) {
      return {
        pass: false,
        message: () => `Expected element to be in the document`,
      }
    }
    return {
      pass: true,
      message: () => `Expected element not to be in the document`,
    }
  },
})
