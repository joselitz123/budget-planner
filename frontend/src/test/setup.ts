/**
 * Vitest setup file for testing
 * Configures global test environment, mocks, and utilities
 */

import fakeIndexedDB from "fake-indexeddb";
import { beforeEach, vi } from "vitest";

// Mock IndexedDB globally - fake-indexeddb exports an indexedDB instance directly
global.indexedDB = fakeIndexedDB;
global.IDBKeyRange = fakeIndexedDB.IDBKeyRange;

// Mock crypto.randomUUID for deterministic tests
vi.stubGlobal("crypto", {
  randomUUID: () => `test-${Math.random().toString(36).substring(2, 15)}`,
});

// Mock environment variables
vi.stubGlobal("import.meta", {
  env: {
    PUBLIC_API_URL: "http://localhost:8080/api",
    PUBLIC_CLERK_PUBLISHABLE_KEY: "pk_test_mock_key",
  },
});

// Mock fetch globally
let mockFetchResponses = new Map();

export function mockFetch(
  url: string,
  response: Partial<Response> & { data?: unknown }
) {
  const { data, ...responseInit } = response;
  mockFetchResponses.set(url, {
    ...responseInit,
    json: async () => data,
    ok: responseInit.status ? responseInit.status < 400 : true,
    status: responseInit.status || 200,
  } as Response);
}

export function clearMockFetch() {
  mockFetchResponses.clear();
}

global.fetch = vi.fn(async (url: string) => {
  const mock = mockFetchResponses.get(url);
  if (mock) {
    return mock;
  }
  throw new Error(`Unexpected fetch call to: ${url}`);
}) as unknown as typeof fetch;

// Create a mock writable store helper
export function createMockStore<T>(initialValue: T) {
  let value = initialValue;
  const listeners = new Set<(val: T) => void>();
  return {
    set: vi.fn((val: T) => {
      value = val;
      listeners.forEach((l) => l(val));
    }),
    get: vi.fn(() => value),
    subscribe: vi.fn((listener: (val: T) => void) => {
      listeners.add(listener);
      listener(value);
      return () => listeners.delete(listener);
    }),
    update: vi.fn((updater: (val: T) => T) => {
      const newVal = updater(value);
      value = newVal;
      listeners.forEach((l) => l(newVal));
    }),
  };
}

// Clean up after each test
beforeEach(() => {
  vi.clearAllMocks();
  clearMockFetch();
});
