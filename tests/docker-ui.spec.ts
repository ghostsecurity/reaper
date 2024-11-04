import { test, expect } from '@playwright/test';

test('sign in button', async ({ page }) => {
  await page.goto('http://localhost:8000');

  // Click the get started link.
  await page.getByRole('button', { name: 'Sign in' }).click();

  // Expects page to have a heading with the name of Installation.
  await expect(page.getByRole('heading', { name: 'Welcome to Reaper' })).toBeVisible();
});
