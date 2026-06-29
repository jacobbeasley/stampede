import { request, expect } from '@playwright/test';

async function globalSetup() {
  const requestContext = await request.newContext({
    baseURL: 'http://localhost:3000',
  });

  // Fetch the registration page to get the CSRF token
  const response = await requestContext.get('/register');
  const responseText = await response.text();

  // Extract CSRF token from the meta tag
  const csrfMatch = responseText.match(/<meta name="csrf-token" content="(.*?)"/);
  if (!csrfMatch) {
    console.error("Could not find CSRF token");
    return;
  }
  const csrfToken = csrfMatch[1];

  // Submit registration form
  const registerResponse = await requestContext.post('/register', {
    form: {
      authenticity_token: csrfToken,
      Email: 'admin@example.com',
      Password: 'password',
      PasswordConfirmation: 'password',
    },
  });

  // Check valid registration
  if (!registerResponse.ok()) {
     console.log("Status: " + registerResponse.status());
     console.log(await registerResponse.text());
  }
  expect(registerResponse.status() === 200 || registerResponse.status() === 303).toBeTruthy();

  await requestContext.dispose();
}

export default globalSetup;
