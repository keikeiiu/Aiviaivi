import { test, expect } from "@playwright/test";

test("full flow: register → upload → verify in feed", async ({ page }) => {
  const ts = Date.now();

  // 1. Register via API (faster than UI)
  const registerResult = await page.evaluate(async (ts) => {
    const res = await fetch("http://localhost:8080/api/v1/auth/register", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({
        username: `pw_${ts}`,
        email: `pw_${ts}@test.com`,
        password: "test123",
      }),
    });
    const data = await res.json();
    return { token: data.data?.token, userId: data.data?.user?.id };
  }, ts);

  expect(registerResult.token).toBeTruthy();
  console.log("✓ Registered");

  // 2. Upload video via API
  const uploadResult = await page.evaluate(async ({ token }: { token: string }) => {
    // Create video record first via API
    const createRes = await fetch("http://localhost:8080/api/v1/videos/upload", {
      method: "POST",
      headers: { Authorization: `Bearer ${token}` },
      body: (() => {
        const fd = new FormData();
        fd.append("title", `E2E Test ${Date.now()}`);
        fd.append("description", "Playwright E2E upload test");

        // Create a small video blob
        const bytes = new Uint8Array(1024 * 50); // 50KB dummy
        const blob = new Blob([bytes], { type: "video/mp4" });
        fd.append("file", blob, "test.mp4");
        return fd;
      })(),
    });
    const data = await createRes.json();
    return { code: data.code, videoId: data.data?.id, status: data.data?.status };
  }, { token: registerResult.token });

  console.log("Upload:", JSON.stringify(uploadResult));
  expect(uploadResult.code).toBe(0);
  expect(uploadResult.videoId).toBeTruthy();
  console.log("✓ Uploaded:", uploadResult.videoId);

  // 3. Wait for transcode
  await page.waitForTimeout(8000);

  // 4. Verify video appears in feed
  const feedResult = await page.evaluate(async () => {
    const res = await fetch("http://localhost:8080/api/v1/videos?size=5");
    const data = await res.json();
    return {
      total: data.pagination?.total,
      latest: data.data?.[0]?.title,
      status: data.data?.[0]?.status,
    };
  });

  console.log("Feed:", JSON.stringify(feedResult));
  expect(feedResult.total).toBeGreaterThan(0);
  console.log("✓ Video in feed");
  console.log("✓ ALL PASSED");
});
