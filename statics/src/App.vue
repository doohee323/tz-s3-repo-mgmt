<template>
  <div id="app">
    <div class="align-left loading">
      <img alt="" src="./assets/loading.gif"/>
    </div>
    <nav class="navbar navbar-expand-lg navbar-light bg-light">
      <ul class="navbar-nav mr-auto">
        <li v-if="hasAuths()" class="nav-item active">
          <router-link :to="{name: 'Home'}" class="nav-link">Home</router-link>
        </li>
        <li v-if="hasAuths()" class="nav-item">
          <router-link :to="{name: 'Test'}" class="nav-link">Test</router-link>
        </li>
        <li v-if="hasAuths() && isAdmin" class="nav-item">
          <router-link :to="{name: 'S3Repos'}" class="nav-link">S3 Repo</router-link>
        </li>
        <li v-if="hasAuths()" class="nav-item">
          <router-link :to="{name: 'Users'}" class="nav-link">Users</router-link>
        </li>
      </ul>
      <form v-if="hasAuths()" id="logout" class="form-inline my-2 my-lg-0">
        <button class="btn btn-outline-success my-2 my-sm-0" @click="logout()">Logout</button>
      </form>
    </nav>
    <router-view></router-view>
  </div>
</template>

<script>

import {getCache, initSCacheEx, setCache} from '@/assets/cache.js'
import $ from "jquery";

export default {
  components: {},
  data() {
    this.$root.hasAuth = getCache('session', 'hasAuth') === 'true';
    return {
      hasAuth: this.$root.hasAuth,
      isAdmin: getCache('session', 'admin') === 'true',
    };
  },
  methods: {
    hasAuths() {
      return this.$root.hasAuth;
    },
    logout() {
      $('.loading').show();
      let apiUrl = getCache('config', 'apiUrl');
      let vData = {
        token: getCache('session', 'token'),
      }
      $.ajax({
        method: "POST",
        url: apiUrl + '/auth/logout',
        data: JSON.stringify(vData),
      }).done(function (e) {
        console.log(e);
        $('.loading').hide();
        initSCacheEx('session', 'hasAuth');
        initSCacheEx('session', 'token');
        initSCacheEx('session', 'user');
        initSCacheEx('session', 'admin');
      });
    },
  },
  mounted() {
    let apiUrl = getCache('config', 'apiUrl');
    if (apiUrl === '') {
      apiUrl = document.location.protocol + '//' + document.location.hostname;
      if (document.location.hostname === 'localhost') {
        apiUrl += ':8080/api';
      } else {
        apiUrl += '/api';
      }
      setCache('config', 'apiUrl', apiUrl);
    }
    if (getCache('session', 'hasAuth') !== 'true') {
      if (this.$route.name !== 'Login') {
        if (this.$route.name !== 'Login') {
          this.$router.push('/login');
        }
      }
    }
    this.$root.hasAuth = getCache('session', 'hasAuth') === 'true';
  },
}
</script>
