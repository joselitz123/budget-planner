

## Login
```html
<!DOCTYPE html>

<html lang="en"><head>

<meta charset="utf-8"/>

<meta content="width=device-width, initial-scale=1.0" name="viewport"/>

<title>Budget Planner Login</title>

<script src="https://cdn.tailwindcss.com?plugins=forms,typography"></script>

<link href="https://fonts.googleapis.com" rel="preconnect"/>

<link crossorigin="" href="https://fonts.gstatic.com" rel="preconnect"/>

<link href="https://fonts.googleapis.com/css2?family=Caveat:wght@500;700&amp;family=Inter:wght@300;400;500;600&amp;family=Playfair+Display:ital,wght@0,600;1,600&amp;display=swap" rel="stylesheet"/>

<link href="https://fonts.googleapis.com/icon?family=Material+Icons+Outlined" rel="stylesheet"/>

<script>

tailwind.config = {

darkMode: "class",

theme: {

extend: {

colors: {

primary: "#333333", // Dark grey/black ink color for primary actions

"paper-light": "#fdfbf7", // Warm creamy paper color

"paper-dark": "#1f1f1f", // Dark mode paper

"line-light": "#e5e7eb", // Light grey lines

"line-dark": "#374151", // Dark mode lines

"accent-gold": "#d4af37", // Gold for spiral binding aesthetic

},

fontFamily: {

display: ["'Playfair Display'", "serif"],

body: ["'Inter'", "sans-serif"],

handwriting: ["'Caveat'", "cursive"],

},

boxShadow: {

'paper': '0 4px 6px -1px rgba(0, 0, 0, 0.1), 0 2px 4px -1px rgba(0, 0, 0, 0.06), 2px 0 10px rgba(0,0,0,0.05)',

},

backgroundImage: {

'paper-pattern': "url('data:image/svg+xml;base64,PHN2ZyB3aWR0aD0iMjAiIGhlaWdodD0iMjAiIHhtbG5zPSJodHRwOi8vd3d3LnczLm9yZy8yMDAwL3N2ZyI+PGNpcmNsZSBjeD0iMSIgY3k9IjEiIHI9IjEiIGZpbGw9IiM5OTkiIGZpbGwtb3BhY2l0eT0iMC4wMyIvPjwvc3ZnPg==')",

}

},

},

};

</script>

<style>.no-scrollbar::-webkit-scrollbar {

display: none;

}

.no-scrollbar {

-ms-overflow-style: none;

scrollbar-width: none;

}.binding-holes {

position: absolute;

left: 12px;

top: 0;

bottom: 0;

width: 20px;

display: flex;

flex-direction: column;

justify-content: space-evenly;

z-index: 10;

}

.binding-hole {

width: 12px;

height: 12px;

background-color: #f3f4f6;border-radius: 50%;

box-shadow: inset 1px 1px 2px rgba(0,0,0,0.2);

position: relative;

}.dark .binding-hole {

background-color: #374151;

box-shadow: inset 1px 1px 3px rgba(0,0,0,0.5);

}

.binding-coil {

position: absolute;

left: -4px;

top: 50%;

transform: translateY(-50%);

width: 24px;

height: 8px;

background: linear-gradient(90deg, #b8860b, #ffd700, #b8860b);

border-radius: 4px;

box-shadow: 1px 1px 2px rgba(0,0,0,0.3);

z-index: 20;

}

</style>

<style>

body {

min-height: max(884px, 100dvh);

}

</style>

</head>

<body class="bg-gray-200 dark:bg-gray-900 min-h-screen flex items-center justify-center font-body p-4 sm:p-0">

<div class="relative w-full max-w-md bg-paper-light dark:bg-paper-dark rounded-r-2xl rounded-l-md shadow-2xl overflow-hidden min-h-[750px] flex">

<div class="w-12 bg-gray-100 dark:bg-zinc-800 border-r border-gray-300 dark:border-zinc-700 relative flex-shrink-0">

<div class="binding-holes py-4">

<div class="relative h-8 w-full"><div class="binding-hole"></div><div class="binding-coil"></div></div>

<div class="relative h-8 w-full"><div class="binding-hole"></div><div class="binding-coil"></div></div>

<div class="relative h-8 w-full"><div class="binding-hole"></div><div class="binding-coil"></div></div>

<div class="relative h-8 w-full"><div class="binding-hole"></div><div class="binding-coil"></div></div>

<div class="relative h-8 w-full"><div class="binding-hole"></div><div class="binding-coil"></div></div>

<div class="relative h-8 w-full"><div class="binding-hole"></div><div class="binding-coil"></div></div>

<div class="relative h-8 w-full"><div class="binding-hole"></div><div class="binding-coil"></div></div>

<div class="relative h-8 w-full"><div class="binding-hole"></div><div class="binding-coil"></div></div>

<div class="relative h-8 w-full"><div class="binding-hole"></div><div class="binding-coil"></div></div>

<div class="relative h-8 w-full"><div class="binding-hole"></div><div class="binding-coil"></div></div>

<div class="relative h-8 w-full"><div class="binding-hole"></div><div class="binding-coil"></div></div>

<div class="relative h-8 w-full"><div class="binding-hole"></div><div class="binding-coil"></div></div>

<div class="relative h-8 w-full"><div class="binding-hole"></div><div class="binding-coil"></div></div>

<div class="relative h-8 w-full"><div class="binding-hole"></div><div class="binding-coil"></div></div>

<div class="relative h-8 w-full"><div class="binding-hole"></div><div class="binding-coil"></div></div>

<div class="relative h-8 w-full"><div class="binding-hole"></div><div class="binding-coil"></div></div>

</div>

</div>

<div class="flex-1 flex flex-col relative bg-paper-pattern">

<div class="pt-10 px-8 pb-4 border-b-2 border-primary/10 dark:border-white/10">

<div class="flex justify-between items-start mb-2">

<div>

<h1 class="font-display text-4xl text-primary dark:text-white mb-1 tracking-tight">Budget Planner</h1>

<p class="font-handwriting text-2xl text-gray-500 dark:text-gray-400 -mt-1 ml-1 transform -rotate-1">Welcome Back</p>

</div>

<div class="text-right hidden sm:block">

<span class="block text-xs font-mono text-gray-400 dark:text-gray-500 uppercase tracking-widest">Date</span>

<span class="block text-sm font-handwriting text-gray-600 dark:text-gray-300">Today</span>

</div>

</div>

</div>

<div class="flex-1 px-8 py-8 flex flex-col">

<div class="mb-8">

<h2 class="text-sm font-bold uppercase tracking-widest text-gray-800 dark:text-gray-200 mb-1">User Credentials</h2>

<p class="text-xs text-gray-500 dark:text-gray-400 font-light">Please enter your details to access your records</p>

</div>

<form action="#" class="space-y-6">

<div class="relative">

<label class="block text-xs font-bold text-gray-500 dark:text-gray-400 uppercase mb-1 tracking-wider pl-1" for="email">

Email Address

</label>

<div class="flex items-center border-b border-gray-300 dark:border-gray-600 pb-1 focus-within:border-primary dark:focus-within:border-white transition-colors">

<span class="material-icons-outlined text-gray-400 dark:text-gray-500 mr-2 text-xl">mail</span>

<input class="w-full bg-transparent border-none p-2 text-gray-800 dark:text-gray-100 placeholder-gray-300 dark:placeholder-gray-600 focus:ring-0 font-handwriting text-xl" id="email" placeholder="user@example.com" type="email"/>

</div>

<div class="absolute inset-0 pointer-events-none -z-10" style="background-image: repeating-linear-gradient(transparent, transparent 39px, #e5e7eb 40px); background-size: 100% 40px; opacity: 0.3;"></div>

</div>

<div class="relative mt-6">

<label class="block text-xs font-bold text-gray-500 dark:text-gray-400 uppercase mb-1 tracking-wider pl-1" for="password">

Password

</label>

<div class="flex items-center border-b border-gray-300 dark:border-gray-600 pb-1 focus-within:border-primary dark:focus-within:border-white transition-colors">

<span class="material-icons-outlined text-gray-400 dark:text-gray-500 mr-2 text-xl">lock</span>

<input class="w-full bg-transparent border-none p-2 text-gray-800 dark:text-gray-100 placeholder-gray-300 dark:placeholder-gray-600 focus:ring-0 font-handwriting text-xl" id="password" placeholder="••••••••" type="password"/>

</div>

</div>

<div class="flex items-center justify-between pt-2">

<div class="flex items-center">

<input class="h-4 w-4 text-primary focus:ring-gray-500 border-gray-300 rounded" id="remember-me" name="remember-me" type="checkbox"/>

<label class="ml-2 block text-sm text-gray-600 dark:text-gray-400 font-handwriting text-lg" for="remember-me">

Remember me

</label>

</div>

<div class="text-sm">

<a class="font-medium text-gray-600 dark:text-gray-400 hover:text-primary dark:hover:text-white underline decoration-dotted underline-offset-4" href="#">

Forgot password?

</a>

</div>

</div>

<div class="pt-6">

<button class="w-full flex justify-center py-3 px-4 border border-transparent text-sm font-bold uppercase tracking-widest text-white bg-primary hover:bg-gray-800 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-gray-500 shadow-md transform hover:-translate-y-0.5 transition-all duration-200" type="submit">

Access Planner

</button>

</div>

</form>

<div class="mt-auto pt-8 pb-4">

<div class="border-t-2 border-dashed border-gray-300 dark:border-gray-600 pt-4 text-center">

<p class="text-xs text-gray-500 dark:text-gray-400 mb-3">Don't have a planner yet?</p>

<a class="inline-block px-6 py-2 border border-primary dark:border-gray-400 text-primary dark:text-gray-200 font-bold text-xs uppercase tracking-widest hover:bg-gray-100 dark:hover:bg-zinc-800 transition-colors" href="#">

Create New Account

</a>

</div>

</div>

</div>

<div class="absolute top-0 right-0 w-16 h-16 overflow-hidden">

<div class="absolute top-0 right-0 bg-gray-200 dark:bg-zinc-700 w-8 h-8 transform rotate-45 translate-x-4 -translate-y-4 shadow-sm"></div>

</div>

<div class="bg-gray-100 dark:bg-zinc-800/50 py-3 px-8 text-center border-t border-gray-200 dark:border-zinc-700">

<p class="font-display italic text-gray-500 dark:text-gray-400 text-sm">"Classify and summarize expenditures."</p>

</div>

</div>

</div>

  

</body></html>
```

