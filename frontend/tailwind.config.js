/** @type {import('tailwindcss').Config} */
export default {
	content: ['./src/**/*.{html,js,svelte,ts}'],
	darkMode: 'class',
	theme: {
		extend: {
			colors: {
				primary: '#333333', // Dark ink color
				'paper-light': '#fdfbf7', // Warm creamy paper
				'paper-dark': '#1f1f1f', // Dark mode paper
				'line-light': '#e5e7eb', // Light notebook lines
				'line-dark': '#374151', // Dark mode lines
				'accent-gold': '#d4af37', // Spiral binding gold
				'accent-highlight': '#FEF3C7' // Highlighter yellow
			},
			fontFamily: {
				display: ['"Playfair Display"', 'serif'],
				body: ['"Inter"', 'sans-serif'],
				handwriting: ['"Caveat"', 'cursive'],
				hand: ['"Kalam"', 'cursive']
			},
			boxShadow: {
				'paper': '0 4px 6px -1px rgba(0, 0, 0, 0.1), 0 2px 4px -1px rgba(0, 0, 0, 0.06), 2px 0 10px rgba(0,0,0,0.05)'
			},
			backgroundImage: {
				'paper-pattern':
					"url('data:image/svg+xml;base64,PHN2ZyB3aWR0aD0iMjAiIGhlaWdodD0iMjAiIHhtbG5zPSJodHRwOi8vd3d3LnczLm9yZy8yMDAwL3N2ZyI+PGNpcmNsZSBjeD0iMSIgY3k9IjEiIHI9IjEiIGZpbGw9IiM5OTkiIGZpbGwtb3BhY2l0eT0iMC4wMyIvPjwvc3ZnPg==')"
			},
			keyframes: {
				shimmer: {
					'0%': { backgroundPosition: '-1000px 0' },
					'100%': { backgroundPosition: '1000px 0' }
				}
			},
			animation: {
				shimmer: 'shimmer 2s infinite linear'
			}
		}
	},
	plugins: [require('@tailwindcss/forms')]
};
