import { Fancybox } from "@fancyapps/ui";
import { createApp } from "vue";
import ArticleCodeBlock from "../components/article/ArticleCodeBlock.vue";
import { preloadLanguages } from "./useCodeHighlighter";

function friendlyLabelFromId(languageId) {
  switch (languageId) {
    case "javascript":
      return "JavaScript";
    case "typescript":
      return "TypeScript";
    case "bash":
      return "Shell";
    case "go":
      return "Go";
    case "json":
      return "JSON";
    case "plaintext":
    default:
      return "Code";
  }
}

function normalizeLanguageId(languageId) {
  if (!languageId) return "plaintext";
  const normalized = String(languageId).trim().toLowerCase();
  if (!normalized || normalized === "text") return "plaintext";
  if (normalized === "shell") return "bash";
  return normalized;
}

function fallbackCodeBlockMeta(preElement, index) {
  const codeElement = preElement.querySelector("code");
  const className = codeElement?.className || "";
  const languageClass = className
    .split(/\s+/)
    .find((token) => token.startsWith("language-"));
  const languageId = normalizeLanguageId(
    languageClass ? languageClass.replace("language-", "") : "plaintext"
  );

  return {
    index,
    languageId,
    languageLabel: friendlyLabelFromId(languageId),
    showLineNumbers: true,
  };
}

function setupGallery(root) {
  root.querySelectorAll("img").forEach((img) => {
    img.loading = "lazy";
    const parent = img.parentElement;
    if (!parent) return;
    if (parent.tagName !== "A") {
      const wrapper = document.createElement("a");
      wrapper.href = img.src;
      wrapper.target = "_blank";
      wrapper.dataset.fancybox = "gallery";
      parent.replaceChild(wrapper, img);
      wrapper.appendChild(img);
      return;
    }
    parent.href = img.src;
    parent.target = "_blank";
    parent.dataset.fancybox = "gallery";
  });

  if (typeof Fancybox.unbind === "function") {
    Fancybox.unbind("[data-fancybox]");
  }
  Fancybox.bind("[data-fancybox]", {});
}

function enhanceImageGroups(root) {
  root.querySelectorAll("p, figure, div").forEach((container) => {
    const images = container.querySelectorAll(":scope > a[data-fancybox] > img, :scope > img");
    if (images.length >= 2) {
      container.classList.add("image-gallery-group");
    }
  });
}

export function useArticleEnhancements() {
  const codeBlockApps = [];

  const cleanupCodeBlocks = () => {
    while (codeBlockApps.length) {
      const mountedApp = codeBlockApps.pop();
      mountedApp?.app?.unmount();
    }
  };

  const enhanceCodeBlocks = async ({ root, codeBlocks = [] }) => {
    cleanupCodeBlocks();

    const preElements = root.querySelectorAll("pre");
    const validBlocks = [];
    preElements.forEach((preElement, index) => {
      const codeElement = preElement.querySelector("code");
      if (!codeElement) return;
      const meta = codeBlocks[index] || fallbackCodeBlockMeta(preElement, index);
      validBlocks.push({ preElement, codeElement, meta });
    });

    const languageIds = validBlocks.map((b) => normalizeLanguageId(b.meta.languageId));
    await preloadLanguages(languageIds);

    validBlocks.forEach(({ preElement, codeElement, meta }) => {
      const mountPoint = document.createElement("div");
      mountPoint.className = "article-code-block-host";

      const app = createApp(ArticleCodeBlock, {
        code: codeElement.textContent ?? "",
        languageId: normalizeLanguageId(meta.languageId),
        languageLabel: meta.languageLabel || friendlyLabelFromId(meta.languageId),
        showLineNumbers: Boolean(meta.showLineNumbers),
      });

      preElement.replaceWith(mountPoint);
      app.mount(mountPoint);
      codeBlockApps.push({ app, mountPoint });
    });
  };

  const cleanupEnhancements = () => {
    cleanupCodeBlocks();
  };

  const enhanceArticleContent = async ({ root, codeBlocks = [] }) => {
    if (!root) return;
    await enhanceCodeBlocks({ root, codeBlocks });
    setupGallery(root);
    enhanceImageGroups(root);
  };

  return {
    enhanceArticleContent,
    cleanupEnhancements,
  };
}
