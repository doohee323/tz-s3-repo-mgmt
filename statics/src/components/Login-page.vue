<template>

  <div class="container-fluid">
    <div class="row">
      <div class="col-md-12 min-vh-100 d-flex flex-column justify-content-center">
        <div class="col-lg-5 col-md-6 mx-auto">
          <div class="card rounded shadow shadow-sm">
            <div class="card-header">
              <h3 class="mb-0">Login</h3>
            </div>
            <div class="card-body">
              <form autocomplete="off" class="require" data-parsley-validate method="post">
                <div class="form-group error-handler">
                  <label for="username">Username</label>
                  <input id="username" v-model="username" class="form-control" data-parsley-errors-messages-disabled="true" data-parsley-required
                         data-parsley-trigger-after-failure="keyup"
                         data-parsley-validation-threshold="0" name="username" required
                         type="text"
                         @keyup="updated">
                </div>
                <div class="form-group error-handler">
                  <label for="password">Password</label>
                  <input id="password" v-model="password" autocomplete="off" class="form-control" name="password" required
                         type="password" @keyup="updated">
                </div>
                <div v-show="wrongAuth" class="form-group error-handler">
                  <ul class="list-unstyled help-block is-invalid col-lg-12 filled">
                    <li class="parsley-required">Wrong Password.</li>
                  </ul>
                </div>
                <button class="btn btn-lg btn-success btn-block" type="submit" @click="login()">
                  Log In
                </button>
              </form>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>

</template>

<script>

import {getCache, setCache} from '@/assets/cache'
import {resetParsley} from "@/assets/forms";
import $ from "jquery";

export default {
  components: {},
  data() {
    this.$root.hasAuth = getCache('session', 'hasAuth') === 'true';
    return {
      hasAuth: this.$root.hasAuth,
      username: '',
      password: '',
      wrongAuth: false,
    };
  },
  computed: {
    $form() {
      return $(this.$el).find('form');
    },
  },
  methods: {
    login() {
      this.getData((res) => {
        if (res.token) {
          setCache('session', 'hasAuth', 'true');
          setCache('session', 'token', res.token);
          setCache('session', 'user', JSON.stringify(res.user));
          setCache('session', 'admin', res.user.role === 'admin' ? 'true' : 'false');
          this.$root.hasAuth = getCache('session', 'hasAuth') === 'true';
          if (res.config.awsRegions) {
            setCache('config', 'awsRegions', res.config.awsRegions);
          }
          if (res.config.s3Repos) {
            setCache('config', 's3Repos', res.config.s3Repos);
          }
          if (this.$route.name !== 'Home') {
            this.$router.push('/');
          }
          $('#logout').show();
        } else {
          this.wrongAuth = true;
        }
      })
    },
    getData(cb) {
      $('.loading').show();
      let apiUrl = getCache('config', 'apiUrl');
      let vData = {
        username: this.username,
        password: this.password,
      }
      $.ajax({
        url: apiUrl + '/auth/login',
        type: "POST",
        data: JSON.stringify(vData),
        success: function (res) {
          $('.loading').hide();
          cb(res);
        },
        error: function (res) {
          $('.loading').hide();
          console.log(res.responseJSON ? res.responseJSON.error : res.statusText);
          cb(res);
        },
      })
    },
    updated() {
      this.wrongAuth = false;
    }
  },
  mounted() {
    resetParsley(this.$form);
  },
}

</script>