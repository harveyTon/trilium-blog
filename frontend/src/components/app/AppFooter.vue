<template>
  <el-footer class="app-footer">
    <div class="footer-content">
      <div class="footer-text">
        <p>
          © {{ new Date().getFullYear() }} {{ site.title }}.
          {{ t('footer.copyright') }}
        </p>
        <p>
          {{ t('footer.poweredByPrefix') }}
          <a href="https://github.com/harveyTon/trilium-blog" target="_blank"
            >Trilium Blog</a
          >
          &
          <a href="https://github.com/TriliumNext/Trilium" target="_blank"
            >Trilium</a
          >
          {{ t('footer.poweredBySuffix') }}
        </p>
      </div>
    </div>
  </el-footer>
</template>

<script>
import { ElFooter } from "element-plus";
import { t } from "../../i18n";

export default {
  name: "AppFooter",
  methods: { t },
  components: {
    ElFooter,
  },
  props: {
    site: {
      type: Object,
      required: true,
    },
  },
};
</script>
