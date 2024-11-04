import { test, expect } from '@playwright/test';

test('sign in button', async ({ page }) => {
  await page.goto('http://localhost:8000');

  // Click the sign in button.
  await page.getByRole('button', { name: 'Sign in' }).click();

  // Expects page to have a heading with Scan title.
  await expect(page.getByRole('heading', { name: 'Scan' })).toBeVisible();
});
