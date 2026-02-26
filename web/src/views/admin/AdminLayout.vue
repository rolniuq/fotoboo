<template>
  <div class="admin-layout">
    <nav class="admin-nav">
      <router-link to="/" class="nav-brand">FotoBoo</router-link>
      <div class="nav-links">
        <router-link :to="{ name: 'admin-dashboard' }" class="nav-link">Dashboard</router-link>
        <router-link :to="{ name: 'admin-photos' }" class="nav-link">Photos</router-link>
        <router-link :to="{ name: 'admin-sessions' }" class="nav-link">Sessions</router-link>
        <router-link :to="{ name: 'admin-devices' }" class="nav-link">Devices</router-link>
        <router-link :to="{ name: 'admin-config' }" class="nav-link">Config</router-link>
      </div>
    </nav>

    <main class="admin-content">
      <router-view />
    </main>
  </div>
</template>

<script setup>
</script>

<style scoped>
.admin-layout {
  min-height: 100vh;
  background: #f5f7fa;
}

.admin-nav {
  background: white;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.08);
  padding: 0 2rem;
  display: flex;
  align-items: center;
  height: 60px;
  gap: 2rem;
}

.nav-brand {
  font-size: 1.4rem;
  font-weight: 700;
  color: var(--color-primary);
  text-decoration: none;
}

.nav-links {
  display: flex;
  gap: 0.25rem;
}

.nav-link {
  padding: 0.5rem 1rem;
  border-radius: var(--radius-sm);
  color: var(--color-muted);
  text-decoration: none;
  font-weight: 500;
  font-size: 0.9rem;
  transition: var(--transition);
}

.nav-link:hover {
  background: #eef1f6;
  color: var(--color-dark);
}

.nav-link.router-link-active {
  background: var(--color-primary);
  color: white;
}

.admin-content {
  padding: 2rem;
  max-width: 1200px;
  margin: 0 auto;
}

@media (max-width: 768px) {
  .admin-nav {
    padding: 0 1rem;
    gap: 1rem;
  }

  .nav-links {
    gap: 0;
    overflow-x: auto;
  }

  .nav-link {
    padding: 0.4rem 0.6rem;
    font-size: 0.8rem;
    white-space: nowrap;
  }

  .admin-content {
    padding: 1rem;
  }
}
</style>