## Sign-up 
```html
<!DOCTYPE html>

<html lang="en"><head>

<meta charset="utf-8"/>

<meta content="width=device-width, initial-scale=1.0" name="viewport"/>

<title>Budget Planner Sign Up</title>

<script src="https://cdn.tailwindcss.com?plugins=forms,container-queries"></script>

<link href="https://fonts.googleapis.com" rel="preconnect"/>

<link crossorigin="" href="https://fonts.gstatic.com" rel="preconnect"/>

<link href="https://fonts.googleapis.com/css2?family=Caveat:wght@500;700&amp;family=Inter:wght@300;400;500;600&amp;family=Playfair+Display:ital,wght@0,600;1,600&amp;display=swap" rel="stylesheet"/>

<link href="https://fonts.googleapis.com/icon?family=Material+Icons+Outlined" rel="stylesheet"/>

<script>

tailwind.config = {

darkMode: "class",

theme: {

extend: {

colors: {

primary: "#333333", // Dark grey/black ink color for primary actions

"paper-light": "#fdfbf7", // Warm creamy paper color

"paper-dark": "#1f1f1f", // Dark mode paper

"line-light": "#e5e7eb", // Light grey lines

"line-dark": "#374151", // Dark mode lines

"accent-gold": "#d4af37", // Gold for spiral binding aesthetic

},

fontFamily: {

display: ["'Playfair Display'", "serif"],

body: ["'Inter'", "sans-serif"],

handwriting: ["'Caveat'", "cursive"],

},

boxShadow: {

'paper': '0 4px 6px -1px rgba(0, 0, 0, 0.1), 0 2px 4px -1px rgba(0, 0, 0, 0.06), 2px 0 10px rgba(0,0,0,0.05)',

},

backgroundImage: {

'paper-pattern': "url('data:image/svg+xml;base64,PHN2ZyB3aWR0aD0iMjAiIGhlaWdodD0iMjAiIHhtbG5zPSJodHRwOi8vd3d3LnczLm9yZy8yMDAwL3N2ZyI+PGNpcmNsZSBjeD0iMSIgY3k9IjEiIHI9IjEiIGZpbGw9IiM5OTkiIGZpbGwtb3BhY2l0eT0iMC4wMyIvPjwvc3ZnPg==')",

}

},

},

};

</script>

<style>.no-scrollbar::-webkit-scrollbar {

display: none;

}

.no-scrollbar {

-ms-overflow-style: none;

scrollbar-width: none;

}.binding-holes {

position: absolute;

left: 12px;

top: 0;

bottom: 0;

width: 20px;

display: flex;

flex-direction: column;

justify-content: space-evenly;

z-index: 10;

}

.binding-hole {

width: 12px;

height: 12px;

background-color: #f3f4f6;border-radius: 50%;

box-shadow: inset 1px 1px 2px rgba(0,0,0,0.2);

position: relative;

}.dark .binding-hole {

background-color: #374151;

box-shadow: inset 1px 1px 3px rgba(0,0,0,0.5);

}

.binding-coil {

position: absolute;

left: -4px;

top: 50%;

transform: translateY(-50%);

width: 24px;

height: 8px;

background: linear-gradient(90deg, #b8860b, #ffd700, #b8860b);

border-radius: 4px;

box-shadow: 1px 1px 2px rgba(0,0,0,0.3);

z-index: 20;

}

</style>

<style>

body {

min-height: max(884px, 100dvh);

}

</style>

<style>

body {

min-height: max(884px, 100dvh);

}

</style>

</head>

<body class="bg-gray-200 dark:bg-gray-900 min-h-screen flex items-center justify-center font-body p-4 sm:p-0">

<div class="relative w-full max-w-md bg-paper-light dark:bg-paper-dark rounded-r-2xl rounded-l-md shadow-2xl overflow-hidden min-h-[750px] flex">

<div class="w-12 bg-gray-100 dark:bg-zinc-800 border-r border-gray-300 dark:border-zinc-700 relative flex-shrink-0">

<div class="binding-holes py-4">

<div class="relative h-8 w-full"><div class="binding-hole"></div><div class="binding-coil"></div></div>

<div class="relative h-8 w-full"><div class="binding-hole"></div><div class="binding-coil"></div></div>

<div class="relative h-8 w-full"><div class="binding-hole"></div><div class="binding-coil"></div></div>

<div class="relative h-8 w-full"><div class="binding-hole"></div><div class="binding-coil"></div></div>

<div class="relative h-8 w-full"><div class="binding-hole"></div><div class="binding-coil"></div></div>

<div class="relative h-8 w-full"><div class="binding-hole"></div><div class="binding-coil"></div></div>

<div class="relative h-8 w-full"><div class="binding-hole"></div><div class="binding-coil"></div></div>

<div class="relative h-8 w-full"><div class="binding-hole"></div><div class="binding-coil"></div></div>

<div class="relative h-8 w-full"><div class="binding-hole"></div><div class="binding-coil"></div></div>

<div class="relative h-8 w-full"><div class="binding-hole"></div><div class="binding-coil"></div></div>

<div class="relative h-8 w-full"><div class="binding-hole"></div><div class="binding-coil"></div></div>

<div class="relative h-8 w-full"><div class="binding-hole"></div><div class="binding-coil"></div></div>

<div class="relative h-8 w-full"><div class="binding-hole"></div><div class="binding-coil"></div></div>

<div class="relative h-8 w-full"><div class="binding-hole"></div><div class="binding-coil"></div></div>

<div class="relative h-8 w-full"><div class="binding-hole"></div><div class="binding-coil"></div></div>

<div class="relative h-8 w-full"><div class="binding-hole"></div><div class="binding-coil"></div></div>

</div>

</div>

<div class="flex-1 flex flex-col relative bg-paper-pattern">

<div class="pt-10 px-8 pb-4 border-b-2 border-primary/10 dark:border-white/10">

<div class="flex justify-between items-start mb-2">

<div>

<h1 class="font-display text-4xl text-primary dark:text-white mb-1 tracking-tight">Budget Planner</h1>

<p class="font-handwriting text-2xl text-gray-500 dark:text-gray-400 -mt-1 ml-1 transform -rotate-1">Join the Planner</p>

</div>

<div class="text-right hidden sm:block">

<span class="block text-xs font-mono text-gray-400 dark:text-gray-500 uppercase tracking-widest">Date</span>

<span class="block text-sm font-handwriting text-gray-600 dark:text-gray-300">Today</span>

</div>

</div>

</div>

<div class="flex-1 px-8 py-8 flex flex-col">

<div class="mb-8">

<h2 class="text-sm font-bold uppercase tracking-widest text-gray-800 dark:text-gray-200 mb-1">New Registration</h2>

<p class="text-xs text-gray-500 dark:text-gray-400 font-light">Please fill in your details to create an account</p>

</div>

<form action="#" class="space-y-6">

<div class="relative">

<label class="block text-xs font-bold text-gray-500 dark:text-gray-400 uppercase mb-1 tracking-wider pl-1" for="email">

Email Address

</label>

<div class="flex items-center border-b border-gray-300 dark:border-gray-600 pb-1 focus-within:border-primary dark:focus-within:border-white transition-colors">

<span class="material-icons-outlined text-gray-400 dark:text-gray-500 mr-2 text-xl">mail</span>

<input class="w-full bg-transparent border-none p-2 text-gray-800 dark:text-gray-100 placeholder-gray-300 dark:placeholder-gray-600 focus:ring-0 font-handwriting text-xl" id="email" placeholder="user@example.com" type="email"/>

</div>

<div class="absolute inset-0 pointer-events-none -z-10" style="background-image: repeating-linear-gradient(transparent, transparent 39px, #e5e7eb 40px); background-size: 100% 40px; opacity: 0.3;"></div>

</div>

<div class="relative mt-6">

<label class="block text-xs font-bold text-gray-500 dark:text-gray-400 uppercase mb-1 tracking-wider pl-1" for="password">

Password

</label>

<div class="flex items-center border-b border-gray-300 dark:border-gray-600 pb-1 focus-within:border-primary dark:focus-within:border-white transition-colors">

<span class="material-icons-outlined text-gray-400 dark:text-gray-500 mr-2 text-xl">lock</span>

<input class="w-full bg-transparent border-none p-2 text-gray-800 dark:text-gray-100 placeholder-gray-300 dark:placeholder-gray-600 focus:ring-0 font-handwriting text-xl" id="password" placeholder="••••••••" type="password"/>

</div>

</div>

<div class="relative mt-6">

<label class="block text-xs font-bold text-gray-500 dark:text-gray-400 uppercase mb-1 tracking-wider pl-1" for="confirm_password">

Confirm Password

</label>

<div class="flex items-center border-b border-gray-300 dark:border-gray-600 pb-1 focus-within:border-primary dark:focus-within:border-white transition-colors">

<span class="material-icons-outlined text-gray-400 dark:text-gray-500 mr-2 text-xl">lock</span>

<input class="w-full bg-transparent border-none p-2 text-gray-800 dark:text-gray-100 placeholder-gray-300 dark:placeholder-gray-600 focus:ring-0 font-handwriting text-xl" id="confirm_password" placeholder="••••••••" type="password"/>

</div>

</div>

<div class="pt-8">

<button class="w-full flex justify-center py-3 px-4 border border-transparent text-sm font-bold uppercase tracking-widest text-white bg-primary hover:bg-gray-800 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-gray-500 shadow-md transform hover:-translate-y-0.5 transition-all duration-200" type="submit">

Create Account

</button>

</div>

</form>

<div class="mt-auto pt-8 pb-4">

<div class="border-t-2 border-dashed border-gray-300 dark:border-gray-600 pt-4 text-center">

<p class="text-xs text-gray-500 dark:text-gray-400 mb-3">Already have a planner?</p>

<a class="inline-block px-6 py-2 border border-primary dark:border-gray-400 text-primary dark:text-gray-200 font-bold text-xs uppercase tracking-widest hover:bg-gray-100 dark:hover:bg-zinc-800 transition-colors" href="#">

Return to Login

</a>

</div>

</div>

</div>

<div class="absolute top-0 right-0 w-16 h-16 overflow-hidden">

<div class="absolute top-0 right-0 bg-gray-200 dark:bg-zinc-700 w-8 h-8 transform rotate-45 translate-x-4 -translate-y-4 shadow-sm"></div>

</div>

<div class="bg-gray-100 dark:bg-zinc-800/50 py-3 px-8 text-center border-t border-gray-200 dark:border-zinc-700">

<p class="font-display italic text-gray-500 dark:text-gray-400 text-sm">"A journey of a thousand miles begins with a single step."</p>

</div>

</div>

</div>

  

</body></html>
```


