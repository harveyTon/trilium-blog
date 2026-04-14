<template>
  <div v-if="pageUrl" class="article-source">
    {{ t('sourceLink.label') }}
    <a :href="pageUrl" target="_blank" rel="noopener noreferrer">
      {{ pageUrl }}
    </a>
  </div>
</template>

<script>
import { t } from "../../i18n";

export default {
  name: "SourceLinkBlock",
  props: {
    pageUrl: {
      type: String,
      default: "",
    },
  },
};
</script>
