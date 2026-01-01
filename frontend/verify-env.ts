// Run this in the browser console to verify environment variables
console.log('=== Environment Variables Check ===');
console.log('PUBLIC_CLERK_PUBLISHABLE_KEY:', import.meta.env.PUBLIC_CLERK_PUBLISHABLE_KEY ? '✅ FOUND' : '❌ MISSING');
console.log('PUBLIC_DEV_AUTH_PROVIDER:', import.meta.env.PUBLIC_DEV_AUTH_PROVIDER);
console.log('PUBLIC_API_URL:', import.meta.env.PUBLIC_API_URL);
console.log('PUBLIC_APP_NAME:', import.meta.env.PUBLIC_APP_NAME);
console.log('PUBLIC_APP_SHORT_NAME:', import.meta.env.PUBLIC_APP_SHORT_NAME);
console.log('MODE:', import.meta.env.MODE);
console.log('DEV:', import.meta.env.DEV);
console.log('PROD:', import.meta.env.PROD);
console.log('All PUBLIC_ env vars:', Object.keys(import.meta.env).filter(k => k.startsWith('PUBLIC_')));