## Budget Overview
```html
<!DOCTYPE html>
<html lang="en"><head>
<meta charset="utf-8"/>
<meta content="width=device-width, initial-scale=1.0" name="viewport"/>
<title>Budget Overview</title>
<script src="https://cdn.tailwindcss.com?plugins=forms,typography"></script>
<link href="https://fonts.googleapis.com/css2?family=Playfair+Display:ital,wght@0,400;0,600;0,700;1,400&amp;family=Inter:wght@300;400;500;600&amp;display=swap" rel="stylesheet"/>
<link href="https://fonts.googleapis.com/icon?family=Material+Icons" rel="stylesheet"/>
<script>
        tailwind.config = {
            darkMode: "class",
            theme: {
                extend: {
                    colors: {
                        primary: "#374151", // Dark gray for primary text/headers
                        "background-light": "#F9FAFB", // Very light gray paper feel
                        "background-dark": "#1F2937", // Dark mode background
                        "paper-light": "#FFFFFF",
                        "paper-dark": "#111827",
                        "accent-blue": "#60A5FA",
                        "accent-purple": "#A78BFA",
                        "accent-green": "#34D399",
                        "accent-yellow": "#FBBF24",
                        "table-header-light": "#374151",
                        "table-header-dark": "#4B5563",
                    },
                    fontFamily: {
                        display: ["'Playfair Display'", "serif"],
                        body: ["'Inter'", "sans-serif"],
                    },
                    borderRadius: {
                        DEFAULT: "0.5rem",
                        "xl": "1rem",
                    },
                    boxShadow: {
                        'paper': '0 4px 6px -1px rgba(0, 0, 0, 0.05), 0 2px 4px -1px rgba(0, 0, 0, 0.03)',
                    }
                },
            },
        };
    </script>
<style>
        .notebook-lines {
            background-image: linear-gradient(#e5e7eb 1px, transparent 1px);
            background-size: 100% 2rem;
        }
        .dark .notebook-lines {
            background-image: linear-gradient(#374151 1px, transparent 1px);
        }
    </style>
<style>
    body {
      min-height: max(884px, 100dvh);
    }
  </style>
  </head>
<body class="bg-background-light dark:bg-background-dark text-gray-800 dark:text-gray-200 font-body min-h-screen pb-10 antialiased selection:bg-accent-blue selection:text-white">
<header class="sticky top-0 z-50 bg-paper-light/90 dark:bg-paper-dark/90 backdrop-blur-sm border-b border-gray-200 dark:border-gray-700 px-4 py-4 flex justify-between items-center shadow-sm">
<div class="flex items-center space-x-2">
<span class="material-icons text-primary dark:text-gray-300">menu_book</span>
<h1 class="text-xl font-display font-bold text-primary dark:text-gray-100">Budget Planner</h1>
</div>
<div class="flex items-center space-x-3">
<span class="text-sm font-medium text-gray-500 dark:text-gray-400">Jan 2024</span>
<button class="p-2 rounded-full hover:bg-gray-100 dark:hover:bg-gray-800 transition-colors">
<span class="material-icons text-primary dark:text-gray-300">calendar_today</span>
</button>
</div>
</header>
<main class="max-w-md mx-auto px-4 mt-6 space-y-6">
<section class="bg-paper-light dark:bg-paper-dark rounded-xl shadow-paper border border-gray-100 dark:border-gray-800 overflow-hidden">
<div class="bg-table-header-light dark:bg-table-header-dark text-white p-3 flex justify-between items-center">
<h2 class="font-display font-semibold tracking-wide uppercase text-sm">Budget Review</h2>
<span class="material-icons text-sm opacity-80">analytics</span>
</div>
<div class="p-4 grid grid-cols-2 gap-3 text-center">
<div class="border border-gray-200 dark:border-gray-700 p-2 rounded-lg bg-gray-50 dark:bg-gray-900/50">
<p class="text-xs text-gray-500 dark:text-gray-400 uppercase tracking-wider mb-1">Opening</p>
<p class="font-bold text-gray-800 dark:text-gray-100 text-lg">$3,000</p>
</div>
<div class="border border-gray-200 dark:border-gray-700 p-2 rounded-lg bg-gray-50 dark:bg-gray-900/50">
<p class="text-xs text-gray-500 dark:text-gray-400 uppercase tracking-wider mb-1">Total Income</p>
<p class="font-bold text-green-600 dark:text-green-400 text-lg">$12,000</p>
</div>
<div class="border border-gray-200 dark:border-gray-700 p-2 rounded-lg bg-gray-50 dark:bg-gray-900/50">
<p class="text-xs text-gray-500 dark:text-gray-400 uppercase tracking-wider mb-1">Total Expenses</p>
<p class="font-bold text-red-500 dark:text-red-400 text-lg">$8,000</p>
</div>
<div class="border-2 border-accent-yellow p-2 rounded-lg bg-yellow-50 dark:bg-yellow-900/20">
<p class="text-xs text-yellow-700 dark:text-yellow-400 uppercase tracking-wider mb-1 font-semibold">Total Savings</p>
<p class="font-bold text-yellow-800 dark:text-yellow-300 text-lg">$7,000</p>
</div>
</div>
<div class="px-4 pb-4">
<div class="flex justify-between items-center text-sm border-t border-gray-200 dark:border-gray-700 pt-3">
<span class="text-gray-500 dark:text-gray-400">Balance Forward</span>
<span class="font-bold text-primary dark:text-white">$4,000</span>
</div>
</div>
</section>
<section class="bg-paper-light dark:bg-paper-dark rounded-xl shadow-paper border border-gray-100 dark:border-gray-800 overflow-hidden">
<div class="bg-table-header-light dark:bg-table-header-dark text-white p-3 flex justify-between items-center">
<h2 class="font-display font-semibold tracking-wide uppercase text-sm">Spending</h2>
<span class="material-icons text-sm opacity-80">pie_chart</span>
</div>
<div class="p-4">
<div class="flex justify-center mb-6 relative">
<div class="w-40 h-40 rounded-full border-4 border-white dark:border-gray-800 shadow-lg relative overflow-hidden" style="background: conic-gradient(
                             #A78BFA 0% 27.5%, 
                             #60A5FA 27.5% 40%, 
                             #34D399 40% 48%,
                             #FBBF24 48% 54.25%,
                             #F87171 54.25% 73%,
                             #9CA3AF 73% 100%
                         );">
</div>
<div class="absolute inset-0 flex items-center justify-center">
<div class="bg-white dark:bg-gray-800 rounded-full w-16 h-16 flex items-center justify-center shadow-sm">
<span class="text-[10px] font-bold text-gray-500 dark:text-gray-400 text-center leading-tight">TOTAL<br/>SPEND</span>
</div>
</div>
</div>
<div class="overflow-x-auto">
<table class="w-full text-sm text-left">
<thead class="text-xs text-gray-500 dark:text-gray-400 uppercase bg-gray-50 dark:bg-gray-800 border-b dark:border-gray-700">
<tr>
<th class="px-3 py-2 font-medium">Category</th>
<th class="px-3 py-2 font-medium text-right">Amount</th>
<th class="px-3 py-2 font-medium text-right w-12">%</th>
</tr>
</thead>
<tbody class="divide-y divide-gray-100 dark:divide-gray-800">
<tr>
<td class="px-3 py-2 flex items-center gap-2">
<span class="w-2 h-2 rounded-full bg-accent-purple"></span>
<span class="text-gray-700 dark:text-gray-300">Housing</span>
</td>
<td class="px-3 py-2 text-right font-medium text-gray-900 dark:text-gray-100">2,200</td>
<td class="px-3 py-2 text-right text-gray-500 dark:text-gray-400">27.5</td>
</tr>
<tr>
<td class="px-3 py-2 flex items-center gap-2">
<span class="w-2 h-2 rounded-full bg-accent-blue"></span>
<span class="text-gray-700 dark:text-gray-300">Food</span>
</td>
<td class="px-3 py-2 text-right font-medium text-gray-900 dark:text-gray-100">1,000</td>
<td class="px-3 py-2 text-right text-gray-500 dark:text-gray-400">12.5</td>
</tr>
<tr>
<td class="px-3 py-2 flex items-center gap-2">
<span class="w-2 h-2 rounded-full bg-accent-green"></span>
<span class="text-gray-700 dark:text-gray-300">Health Care</span>
</td>
<td class="px-3 py-2 text-right font-medium text-gray-900 dark:text-gray-100">800</td>
<td class="px-3 py-2 text-right text-gray-500 dark:text-gray-400">8</td>
</tr>
<tr>
<td class="px-3 py-2 flex items-center gap-2">
<span class="w-2 h-2 rounded-full bg-accent-yellow"></span>
<span class="text-gray-700 dark:text-gray-300">Transportation</span>
</td>
<td class="px-3 py-2 text-right font-medium text-gray-900 dark:text-gray-100">500</td>
<td class="px-3 py-2 text-right text-gray-500 dark:text-gray-400">6.25</td>
</tr>
<tr>
<td class="px-3 py-2 flex items-center gap-2">
<span class="w-2 h-2 rounded-full bg-red-400"></span>
<span class="text-gray-700 dark:text-gray-300">Personal</span>
</td>
<td class="px-3 py-2 text-right font-medium text-gray-900 dark:text-gray-100">1,500</td>
<td class="px-3 py-2 text-right text-gray-500 dark:text-gray-400">18.75</td>
</tr>
<tr>
<td class="px-3 py-2 flex items-center gap-2">
<span class="w-2 h-2 rounded-full bg-gray-400"></span>
<span class="text-gray-700 dark:text-gray-300">Entertainment</span>
</td>
<td class="px-3 py-2 text-right font-medium text-gray-900 dark:text-gray-100">2,000</td>
<td class="px-3 py-2 text-right text-gray-500 dark:text-gray-400">25</td>
</tr>
</tbody>
</table>
</div>
</div>
</section>
<section class="relative bg-paper-light dark:bg-paper-dark rounded-xl shadow-paper border border-gray-100 dark:border-gray-800 p-5 overflow-hidden">
<div class="absolute top-0 left-0 bottom-0 w-6 flex flex-col justify-evenly py-4 pl-1 pointer-events-none">
<div class="w-3 h-3 rounded-full bg-gray-300 dark:bg-gray-600 shadow-inner mb-2"></div>
<div class="w-3 h-3 rounded-full bg-gray-300 dark:bg-gray-600 shadow-inner mb-2"></div>
<div class="w-3 h-3 rounded-full bg-gray-300 dark:bg-gray-600 shadow-inner mb-2"></div>
<div class="w-3 h-3 rounded-full bg-gray-300 dark:bg-gray-600 shadow-inner mb-2"></div>
<div class="w-3 h-3 rounded-full bg-gray-300 dark:bg-gray-600 shadow-inner mb-2"></div>
<div class="w-3 h-3 rounded-full bg-gray-300 dark:bg-gray-600 shadow-inner mb-2"></div>
</div>
<div class="pl-6">
<h3 class="font-display text-lg font-bold text-primary dark:text-gray-100 mb-4 border-b-2 border-primary dark:border-gray-500 inline-block">Monthly Reflection</h3>
<div class="mb-5">
<label class="block text-xs font-bold uppercase tracking-wide text-white bg-primary dark:bg-gray-600 py-1 px-2 mb-1 w-full rounded-sm">
                        My biggest wins this month
                    </label>
<div class="notebook-lines min-h-[4rem] text-sm text-blue-600 dark:text-blue-300 font-handwriting leading-8 pl-1">
                        Save 3000 more than last month.
                    </div>
</div>
<div class="mb-5">
<label class="block text-xs font-bold uppercase tracking-wide text-white bg-primary dark:bg-gray-600 py-1 px-2 mb-1 w-full rounded-sm">
                        Did I meet my budget? If not, why not?
                    </label>
<div class="notebook-lines min-h-[4rem] text-sm text-gray-700 dark:text-gray-300 font-handwriting leading-8 pl-1">
                        Yes, mostly. Food budget went slightly over due to the birthday party.
                    </div>
</div>
<div class="mb-2">
<label class="block text-xs font-bold uppercase tracking-wide text-white bg-primary dark:bg-gray-600 py-1 px-2 mb-1 w-full rounded-sm">
                        I will do this within 1 month to improve
                    </label>
<div class="notebook-lines min-h-[4rem] text-sm text-gray-700 dark:text-gray-300 font-handwriting leading-8 pl-1">
                        Cook at home at least 5 days a week.
                    </div>
</div>
</div>
</section>
<button class="fixed bottom-6 right-6 bg-primary text-white p-4 rounded-full shadow-lg hover:bg-gray-800 transition-transform hover:scale-105 active:scale-95 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-primary dark:focus:ring-offset-gray-900 z-50">
<span class="material-icons text-2xl">add</span>
</button>
</main>
<nav class="fixed bottom-0 w-full bg-paper-light dark:bg-paper-dark border-t border-gray-200 dark:border-gray-700 flex justify-around py-3 z-40 pb-safe">
<a class="flex flex-col items-center text-primary dark:text-white" href="#">
<span class="material-icons">dashboard</span>
<span class="text-[10px] mt-1 font-medium">Overview</span>
</a>
<a class="flex flex-col items-center text-gray-400 dark:text-gray-500 hover:text-primary dark:hover:text-gray-300 transition-colors" href="#">
<span class="material-icons">receipt_long</span>
<span class="text-[10px] mt-1 font-medium">Transactions</span>
</a>
<a class="flex flex-col items-center text-gray-400 dark:text-gray-500 hover:text-primary dark:hover:text-gray-300 transition-colors" href="#">
<span class="material-icons">account_balance_wallet</span>
<span class="text-[10px] mt-1 font-medium">Bills</span>
</a>
<a class="flex flex-col items-center text-gray-400 dark:text-gray-500 hover:text-primary dark:hover:text-gray-300 transition-colors" href="#">
<span class="material-icons">settings</span>
<span class="text-[10px] mt-1 font-medium">Settings</span>
</a>
</nav>
<div class="h-20"></div>

</body></html>
```


