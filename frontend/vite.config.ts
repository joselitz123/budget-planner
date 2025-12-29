import { sveltekit } from '@sveltejs/kit/vite';
import { defineConfig } from 'vite';
import { VitePWA } from 'vite-plugin-pwa';

export default defineConfig({
	plugins: [
		sveltekit(),
		VitePWA({
			srcDir: 'static',
			manifest: {
				name: 'Budget Planner',
				short_name: 'Budget',
				description: 'Offline-first personal budget planning with notebook aesthetic',
				theme_color: '#333333',
				background_color: '#fdfbf7',
				display: 'standalone',
				icons: [
					{
						src: '/icons/icon-192x192.png',
						sizes: '192x192',
						type: 'image/png'
					},
					{
						src: '/icons/icon-512x512.png',
						sizes: '512x512',
						type: 'image/png'
					}
				]
			},
			strategies: 'generateSW',
			devOptions: {
				enabled: true
			},
			workbox: {
				// Increase limit to handle Clerk SDK (3MB+)
				maximumFileSizeToCacheInBytes: 5 * 1024 * 1024 // 5 MB
			}
		})
	],
	resolve: {
		alias: {
			$lib: '/src/lib'
		}
	}
});
