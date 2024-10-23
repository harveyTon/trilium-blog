<template>
  <div class="not-found">
    <h1>404</h1>
    <p>抱歉，您访问的页面不存在。</p>
    <p>{{ countdown }}秒后自动返回首页</p>
    <router-link to="/">立即返回首页</router-link>
  </div>
</template>

<script>
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
  font-size: 72px;
  color: #409eff;
}

a {
  color: #409eff;
  text-decoration: none;
}
</style>