## Expense Entry and Tracker

```html
<!DOCTYPE html>
<html lang="en"><head>
<meta charset="utf-8"/>
<meta content="width=device-width, initial-scale=1.0" name="viewport"/>
<title>Expense Entry &amp; Tracker</title>
<script src="https://cdn.tailwindcss.com?plugins=forms,typography"></script>
<link href="https://fonts.googleapis.com/css2?family=Playfair+Display:ital,wght@0,400;0,600;0,700;1,400&amp;family=Inter:wght@300;400;500;600&amp;family=Kalam:wght@300;400&amp;display=swap" rel="stylesheet"/>
<link href="https://fonts.googleapis.com/icon?family=Material+Icons+Outlined" rel="stylesheet"/>
<script>
      tailwind.config = {
        darkMode: "class",
        theme: {
          extend: {
            colors: {
              primary: "#374151", // Dark Slate/Charcoal for a classic ink look
              "background-light": "#F9F9F7", // Paper-like off-white
              "background-dark": "#1F2937", // Dark gray for dark mode
              "paper-light": "#FFFFFF",
              "paper-dark": "#374151",
              "line-light": "#E5E7EB",
              "line-dark": "#4B5563",
              "accent-highlight": "#FEF3C7", // Highlighter yellow/cream
            },
            fontFamily: {
              display: ["'Playfair Display'", "serif"],
              body: ["'Inter'", "sans-serif"],
              hand: ["'Kalam'", "cursive"],
            },
            borderRadius: {
              DEFAULT: "12px",
            },
            boxShadow: {
              'paper': '0 4px 6px -1px rgba(0, 0, 0, 0.05), 0 2px 4px -1px rgba(0, 0, 0, 0.03)',
            }
          },
        },
      };
    </script>
<style>.no-scrollbar::-webkit-scrollbar {
            display: none;
        }
        .no-scrollbar {
            -ms-overflow-style: none;
            scrollbar-width: none;
        }
        .notebook-lines {
            background-image: linear-gradient(#E5E7EB 1px, transparent 1px);
            background-size: 100% 3rem;
        }
        .dark .notebook-lines {
            background-image: linear-gradient(#4B5563 1px, transparent 1px);
        }
    </style>
<style>
    body {
      min-height: max(884px, 100dvh);
    }
  </style>
  </head>
<body class="bg-background-light dark:bg-background-dark text-gray-800 dark:text-gray-100 font-body transition-colors duration-300 antialiased min-h-screen pb-20">
<header class="pt-10 pb-4 px-6 sticky top-0 bg-background-light/95 dark:bg-background-dark/95 backdrop-blur-sm z-20 border-b border-gray-200 dark:border-gray-700">
<div class="flex justify-between items-start mb-4">
<div>
<h2 class="text-sm font-medium text-gray-500 dark:text-gray-400 uppercase tracking-widest">January 2024</h2>
<h1 class="text-3xl font-display font-bold text-gray-900 dark:text-white mt-1">Expense Tracker</h1>
</div>
<button class="p-2 rounded-full hover:bg-gray-200 dark:hover:bg-gray-700 transition-colors">
<span class="material-icons-outlined text-gray-600 dark:text-gray-300">calendar_today</span>
</button>
</div>
<div class="grid grid-cols-2 gap-3 mt-4">
<div class="bg-paper-light dark:bg-paper-dark p-3 rounded-DEFAULT shadow-sm border border-gray-100 dark:border-gray-600">
<p class="text-xs text-gray-500 dark:text-gray-400 mb-1">Total Spent</p>
<p class="text-xl font-bold font-display text-gray-900 dark:text-white">$1,245.50</p>
</div>
<div class="bg-paper-light dark:bg-paper-dark p-3 rounded-DEFAULT shadow-sm border border-gray-100 dark:border-gray-600">
<p class="text-xs text-gray-500 dark:text-gray-400 mb-1">Budget Left</p>
<p class="text-xl font-bold font-display text-emerald-600 dark:text-emerald-400">$754.50</p>
</div>
</div>
</header>
<main class="px-4 mt-2">
<div class="flex justify-between items-center py-4">
<div class="flex space-x-2 overflow-x-auto no-scrollbar pb-2">
<span class="px-3 py-1 bg-primary text-white text-xs rounded-full whitespace-nowrap">All</span>
<span class="px-3 py-1 bg-gray-200 dark:bg-gray-700 text-gray-600 dark:text-gray-300 text-xs rounded-full whitespace-nowrap">Food</span>
<span class="px-3 py-1 bg-gray-200 dark:bg-gray-700 text-gray-600 dark:text-gray-300 text-xs rounded-full whitespace-nowrap">Transport</span>
<span class="px-3 py-1 bg-gray-200 dark:bg-gray-700 text-gray-600 dark:text-gray-300 text-xs rounded-full whitespace-nowrap">Bills</span>
</div>
<button class="text-primary dark:text-white">
<span class="material-icons-outlined">filter_list</span>
</button>
</div>
<div class="bg-paper-light dark:bg-paper-dark rounded-DEFAULT shadow-paper overflow-hidden relative border border-gray-200 dark:border-gray-600">
<div class="absolute left-0 top-0 bottom-0 w-6 flex flex-col items-center justify-start py-4 gap-3 bg-gray-100 dark:bg-gray-800 border-r border-gray-300 dark:border-gray-600 z-10">
<div class="w-2 h-4 rounded-full bg-gradient-to-b from-gray-300 to-gray-500 shadow-sm"></div>
<div class="w-2 h-4 rounded-full bg-gradient-to-b from-gray-300 to-gray-500 shadow-sm"></div>
<div class="w-2 h-4 rounded-full bg-gradient-to-b from-gray-300 to-gray-500 shadow-sm"></div>
<div class="w-2 h-4 rounded-full bg-gradient-to-b from-gray-300 to-gray-500 shadow-sm"></div>
<div class="w-2 h-4 rounded-full bg-gradient-to-b from-gray-300 to-gray-500 shadow-sm"></div>
<div class="w-2 h-4 rounded-full bg-gradient-to-b from-gray-300 to-gray-500 shadow-sm"></div>
<div class="w-2 h-4 rounded-full bg-gradient-to-b from-gray-300 to-gray-500 shadow-sm"></div>
<div class="w-2 h-4 rounded-full bg-gradient-to-b from-gray-300 to-gray-500 shadow-sm"></div>
<div class="w-2 h-4 rounded-full bg-gradient-to-b from-gray-300 to-gray-500 shadow-sm"></div>
<div class="w-2 h-4 rounded-full bg-gradient-to-b from-gray-300 to-gray-500 shadow-sm"></div>
<div class="w-2 h-4 rounded-full bg-gradient-to-b from-gray-300 to-gray-500 shadow-sm"></div>
<div class="w-2 h-4 rounded-full bg-gradient-to-b from-gray-300 to-gray-500 shadow-sm"></div>
</div>
<div class="pl-8 overflow-x-auto">
<table class="w-full min-w-[500px] text-left border-collapse">
<thead>
<tr class="border-b-2 border-primary">
<th class="py-3 pl-4 pr-2 text-xs font-bold uppercase tracking-wider text-gray-500 dark:text-gray-400 w-16">Date</th>
<th class="py-3 px-2 text-xs font-bold uppercase tracking-wider text-gray-500 dark:text-gray-400 w-32">Description</th>
<th class="py-3 px-2 text-xs font-bold uppercase tracking-wider text-gray-500 dark:text-gray-400 w-24">Category</th>
<th class="py-3 px-2 text-xs font-bold uppercase tracking-wider text-gray-500 dark:text-gray-400 w-20 text-right">Amt</th>
<th class="py-3 px-2 text-xs font-bold uppercase tracking-wider text-gray-500 dark:text-gray-400 w-10 text-center">Paid</th>
</tr>
</thead>
<tbody class="text-sm font-hand text-lg">
<tr class="border-b border-gray-200 dark:border-gray-700 hover:bg-gray-50 dark:hover:bg-gray-700/50 transition-colors h-12">
<td class="pl-4 py-2 text-gray-600 dark:text-gray-300">01/15</td>
<td class="px-2 py-2 font-medium text-gray-800 dark:text-white">Grocery Run</td>
<td class="px-2 py-2">
<span class="bg-orange-100 dark:bg-orange-900 text-orange-800 dark:text-orange-200 text-xs font-sans px-2 py-0.5 rounded-full">Food</span>
</td>
<td class="px-2 py-2 text-right font-bold text-gray-800 dark:text-white">$120.50</td>
<td class="px-2 py-2 text-center text-emerald-500">
<span class="material-icons-outlined text-sm">check_circle</span>
</td>
</tr>
<tr class="border-b border-gray-200 dark:border-gray-700 hover:bg-gray-50 dark:hover:bg-gray-700/50 transition-colors h-12 bg-accent-highlight/30 dark:bg-yellow-900/10">
<td class="pl-4 py-2 text-gray-600 dark:text-gray-300">01/16</td>
<td class="px-2 py-2 font-medium text-gray-800 dark:text-white">Electricity Bill</td>
<td class="px-2 py-2">
<span class="bg-blue-100 dark:bg-blue-900 text-blue-800 dark:text-blue-200 text-xs font-sans px-2 py-0.5 rounded-full">Bills</span>
</td>
<td class="px-2 py-2 text-right font-bold text-gray-800 dark:text-white">$85.00</td>
<td class="px-2 py-2 text-center text-emerald-500">
<span class="material-icons-outlined text-sm">check_circle</span>
</td>
</tr>
<tr class="border-b border-gray-200 dark:border-gray-700 hover:bg-gray-50 dark:hover:bg-gray-700/50 transition-colors h-12">
<td class="pl-4 py-2 text-gray-600 dark:text-gray-300">01/18</td>
<td class="px-2 py-2 font-medium text-gray-800 dark:text-white">Uber Ride</td>
<td class="px-2 py-2">
<span class="bg-purple-100 dark:bg-purple-900 text-purple-800 dark:text-purple-200 text-xs font-sans px-2 py-0.5 rounded-full">Transport</span>
</td>
<td class="px-2 py-2 text-right font-bold text-gray-800 dark:text-white">$24.00</td>
<td class="px-2 py-2 text-center text-gray-300 dark:text-gray-600">
<span class="material-icons-outlined text-sm">radio_button_unchecked</span>
</td>
</tr>
<tr class="border-b border-gray-200 dark:border-gray-700 hover:bg-gray-50 dark:hover:bg-gray-700/50 transition-colors h-12">
<td class="pl-4 py-2 text-gray-600 dark:text-gray-300">01/20</td>
<td class="px-2 py-2 font-medium text-gray-800 dark:text-white">Coffee Date</td>
<td class="px-2 py-2">
<span class="bg-pink-100 dark:bg-pink-900 text-pink-800 dark:text-pink-200 text-xs font-sans px-2 py-0.5 rounded-full">Social</span>
</td>
<td class="px-2 py-2 text-right font-bold text-gray-800 dark:text-white">$15.75</td>
<td class="px-2 py-2 text-center text-emerald-500">
<span class="material-icons-outlined text-sm">check_circle</span>
</td>
</tr>
<tr class="border-b border-gray-200 dark:border-gray-700 hover:bg-gray-50 dark:hover:bg-gray-700/50 transition-colors h-12">
<td class="pl-4 py-2 text-gray-600 dark:text-gray-300">01/22</td>
<td class="px-2 py-2 font-medium text-gray-800 dark:text-white">Internet Bill</td>
<td class="px-2 py-2">
<span class="bg-blue-100 dark:bg-blue-900 text-blue-800 dark:text-blue-200 text-xs font-sans px-2 py-0.5 rounded-full">Bills</span>
</td>
<td class="px-2 py-2 text-right font-bold text-gray-800 dark:text-white">$60.00</td>
<td class="px-2 py-2 text-center text-gray-300 dark:text-gray-600">
<span class="material-icons-outlined text-sm">radio_button_unchecked</span>
</td>
</tr>
<tr class="border-b border-gray-200 dark:border-gray-700 hover:bg-gray-50 dark:hover:bg-gray-700/50 transition-colors h-12">
<td class="pl-4 py-2 text-gray-600 dark:text-gray-300">01/24</td>
<td class="px-2 py-2 font-medium text-gray-800 dark:text-white">Gym Membership</td>
<td class="px-2 py-2">
<span class="bg-green-100 dark:bg-green-900 text-green-800 dark:text-green-200 text-xs font-sans px-2 py-0.5 rounded-full">Health</span>
</td>
<td class="px-2 py-2 text-right font-bold text-gray-800 dark:text-white">$45.00</td>
<td class="px-2 py-2 text-center text-emerald-500">
<span class="material-icons-outlined text-sm">check_circle</span>
</td>
</tr>
<tr class="border-b border-gray-200 dark:border-gray-700 h-12">
<td class="pl-4 py-2"></td>
<td class="px-2 py-2"></td>
<td class="px-2 py-2"></td>
<td class="px-2 py-2"></td>
<td class="px-2 py-2"></td>
</tr>
<tr class="border-b border-gray-200 dark:border-gray-700 h-12">
<td class="pl-4 py-2"></td>
<td class="px-2 py-2"></td>
<td class="px-2 py-2"></td>
<td class="px-2 py-2"></td>
<td class="px-2 py-2"></td>
</tr>
</tbody>
</table>
</div>
<div class="h-2 bg-gray-100 dark:bg-gray-800 w-full border-t border-gray-200 dark:border-gray-700"></div>
</div>
<button class="fixed bottom-24 right-6 bg-primary text-white p-4 rounded-full shadow-lg hover:shadow-xl hover:bg-gray-700 transition-all transform hover:-translate-y-1 z-30">
<span class="material-icons-outlined text-2xl">add</span>
</button>
</main>
<nav class="fixed bottom-0 left-0 right-0 bg-white dark:bg-gray-900 border-t border-gray-200 dark:border-gray-700 pb-6 pt-3 px-8 z-40">
<div class="flex justify-between items-center">
<a class="flex flex-col items-center text-gray-400 hover:text-primary dark:hover:text-white transition-colors" href="#">
<span class="material-icons-outlined">home</span>
<span class="text-[10px] mt-1">Home</span>
</a>
<a class="flex flex-col items-center text-primary dark:text-white" href="#">
<span class="material-icons-outlined">edit_note</span>
<span class="text-[10px] mt-1 font-medium">Tracker</span>
</a>
<a class="flex flex-col items-center text-gray-400 hover:text-primary dark:hover:text-white transition-colors" href="#">
<span class="material-icons-outlined">pie_chart</span>
<span class="text-[10px] mt-1">Budget</span>
</a>
<a class="flex flex-col items-center text-gray-400 hover:text-primary dark:hover:text-white transition-colors" href="#">
<span class="material-icons-outlined">person</span>
<span class="text-[10px] mt-1">Profile</span>
</a>
</div>
</nav>
<script>
        // Simple script to toggle dark mode for demonstration
        // In a real app, this would be handled by system preference or user setting
        if (window.matchMedia && window.matchMedia('(prefers-color-scheme: dark)').matches) {
            document.documentElement.classList.add('dark');
        }
    </script>

</body></html>
```


