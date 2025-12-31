import { defineConfig } from 'vitest/config';
import { svelte } from '@sveltejs/vite-plugin-svelte';
import { resolve } from 'node:path';

export default defineConfig({
	plugins: [svelte({ hot: !process.env.VITEST })],
	test: {
		globals: true,
		environment: 'jsdom',
		setupFiles: ['./src/test/setup.ts'],
		include: ['src/**/*.{test,spec}.{js,ts}'],
		coverage: {
			provider: 'v8',
			reporter: ['text', 'json', 'html'],
			exclude: [
				'node_modules/',
				'src/test/',
				'*.config.ts',
				'*.config.js',
				'src/lib/components/ui/'
			]
		},
		testTimeout: 10000
	},
	resolve: {
		alias: {
			$lib: resolve(__dirname, './src/lib'),
			'$test': resolve(__dirname, './src/test')
		}
	}
});
