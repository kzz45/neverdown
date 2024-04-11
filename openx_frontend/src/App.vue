<template>
  <router-view v-slot="{ Component }">
    <component :is="Component" />
  </router-view>
</template>

<script setup lang="ts"></script>

<style lang="scss">
$mousepos: 0 0;
#app {
  font-family: Helvetica, Tahoma, Arial;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  text-align: center;
  color: #2c3e50;
}

#nav {
  padding: 30px;
}
#nav a {
  font-weight: bold;
  color: #2c3e50;
}
#nav a.router-link-exact-active {
  color: #42b983;
}

.fade-enter-active {
  animation: bounce-in 0.6s;
}
.fade-leave-active {
  animation: bounce-out 0.3s;
}
.slide-left-leave-active {
  animation: left-leave 0.6s;
}
@keyframes bounce-in {
  0% {
    transform: translateX(-100%);
  }
  50% {
    transform-origin: $mousepos;
    transform: scale(0);
  }
  100% {
    transform-origin: $mousepos;
    transform: scale(1);
  }
}
@keyframes bounce-out {
  0% {
    transform-origin: $mousepos;
    transform: scale(1);
  }
  100% {
    transform-origin: $mousepos;
    transform: scale(0);
  }
}
@keyframes left-leave {
  0% {
    transform-origin: $mousepos;
    transform: scale(1);
  }
  100% {
    transform-origin: $mousepos;
    transform: scale(0);
  }
}
</style>