## Expense Entry and Tracker (Showing modal for expense entry)
```html
<!DOCTYPE html>
<html lang="en"><head>
<meta charset="utf-8"/>
<meta content="width=device-width, initial-scale=1.0" name="viewport"/>
<title>Expense Entry &amp; Tracker</title>
<script src="https://cdn.tailwindcss.com?plugins=forms,typography"></script>
<link href="https://fonts.googleapis.com/css2?family=Playfair+Display:ital,wght@0,400;0,600;0,700;1,400&amp;family=Inter:wght@300;400;500;600&amp;family=Kalam:wght@300;400&amp;display=swap" rel="stylesheet"/>
<link href="https://fonts.googleapis.com/icon?family=Material+Icons+Outlined" rel="stylesheet"/>
<script>
      tailwind.config = {
        darkMode: "class",
        theme: {
          extend: {
            colors: {
              primary: "#374151", // Dark Slate/Charcoal for a classic ink look
              "background-light": "#F9F9F7", // Paper-like off-white
              "background-dark": "#1F2937", // Dark gray for dark mode
              "paper-light": "#FFFFFF",
              "paper-dark": "#374151",
              "line-light": "#E5E7EB",
              "line-dark": "#4B5563",
              "accent-highlight": "#FEF3C7", // Highlighter yellow/cream
            },
            fontFamily: {
              display: ["'Playfair Display'", "serif"],
              body: ["'Inter'", "sans-serif"],
              hand: ["'Kalam'", "cursive"],
            },
            borderRadius: {
              DEFAULT: "12px",
            },
            boxShadow: {
              'paper': '0 4px 6px -1px rgba(0, 0, 0, 0.05), 0 2px 4px -1px rgba(0, 0, 0, 0.03)',
            }
          },
        },
      };
    </script>
<style>.no-scrollbar::-webkit-scrollbar {
            display: none;
        }
        .no-scrollbar {
            -ms-overflow-style: none;
            scrollbar-width: none;
        }
        .notebook-lines {
            background-image: linear-gradient(#E5E7EB 1px, transparent 1px);
            background-size: 100% 2rem;
            line-height: 2rem;
            background-position-y: 1.9rem;
        }
        .dark .notebook-lines {
            background-image: linear-gradient(#4B5563 1px, transparent 1px);
        }
    </style>
<style>
        body {
          min-height: max(884px, 100dvh);
        }
    </style>
<style>
    body {
      min-height: max(884px, 100dvh);
    }
  </style>
  </head>
<body class="bg-background-light dark:bg-background-dark text-gray-800 dark:text-gray-100 font-body transition-colors duration-300 antialiased min-h-screen pb-20">
<div class="fixed inset-0 z-50 flex items-end sm:items-center justify-center px-4 pb-4 pt-12">
<div class="absolute inset-0 bg-gray-900/40 backdrop-blur-sm transition-opacity"></div>
<div class="relative w-full max-w-lg transform overflow-hidden rounded-2xl bg-paper-light dark:bg-paper-dark text-left shadow-xl transition-all border border-gray-200 dark:border-gray-600 flex flex-col max-h-[90vh]">
<div class="px-6 py-4 border-b border-gray-100 dark:border-gray-600 flex justify-between items-center bg-background-light/50 dark:bg-gray-800/50">
<div>
<h3 class="text-xl font-display font-bold text-gray-900 dark:text-white">Add New Expense</h3>
<p class="text-xs text-gray-500 dark:text-gray-400 mt-0.5">Record a new transaction in your journal</p>
</div>
<button class="text-gray-400 hover:text-gray-500 dark:text-gray-500 dark:hover:text-gray-300 transition-colors">
<span class="material-icons-outlined">close</span>
</button>
</div>
<div class="p-6 overflow-y-auto custom-scrollbar space-y-6">
<div class="grid grid-cols-2 gap-6">
<div class="relative">
<label class="block text-xs font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider mb-1" for="amount">Amount</label>
<div class="relative rounded-md shadow-sm">
<div class="pointer-events-none absolute inset-y-0 left-0 flex items-center pl-3">
<span class="text-gray-500 sm:text-sm font-display font-bold">$</span>
</div>
<input class="block w-full rounded-md border-0 py-2.5 pl-7 pr-4 text-gray-900 ring-1 ring-inset ring-gray-300 placeholder:text-gray-400 focus:ring-2 focus:ring-inset focus:ring-primary dark:bg-gray-800 dark:ring-gray-600 dark:text-white sm:text-lg font-bold font-display" id="amount" name="amount" placeholder="0.00" type="number"/>
</div>
</div>
<div>
<label class="block text-xs font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider mb-1" for="date">Date</label>
<input class="block w-full rounded-md border-0 py-2.5 text-gray-900 ring-1 ring-inset ring-gray-300 placeholder:text-gray-400 focus:ring-2 focus:ring-inset focus:ring-primary dark:bg-gray-800 dark:ring-gray-600 dark:text-white sm:text-sm" id="date" name="date" type="date"/>
</div>
</div>
<div>
<label class="block text-xs font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider mb-1" for="description">Description</label>
<input class="block w-full border-0 border-b-2 border-gray-200 dark:border-gray-600 bg-transparent py-2 px-0 text-gray-900 dark:text-white placeholder:text-gray-400 focus:border-primary focus:ring-0 sm:text-base font-hand" id="description" name="description" placeholder="e.g. Coffee at Sarah's" type="text"/>
</div>
<div>
<label class="block text-xs font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider mb-3">Category</label>
<div class="flex flex-wrap gap-2">
<label class="cursor-pointer">
<input checked="" class="peer sr-only" name="category" type="radio" value="food"/>
<span class="inline-flex items-center rounded-full px-3 py-1 text-xs font-medium ring-1 ring-inset ring-gray-200 dark:ring-gray-600 text-gray-600 dark:text-gray-300 peer-checked:bg-orange-100 peer-checked:text-orange-800 peer-checked:ring-orange-200 dark:peer-checked:bg-orange-900/50 dark:peer-checked:text-orange-200 dark:peer-checked:ring-orange-800 transition-all">Food</span>
</label>
<label class="cursor-pointer">
<input class="peer sr-only" name="category" type="radio" value="transport"/>
<span class="inline-flex items-center rounded-full px-3 py-1 text-xs font-medium ring-1 ring-inset ring-gray-200 dark:ring-gray-600 text-gray-600 dark:text-gray-300 peer-checked:bg-purple-100 peer-checked:text-purple-800 peer-checked:ring-purple-200 dark:peer-checked:bg-purple-900/50 dark:peer-checked:text-purple-200 dark:peer-checked:ring-purple-800 transition-all">Transport</span>
</label>
<label class="cursor-pointer">
<input class="peer sr-only" name="category" type="radio" value="bills"/>
<span class="inline-flex items-center rounded-full px-3 py-1 text-xs font-medium ring-1 ring-inset ring-gray-200 dark:ring-gray-600 text-gray-600 dark:text-gray-300 peer-checked:bg-blue-100 peer-checked:text-blue-800 peer-checked:ring-blue-200 dark:peer-checked:bg-blue-900/50 dark:peer-checked:text-blue-200 dark:peer-checked:ring-blue-800 transition-all">Bills</span>
</label>
<label class="cursor-pointer">
<input class="peer sr-only" name="category" type="radio" value="social"/>
<span class="inline-flex items-center rounded-full px-3 py-1 text-xs font-medium ring-1 ring-inset ring-gray-200 dark:ring-gray-600 text-gray-600 dark:text-gray-300 peer-checked:bg-pink-100 peer-checked:text-pink-800 peer-checked:ring-pink-200 dark:peer-checked:bg-pink-900/50 dark:peer-checked:text-pink-200 dark:peer-checked:ring-pink-800 transition-all">Social</span>
</label>
<label class="cursor-pointer">
<input class="peer sr-only" name="category" type="radio" value="health"/>
<span class="inline-flex items-center rounded-full px-3 py-1 text-xs font-medium ring-1 ring-inset ring-gray-200 dark:ring-gray-600 text-gray-600 dark:text-gray-300 peer-checked:bg-green-100 peer-checked:text-green-800 peer-checked:ring-green-200 dark:peer-checked:bg-green-900/50 dark:peer-checked:text-green-200 dark:peer-checked:ring-green-800 transition-all">Health</span>
</label>
</div>
</div>
<div>
<label class="block text-xs font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider mb-1" for="notes">Notes</label>
<div class="relative bg-white dark:bg-gray-800 border border-gray-200 dark:border-gray-600 rounded-md overflow-hidden">
<textarea class="block w-full border-0 bg-transparent text-gray-900 dark:text-white placeholder:text-gray-400 focus:ring-0 sm:text-sm notebook-lines font-hand pl-3" id="notes" name="notes" placeholder="Add extra details here..." rows="3"></textarea>
</div>
</div>
</div>
<div class="px-6 py-4 bg-gray-50 dark:bg-gray-800/50 flex items-center justify-end gap-3 border-t border-gray-100 dark:border-gray-600">
<button class="px-4 py-2 text-sm font-medium text-gray-700 dark:text-gray-300 bg-white dark:bg-gray-700 border border-gray-300 dark:border-gray-600 rounded-lg hover:bg-gray-50 dark:hover:bg-gray-600 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-primary transition-colors" type="button">
                Cancel
            </button>
<button class="px-4 py-2 text-sm font-medium text-white bg-primary rounded-lg hover:bg-gray-800 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-primary shadow-sm transition-colors" type="button">
                Save Entry
            </button>
</div>
</div>
</div>
<header class="pt-10 pb-4 px-6 sticky top-0 bg-background-light/95 dark:bg-background-dark/95 backdrop-blur-sm z-20 border-b border-gray-200 dark:border-gray-700">
<div class="flex justify-between items-start mb-4">
<div>
<h2 class="text-sm font-medium text-gray-500 dark:text-gray-400 uppercase tracking-widest">January 2024</h2>
<h1 class="text-3xl font-display font-bold text-gray-900 dark:text-white mt-1">Expense Tracker</h1>
</div>
<button class="p-2 rounded-full hover:bg-gray-200 dark:hover:bg-gray-700 transition-colors">
<span class="material-icons-outlined text-gray-600 dark:text-gray-300">calendar_today</span>
</button>
</div>
<div class="grid grid-cols-2 gap-3 mt-4">
<div class="bg-paper-light dark:bg-paper-dark p-3 rounded-DEFAULT shadow-sm border border-gray-100 dark:border-gray-600">
<p class="text-xs text-gray-500 dark:text-gray-400 mb-1">Total Spent</p>
<p class="text-xl font-bold font-display text-gray-900 dark:text-white">$1,245.50</p>
</div>
<div class="bg-paper-light dark:bg-paper-dark p-3 rounded-DEFAULT shadow-sm border border-gray-100 dark:border-gray-600">
<p class="text-xs text-gray-500 dark:text-gray-400 mb-1">Budget Left</p>
<p class="text-xl font-bold font-display text-emerald-600 dark:text-emerald-400">$754.50</p>
</div>
</div>
</header>
<main class="px-4 mt-2">
<div class="flex justify-between items-center py-4">
<div class="flex space-x-2 overflow-x-auto no-scrollbar pb-2">
<span class="px-3 py-1 bg-primary text-white text-xs rounded-full whitespace-nowrap">All</span>
<span class="px-3 py-1 bg-gray-200 dark:bg-gray-700 text-gray-600 dark:text-gray-300 text-xs rounded-full whitespace-nowrap">Food</span>
<span class="px-3 py-1 bg-gray-200 dark:bg-gray-700 text-gray-600 dark:text-gray-300 text-xs rounded-full whitespace-nowrap">Transport</span>
<span class="px-3 py-1 bg-gray-200 dark:bg-gray-700 text-gray-600 dark:text-gray-300 text-xs rounded-full whitespace-nowrap">Bills</span>
</div>
<button class="text-primary dark:text-white">
<span class="material-icons-outlined">filter_list</span>
</button>
</div>
<div class="bg-paper-light dark:bg-paper-dark rounded-DEFAULT shadow-paper overflow-hidden relative border border-gray-200 dark:border-gray-600">
<div class="absolute left-0 top-0 bottom-0 w-6 flex flex-col items-center justify-start py-4 gap-3 bg-gray-100 dark:bg-gray-800 border-r border-gray-300 dark:border-gray-600 z-10">
<div class="w-2 h-4 rounded-full bg-gradient-to-b from-gray-300 to-gray-500 shadow-sm"></div>
<div class="w-2 h-4 rounded-full bg-gradient-to-b from-gray-300 to-gray-500 shadow-sm"></div>
<div class="w-2 h-4 rounded-full bg-gradient-to-b from-gray-300 to-gray-500 shadow-sm"></div>
<div class="w-2 h-4 rounded-full bg-gradient-to-b from-gray-300 to-gray-500 shadow-sm"></div>
<div class="w-2 h-4 rounded-full bg-gradient-to-b from-gray-300 to-gray-500 shadow-sm"></div>
<div class="w-2 h-4 rounded-full bg-gradient-to-b from-gray-300 to-gray-500 shadow-sm"></div>
<div class="w-2 h-4 rounded-full bg-gradient-to-b from-gray-300 to-gray-500 shadow-sm"></div>
<div class="w-2 h-4 rounded-full bg-gradient-to-b from-gray-300 to-gray-500 shadow-sm"></div>
<div class="w-2 h-4 rounded-full bg-gradient-to-b from-gray-300 to-gray-500 shadow-sm"></div>
<div class="w-2 h-4 rounded-full bg-gradient-to-b from-gray-300 to-gray-500 shadow-sm"></div>
<div class="w-2 h-4 rounded-full bg-gradient-to-b from-gray-300 to-gray-500 shadow-sm"></div>
<div class="w-2 h-4 rounded-full bg-gradient-to-b from-gray-300 to-gray-500 shadow-sm"></div>
</div>
<div class="pl-8 overflow-x-auto">
<table class="w-full min-w-[500px] text-left border-collapse">
<thead>
<tr class="border-b-2 border-primary">
<th class="py-3 pl-4 pr-2 text-xs font-bold uppercase tracking-wider text-gray-500 dark:text-gray-400 w-16">Date</th>
<th class="py-3 px-2 text-xs font-bold uppercase tracking-wider text-gray-500 dark:text-gray-400 w-32">Description</th>
<th class="py-3 px-2 text-xs font-bold uppercase tracking-wider text-gray-500 dark:text-gray-400 w-24">Category</th>
<th class="py-3 px-2 text-xs font-bold uppercase tracking-wider text-gray-500 dark:text-gray-400 w-20 text-right">Amt</th>
<th class="py-3 px-2 text-xs font-bold uppercase tracking-wider text-gray-500 dark:text-gray-400 w-10 text-center">Paid</th>
</tr>
</thead>
<tbody class="text-sm font-hand text-lg">
<tr class="border-b border-gray-200 dark:border-gray-700 hover:bg-gray-50 dark:hover:bg-gray-700/50 transition-colors h-12">
<td class="pl-4 py-2 text-gray-600 dark:text-gray-300">01/15</td>
<td class="px-2 py-2 font-medium text-gray-800 dark:text-white">Grocery Run</td>
<td class="px-2 py-2">
<span class="bg-orange-100 dark:bg-orange-900 text-orange-800 dark:text-orange-200 text-xs font-sans px-2 py-0.5 rounded-full">Food</span>
</td>
<td class="px-2 py-2 text-right font-bold text-gray-800 dark:text-white">$120.50</td>
<td class="px-2 py-2 text-center text-emerald-500">
<span class="material-icons-outlined text-sm">check_circle</span>
</td>
</tr>
<tr class="border-b border-gray-200 dark:border-gray-700 hover:bg-gray-50 dark:hover:bg-gray-700/50 transition-colors h-12 bg-accent-highlight/30 dark:bg-yellow-900/10">
<td class="pl-4 py-2 text-gray-600 dark:text-gray-300">01/16</td>
<td class="px-2 py-2 font-medium text-gray-800 dark:text-white">Electricity Bill</td>
<td class="px-2 py-2">
<span class="bg-blue-100 dark:bg-blue-900 text-blue-800 dark:text-blue-200 text-xs font-sans px-2 py-0.5 rounded-full">Bills</span>
</td>
<td class="px-2 py-2 text-right font-bold text-gray-800 dark:text-white">$85.00</td>
<td class="px-2 py-2 text-center text-emerald-500">
<span class="material-icons-outlined text-sm">check_circle</span>
</td>
</tr>
<tr class="border-b border-gray-200 dark:border-gray-700 hover:bg-gray-50 dark:hover:bg-gray-700/50 transition-colors h-12">
<td class="pl-4 py-2 text-gray-600 dark:text-gray-300">01/18</td>
<td class="px-2 py-2 font-medium text-gray-800 dark:text-white">Uber Ride</td>
<td class="px-2 py-2">
<span class="bg-purple-100 dark:bg-purple-900 text-purple-800 dark:text-purple-200 text-xs font-sans px-2 py-0.5 rounded-full">Transport</span>
</td>
<td class="px-2 py-2 text-right font-bold text-gray-800 dark:text-white">$24.00</td>
<td class="px-2 py-2 text-center text-gray-300 dark:text-gray-600">
<span class="material-icons-outlined text-sm">radio_button_unchecked</span>
</td>
</tr>
<tr class="border-b border-gray-200 dark:border-gray-700 hover:bg-gray-50 dark:hover:bg-gray-700/50 transition-colors h-12">
<td class="pl-4 py-2 text-gray-600 dark:text-gray-300">01/20</td>
<td class="px-2 py-2 font-medium text-gray-800 dark:text-white">Coffee Date</td>
<td class="px-2 py-2">
<span class="bg-pink-100 dark:bg-pink-900 text-pink-800 dark:text-pink-200 text-xs font-sans px-2 py-0.5 rounded-full">Social</span>
</td>
<td class="px-2 py-2 text-right font-bold text-gray-800 dark:text-white">$15.75</td>
<td class="px-2 py-2 text-center text-emerald-500">
<span class="material-icons-outlined text-sm">check_circle</span>
</td>
</tr>
<tr class="border-b border-gray-200 dark:border-gray-700 hover:bg-gray-50 dark:hover:bg-gray-700/50 transition-colors h-12">
<td class="pl-4 py-2 text-gray-600 dark:text-gray-300">01/22</td>
<td class="px-2 py-2 font-medium text-gray-800 dark:text-white">Internet Bill</td>
<td class="px-2 py-2">
<span class="bg-blue-100 dark:bg-blue-900 text-blue-800 dark:text-blue-200 text-xs font-sans px-2 py-0.5 rounded-full">Bills</span>
</td>
<td class="px-2 py-2 text-right font-bold text-gray-800 dark:text-white">$60.00</td>
<td class="px-2 py-2 text-center text-gray-300 dark:text-gray-600">
<span class="material-icons-outlined text-sm">radio_button_unchecked</span>
</td>
</tr>
<tr class="border-b border-gray-200 dark:border-gray-700 hover:bg-gray-50 dark:hover:bg-gray-700/50 transition-colors h-12">
<td class="pl-4 py-2 text-gray-600 dark:text-gray-300">01/24</td>
<td class="px-2 py-2 font-medium text-gray-800 dark:text-white">Gym Membership</td>
<td class="px-2 py-2">
<span class="bg-green-100 dark:bg-green-900 text-green-800 dark:text-green-200 text-xs font-sans px-2 py-0.5 rounded-full">Health</span>
</td>
<td class="px-2 py-2 text-right font-bold text-gray-800 dark:text-white">$45.00</td>
<td class="px-2 py-2 text-center text-emerald-500">
<span class="material-icons-outlined text-sm">check_circle</span>
</td>
</tr>
<tr class="border-b border-gray-200 dark:border-gray-700 h-12">
<td class="pl-4 py-2"></td>
<td class="px-2 py-2"></td>
<td class="px-2 py-2"></td>
<td class="px-2 py-2"></td>
<td class="px-2 py-2"></td>
</tr>
<tr class="border-b border-gray-200 dark:border-gray-700 h-12">
<td class="pl-4 py-2"></td>
<td class="px-2 py-2"></td>
<td class="px-2 py-2"></td>
<td class="px-2 py-2"></td>
<td class="px-2 py-2"></td>
</tr>
</tbody>
</table>
</div>
<div class="h-2 bg-gray-100 dark:bg-gray-800 w-full border-t border-gray-200 dark:border-gray-700"></div>
</div>
<button class="fixed bottom-24 right-6 bg-primary text-white p-4 rounded-full shadow-lg hover:shadow-xl hover:bg-gray-700 transition-all transform hover:-translate-y-1 z-30">
<span class="material-icons-outlined text-2xl">add</span>
</button>
</main>
<nav class="fixed bottom-0 left-0 right-0 bg-white dark:bg-gray-900 border-t border-gray-200 dark:border-gray-700 pb-6 pt-3 px-8 z-40">
<div class="flex justify-between items-center">
<a class="flex flex-col items-center text-gray-400 hover:text-primary dark:hover:text-white transition-colors" href="#">
<span class="material-icons-outlined">home</span>
<span class="text-[10px] mt-1">Home</span>
</a>
<a class="flex flex-col items-center text-primary dark:text-white" href="#">
<span class="material-icons-outlined">edit_note</span>
<span class="text-[10px] mt-1 font-medium">Tracker</span>
</a>
<a class="flex flex-col items-center text-gray-400 hover:text-primary dark:hover:text-white transition-colors" href="#">
<span class="material-icons-outlined">pie_chart</span>
<span class="text-[10px] mt-1">Budget</span>
</a>
<a class="flex flex-col items-center text-gray-400 hover:text-primary dark:hover:text-white transition-colors" href="#">
<span class="material-icons-outlined">person</span>
<span class="text-[10px] mt-1">Profile</span>
</a>
</div>
</nav>
<script>
        // Simple script to toggle dark mode for demonstration
        // In a real app, this would be handled by system preference or user setting
        if (window.matchMedia && window.matchMedia('(prefers-color-scheme: dark)').matches) {
            document.documentElement.classList.add('dark');
        }
    </script>
</body></html>
```

