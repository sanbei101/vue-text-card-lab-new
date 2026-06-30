import path from "path";

import { Hono } from "hono";
import { Window, FontLibrary } from "skia-canvas";
import { createSSRApp } from "vue";
import { renderToString } from "vue/server-renderer";

import CardSvg from "@/card/CardSvg.vue";
import { templates } from "@/templates";

const FONTS_DIR = path.resolve(process.cwd(), "./server/fonts");

const win = new Window(1, 1) as unknown as {
  document: Document;
  OffscreenCanvas: typeof OffscreenCanvas;
};
(globalThis as unknown as Record<string, unknown>).document = win.document;
(globalThis as unknown as Record<string, unknown>).OffscreenCanvas = win.OffscreenCanvas;

// ── 注册 Google Fonts ────────────────────────────────────────────────
FontLibrary.use("ZCOOL XiaoWei", path.join(FONTS_DIR, "zcool-xiaowei.ttf"));
FontLibrary.use("Ma Shan Zheng", path.join(FONTS_DIR, "ma-shan-zheng.ttf"));
FontLibrary.use("ZCOOL KuaiLe", path.join(FONTS_DIR, "zcool-kuaile.ttf"));
FontLibrary.use("Long Cang", path.join(FONTS_DIR, "long-cang.ttf"));
FontLibrary.use("ZCOOL QingKe HuangYou", path.join(FONTS_DIR, "zcool-qingke-huangyou.ttf"));

// ── 路由 ────────────────────────────────────────────────────────────
const app = new Hono();

app.get("/api/cards", async (c) => {
  try {
    const title = (c.req.query("title") as string) || "默认标题";
    const keyword = (c.req.query("keyword") as string) || "";

    const results: Array<{ id: string; name: string; svg: string }> = [];

    for (const tpl of templates) {
      const ssrApp = createSSRApp(CardSvg, {
        template: tpl,
        title,
        keyword,
      });

      const svgString = await renderToString(ssrApp);
      results.push({ id: tpl.id, name: tpl.name, svg: svgString });
    }

    return c.json({ templates: results });
  } catch (error) {
    console.error("生成卡片失败:", error);
    return c.json({ error: "服务器内部错误" }, 500);
  }
});

app.get("/health", (c) => c.text("ok"));

// ── 启动 ────────────────────────────────────────────────────────────
const PORT = Number(process.env.PORT ?? 5174);

console.log(`Card Engine Server running at http://localhost:${PORT}`);

export default {
  port: PORT,
  fetch: app.fetch,
};
