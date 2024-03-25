<template>
  <div class="container-fluid">
    <div class="row">
      <div class="col-lg-12">
        <h2 class="page-header">
          {{ $route.name }}
        </h2>
      </div>
    </div>
    <div class="row form-edit-main">
      <div class="col-lg-6 offset-lg-1">
        <form method="post" class="require" data-parsley-validate autocomplete="off">
          <div class="form-group row error-handler">
            <label class="col-lg-4">Username</label>
            <div class="col-lg-6">
              <input name="username" type="text" class="form-control" required v-model="vData.username" disabled>
            </div>
          </div>
          <div class="form-group row error-handler">
            <label class="col-lg-4">Password</label>
            <div class="col-lg-6">
              <input name="password" type="password" class="form-control" required v-model="vData.password"
                     data-parsley-validation-threshold="0">
            </div>
          </div>
          <div class="form-group row error-handler" v-show="isAdmin">
            <label class="col-lg-4">Name</label>
            <div class="col-lg-6">
              <input name="name" type="text" class="form-control" required v-model="vData.name">
            </div>
          </div>
          <div class="form-group row error-handler" v-show="isAdmin">
            <label class="col-lg-4">Email</label>
            <div class="col-lg-6">
              <input name="email" type="email" class="form-control" data-parsley-type="email" required
                     v-model="vData.email">
            </div>
          </div>
          <div class="form-group row error-handler" v-show="isAdmin">
            <label class="col-lg-4">Description</label>
            <div class="col-lg-6">
              <input name="description" type="text" class="form-control" v-model="vData.description">
            </div>
          </div>
          <div class="form-group row error-handler" v-show="isAdmin">
            <label class="col-lg-4">Role</label>
            <div class="col-lg-6">
              <select name="role" id="role" class="form-control" v-model="vData.role">
                <option v-for="item in roles" :value="item" :key="item.id" :selected="item === vData.role">
                  {{ item }}
                </option>
              </select>
            </div>
          </div>
          <div class="form-group row">
            <div class="btn-group offset-lg-8">
              <router-link :to="{name: 'Users'}" class="btn btn-secondary">Cancel</router-link>
              <button type="submit" class="btn btn-primary once" @click="saveData()">Save</button>
            </div>
          </div>
        </form>
      </div>
    </div>
  </div>
</template>
<script>

import $ from 'jquery';
import {checkLogin, resetParsley} from '@/assets/forms.js'
import {getCache} from "@/assets/cache";

export default {
  data: function () {
    return {
      isAdmin: getCache('session', 'admin') === 'true',
      roles: ["developer", "admin"],
      vData: {
        username: this.$route.query.user,
      }
    }
  },
  props: [
    'title',
    'namesInUse',
    'nameMaxLength',
  ],
  computed: {
    $form() {
      return $(this.$el).find('form');
    },
  },
  methods: {
    getData() {
      let $this = this;
      $('.loading').show();
      let apiUrl = getCache('config', 'apiUrl');
      $.ajax({
        url: apiUrl + '/user?username=' + this.vData.username,
        type: "GET",
        success: function (res) {
          $this.vData = res;
          $('.loading').hide();
        },
        error: function (e) {
          console.log(e);
        }
      })
    },
    saveData() {
      let $this = this;
      $('.loading').show();
      let apiUrl = getCache('config', 'apiUrl');
      $.ajax({
        url: apiUrl + '/user',
        type: "POST",
        data: JSON.stringify(this.vData),
        success: function (e) {
          console.log(e.message);
          $('.loading').hide();
          $this.$router.push('/users');
        },
        error: function (e) {
          console.log(e);
        }
      })
    },
  },
  mounted() {
    checkLogin(this);
    this.getData(() => {
    });
    resetParsley(this.$form);
  },
  updated() {
    resetParsley(this.$form);
  },
};

</script>