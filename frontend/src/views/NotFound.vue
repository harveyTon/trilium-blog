<template>
  <div class="not-found">
    <h1>404</h1>
    <p>{{ t('notFound.message') }}</p>
    <p>{{ t('notFound.countdown', { seconds: countdown }) }}</p>
    <router-link to="/">{{ t('notFound.goHome') }}</router-link>
  </div>
</template>

<script>
import { t } from "../i18n";

export default {
  name: "NotFound",
  data() {
    return {
      countdown: 3,
    };
  },
  mounted() {
    this.startCountdown();
  },
  methods: {
    startCountdown() {
      const timer = setInterval(() => {
        this.countdown--;
        if (this.countdown <= 0) {
          clearInterval(timer);
          this.$router.push("/");
        }
      }, 1000);
    },
  },
};
</script>

<style scoped>
.not-found {
  text-align: center;
  padding: 50px;
}

h1 {
  font-size: clamp(3rem, 12vw, 72px);
  color: var(--accent);
}

a {
  color: var(--accent);
  text-decoration: none;
}

a:hover,
a:focus-visible {
  text-decoration: underline;
}
</style>
