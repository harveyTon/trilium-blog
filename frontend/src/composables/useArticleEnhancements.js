import { Fancybox } from "@fancyapps/ui";

function injectCodeToolbar(codeElement, language) {
  const pre = codeElement.closest("pre");
  if (!pre || pre.dataset.enhanced === "true") return;

  const toolbar = document.createElement("div");
  toolbar.className = "code-toolbar";

  const langLabel = document.createElement("span");
  langLabel.className = "code-language";
  langLabel.textContent = language || "text";

  const actions = document.createElement("div");
  actions.className = "code-actions";

  const copyButton = document.createElement("button");
  copyButton.type = "button";
  copyButton.className = "code-action-button";
  copyButton.textContent = "复制";
  copyButton.addEventListener("click", async () => {
    try {
      await navigator.clipboard.writeText(codeElement.textContent || "");
      copyButton.textContent = "已复制";
      window.setTimeout(() => {
        copyButton.textContent = "复制";
      }, 1200);
    } catch {
      copyButton.textContent = "失败";
      window.setTimeout(() => {
        copyButton.textContent = "复制";
      }, 1200);
    }
  });

  const collapseButton = document.createElement("button");
  collapseButton.type = "button";
  collapseButton.className = "code-action-button";
  collapseButton.textContent = "折叠";
  collapseButton.addEventListener("click", () => {
    const collapsed = pre.dataset.collapsed === "true";
    pre.dataset.collapsed = collapsed ? "false" : "true";
    collapseButton.textContent = collapsed ? "折叠" : "展开";
  });

  actions.appendChild(copyButton);
  actions.appendChild(collapseButton);
  toolbar.appendChild(langLabel);
  toolbar.appendChild(actions);
  pre.parentNode.insertBefore(toolbar, pre);
  pre.dataset.enhanced = "true";
}

function setupGallery() {
  document.querySelectorAll(".article-content img").forEach((img) => {
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
  Fancybox.bind("[data-fancybox]", {});
}

function enhanceImageGroups() {
  document.querySelectorAll(".article-content p, .article-content figure, .article-content div").forEach((container) => {
    const images = container.querySelectorAll(":scope > a[data-fancybox] > img, :scope > img");
    if (images.length >= 2) {
      container.classList.add("image-gallery-group");
    }
  });
}

export function useArticleEnhancements({ hljs, applyHighlightTheme }) {
  const enhanceCodeBlocks = () => {
    applyHighlightTheme();
    document.querySelectorAll("pre code").forEach((el) => {
      const code = el.textContent ?? "";
      const languageMatch = el.className.match(/language-(\S+)/);
      let resolvedLanguage = "text";
      if (languageMatch) {
        resolvedLanguage = languageMatch[1];
        el.innerHTML = hljs.highlight(code, {
          language: languageMatch[1],
        }).value;
      } else {
        const result = hljs.highlightAuto(code);
        resolvedLanguage = result.language || "text";
        el.innerHTML = result.value;
      }
      el.classList.add("hljs");
      injectCodeToolbar(el, resolvedLanguage);
    });
  };

  const enhanceArticleContent = () => {
    enhanceCodeBlocks();
    setupGallery();
    enhanceImageGroups();
  };

  return {
    enhanceArticleContent,
  };
}