## Bill Payment Tracker
```html
<!DOCTYPE html>
<html lang="en"><head>
<meta charset="utf-8"/>
<meta content="width=device-width, initial-scale=1.0" name="viewport"/>
<title>Bill Payment Tracker</title>
<script src="https://cdn.tailwindcss.com?plugins=forms,typography"></script>
<link href="https://fonts.googleapis.com/css2?family=Playfair+Display:ital,wght@0,400;0,600;0,700;1,400&amp;family=Inter:wght@300;400;500;600&amp;display=swap" rel="stylesheet"/>
<link href="https://fonts.googleapis.com/icon?family=Material+Icons+Outlined" rel="stylesheet"/>
<script>
        tailwind.config = {
            darkMode: "class",
            theme: {
                extend: {
                    colors: {
                        primary: "#374151", // Dark gray/slate for the "ink" look
                        "primary-accent": "#E11D48", // Rose for alerts/overdue
                        "background-light": "#F9FAFB", // Paper white
                        "background-dark": "#1F2937", // Dark slate for dark mode
                        "paper-light": "#FFFFFF",
                        "paper-dark": "#374151",
                        "line-light": "#E5E7EB",
                        "line-dark": "#4B5563",
                    },
                    fontFamily: {
                        display: ["'Playfair Display'", "serif"],
                        sans: ["'Inter'", "sans-serif"],
                    },
                    borderRadius: {
                        DEFAULT: "12px",
                        'lg': "16px",
                        'xl': "24px"
                    },
                    boxShadow: {
                        'paper': '0 4px 6px -1px rgba(0, 0, 0, 0.05), 0 2px 4px -1px rgba(0, 0, 0, 0.03)',
                    }
                },
            },
        };
    </script>
<style>
    body {
      min-height: max(884px, 100dvh);
    }
  </style>
  </head>
<body class="bg-background-light dark:bg-background-dark text-primary dark:text-gray-200 transition-colors duration-300 font-sans antialiased min-h-screen pb-20">
<header class="sticky top-0 z-50 bg-background-light/90 dark:bg-background-dark/90 backdrop-blur-md border-b border-line-light dark:border-line-dark px-4 py-3 flex items-center justify-between">
<div class="flex items-center space-x-2">
<button class="p-2 -ml-2 rounded-full hover:bg-gray-100 dark:hover:bg-gray-700 transition">
<span class="material-icons-outlined text-2xl text-primary dark:text-white">arrow_back</span>
</button>
<h1 class="font-display text-xl font-bold tracking-tight text-primary dark:text-white">Bill Payment</h1>
</div>
<div class="flex items-center space-x-3">
<button class="p-2 rounded-full hover:bg-gray-100 dark:hover:bg-gray-700 transition">
<span class="material-icons-outlined text-2xl text-primary dark:text-white">calendar_today</span>
</button>
<button class="p-2 rounded-full hover:bg-gray-100 dark:hover:bg-gray-700 transition relative">
<span class="material-icons-outlined text-2xl text-primary dark:text-white">notifications</span>
<span class="absolute top-2 right-2 w-2 h-2 bg-red-500 rounded-full"></span>
</button>
</div>
</header>
<main class="px-4 py-6 max-w-lg mx-auto space-y-8">
<div class="flex flex-col items-center justify-center space-y-1">
<h2 class="font-display text-3xl font-bold text-primary dark:text-white">October 2023</h2>
<p class="text-sm text-gray-500 dark:text-gray-400 font-serif italic">"Financial freedom starts with a plan."</p>
</div>
<section class="grid grid-cols-2 gap-4">
<div class="bg-paper-light dark:bg-paper-dark p-4 rounded-xl shadow-paper border border-line-light dark:border-line-dark flex flex-col justify-between h-32">
<div class="flex justify-between items-start">
<span class="material-icons-outlined text-gray-400 dark:text-gray-300">receipt_long</span>
<span class="text-xs font-semibold uppercase tracking-wider text-gray-500 dark:text-gray-400">Total Due</span>
</div>
<div>
<h3 class="text-2xl font-bold text-primary dark:text-white font-display">$4,250</h3>
<p class="text-xs text-red-500 mt-1 flex items-center">
<span class="material-icons-outlined text-[14px] mr-1">warning</span> 2 Overdue
                    </p>
</div>
</div>
<div class="bg-paper-light dark:bg-paper-dark p-4 rounded-xl shadow-paper border border-line-light dark:border-line-dark flex flex-col justify-between h-32">
<div class="flex justify-between items-start">
<span class="material-icons-outlined text-green-500">check_circle</span>
<span class="text-xs font-semibold uppercase tracking-wider text-gray-500 dark:text-gray-400">Paid</span>
</div>
<div>
<h3 class="text-2xl font-bold text-green-600 dark:text-green-400 font-display">$2,100</h3>
<p class="text-xs text-gray-500 dark:text-gray-400 mt-1">49% Completed</p>
</div>
</div>
</section>
<section class="bg-paper-light dark:bg-paper-dark rounded-xl shadow-paper border border-line-light dark:border-line-dark overflow-hidden relative">
<div class="absolute left-0 top-0 bottom-0 w-8 bg-gray-100 dark:bg-gray-800 border-r border-line-light dark:border-line-dark flex flex-col items-center py-4 space-y-4 z-10">
<div class="w-4 h-2 bg-gradient-to-r from-gray-400 to-gray-200 rounded-full shadow-sm"></div>
<div class="w-4 h-2 bg-gradient-to-r from-gray-400 to-gray-200 rounded-full shadow-sm"></div>
<div class="w-4 h-2 bg-gradient-to-r from-gray-400 to-gray-200 rounded-full shadow-sm"></div>
<div class="w-4 h-2 bg-gradient-to-r from-gray-400 to-gray-200 rounded-full shadow-sm"></div>
<div class="w-4 h-2 bg-gradient-to-r from-gray-400 to-gray-200 rounded-full shadow-sm"></div>
<div class="w-4 h-2 bg-gradient-to-r from-gray-400 to-gray-200 rounded-full shadow-sm"></div>
<div class="w-4 h-2 bg-gradient-to-r from-gray-400 to-gray-200 rounded-full shadow-sm"></div>
<div class="w-4 h-2 bg-gradient-to-r from-gray-400 to-gray-200 rounded-full shadow-sm"></div>
<div class="w-4 h-2 bg-gradient-to-r from-gray-400 to-gray-200 rounded-full shadow-sm"></div>
</div>
<div class="pl-12 pr-4 py-4 border-b border-line-light dark:border-line-dark bg-gray-50 dark:bg-gray-800/50 flex justify-between items-center">
<h3 class="font-display font-bold text-lg text-primary dark:text-white uppercase tracking-widest">Bill Details</h3>
<button class="text-primary dark:text-white hover:bg-gray-200 dark:hover:bg-gray-700 p-1 rounded">
<span class="material-icons-outlined">filter_list</span>
</button>
</div>
<div class="pl-12 divide-y divide-line-light dark:divide-line-dark">
<div class="p-4 flex items-center justify-between group bg-red-50/50 dark:bg-red-900/10">
<div class="flex items-center space-x-3">
<div class="h-10 w-10 rounded-full bg-white dark:bg-gray-700 border border-gray-200 dark:border-gray-600 flex items-center justify-center text-xl shadow-sm">
                            💡
                        </div>
<div>
<p class="font-semibold text-primary dark:text-white leading-tight">Electricity Bill</p>
<p class="text-xs text-red-500 font-medium">Due: Oct 15 (Overdue)</p>
</div>
</div>
<div class="text-right">
<p class="font-bold text-primary dark:text-white">$120.50</p>
<span class="inline-block px-2 py-0.5 text-[10px] font-bold uppercase tracking-wider text-red-600 bg-red-100 dark:bg-red-900/40 dark:text-red-300 rounded-full">Unpaid</span>
</div>
</div>
<div class="p-4 flex items-center justify-between group relative overflow-hidden">
<div class="absolute right-12 top-2 opacity-10 pointer-events-none transform rotate-12">
<span class="text-green-600 dark:text-green-400 text-6xl font-display">✓</span>
</div>
<div class="flex items-center space-x-3 opacity-60">
<div class="h-10 w-10 rounded-full bg-white dark:bg-gray-700 border border-gray-200 dark:border-gray-600 flex items-center justify-center text-xl shadow-sm">
                            💧
                        </div>
<div>
<p class="font-semibold text-primary dark:text-white line-through decoration-gray-400">Water Bill</p>
<p class="text-xs text-gray-400 dark:text-gray-500">Paid: Oct 12</p>
</div>
</div>
<div class="text-right opacity-60">
<p class="font-bold text-primary dark:text-white line-through decoration-gray-400">$45.00</p>
<span class="inline-block px-2 py-0.5 text-[10px] font-bold uppercase tracking-wider text-green-600 bg-green-100 dark:bg-green-900/40 dark:text-green-300 rounded-full">Paid</span>
</div>
</div>
<div class="p-4 flex items-center justify-between group">
<div class="flex items-center space-x-3">
<div class="h-10 w-10 rounded-full bg-white dark:bg-gray-700 border border-gray-200 dark:border-gray-600 flex items-center justify-center text-xl shadow-sm">
                            🏠
                        </div>
<div>
<p class="font-semibold text-primary dark:text-white leading-tight">Rent</p>
<p class="text-xs text-gray-500 dark:text-gray-400">Due: Oct 28</p>
</div>
</div>
<div class="text-right">
<p class="font-bold text-primary dark:text-white">$2,500.00</p>
<button class="mt-1 text-xs text-primary underline decoration-dotted hover:text-blue-600 dark:hover:text-blue-400">Mark Paid</button>
</div>
</div>
<div class="p-4 flex items-center justify-between group">
<div class="flex items-center space-x-3">
<div class="h-10 w-10 rounded-full bg-white dark:bg-gray-700 border border-gray-200 dark:border-gray-600 flex items-center justify-center text-xl shadow-sm">
                            🚗
                        </div>
<div>
<p class="font-semibold text-primary dark:text-white leading-tight">Car Loan</p>
<p class="text-xs text-gray-500 dark:text-gray-400">Due: Oct 30</p>
</div>
</div>
<div class="text-right">
<p class="font-bold text-primary dark:text-white">$400.00</p>
<button class="mt-1 text-xs text-primary underline decoration-dotted hover:text-blue-600 dark:hover:text-blue-400">Mark Paid</button>
</div>
</div>
<div class="p-4 flex items-center justify-between group">
<div class="flex items-center space-x-3">
<div class="h-10 w-10 rounded-full bg-white dark:bg-gray-700 border border-gray-200 dark:border-gray-600 flex items-center justify-center text-xl shadow-sm">
                            📺
                        </div>
<div>
<p class="font-semibold text-primary dark:text-white leading-tight">Netflix</p>
<p class="text-xs text-gray-500 dark:text-gray-400">Due: Oct 31</p>
</div>
</div>
<div class="text-right">
<p class="font-bold text-primary dark:text-white">$15.99</p>
<button class="mt-1 text-xs text-primary underline decoration-dotted hover:text-blue-600 dark:hover:text-blue-400">Mark Paid</button>
</div>
</div>
</div>
<div class="pl-12 bg-paper-light dark:bg-paper-dark border-t border-dashed border-gray-300 dark:border-gray-600">
<div class="h-10 border-b border-line-light dark:border-line-dark w-full"></div>
<div class="h-10 border-b border-line-light dark:border-line-dark w-full"></div>
<div class="h-10 border-b border-line-light dark:border-line-dark w-full"></div>
</div>
</section>
<section class="space-y-4">
<h3 class="font-display font-bold text-xl text-primary dark:text-white pl-2 border-l-4 border-primary dark:border-white">Monthly Insights</h3>
<div class="bg-paper-light dark:bg-paper-dark p-6 rounded-xl shadow-paper border border-line-light dark:border-line-dark">
<div class="flex justify-between items-center mb-4">
<p class="text-sm font-semibold text-gray-600 dark:text-gray-300">Remaining Budget</p>
<span class="text-sm text-green-600 dark:text-green-400 font-bold">$1,250 left</span>
</div>
<div class="relative w-full h-4 bg-gray-200 dark:bg-gray-700 rounded-full overflow-hidden border border-gray-300 dark:border-gray-600">
<div class="absolute top-0 left-0 h-full w-[70%] bg-primary dark:bg-gray-400" style="background-image: repeating-linear-gradient(45deg, transparent, transparent 5px, rgba(255,255,255,0.1) 5px, rgba(255,255,255,0.1) 10px);"></div>
</div>
<div class="flex justify-between mt-2 text-xs text-gray-500 dark:text-gray-400">
<span>Spent: $3,000</span>
<span>Limit: $4,250</span>
</div>
<div class="mt-6 pt-4 border-t border-dashed border-gray-300 dark:border-gray-600">
<p class="text-sm text-gray-500 dark:text-gray-400 italic">"You are spending 10% less on utilities compared to last month. Keep it up!"</p>
</div>
</div>
</section>
</main>
<div class="fixed bottom-6 right-6 z-50">
<button class="bg-primary text-white p-4 rounded-full shadow-lg hover:bg-gray-800 transition transform hover:scale-105 flex items-center justify-center">
<span class="material-icons-outlined text-3xl">add</span>
</button>
</div>

</body></html>
```