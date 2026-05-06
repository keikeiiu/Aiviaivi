import { test, expect } from "@playwright/test";

test("upload via API with FormData", async ({ page }) => {
  // Login first
  const loginResult = await page.evaluate(async () => {
    const res = await fetch("http://localhost:8080/api/v1/auth/login", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ email: "demo@ailivili.dev", password: "demo123" }),
    });
    return (await res.json()).data?.token;
  });
  expect(loginResult).toBeTruthy();
  console.log("✓ Logged in");

  // Upload with FormData
  const uploadResult = await page.evaluate(async (token) => {
    const fd = new FormData();
    fd.append("title", "Playwright Multipart Test");

    // Create a real video-like blob
    const bytes = new Uint8Array(50 * 1024);
    const blob = new Blob([bytes], { type: "video/mp4" });
    fd.append("file", blob, "test.mp4");

    const res = await fetch("http://localhost:8080/api/v1/videos/upload", {
      method: "POST",
      headers: { Authorization: `Bearer ${token}` },
      body: fd,
    });
    const data = await res.json();
    return { status: res.status, code: data.code, videoId: data.data?.id };
  }, loginResult);

  console.log("Result:", JSON.stringify(uploadResult));
  expect(uploadResult.status).toBe(200);
  expect(uploadResult.code).toBe(0);
  console.log("✓ Upload works with FormData!");
});
